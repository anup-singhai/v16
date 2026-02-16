// Tab switching
document.querySelectorAll('.tab').forEach(tab => {
    tab.addEventListener('click', () => {
        const tabName = tab.dataset.tab;

        // Update active tab
        document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
        tab.classList.add('active');

        // Update active content
        document.querySelectorAll('.tab-content').forEach(content => {
            content.classList.remove('active');
        });
        document.getElementById(`${tabName}-tab`).classList.add('active');
    });
});

let currentConfig = {};

// Load configuration
async function loadConfig() {
    try {
        const response = await fetch('/api/config');
        if (!response.ok) throw new Error('Failed to load config');

        currentConfig = await response.json();
        populateForm(currentConfig);
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
        document.getElementById('status-provider').textContent = status.provider || '-';
        document.getElementById('status-model').textContent = status.model || '-';

        // Update channel status
        updateChannelStatus('telegram', status.channels?.telegram);
        updateChannelStatus('discord', status.channels?.discord);
        updateChannelStatus('slack', status.channels?.slack);

        showMessage('Status refreshed', 'success');
    } catch (error) {
        showMessage('Error loading status: ' + error.message, 'error');
    }
}

function updateChannelStatus(channel, enabled) {
    const element = document.getElementById(`status-${channel}`);
    const indicator = element.querySelector('.status-indicator');

    if (enabled) {
        indicator.className = 'status-indicator enabled';
        element.innerHTML = `<span class="status-indicator enabled"></span> Enabled`;
    } else {
        indicator.className = 'status-indicator disabled';
        element.innerHTML = `<span class="status-indicator disabled"></span> Disabled`;
    }
}

// Populate form with config data
function populateForm(config) {
    // Agents
    document.getElementById('agent-workspace').value = config.agents?.defaults?.workspace || '';
    document.getElementById('agent-provider').value = config.agents?.defaults?.provider || 'openrouter';
    document.getElementById('agent-model').value = config.agents?.defaults?.model || '';
    document.getElementById('agent-max-tokens').value = config.agents?.defaults?.max_tokens || 4096;
    document.getElementById('agent-temperature').value = config.agents?.defaults?.temperature || 0.7;
    document.getElementById('agent-max-iterations').value = config.agents?.defaults?.max_tool_iterations || 50;
    document.getElementById('agent-restrict').checked = config.agents?.defaults?.restrict_to_workspace || false;

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
    document.getElementById('provider-openrouter').value = config.providers?.openrouter?.api_key || '';
    document.getElementById('provider-anthropic').value = config.providers?.anthropic?.api_key || '';
    document.getElementById('provider-openai').value = config.providers?.openai?.api_key || '';
    document.getElementById('provider-gemini').value = config.providers?.gemini?.api_key || '';
    document.getElementById('provider-groq').value = config.providers?.groq?.api_key || '';
}

// Collect form data
function collectFormData() {
    const config = JSON.parse(JSON.stringify(currentConfig)); // Deep clone

    // Agents
    config.agents = config.agents || {};
    config.agents.defaults = {
        workspace: document.getElementById('agent-workspace').value,
        restrict_to_workspace: document.getElementById('agent-restrict').checked,
        provider: document.getElementById('agent-provider').value,
        model: document.getElementById('agent-model').value,
        max_tokens: parseInt(document.getElementById('agent-max-tokens').value) || 4096,
        temperature: parseFloat(document.getElementById('agent-temperature').value) || 0.7,
        max_tool_iterations: parseInt(document.getElementById('agent-max-iterations').value) || 50
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

    config.providers.openrouter = config.providers.openrouter || {};
    config.providers.openrouter.api_key = document.getElementById('provider-openrouter').value;

    config.providers.anthropic = config.providers.anthropic || {};
    config.providers.anthropic.api_key = document.getElementById('provider-anthropic').value;

    config.providers.openai = config.providers.openai || {};
    config.providers.openai.api_key = document.getElementById('provider-openai').value;

    config.providers.gemini = config.providers.gemini || {};
    config.providers.gemini.api_key = document.getElementById('provider-gemini').value;

    config.providers.groq = config.providers.groq || {};
    config.providers.groq.api_key = document.getElementById('provider-groq').value;

    return config;
}

// Save configuration
async function saveConfig() {
    try {
        const config = collectFormData();

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
        showMessage(result.message || 'Configuration saved successfully! Restart v16 gateway to apply changes.', 'success');

        // Reload config to ensure sync
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

// Load initial data
window.addEventListener('DOMContentLoaded', () => {
    loadConfig();
    loadStatus();
});
