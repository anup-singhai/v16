package v16

import (
	"fmt"
)

// CommandMap maps v16.ai platform commands to internal tool names
var CommandMap = map[string]string{
	// File operations
	"read-file":  "read_file",
	"write-file": "write_file",
	"list-dir":   "list_dir",
	"edit-file":  "edit",
	"append-file": "append_file",

	// Execution
	"execute":       "exec",
	"terminal-open": "terminal",
	"terminal-input": "terminal",
	"terminal-output": "terminal",
	"terminal-close": "terminal",

	// Desktop
	"screenshot":     "desktop",
	"mouse-move":     "desktop",
	"mouse-click":    "desktop",
	"keyboard-type":  "desktop",
	"get-mouse-pos":  "desktop",
	"screen-size":    "desktop",
	"window-list":    "desktop",
	"active-window":  "desktop",

	// Browser
	"browser-navigate":   "browser",
	"browser-screenshot": "browser",
	"browser-click":      "browser",
	"browser-type":       "browser",
	"browser-content":    "browser",
	"browser-text":       "browser",
	"browser-eval":       "browser",
	"browser-wait":       "browser",

	// Web
	"web-search": "web_search",
	"web-fetch":  "web_fetch",

	// Message
	"send-message": "message",

	// Cron
	"cron-add":    "cron",
	"cron-remove": "cron",
	"cron-list":   "cron",

	// Subagent
	"spawn": "spawn",
}

// ActionMap maps v16 commands to tool-specific actions
var ActionMap = map[string]string{
	// Desktop actions
	"screenshot":    "screenshot",
	"mouse-move":    "mouse_move",
	"mouse-click":   "mouse_click",
	"keyboard-type": "keyboard_type",
	"get-mouse-pos": "get_mouse_pos",
	"screen-size":   "screen_size",
	"window-list":   "window_list",
	"active-window": "get_active_window",

	// Browser actions
	"browser-navigate":   "navigate",
	"browser-screenshot": "screenshot",
	"browser-click":      "click",
	"browser-type":       "type",
	"browser-content":    "get_content",
	"browser-text":       "get_text",
	"browser-eval":       "evaluate",
	"browser-wait":       "wait_visible",

	// Terminal actions
	"terminal-open":   "open",
	"terminal-input":  "input",
	"terminal-output": "get_output",
	"terminal-close":  "close",
}

// TranslateCommand translates a v16.ai command to internal tool name and arguments
func TranslateCommand(v16Command string, args map[string]interface{}) (toolName string, toolArgs map[string]interface{}, err error) {
	// Check if command is mapped
	toolName, ok := CommandMap[v16Command]
	if !ok {
		return "", nil, fmt.Errorf("unknown command: %s", v16Command)
	}

	// Create tool arguments
	toolArgs = make(map[string]interface{})

	// Copy all args
	for k, v := range args {
		toolArgs[k] = v
	}

	// Add action parameter if needed
	if action, hasAction := ActionMap[v16Command]; hasAction {
		toolArgs["action"] = action
	}

	// Special handling for specific commands

	// Desktop commands
	switch v16Command {
	case "screenshot":
		// screenshot might need path mapping
		if _, hasPath := toolArgs["path"]; !hasPath {
			toolArgs["path"] = "/tmp/screenshot.png"
		}

	case "mouse-click":
		// Ensure button is set
		if _, hasButton := toolArgs["button"]; !hasButton {
			toolArgs["button"] = "left"
		}
	}

	// Browser commands
	switch v16Command {
	case "browser-screenshot":
		if _, hasPath := toolArgs["path"]; !hasPath {
			toolArgs["path"] = "/tmp/browser-screenshot.png"
		}

	case "browser-navigate":
		// URL is required, but may be in different param names
		if url, ok := args["url"].(string); ok {
			toolArgs["url"] = url
		}
	}

	// Terminal commands
	switch v16Command {
	case "terminal-open":
		if _, hasShell := toolArgs["shell"]; !hasShell {
			toolArgs["shell"] = "/bin/bash"
		}
	}

	// File operations - ensure paths are provided
	switch v16Command {
	case "read-file", "write-file", "edit-file", "append-file":
		if _, hasPath := toolArgs["path"]; !hasPath {
			return "", nil, fmt.Errorf("path parameter required for %s", v16Command)
		}
	}

	return toolName, toolArgs, nil
}

// ReverseTranslateCommand translates internal tool result back to v16 format (for future use)
func ReverseTranslateCommand(toolName string, result interface{}) map[string]interface{} {
	response := make(map[string]interface{})
	response["tool"] = toolName
	response["result"] = result
	return response
}
