// WebSocket connection
let ws = null;
let reconnectInterval = null;

// DOM elements
const elements = {
    connectionStatus: document.getElementById('connection-status'),
    deviceInfo: document.getElementById('device-info'),
    startButton: document.getElementById('start-workout'),
    stopButton: document.getElementById('stop-workout'),
    log: document.getElementById('log'),

    // Stats
    time: document.getElementById('time'),
    distance: document.getElementById('distance-val'),
    pace: document.getElementById('pace'),
    avgPace: document.getElementById('avg-pace'),
    power: document.getElementById('power'),
    strokeRate: document.getElementById('stroke-rate'),
    calories: document.getElementById('calories'),
    heartRate: document.getElementById('heart-rate'),
    workoutState: document.getElementById('workout-state'),
    rowingState: document.getElementById('rowing-state'),
};

// Initialize
function init() {
    // Setup workout type change handlers
    document.querySelectorAll('input[name="workout-type"]').forEach(radio => {
        radio.addEventListener('change', handleWorkoutTypeChange);
    });

    // Setup button handlers
    elements.startButton.addEventListener('click', startWorkout);
    elements.stopButton.addEventListener('click', stopWorkout);

    // Connect to WebSocket
    connect();
}

// Handle workout type change
function handleWorkoutTypeChange(e) {
    const type = e.target.value;

    document.getElementById('distance-options').style.display =
        type === 'fixed_distance' ? 'block' : 'none';

    document.getElementById('time-options').style.display =
        type === 'fixed_time' ? 'block' : 'none';
}

// WebSocket connection
function connect() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;

    addLog('Connecting to ' + wsUrl);
    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
        addLog('Connected to ergometer.live');
        updateConnectionStatus(true);
        getStatus();

        if (reconnectInterval) {
            clearInterval(reconnectInterval);
            reconnectInterval = null;
        }
    };

    ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        handleMessage(msg);
    };

    ws.onclose = () => {
        addLog('Disconnected from server');
        updateConnectionStatus(false);
        scheduleReconnect();
    };

    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        addLog('Connection error');
    };
}

// Schedule reconnection
function scheduleReconnect() {
    if (reconnectInterval) return;

    addLog('Reconnecting in 5 seconds...');
    reconnectInterval = setInterval(() => {
        connect();
    }, 5000);
}

// Handle incoming messages
function handleMessage(msg) {
    switch (msg.type) {
        case 'workout_stats':
            updateStats(msg.data);
            break;

        case 'workout_state':
            updateWorkoutState(msg.data);
            break;

        case 'status':
            updateDeviceInfo(msg.data);
            break;

        case 'error':
            addLog('Error: ' + msg.data.message, 'error');
            break;

        case 'success':
            addLog(msg.data.message, 'success');
            if (msg.data.action === 'start_workout') {
                enableStopButton();
            } else if (msg.data.action === 'stop_workout') {
                enableStartButton();
            }
            break;

        case 'workout_started':
            addLog('Workout started!', 'success');
            enableStopButton();
            break;

        case 'workout_ended':
            addLog('Workout ended', 'info');
            enableStartButton();
            break;
    }
}

// Update stats display
function updateStats(data) {
    elements.time.textContent = formatTime(data.elapsed_time);
    elements.distance.textContent = Math.round(data.distance) + ' m';
    elements.pace.textContent = formatPace(data.pace);
    elements.avgPace.textContent = formatPace(data.avg_pace);
    elements.power.textContent = data.power;
    elements.strokeRate.textContent = data.stroke_rate;
    elements.calories.textContent = data.calories;
    elements.heartRate.textContent = data.heart_rate > 0 ? data.heart_rate : '--';
    elements.workoutState.textContent = data.workout_state;
    elements.rowingState.textContent = data.rowing_state;

    // Update state badge colors
    updateStateBadge(elements.workoutState, data.workout_state);
    updateStateBadge(elements.rowingState, data.rowing_state);
}

// Update workout state
function updateWorkoutState(data) {
    if (data.workout_state) {
        elements.workoutState.textContent = data.workout_state;
        updateStateBadge(elements.workoutState, data.workout_state);
    }

    if (data.is_active !== undefined) {
        if (data.is_active) {
            enableStopButton();
        } else {
            enableStartButton();
        }
    }
}

// Update device info
function updateDeviceInfo(data) {
    if (data.connected) {
        elements.deviceInfo.textContent = `${data.erg_type} (${data.serial}) - Battery: ${data.battery}%`;
        addLog(`Connected to ${data.erg_type}`, 'success');
    } else {
        elements.deviceInfo.textContent = 'No device connected';
    }
}

// Update connection status
function updateConnectionStatus(connected) {
    if (connected) {
        elements.connectionStatus.textContent = 'Connected';
        elements.connectionStatus.className = 'status-connected';
    } else {
        elements.connectionStatus.textContent = 'Disconnected';
        elements.connectionStatus.className = 'status-disconnected';
        elements.deviceInfo.textContent = '';
        enableStartButton();
    }
}

// Update state badge color
function updateStateBadge(element, state) {
    element.className = 'state-badge';

    if (state.includes('Workout') || state.includes('Active') || state.includes('Driving') || state.includes('Recovery')) {
        element.classList.add('state-active');
    } else if (state.includes('Ready') || state.includes('Inactive')) {
        element.classList.add('state-ready');
    } else {
        element.classList.add('state-idle');
    }
}

// Format time (seconds -> M:SS.s)
function formatTime(seconds) {
    if (!seconds || seconds === 0) return '0:00.0';

    const mins = Math.floor(seconds / 60);
    const secs = (seconds % 60).toFixed(1);
    return `${mins}:${secs.padStart(4, '0')}`;
}

// Format pace (seconds per 500m -> M:SS)
function formatPace(seconds) {
    if (!seconds || seconds === 0) return '--:--';

    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, '0')}`;
}

// Enable/disable buttons
function enableStopButton() {
    elements.startButton.disabled = true;
    elements.stopButton.disabled = false;
}

function enableStartButton() {
    elements.startButton.disabled = false;
    elements.stopButton.disabled = true;
}

// Start workout
function startWorkout() {
    const type = document.querySelector('input[name="workout-type"]:checked').value;

    const data = {
        workout_type: type
    };

    if (type === 'fixed_distance') {
        data.distance = parseInt(document.getElementById('distance').value);
        data.split_distance = parseInt(document.getElementById('split-distance').value) || 0;
    } else if (type === 'fixed_time') {
        const mins = parseInt(document.getElementById('time-minutes').value) || 0;
        const secs = parseInt(document.getElementById('time-seconds').value) || 0;
        data.time = (mins * 60) + secs;
        data.split_time = parseInt(document.getElementById('split-time').value) || 0;
    }

    send({ type: 'start_workout', data: data });
    addLog(`Starting ${type} workout...`, 'info');
}

// Stop workout
function stopWorkout() {
    send({ type: 'stop_workout' });
    addLog('Stopping workout...', 'info');
}

// Get status
function getStatus() {
    send({ type: 'get_status' });
}

// Send message
function send(msg) {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify(msg));
    } else {
        addLog('Not connected to server', 'error');
    }
}

// Add log message
function addLog(message, type = 'info') {
    const logEntry = document.createElement('div');
    logEntry.className = `log-entry log-${type}`;

    const timestamp = new Date().toLocaleTimeString();
    logEntry.textContent = `[${timestamp}] ${message}`;

    elements.log.appendChild(logEntry);
    elements.log.scrollTop = elements.log.scrollHeight;

    // Keep only last 50 messages
    while (elements.log.children.length > 50) {
        elements.log.removeChild(elements.log.firstChild);
    }
}

// Initialize on load
document.addEventListener('DOMContentLoaded', init);
