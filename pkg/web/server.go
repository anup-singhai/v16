package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/v16ai/v16-client/pkg/config"
	"github.com/v16ai/v16-client/pkg/logger"
)

//go:embed static/*
var staticFiles embed.FS

// Server provides web UI for configuration
type Server struct {
	config     *config.Config
	configPath string
	addr       string
}

// NewServer creates a new web server
func NewServer(cfg *config.Config, configPath string, addr string) *Server {
	return &Server{
		config:     cfg,
		configPath: configPath,
		addr:       addr,
	}
}

// Start starts the web server
func (s *Server) Start() error {
	// Serve static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("failed to get static files: %w", err)
	}

	http.Handle("/", http.FileServer(http.FS(staticFS)))

	// API endpoints
	http.HandleFunc("/api/config", s.handleConfig)
	http.HandleFunc("/api/config/save", s.handleConfigSave)
	http.HandleFunc("/api/status", s.handleStatus)
	http.HandleFunc("/api/browse", s.handleBrowse)

	logger.InfoCF("web", "Starting web UI", map[string]interface{}{
		"addr": s.addr,
	})

	fmt.Printf("\n🌐 Web UI available at: http://%s\n", s.addr)

	return http.ListenAndServe(s.addr, nil)
}

// handleConfig returns the current configuration
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.config)
}

// handleConfigSave saves the configuration
func (s *Server) handleConfigSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newConfig config.Config
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Save to file
	if err := config.SaveConfig(s.configPath, &newConfig); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save config: %v", err), http.StatusInternalServerError)
		return
	}

	// Update in-memory config
	s.config = &newConfig

	logger.InfoC("web", "Configuration saved")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Configuration saved successfully",
	})
}

// handleStatus returns the agent status
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	workspace := s.config.WorkspacePath()

	status := map[string]interface{}{
		"config_path": s.configPath,
		"workspace":   workspace,
		"provider":    s.config.Agents.Defaults.Provider,
		"model":       s.config.Agents.Defaults.Model,
		"channels": map[string]bool{
			"telegram": s.config.Channels.Telegram.Enabled,
			"discord":  s.config.Channels.Discord.Enabled,
			"slack":    s.config.Channels.Slack.Enabled,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// handleBrowse lists directories for browsing
func (s *Server) handleBrowse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			http.Error(w, "Failed to get home directory", http.StatusInternalServerError)
			return
		}
		path = home
	}

	// Expand ~ to home directory
	if path == "~" || filepath.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			http.Error(w, "Failed to get home directory", http.StatusInternalServerError)
			return
		}
		if path == "~" {
			path = home
		} else {
			path = filepath.Join(home, path[2:])
		}
	}

	// Clean the path
	path = filepath.Clean(path)

	// Read directory contents
	entries, err := os.ReadDir(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read directory: %v", err), http.StatusBadRequest)
		return
	}

	// Filter to only directories
	var dirs []map[string]interface{}
	for _, entry := range entries {
		if entry.IsDir() {
			// Skip hidden directories (starting with .)
			if entry.Name()[0] == '.' {
				continue
			}

			dirs = append(dirs, map[string]interface{}{
				"name": entry.Name(),
				"path": filepath.Join(path, entry.Name()),
			})
		}
	}

	// Get parent directory
	parent := filepath.Dir(path)
	if parent == path {
		parent = "" // At root
	}

	response := map[string]interface{}{
		"current": path,
		"parent":  parent,
		"dirs":    dirs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetConfigPath returns the config file path
func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".v16", "config.json")
}
