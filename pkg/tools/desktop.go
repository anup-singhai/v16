package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-vgo/robotgo"
)

type DesktopTool struct {
	enabled bool
}

func NewDesktopTool() *DesktopTool {
	return &DesktopTool{
		enabled: true,
	}
}

func (t *DesktopTool) Name() string {
	return "desktop"
}

func (t *DesktopTool) Description() string {
	return "Control desktop: capture screen, move mouse, type, click"
}

func (t *DesktopTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type":     "object",
		"required": []string{"action"},
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type": "string",
				"enum": []string{
					"screenshot",
					"mouse_move",
					"mouse_click",
					"keyboard_type",
					"get_mouse_pos",
					"screen_size",
					"window_list",
					"get_active_window",
				},
				"description": "Desktop action to perform",
			},
			"x": map[string]interface{}{
				"type":        "number",
				"description": "X coordinate (for mouse actions)",
			},
			"y": map[string]interface{}{
				"type":        "number",
				"description": "Y coordinate (for mouse actions)",
			},
			"text": map[string]interface{}{
				"type":        "string",
				"description": "Text to type (for keyboard_type)",
			},
			"button": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"left", "right", "middle"},
				"description": "Mouse button (for mouse_click)",
				"default":     "left",
			},
			"path": map[string]interface{}{
				"type":        "string",
				"description": "File path to save screenshot (optional, defaults to /tmp/screenshot.png)",
			},
		},
	}
}

func (t *DesktopTool) Execute(ctx context.Context, args map[string]interface{}) *ToolResult {
	if !t.enabled {
		return ErrorResult("desktop tool is disabled")
	}

	action, ok := args["action"].(string)
	if !ok {
		return ErrorResult("action parameter is required")
	}

	switch action {
	case "screenshot":
		return t.screenshot(args)
	case "mouse_move":
		return t.mouseMove(args)
	case "mouse_click":
		return t.mouseClick(args)
	case "keyboard_type":
		return t.keyboardType(args)
	case "get_mouse_pos":
		return t.getMousePos()
	case "screen_size":
		return t.getScreenSize()
	case "window_list":
		return t.getWindowList()
	case "get_active_window":
		return t.getActiveWindow()
	default:
		return ErrorResult(fmt.Sprintf("unknown action: %s", action))
	}
}

func (t *DesktopTool) screenshot(args map[string]interface{}) *ToolResult {
	// Get path or use default
	path, ok := args["path"].(string)
	if !ok || path == "" {
		path = "/tmp/screenshot.png"
	}

	// Capture screen and save directly to file
	err := robotgo.SaveCapture(path)
	if err != nil {
		return ErrorResult(fmt.Sprintf("failed to capture screen: %v", err))
	}

	return NewToolResult(fmt.Sprintf("Screenshot saved to %s", path))
}

func (t *DesktopTool) mouseMove(args map[string]interface{}) *ToolResult {
	x, xok := args["x"].(float64)
	y, yok := args["y"].(float64)

	if !xok || !yok {
		return ErrorResult("x and y coordinates required")
	}

	robotgo.Move(int(x), int(y))
	return NewToolResult(fmt.Sprintf("Moved mouse to (%d, %d)", int(x), int(y)))
}

func (t *DesktopTool) mouseClick(args map[string]interface{}) *ToolResult {
	button, ok := args["button"].(string)
	if !ok {
		button = "left"
	}

	robotgo.Click(button)
	return NewToolResult(fmt.Sprintf("Clicked %s mouse button", button))
}

func (t *DesktopTool) keyboardType(args map[string]interface{}) *ToolResult {
	text, ok := args["text"].(string)
	if !ok {
		return ErrorResult("text parameter required")
	}

	robotgo.TypeStr(text)
	return NewToolResult(fmt.Sprintf("Typed: %s", text))
}

func (t *DesktopTool) getMousePos() *ToolResult {
	x, y := robotgo.Location()
	return NewToolResult(fmt.Sprintf("Mouse position: (%d, %d)", x, y))
}

func (t *DesktopTool) getScreenSize() *ToolResult {
	width, height := robotgo.GetScreenSize()
	return NewToolResult(fmt.Sprintf("Screen size: %dx%d", width, height))
}

func (t *DesktopTool) getWindowList() *ToolResult {
	pids, err := robotgo.Pids()
	if err != nil {
		return ErrorResult(fmt.Sprintf("failed to get process list: %v", err))
	}

	var windows []string
	for _, pid := range pids {
		title := robotgo.GetTitle(pid)
		if title != "" {
			windows = append(windows, fmt.Sprintf("PID %d: %s", pid, title))
		}
	}

	if len(windows) == 0 {
		return NewToolResult("No windows found")
	}

	return NewToolResult(fmt.Sprintf("Windows:\n%s", strings.Join(windows, "\n")))
}

func (t *DesktopTool) getActiveWindow() *ToolResult {
	pid := robotgo.GetPid()
	title := robotgo.GetTitle()
	return NewToolResult(fmt.Sprintf("Active window: %s (PID: %d)", title, pid))
}
