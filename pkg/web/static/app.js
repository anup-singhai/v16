// Global state
let currentConfig = {};
let agents = [];
let cronJobs = [];

// Tab switching
document.querySelectorAll('.tab').forEach(tab => {
    tab.addEventListener('click', () => {
        const tabName = tab.dataset.tab;
        document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
        tab.classList.add('active');
        document.querySelectorAll('.tab-content').forEach(content => {
            content.classList.remove('active');
        });
        document.getElementById(`${tabName}-tab`).classList.add('active');
    });
});

// Load configuration
async function loadConfig() {
    try {
        const response = await fetch('/api/config');
        if (!response.ok) throw new Error('Failed to load config');

        currentConfig = await response.json();

        // Extract agents from config (current single agent becomes default)
        agents = [];

        // Check if multi-agent config exists
        if (currentConfig.agents && currentConfig.agents.list) {
            agents = currentConfig.agents.list;
        } else {
            // Migrate from single agent to multi-agent
            agents = [{
                id: 'default',
                name: 'assistant',
                personality: 'You are a helpful AI assistant.',
                provider: currentConfig.agents?.defaults?.provider || 'moonshot',
                model: currentConfig.agents?.defaults?.model || 'moonshot-v1-128k',
                max_tokens: currentConfig.agents?.defaults?.max_tokens || 4096,
                temperature: currentConfig.agents?.defaults?.temperature || 0.7,
                enabled: true,
                cron_jobs: []
            }];
        }

        populateForm(currentConfig);
        renderAgents();
        showMessage('Configuration loaded', 'success');
    } catch (error) {
        showMessage('Error loading configuration: ' + error.message, 'error');
    }
}

// Load status
async function loadStatus() {
    try {
        const response = await fetch('/api/status');
        if (!response.ok) throw new Error('Failed to load status');

        const status = await response.json();
        document.getElementById('status-config-path').textContent = status.config_path || '-';
        document.getElementById('status-workspace').textContent = status.workspace || '-';

        showMessage('Status refreshed', 'success');
    } catch (error) {
        showMessage('Error loading status: ' + error.message, 'error');
    }
}

// Render agents grid
function renderAgents() {
    const grid = document.getElementById('agents-grid');
    grid.innerHTML = '';

    if (agents.length === 0) {
        grid.innerHTML = `
            <div class="empty-state" style="grid-column: 1 / -1;">
                <div class="empty-state-icon">🤖</div>
                <div class="empty-state-text">No Agents Yet</div>
                <div class="empty-state-subtext">Create your first autonomous AI agent to get started</div>
                <button onclick="openAgentModal()">Create First Agent</button>
            </div>
        `;
        return;
    }

    // Render agent cards
    agents.forEach(agent => {
        const card = document.createElement('div');
        card.className = 'agent-card';

        const statusClass = agent.enabled ? 'active' : 'inactive';
        const statusText = agent.enabled ? 'Active' : 'Inactive';

        const cronCount = agent.cron_jobs ? agent.cron_jobs.length : 0;
        const cronText = cronCount === 1 ? '1 scheduled task' : `${cronCount} scheduled tasks`;

        const workspacePath = agent.workspace || '(default)';
        const workspaceRestricted = agent.restrict_to_workspace ? ' 🔒' : '';

        card.innerHTML = `
            <div class="agent-header">
                <div class="agent-name">@${agent.name}</div>
                <div class="agent-status ${statusClass}">${statusText}</div>
            </div>
            <div class="agent-description">${agent.personality || 'No description'}</div>
            <div class="agent-details">
                <div class="agent-detail">
                    <span class="detail-label">Provider</span>
                    <span class="detail-value">${agent.provider}</span>
                </div>
                <div class="agent-detail">
                    <span class="detail-label">Model</span>
                    <span class="detail-value">${agent.model}</span>
                </div>
                <div class="agent-detail">
                    <span class="detail-label">Workspace</span>
                    <span class="detail-value">${workspacePath}${workspaceRestricted}</span>
                </div>
                <div class="agent-detail">
                    <span class="detail-label">Autonomous Tasks</span>
                    <span class="detail-value">${cronText}</span>
                </div>
            </div>
            <div class="agent-actions">
                <button class="small secondary" onclick="editAgent('${agent.id}')">Edit</button>
                <button class="small secondary" onclick="toggleAgent('${agent.id}')">${agent.enabled ? 'Disable' : 'Enable'}</button>
            </div>
        `;

        grid.appendChild(card);
    });

    // Add "Create Agent" card
    const addCard = document.createElement('div');
    addCard.className = 'add-agent-card';
    addCard.onclick = () => openAgentModal();
    addCard.innerHTML = `
        <div class="add-icon">+</div>
        <div class="add-text">Create New Agent</div>
    `;
    grid.appendChild(addCard);
}

// Open agent modal (create or edit)
function openAgentModal(agentId = null) {
    const modal = document.getElementById('agent-modal');
    const title = document.getElementById('modal-title');
    const deleteBtn = document.getElementById('delete-agent-btn');

    // Clear form
    document.getElementById('agent-id').value = '';
    document.getElementById('agent-name').value = '';
    document.getElementById('agent-personality').value = '';
    document.getElementById('agent-provider').value = 'anthropic';
    document.getElementById('agent-model').value = '';
    document.getElementById('agent-max-tokens').value = '4096';
    document.getElementById('agent-temperature').value = '0.7';
    document.getElementById('agent-workspace').value = '';
    document.getElementById('agent-restrict-workspace').checked = false;
    document.getElementById('agent-enabled').checked = true;
    cronJobs = [];
    renderCronJobs();

    if (agentId) {
        // Edit mode
        const agent = agents.find(a => a.id === agentId);
        if (agent) {
            title.textContent = 'Edit Agent: @' + agent.name;
            deleteBtn.style.display = 'inline-block';

            document.getElementById('agent-id').value = agent.id;
            document.getElementById('agent-name').value = agent.name;
            document.getElementById('agent-personality').value = agent.personality || '';
            document.getElementById('agent-provider').value = agent.provider;
            document.getElementById('agent-model').value = agent.model;
            document.getElementById('agent-max-tokens').value = agent.max_tokens || 4096;
            document.getElementById('agent-temperature').value = agent.temperature || 0.7;
            document.getElementById('agent-workspace').value = agent.workspace || '';
            document.getElementById('agent-restrict-workspace').checked = agent.restrict_to_workspace || false;
            document.getElementById('agent-enabled').checked = agent.enabled !== false;
            cronJobs = agent.cron_jobs || [];
            renderCronJobs();
        }
    } else {
        // Create mode
        title.textContent = 'Create New Agent';
        deleteBtn.style.display = 'none';
    }

    modal.classList.add('active');
}

function closeAgentModal() {
    document.getElementById('agent-modal').classList.remove('active');
}

function editAgent(agentId) {
    openAgentModal(agentId);
}

function toggleAgent(agentId) {
    const agent = agents.find(a => a.id === agentId);
    if (agent) {
        agent.enabled = !agent.enabled;
        renderAgents();
        showMessage(`Agent @${agent.name} ${agent.enabled ? 'enabled' : 'disabled'}`, 'success');
    }
}

function deleteAgent() {
    const agentId = document.getElementById('agent-id').value;
    if (!agentId) return;

    if (!confirm('Are you sure you want to delete this agent?')) return;

    agents = agents.filter(a => a.id !== agentId);
    closeAgentModal();
    renderAgents();
    showMessage('Agent deleted', 'success');
}

function saveAgent() {
    const agentId = document.getElementById('agent-id').value;
    const name = document.getElementById('agent-name').value.trim();
    const personality = document.getElementById('agent-personality').value.trim();
    const provider = document.getElementById('agent-provider').value;
    const model = document.getElementById('agent-model').value.trim();
    const maxTokens = parseInt(document.getElementById('agent-max-tokens').value) || 4096;
    const temperature = parseFloat(document.getElementById('agent-temperature').value) || 0.7;
    const workspace = document.getElementById('agent-workspace').value.trim();
    const restrictWorkspace = document.getElementById('agent-restrict-workspace').checked;
    const enabled = document.getElementById('agent-enabled').checked;

    if (!name) {
        showMessage('Agent name is required', 'error');
        return;
    }

    if (!model) {
        showMessage('Model is required', 'error');
        return;
    }

    const agent = {
        id: agentId || generateId(),
        name: name,
        personality: personality,
        provider: provider,
        model: model,
        max_tokens: maxTokens,
        temperature: temperature,
        workspace: workspace,
        restrict_to_workspace: restrictWorkspace,
        enabled: enabled,
        cron_jobs: cronJobs
    };

    if (agentId) {
        // Update existing
        const index = agents.findIndex(a => a.id === agentId);
        if (index !== -1) {
            agents[index] = agent;
        }
    } else {
        // Create new
        agents.push(agent);
    }

    closeAgentModal();
    renderAgents();
    showMessage(`Agent @${agent.name} saved`, 'success');
}

// Cron job management
function renderCronJobs() {
    const list = document.getElementById('cron-list');
    list.innerHTML = '';

    cronJobs.forEach((job, index) => {
        const item = document.createElement('div');
        item.className = 'cron-item';
        item.innerHTML = `
            <div class="cron-schedule">${job.schedule}</div>
            <div class="cron-task">${job.task}</div>
            <button class="small danger" onclick="removeCronJob(${index})">×</button>
        `;
        list.appendChild(item);
    });
}

function addCronJob() {
    const schedule = prompt('Enter cron schedule (e.g., "0 9 * * *" for 9am daily, or "*/30 * * * *" for every 30 min):');
    if (!schedule) return;

    const task = prompt('Enter task description:');
    if (!task) return;

    cronJobs.push({ schedule, task });
    renderCronJobs();
}

function removeCronJob(index) {
    cronJobs.splice(index, 1);
    renderCronJobs();
}

// Populate form with config data
function populateForm(config) {
    // Channels - Telegram
    document.getElementById('telegram-enabled').checked = config.channels?.telegram?.enabled || false;
    document.getElementById('telegram-token').value = config.channels?.telegram?.token || '';

    // Channels - Discord
    document.getElementById('discord-enabled').checked = config.channels?.discord?.enabled || false;
    document.getElementById('discord-token').value = config.channels?.discord?.token || '';

    // Channels - Slack
    document.getElementById('slack-enabled').checked = config.channels?.slack?.enabled || false;
    document.getElementById('slack-token').value = config.channels?.slack?.bot_token || '';

    // Providers
    document.getElementById('provider-anthropic').value = config.providers?.anthropic?.api_key || '';
    document.getElementById('provider-openai').value = config.providers?.openai?.api_key || '';
    document.getElementById('provider-openrouter').value = config.providers?.openrouter?.api_key || '';
    document.getElementById('provider-moonshot').value = config.providers?.moonshot?.api_key || '';
    document.getElementById('provider-deepseek').value = config.providers?.deepseek?.api_key || '';
    document.getElementById('provider-groq').value = config.providers?.groq?.api_key || '';
    document.getElementById('provider-gemini').value = config.providers?.gemini?.api_key || '';
    document.getElementById('provider-vllm').value = config.providers?.vllm?.api_base || '';

    // Settings
    document.getElementById('default-workspace').value = config.agents?.defaults?.workspace || '';
    document.getElementById('default-max-iterations').value = config.agents?.defaults?.max_tool_iterations || 50;
    document.getElementById('default-restrict-workspace').checked = config.agents?.defaults?.restrict_to_workspace || false;
}

// Collect form data and save
async function saveGlobalConfig() {
    try {
        const config = JSON.parse(JSON.stringify(currentConfig)); // Deep clone

        // Agents - save multi-agent list
        config.agents = config.agents || {};
        config.agents.list = agents;

        // Keep defaults for backward compatibility
        config.agents.defaults = {
            workspace: document.getElementById('default-workspace').value || currentConfig.agents?.defaults?.workspace,
            restrict_to_workspace: document.getElementById('default-restrict-workspace').checked,
            provider: agents.length > 0 ? agents[0].provider : 'moonshot',
            model: agents.length > 0 ? agents[0].model : 'moonshot-v1-128k',
            max_tokens: currentConfig.agents?.defaults?.max_tokens || 4096,
            temperature: currentConfig.agents?.defaults?.temperature || 0.7,
            max_tool_iterations: parseInt(document.getElementById('default-max-iterations').value) || 50
        };

        // Channels
        config.channels = config.channels || {};

        config.channels.telegram = config.channels.telegram || {};
        config.channels.telegram.enabled = document.getElementById('telegram-enabled').checked;
        config.channels.telegram.token = document.getElementById('telegram-token').value;

        config.channels.discord = config.channels.discord || {};
        config.channels.discord.enabled = document.getElementById('discord-enabled').checked;
        config.channels.discord.token = document.getElementById('discord-token').value;

        config.channels.slack = config.channels.slack || {};
        config.channels.slack.enabled = document.getElementById('slack-enabled').checked;
        config.channels.slack.bot_token = document.getElementById('slack-token').value;

        // Providers
        config.providers = config.providers || {};

        config.providers.anthropic = config.providers.anthropic || {};
        config.providers.anthropic.api_key = document.getElementById('provider-anthropic').value;

        config.providers.openai = config.providers.openai || {};
        config.providers.openai.api_key = document.getElementById('provider-openai').value;

        config.providers.openrouter = config.providers.openrouter || {};
        config.providers.openrouter.api_key = document.getElementById('provider-openrouter').value;

        config.providers.moonshot = config.providers.moonshot || {};
        config.providers.moonshot.api_key = document.getElementById('provider-moonshot').value;

        config.providers.deepseek = config.providers.deepseek || {};
        config.providers.deepseek.api_key = document.getElementById('provider-deepseek').value;

        config.providers.groq = config.providers.groq || {};
        config.providers.groq.api_key = document.getElementById('provider-groq').value;

        config.providers.gemini = config.providers.gemini || {};
        config.providers.gemini.api_key = document.getElementById('provider-gemini').value;

        config.providers.vllm = config.providers.vllm || {};
        config.providers.vllm.api_base = document.getElementById('provider-vllm').value;

        const response = await fetch('/api/config/save', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(config)
        });

        if (!response.ok) {
            const error = await response.text();
            throw new Error(error);
        }

        const result = await response.json();
        showMessage(result.message || 'Configuration saved! Restart v16 gateway to apply changes.', 'success');

        // Reload config
        await loadConfig();
        await loadStatus();
    } catch (error) {
        showMessage('Error saving configuration: ' + error.message, 'error');
    }
}

// Show message
function showMessage(text, type) {
    const messageEl = document.getElementById('message');
    messageEl.textContent = text;
    messageEl.className = `message ${type}`;
    messageEl.style.display = 'block';

    setTimeout(() => {
        messageEl.style.display = 'none';
    }, 5000);
}

// Generate unique ID
function generateId() {
    return 'agent-' + Date.now() + '-' + Math.random().toString(36).substr(2, 9);
}

// Load initial data
window.addEventListener('DOMContentLoaded', () => {
    loadConfig();
    loadStatus();
});
