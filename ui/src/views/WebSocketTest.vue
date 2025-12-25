<template>
  <div class="websocket-test">
    <div class="header">
      <h1>WebSocket Connection Test</h1>
      <p class="subtitle">Testing connection to local PM5 WebSocket server</p>
    </div>

    <div class="connection-panel">
      <div class="status-row">
        <div class="status-indicator" :class="statusClass">
          <span class="status-dot"></span>
          <span class="status-text">{{ statusText }}</span>
        </div>

        <div class="actions">
          <button
            v-if="!connected"
            @click="connect"
            :disabled="connecting"
            class="btn btn-primary"
          >
            {{ connecting ? 'Connecting...' : 'Connect' }}
          </button>
          <button v-else @click="disconnect" class="btn btn-secondary">Disconnect</button>
          <button @click="clearMessages" class="btn btn-secondary" :disabled="messages.length === 0">
            Clear Messages
          </button>
        </div>
      </div>

      <div v-if="error" class="error-message">
        <strong>Error:</strong> {{ error }}
      </div>

      <div class="config-info">
        <p><strong>WebSocket URL:</strong> <code>{{ wsUrl }}</code></p>
        <p class="note">
          ⚠️ Make sure the PM5 server is running on localhost:8080 and CORS is enabled
        </p>
      </div>
    </div>

    <div v-if="lastMessage" class="last-message-panel">
      <h2>Latest Message</h2>
      <div class="message-display">
        <div class="message-header">
          <span class="message-type">{{ lastMessage.type }}</span>
          <span class="message-time">{{ formatTime(lastMessage.timestamp) }}</span>
        </div>
        <pre class="message-data">{{ JSON.stringify(lastMessage.data, null, 2) }}</pre>
      </div>
    </div>

    <div class="messages-panel">
      <div class="panel-header">
        <h2>Message Log</h2>
        <span class="message-count">{{ messages.length }} messages</span>
      </div>

      <div v-if="messages.length === 0" class="empty-state">
        <p>No messages received yet. Connect to the WebSocket server to see messages.</p>
      </div>

      <div v-else class="messages-list">
        <div v-for="(msg, index) in messages" :key="index" class="message-item">
          <div class="message-header">
            <span class="message-type">{{ msg.type }}</span>
            <span class="message-time">{{ formatTime(msg.timestamp) }}</span>
          </div>
          <pre class="message-data">{{ JSON.stringify(msg.data, null, 2) }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useWebSocket } from '@/composables/useWebSocket'

const wsUrl = 'ws://localhost:8080/ws'
const { connected, connecting, error, lastMessage, messages, connect, disconnect, clearMessages } =
  useWebSocket(wsUrl)

const statusClass = computed(() => {
  if (connected.value) return 'status-connected'
  if (connecting.value) return 'status-connecting'
  if (error.value) return 'status-error'
  return 'status-disconnected'
})

const statusText = computed(() => {
  if (connected.value) return 'Connected'
  if (connecting.value) return 'Connecting...'
  if (error.value) return 'Error'
  return 'Disconnected'
})

const formatTime = (timestamp: string) => {
  const date = new Date(timestamp)
  return date.toLocaleTimeString()
}
</script>

<style scoped>
.websocket-test {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.header {
  margin-bottom: 2rem;
}

.header h1 {
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.subtitle {
  color: #666;
  font-size: 1rem;
}

.connection-panel {
  background: #f5f5f5;
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.status-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 600;
}

.status-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #999;
}

.status-connected .status-dot {
  background: #22c55e;
  box-shadow: 0 0 8px rgba(34, 197, 94, 0.5);
}

.status-connecting .status-dot {
  background: #f59e0b;
  animation: pulse 1.5s ease-in-out infinite;
}

.status-error .status-dot {
  background: #ef4444;
}

.status-disconnected .status-dot {
  background: #999;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.actions {
  display: flex;
  gap: 0.5rem;
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #2563eb;
}

.btn-secondary {
  background: #e5e7eb;
  color: #374151;
}

.btn-secondary:hover:not(:disabled) {
  background: #d1d5db;
}

.error-message {
  background: #fee2e2;
  border: 1px solid #fecaca;
  border-radius: 4px;
  padding: 0.75rem;
  margin-bottom: 1rem;
  color: #991b1b;
}

.config-info {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid #e5e7eb;
}

.config-info p {
  margin: 0.5rem 0;
  font-size: 0.875rem;
}

.config-info code {
  background: white;
  padding: 0.25rem 0.5rem;
  border-radius: 3px;
  font-family: monospace;
  color: #1f2937;
}

.note {
  color: #f59e0b;
  font-size: 0.875rem;
  margin-top: 0.5rem;
}

.last-message-panel {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.last-message-panel h2 {
  font-size: 1.25rem;
  margin-bottom: 1rem;
}

.messages-panel {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 1.5rem;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.panel-header h2 {
  font-size: 1.25rem;
  margin: 0;
}

.message-count {
  background: #e5e7eb;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  color: #374151;
}

.empty-state {
  text-align: center;
  padding: 3rem 1rem;
  color: #9ca3af;
}

.messages-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  max-height: 600px;
  overflow-y: auto;
}

.message-item {
  border: 1px solid #e5e7eb;
  border-radius: 4px;
  padding: 1rem;
}

.message-display {
  border: 1px solid #e5e7eb;
  border-radius: 4px;
  padding: 1rem;
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid #e5e7eb;
}

.message-type {
  background: #dbeafe;
  color: #1e40af;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
}

.message-time {
  font-size: 0.75rem;
  color: #9ca3af;
}

.message-data {
  margin: 0;
  padding: 0.75rem;
  background: #f9fafb;
  border-radius: 4px;
  font-size: 0.875rem;
  font-family: monospace;
  overflow-x: auto;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
