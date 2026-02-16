package agent

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/v16ai/v16-client/pkg/bus"
	"github.com/v16ai/v16-client/pkg/config"
	"github.com/v16ai/v16-client/pkg/logger"
	"github.com/v16ai/v16-client/pkg/providers"
	"github.com/v16ai/v16-client/pkg/session"
	"github.com/v16ai/v16-client/pkg/state"
	"github.com/v16ai/v16-client/pkg/tools"
)

// AutonomousAgent represents a persistent agent with scheduled behaviors
type AutonomousAgent struct {
	ID             string
	Name           string
	Personality    string
	Provider       providers.LLMProvider
	Model          string
	MaxTokens      int
	Temperature    float64
	Workspace      string
	Restrict       bool
	Enabled        bool
	TelegramChatID string
	CronJobs       []config.AgentCronJob

	agentLoop   *AgentLoop
	scheduler   *cron.Cron
	running     bool
	mu          sync.Mutex
}

// AutonomousAgentManager manages multiple autonomous agents
type AutonomousAgentManager struct {
	agents   map[string]*AutonomousAgent
	msgBus   *bus.MessageBus
	config   *config.Config
	mu       sync.RWMutex
}

// NewAutonomousAgentManager creates a new manager for autonomous agents
func NewAutonomousAgentManager(cfg *config.Config, msgBus *bus.MessageBus) *AutonomousAgentManager {
	return &AutonomousAgentManager{
		agents: make(map[string]*AutonomousAgent),
		msgBus: msgBus,
		config: cfg,
	}
}

// LoadAgents loads agents from configuration
func (m *AutonomousAgentManager) LoadAgents() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.config.Agents.List == nil || len(m.config.Agents.List) == 0 {
		logger.InfoC("agents", "No autonomous agents configured")
		return nil
	}

	logger.InfoCF("agents", "Loading autonomous agents", map[string]interface{}{
		"count": len(m.config.Agents.List),
	})

	for _, agentCfg := range m.config.Agents.List {
		if !agentCfg.Enabled {
			logger.InfoCF("agents", "Skipping disabled agent", map[string]interface{}{
				"name": agentCfg.Name,
			})
			continue
		}

		// Create provider for this agent
		provider, err := providers.CreateProviderWithOverride(m.config, agentCfg.Provider, agentCfg.Model)
		if err != nil {
			logger.ErrorCF("agents", "Failed to create provider for agent", map[string]interface{}{
				"name":  agentCfg.Name,
				"error": err.Error(),
			})
			continue
		}

		// Determine workspace
		workspace := agentCfg.Workspace
		if workspace == "" {
			workspace = m.config.WorkspacePath()
		}

		// Create autonomous agent
		agent := &AutonomousAgent{
			ID:             agentCfg.ID,
			Name:           agentCfg.Name,
			Personality:    agentCfg.Personality,
			Provider:       provider,
			Model:          agentCfg.Model,
			MaxTokens:      agentCfg.MaxTokens,
			Temperature:    agentCfg.Temperature,
			Workspace:      workspace,
			Restrict:       agentCfg.RestrictToWorkspace,
			Enabled:        agentCfg.Enabled,
			TelegramChatID: agentCfg.TelegramChatID,
			CronJobs:       agentCfg.CronJobs,
		}

		// Create agent loop for this agent
		agentLoop := m.createAgentLoop(agent)
		agent.agentLoop = agentLoop

		m.agents[agent.Name] = agent

		logger.InfoCF("agents", "Loaded autonomous agent", map[string]interface{}{
			"name":      agent.Name,
			"provider":  agentCfg.Provider,
			"model":     agent.Model,
			"cron_jobs": len(agent.CronJobs),
		})
	}

	return nil
}

// createAgentLoop creates an AgentLoop instance for an autonomous agent
func (m *AutonomousAgentManager) createAgentLoop(agent *AutonomousAgent) *AgentLoop {
	// Create tool registry for this agent
	toolsRegistry := createToolRegistry(agent.Workspace, agent.Restrict, m.config, m.msgBus)

	// Create subagent manager with its own tool registry
	subagentManager := tools.NewSubagentManager(agent.Provider, agent.Model, agent.Workspace, m.msgBus)
	subagentTools := createToolRegistry(agent.Workspace, agent.Restrict, m.config, m.msgBus)
	subagentManager.SetTools(subagentTools)

	// Register subagent tool (synchronous execution)
	subagentTool := tools.NewSubagentTool(subagentManager)
	toolsRegistry.Register(subagentTool)

	// Create sessions manager
	sessionsManager := session.NewSessionManager(filepath.Join(agent.Workspace, "sessions"))

	// Create state manager
	stateManager := state.NewManager(agent.Workspace)

	// Create context builder
	contextBuilder := NewContextBuilder(agent.Workspace)
	contextBuilder.SetToolsRegistry(toolsRegistry)

	return &AgentLoop{
		bus:            m.msgBus,
		provider:       agent.Provider,
		workspace:      agent.Workspace,
		model:          agent.Model,
		contextWindow:  100000,
		maxIterations:  m.config.Agents.Defaults.MaxToolIterations,
		sessions:       sessionsManager,
		state:          stateManager,
		contextBuilder: contextBuilder,
		tools:          toolsRegistry,
		summarizing:    sync.Map{},
	}
}

// StartAgents starts all enabled agents and their cron schedulers
func (m *AutonomousAgentManager) StartAgents() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, agent := range m.agents {
		if err := m.startAgent(agent); err != nil {
			logger.ErrorCF("agents", "Failed to start agent", map[string]interface{}{
				"name":  name,
				"error": err.Error(),
			})
			continue
		}
	}

	return nil
}

// startAgent starts a single agent and its cron jobs
func (m *AutonomousAgentManager) startAgent(agent *AutonomousAgent) error {
	agent.mu.Lock()
	defer agent.mu.Unlock()

	if agent.running {
		return fmt.Errorf("agent already running")
	}

	// Create cron scheduler
	agent.scheduler = cron.New(cron.WithSeconds())

	// Schedule all cron jobs for this agent
	for _, job := range agent.CronJobs {
		jobCopy := job // Capture for closure
		_, err := agent.scheduler.AddFunc(jobCopy.Schedule, func() {
			m.executeCronJob(agent, jobCopy)
		})
		if err != nil {
			logger.ErrorCF("agents", "Failed to schedule cron job", map[string]interface{}{
				"agent":    agent.Name,
				"schedule": jobCopy.Schedule,
				"error":    err.Error(),
			})
			continue
		}

		logger.InfoCF("agents", "Scheduled cron job", map[string]interface{}{
			"agent":    agent.Name,
			"schedule": jobCopy.Schedule,
			"task":     jobCopy.Task,
		})
	}

	// Start the scheduler
	agent.scheduler.Start()
	agent.running = true

	logger.InfoCF("agents", "Started autonomous agent", map[string]interface{}{
		"name":      agent.Name,
		"cron_jobs": len(agent.CronJobs),
	})

	return nil
}

// executeCronJob executes a scheduled task for an agent
func (m *AutonomousAgentManager) executeCronJob(agent *AutonomousAgent, job config.AgentCronJob) {
	logger.InfoCF("agents", "Executing cron job", map[string]interface{}{
		"agent": agent.Name,
		"task":  job.Task,
	})

	ctx := context.Background()

	// Build the message with agent's personality and the task
	message := agent.Personality + "\n\n" + "AUTONOMOUS TASK: " + job.Task

	// Execute the task using the agent loop
	// Use configured Telegram chat ID or empty string if not set
	channel := "telegram"
	chatID := agent.TelegramChatID
	if chatID == "" {
		chatID = "0" // Default to prevent errors - will fail at send time
	}

	// Create an inbound message for the agent to process
	inboundMsg := bus.InboundMessage{
		Channel: channel,
		ChatID:  chatID,
		Content: message,
		Media:   []string{},
	}

	response, err := agent.agentLoop.processMessage(ctx, inboundMsg)

	if err != nil {
		logger.ErrorCF("agents", "Cron job execution failed", map[string]interface{}{
			"agent": agent.Name,
			"error": err.Error(),
		})
		// Send error notification to Telegram
		m.msgBus.PublishOutbound(bus.OutboundMessage{
			Channel: channel,
			ChatID:  chatID,
			Content: fmt.Sprintf("⚠️ Agent @%s cron job failed: %v", agent.Name, err),
		})
		return
	}

	logger.InfoCF("agents", "Cron job completed", map[string]interface{}{
		"agent":    agent.Name,
		"response": response,
	})
}

// StopAgents stops all running agents
func (m *AutonomousAgentManager) StopAgents() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, agent := range m.agents {
		m.stopAgent(agent)
	}
}

// stopAgent stops a single agent
func (m *AutonomousAgentManager) stopAgent(agent *AutonomousAgent) {
	agent.mu.Lock()
	defer agent.mu.Unlock()

	if !agent.running {
		return
	}

	if agent.scheduler != nil {
		agent.scheduler.Stop()
	}

	agent.running = false

	logger.InfoCF("agents", "Stopped autonomous agent", map[string]interface{}{
		"name": agent.Name,
	})
}

// GetAgent returns an agent by name
func (m *AutonomousAgentManager) GetAgent(name string) (*AutonomousAgent, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	agent, ok := m.agents[name]
	return agent, ok
}

// GetAgentCount returns the number of loaded agents
func (m *AutonomousAgentManager) GetAgentCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.agents)
}
