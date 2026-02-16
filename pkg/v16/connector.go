package v16

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/v16ai/v16-client/pkg/agent"
	"github.com/v16ai/v16-client/pkg/logger"
)

// V16Connector manages connection to v16.ai platform
type V16Connector struct {
	serverURL string
	token     string
	agentLoop *agent.AgentLoop
	client    *SocketIOClient
	connected bool
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	startTime time.Time
}

// CommandRequest represents a command from v16.ai platform
type CommandMessage struct {
	ID      string                 `json:"id"`
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

// CommandResponse represents the response sent back to v16.ai
type CommandResponse struct {
	ID      string      `json:"id"`
	Payload interface{} `json:"payload"`
}

// StatusInfo represents agent status information
type StatusInfo struct {
	Version   string   `json:"version"`
	Platform  string   `json:"platform"`
	Tools     []string `json:"tools"`
	Connected bool     `json:"connected"`
	Uptime    int64    `json:"uptime"`
}

// NewV16Connector creates a new connector to v16.ai platform
func NewV16Connector(serverURL, token string, agentLoop *agent.AgentLoop) *V16Connector {
	ctx, cancel := context.WithCancel(context.Background())

	client := NewSocketIOClient(serverURL, "/local-agent", token)

	return &V16Connector{
		serverURL: serverURL,
		token:     token,
		agentLoop: agentLoop,
		client:    client,
		connected: false,
		ctx:       ctx,
		cancel:    cancel,
		startTime: time.Now(),
	}
}

// Connect establishes connection to v16.ai platform
func (c *V16Connector) Connect() error {
	c.mu.Lock()
	if c.connected {
		c.mu.Unlock()
		return fmt.Errorf("already connected")
	}
	c.mu.Unlock()

	logger.Info("Connecting to V16 platform...")

	// Setup event handlers before connecting
	c.setupEventHandlers()

	// Connect to Socket.IO server
	if err := c.client.Connect(); err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	// Wait for connection to establish
	time.Sleep(500 * time.Millisecond)

	// Send ready event with capabilities
	if err := c.client.SendReady(); err != nil {
		c.client.Disconnect()
		return fmt.Errorf("failed to send ready event: %v", err)
	}

	c.mu.Lock()
	c.connected = true
	c.mu.Unlock()

	logger.Info(fmt.Sprintf("✅ Connected to V16 platform at %s", c.serverURL))
	logger.Info(fmt.Sprintf("📡 Namespace: /local-agent"))
	logger.Info(fmt.Sprintf("🔧 Tools available: %d", len(c.agentLoop.GetToolNames())))

	return nil
}

// Disconnect closes the connection to v16.ai platform
func (c *V16Connector) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return fmt.Errorf("not connected")
	}

	// Disconnect Socket.IO client
	if err := c.client.Disconnect(); err != nil {
		return err
	}

	// Cancel context
	if c.cancel != nil {
		c.cancel()
	}

	c.connected = false

	logger.Info("Disconnected from v16.ai platform")

	return nil
}

// IsConnected returns whether the connector is currently connected
func (c *V16Connector) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected && c.client.IsConnected()
}

// setupEventHandlers configures Socket.IO event handlers
func (c *V16Connector) setupEventHandlers() {
	// Handle incoming commands
	c.client.On("command", func(data interface{}) {
		c.handleCommand(data)
	})

	// Handle disconnect
	c.client.On("disconnect", func(data interface{}) {
		c.mu.Lock()
		c.connected = false
		c.mu.Unlock()
		logger.Info("Disconnected from V16 platform")
	})

	// Handle errors
	c.client.On("error", func(data interface{}) {
		logger.Info(fmt.Sprintf("Error from platform: %v", data))
	})
}

// handleCommand processes incoming commands from the platform
func (c *V16Connector) handleCommand(data interface{}) {
	// Parse command message
	dataBytes, err := json.Marshal(data)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to marshal command data: %v", err))
		return
	}

	var cmd CommandMessage
	if err := json.Unmarshal(dataBytes, &cmd); err != nil {
		logger.Info(fmt.Sprintf("Failed to unmarshal command: %v", err))
		return
	}

	logger.Info(fmt.Sprintf("📥 Received command: %s (id: %s)", cmd.Type, cmd.ID))

	// Execute command
	go c.executeCommand(cmd)
}

// executeCommand executes a command and sends result back to platform
func (c *V16Connector) executeCommand(cmd CommandMessage) {
	// Map command type to tool
	toolName, toolArgs, err := c.mapCommandToTool(cmd.Type, cmd.Payload)
	if err != nil {
		c.sendResult(cmd.ID, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Get tool from registry
	tool := c.agentLoop.GetTool(toolName)
	if tool == nil {
		c.sendResult(cmd.ID, map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("tool not found: %s", toolName),
		})
		return
	}

	// Execute tool
	ctx := c.ctx
	if timeout, ok := cmd.Payload["timeout"].(float64); ok && timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(c.ctx, time.Duration(timeout)*time.Second)
		defer cancel()
	}

	result := tool.Execute(ctx, toolArgs)

	// Send response
	response := map[string]interface{}{
		"success": !result.IsError,
	}

	if result.IsError {
		response["error"] = result.ForLLM
	} else {
		response["result"] = result.ForLLM
		response["output"] = result.ForUser
	}

	c.sendResult(cmd.ID, response)

	logger.Info(fmt.Sprintf("📤 Sent result for command: %s", cmd.ID))
}

// mapCommandToTool maps backend command types to v16 tools
func (c *V16Connector) mapCommandToTool(cmdType string, payload map[string]interface{}) (string, map[string]interface{}, error) {
	// Direct mapping for commands that match tool actions
	switch cmdType {
	// File operations
	case "read-file":
		return "read_file", payload, nil
	case "write-file":
		return "write_file", payload, nil
	case "list-files":
		return "list_dir", payload, nil
	case "view-file":
		return "read_file", payload, nil

	// Shell execution
	case "execute":
		return "exec", payload, nil

	// Git operations
	case "git-status":
		payload["action"] = "status"
		return "git", payload, nil
	case "git-diff":
		payload["action"] = "diff"
		return "git", payload, nil
	case "git-log":
		payload["action"] = "log"
		return "git", payload, nil
	case "git-commit":
		payload["action"] = "commit"
		return "git", payload, nil

	// Browser operations
	case "browser-navigate":
		payload["action"] = "navigate"
		return "browser", payload, nil
	case "browser-screenshot":
		payload["action"] = "screenshot"
		return "browser", payload, nil
	case "browser-click":
		payload["action"] = "click"
		return "browser", payload, nil
	case "browser-type":
		payload["action"] = "type"
		return "browser", payload, nil
	case "browser-content":
		payload["action"] = "get_html"
		return "browser", payload, nil
	case "browser-evaluate":
		payload["action"] = "execute_js"
		return "browser", payload, nil
	case "browser-elements":
		payload["action"] = "get_html"
		return "browser", payload, nil
	case "browser-click-text":
		payload["action"] = "click"
		return "browser", payload, nil
	case "browser-close":
		payload["action"] = "close"
		return "browser", payload, nil
	case "browser-init":
		payload["action"] = "navigate"
		return "browser", payload, nil
	case "browser-status":
		payload["action"] = "get_html"
		return "browser", payload, nil

	// Desktop operations
	case "screenshot":
		payload["action"] = "screenshot"
		return "desktop", payload, nil
	case "screen-size":
		payload["action"] = "screen_size"
		return "desktop", payload, nil
	case "mouse-position":
		payload["action"] = "mouse_position"
		return "desktop", payload, nil
	case "mouse-move":
		payload["action"] = "mouse_move"
		return "desktop", payload, nil
	case "mouse-click":
		payload["action"] = "mouse_click"
		return "desktop", payload, nil
	case "mouse-double-click":
		payload["action"] = "mouse_click"
		payload["clicks"] = 2.0
		return "desktop", payload, nil
	case "mouse-right-click":
		payload["action"] = "mouse_click"
		payload["button"] = "right"
		return "desktop", payload, nil
	case "mouse-drag":
		payload["action"] = "mouse_move"
		return "desktop", payload, nil
	case "mouse-scroll":
		payload["action"] = "mouse_scroll"
		return "desktop", payload, nil
	case "keyboard-type":
		payload["action"] = "keyboard_type"
		return "desktop", payload, nil
	case "keyboard-press":
		payload["action"] = "keyboard_press"
		return "desktop", payload, nil
	case "keyboard-hotkey":
		payload["action"] = "keyboard_hotkey"
		return "desktop", payload, nil
	case "clipboard-read":
		payload["action"] = "clipboard_read"
		return "desktop", payload, nil
	case "clipboard-write":
		payload["action"] = "clipboard_write"
		return "desktop", payload, nil
	case "window-list":
		payload["action"] = "window_list"
		return "desktop", payload, nil
	case "window-focus":
		payload["action"] = "window_focus"
		return "desktop", payload, nil
	case "window-bounds":
		payload["action"] = "window_bounds"
		return "desktop", payload, nil

	// Status
	case "status":
		return "status", payload, nil

	default:
		return "", nil, fmt.Errorf("unknown command type: %s", cmdType)
	}
}

// sendResult sends command result back to platform
func (c *V16Connector) sendResult(commandID string, payload interface{}) error {
	response := CommandResponse{
		ID:      commandID,
		Payload: payload,
	}

	return c.client.Emit("result", response)
}

// GetUptime returns connection uptime in seconds
func (c *V16Connector) GetUptime() int64 {
	return int64(time.Since(c.startTime).Seconds())
}
