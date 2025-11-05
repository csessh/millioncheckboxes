// WebSocket connection
let ws = null;
let reconnectTimer = null;
const GRID_SIZE = 10;
const TOTAL_CHECKBOXES = GRID_SIZE * GRID_SIZE;

// DOM elements
const grid = document.getElementById('checkbox-grid');
const statusEl = document.getElementById('connection-status');
const userCountEl = document.getElementById('user-count');

// State
const checkboxStates = new Array(TOTAL_CHECKBOXES).fill(false);

// Initialize the grid
function initGrid() {
    grid.innerHTML = '';

    for (let i = 0; i < TOTAL_CHECKBOXES; i++) {
        const wrapper = document.createElement('div');
        wrapper.className = 'checkbox-wrapper';

        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.id = `cb-${i}`;
        checkbox.dataset.index = i;
        checkbox.checked = checkboxStates[i];

        checkbox.addEventListener('change', handleCheckboxChange);

        wrapper.appendChild(checkbox);
        grid.appendChild(wrapper);
    }
}

// Handle checkbox change
function handleCheckboxChange(event) {
    const checkbox = event.target;
    const index = parseInt(checkbox.dataset.index);
    const isChecked = checkbox.checked;

    // Update local state
    checkboxStates[index] = isChecked;

    // Send to server
    sendCheckboxUpdate(index, isChecked);
}

// Send checkbox update to server
function sendCheckboxUpdate(index, isChecked) {
    if (ws && ws.readyState === WebSocket.OPEN) {
        const message = {
            cmd: 'SET',
            index: index,
            value: isChecked ? 'true' : 'false'
        };
        ws.send(JSON.stringify(message));
    }
}

// Handle incoming WebSocket message
function handleMessage(event) {
    try {
        const message = JSON.parse(event.data);

        if (message.cmd === 'SET') {
            const index = message.index;
            const isChecked = message.value === 'true';

            // Update local state
            checkboxStates[index] = isChecked;

            // Update UI
            const checkbox = document.getElementById(`cb-${index}`);
            if (checkbox && checkbox.checked !== isChecked) {
                checkbox.checked = isChecked;
            }
        }
    } catch (err) {
        console.error('Error parsing message:', err);
    }
}

// Connect to WebSocket server
function connect() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;

    console.log('Connecting to:', wsUrl);

    try {
        ws = new WebSocket(wsUrl);

        ws.onopen = () => {
            console.log('Connected to server');
            statusEl.textContent = 'Connected';
            statusEl.className = 'connected';

            // Enable all checkboxes
            document.querySelectorAll('input[type="checkbox"]').forEach(cb => {
                cb.disabled = false;
            });

            // Clear reconnect timer
            if (reconnectTimer) {
                clearTimeout(reconnectTimer);
                reconnectTimer = null;
            }
        };

        ws.onmessage = handleMessage;

        ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        ws.onclose = () => {
            console.log('Disconnected from server');
            statusEl.textContent = 'Disconnected';
            statusEl.className = 'disconnected';

            // Disable all checkboxes
            document.querySelectorAll('input[type="checkbox"]').forEach(cb => {
                cb.disabled = true;
            });

            // Attempt to reconnect after 3 seconds
            if (!reconnectTimer) {
                reconnectTimer = setTimeout(connect, 3000);
            }
        };
    } catch (err) {
        console.error('Error creating WebSocket:', err);
        // Attempt to reconnect after 3 seconds
        if (!reconnectTimer) {
            reconnectTimer = setTimeout(connect, 3000);
        }
    }
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
    initGrid();
    connect();
});
