package v16

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// SocketIOClient handles Socket.IO protocol over WebSocket
type SocketIOClient struct {
	url       string
	namespace string
	token     string
	conn      *websocket.Conn
	mu        sync.RWMutex
	connected bool
	handlers  map[string]func(interface{})
	stopCh    chan struct{}
}

// Socket.IO message types
const (
	SocketIOConnect    = "0"
	SocketIODisconnect = "1"
	SocketIOEvent      = "2"
	SocketIOAck        = "3"
	SocketIOError      = "4"
)

// NewSocketIOClient creates a new Socket.IO client
func NewSocketIOClient(serverURL, namespace, token string) *SocketIOClient {
	return &SocketIOClient{
		url:       serverURL,
		namespace: namespace,
		token:     token,
		handlers:  make(map[string]func(interface{})),
		stopCh:    make(chan struct{}),
	}
}

// Connect establishes WebSocket connection with Socket.IO protocol
func (c *SocketIOClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return fmt.Errorf("already connected")
	}

	// Parse server URL
	u, err := url.Parse(c.url)
	if err != nil {
		return fmt.Errorf("invalid server URL: %v", err)
	}

	// Convert http(s) to ws(s)
	if u.Scheme == "http" {
		u.Scheme = "ws"
	} else if u.Scheme == "https" {
		u.Scheme = "wss"
	}

	// Add Socket.IO path and namespace
	u.Path = "/socket.io/" + strings.TrimPrefix(c.namespace, "/")

	// Add query parameters for Socket.IO
	query := u.Query()
	query.Set("transport", "websocket")
	query.Set("EIO", "4") // Engine.IO protocol version 4
	u.RawQuery = query.Encode()

	// Set up WebSocket headers with auth token
	headers := make(map[string][]string)
	headers["Authorization"] = []string{c.token}

	// Connect to WebSocket
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.Dial(u.String(), headers)
	if err != nil {
		return fmt.Errorf("websocket connection failed: %v", err)
	}

	c.conn = conn
	c.connected = true

	// Send Socket.IO connect packet
	connectPacket := fmt.Sprintf("40%s", c.namespace)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(connectPacket)); err != nil {
		conn.Close()
		c.connected = false
		return fmt.Errorf("failed to send connect packet: %v", err)
	}

	// Start message listener
	go c.listen()

	// Start ping/pong handler
	go c.pingHandler()

	return nil
}

// Disconnect closes the WebSocket connection
func (c *SocketIOClient) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return nil
	}

	close(c.stopCh)

	if c.conn != nil {
		// Send disconnect packet
		disconnectPacket := fmt.Sprintf("41%s", c.namespace)
		c.conn.WriteMessage(websocket.TextMessage, []byte(disconnectPacket))
		c.conn.Close()
		c.conn = nil
	}

	c.connected = false
	return nil
}

// On registers an event handler
func (c *SocketIOClient) On(event string, handler func(interface{})) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[event] = handler
}

// Emit sends an event to the server
func (c *SocketIOClient) Emit(event string, data interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected || c.conn == nil {
		return fmt.Errorf("not connected")
	}

	// Create Socket.IO event packet: 42["event",data]
	eventData, err := json.Marshal([]interface{}{event, data})
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %v", err)
	}

	packet := fmt.Sprintf("42%s%s", c.namespace, string(eventData))

	return c.conn.WriteMessage(websocket.TextMessage, []byte(packet))
}

// listen continuously reads messages from WebSocket
func (c *SocketIOClient) listen() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Socket.IO listener panic: %v\n", r)
		}
	}()

	for {
		select {
		case <-c.stopCh:
			return
		default:
			c.mu.RLock()
			conn := c.conn
			c.mu.RUnlock()

			if conn == nil {
				return
			}

			messageType, message, err := conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					fmt.Printf("WebSocket read error: %v\n", err)
				}
				c.handleDisconnect()
				return
			}

			if messageType == websocket.TextMessage {
				c.handleMessage(string(message))
			}
		}
	}
}

// handleMessage processes incoming Socket.IO messages
func (c *SocketIOClient) handleMessage(message string) {
	if len(message) == 0 {
		return
	}

	// Socket.IO packet format: <type>[namespace][data]
	msgType := string(message[0])
	payload := message[1:]

	switch msgType {
	case "0": // Connect
		fmt.Println("[SocketIO] Connected to server")

	case "1": // Disconnect
		fmt.Println("[SocketIO] Disconnected from server")
		c.handleDisconnect()

	case "2": // Event (42 for namespace events)
		if strings.HasPrefix(message, "42") {
			c.handleEvent(payload)
		}

	case "3": // Ack
		// Handle ack if needed

	case "4": // Error
		fmt.Printf("[SocketIO] Error: %s\n", payload)
	}
}

// handleEvent processes Socket.IO events
func (c *SocketIOClient) handleEvent(payload string) {
	// Remove namespace prefix if present
	payload = strings.TrimPrefix(payload, c.namespace)

	// Parse event data: ["event", data]
	var eventData []json.RawMessage
	if err := json.Unmarshal([]byte(payload), &eventData); err != nil {
		fmt.Printf("[SocketIO] Failed to parse event: %v\n", err)
		return
	}

	if len(eventData) < 1 {
		return
	}

	// Extract event name
	var eventName string
	if err := json.Unmarshal(eventData[0], &eventName); err != nil {
		fmt.Printf("[SocketIO] Failed to parse event name: %v\n", err)
		return
	}

	// Extract event data
	var data interface{}
	if len(eventData) > 1 {
		if err := json.Unmarshal(eventData[1], &data); err != nil {
			fmt.Printf("[SocketIO] Failed to parse event data: %v\n", err)
			return
		}
	}

	// Call handler
	c.mu.RLock()
	handler, exists := c.handlers[eventName]
	c.mu.RUnlock()

	if exists && handler != nil {
		go handler(data)
	}
}

// handleDisconnect handles connection loss
func (c *SocketIOClient) handleDisconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	c.connected = false

	// Call disconnect handler if registered
	if handler, exists := c.handlers["disconnect"]; exists {
		go handler(nil)
	}
}

// pingHandler sends periodic pings to keep connection alive
func (c *SocketIOClient) pingHandler() {
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopCh:
			return
		case <-ticker.C:
			c.mu.RLock()
			conn := c.conn
			connected := c.connected
			c.mu.RUnlock()

			if !connected || conn == nil {
				return
			}

			// Send ping packet (Engine.IO ping)
			if err := conn.WriteMessage(websocket.TextMessage, []byte("2")); err != nil {
				fmt.Printf("[SocketIO] Ping failed: %v\n", err)
				c.handleDisconnect()
				return
			}
		}
	}
}

// IsConnected returns connection status
func (c *SocketIOClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// SendReady sends ready event with client capabilities
func (c *SocketIOClient) SendReady() error {
	hostname, _ := os.Hostname()
	cwd, _ := os.Getwd()

	capabilities := []string{
		"file-ops",
		"shell-exec",
		"browser",
		"desktop",
		"terminal",
		"git",
		"grep",
		"glob",
	}

	payload := map[string]interface{}{
		"platform":     runtime.GOOS,
		"arch":         runtime.GOARCH,
		"hostname":     hostname,
		"cwd":          cwd,
		"version":      "v0.1.0",
		"capabilities": capabilities,
		"browserConfig": map[string]interface{}{
			"mode":                 "persistent",
			"chromeDebugAvailable": false,
			"profileDir":           "~/.v16/browser",
		},
	}

	return c.Emit("ready", payload)
}
