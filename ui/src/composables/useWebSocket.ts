import { ref, onUnmounted, shallowRef } from 'vue'

export interface WebSocketMessage {
  type: string
  data: any
  timestamp: string
}

// Use globalThis to persist across HMR updates
declare global {
  var __WS_SHARED__: WebSocket | null
  var __WS_CONNECTION_COUNT__: number
  var __WS_IS_CONNECTING__: boolean
  var __WS_INSTANCE_COUNT__: number
}

// Shared WebSocket instance to prevent multiple connections (persists across HMR)
const sharedWs = {
  get value() {
    return globalThis.__WS_SHARED__ ?? null
  },
  set value(ws: WebSocket | null) {
    globalThis.__WS_SHARED__ = ws
  }
}

const connectionCount = {
  get value() {
    return globalThis.__WS_CONNECTION_COUNT__ ?? 0
  },
  set value(count: number) {
    globalThis.__WS_CONNECTION_COUNT__ = count
  }
}

const isConnecting = {
  get value() {
    return globalThis.__WS_IS_CONNECTING__ ?? false
  },
  set value(connecting: boolean) {
    globalThis.__WS_IS_CONNECTING__ = connecting
  }
}

const instanceCount = {
  get value() {
    return globalThis.__WS_INSTANCE_COUNT__ ?? 0
  },
  set value(count: number) {
    globalThis.__WS_INSTANCE_COUNT__ = count
  }
}

export function useWebSocket(url: string) {
  instanceCount.value++
  const instanceId = instanceCount.value
  console.log(`[useWebSocket] Instance #${instanceId} created (HMR-safe)`)

  const ws = shallowRef<WebSocket | null>(null)
  const connected = ref(false)
  const connecting = ref(false)
  const error = ref<string | null>(null)
  const lastMessage = ref<WebSocketMessage | null>(null)
  const messages = ref<WebSocketMessage[]>([])

  const connect = () => {
    connectionCount.value++
    console.log(`[useWebSocket #${instanceId}] connect() called, count: ${connectionCount.value}, isConnecting: ${isConnecting.value}`)

    // Prevent concurrent connection attempts
    if (isConnecting.value) {
      console.log('[useWebSocket] Connection attempt already in progress, skipping')
      return
    }

    // Use shared instance if it exists and is connected/connecting
    if (sharedWs.value) {
      if (sharedWs.value.readyState === WebSocket.OPEN) {
        console.log('[useWebSocket] Using existing connected WebSocket (HMR reuse)')
        ws.value = sharedWs.value
        connected.value = true
        return
      }
      if (sharedWs.value.readyState === WebSocket.CONNECTING) {
        console.log('[useWebSocket] Using existing connecting WebSocket (HMR reuse)')
        ws.value = sharedWs.value
        connecting.value = true
        return
      }
      // Clean up dead connection
      console.log('[useWebSocket] Cleaning up dead connection')
      sharedWs.value = null
    }

    isConnecting.value = true
    connecting.value = true
    error.value = null

    try {
      console.log('[useWebSocket] Creating new WebSocket connection to', url)
      sharedWs.value = new WebSocket(url)
      ws.value = sharedWs.value

      ws.value.onopen = () => {
        console.log('[useWebSocket] WebSocket connected to', url)
        isConnecting.value = false
        connected.value = true
        connecting.value = false
        error.value = null
      }

      ws.value.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data) as WebSocketMessage
          lastMessage.value = message
          messages.value.push(message)

          // Keep only last 50 messages
          if (messages.value.length > 50) {
            messages.value.shift()
          }
        } catch (err) {
          console.error('Failed to parse WebSocket message:', err)
        }
      }

      ws.value.onerror = (event) => {
        console.error('[useWebSocket] WebSocket error:', event)
        isConnecting.value = false
        error.value = 'Connection error occurred'
        connecting.value = false
      }

      ws.value.onclose = (event) => {
        console.log('[useWebSocket] WebSocket closed:', event.code, event.reason)
        isConnecting.value = false
        connected.value = false
        connecting.value = false
        sharedWs.value = null

        if (!event.wasClean) {
          error.value = `Connection lost (code: ${event.code})`
        }
      }
    } catch (err) {
      console.error('[useWebSocket] Failed to create WebSocket:', err)
      isConnecting.value = false
      error.value = err instanceof Error ? err.message : 'Failed to connect'
      connecting.value = false
    }
  }

  const disconnect = () => {
    console.log(`[useWebSocket #${instanceId}] disconnect() called`)
    if (ws.value) {
      ws.value.close()
      ws.value = null
    }
    if (sharedWs.value) {
      sharedWs.value.close()
      sharedWs.value = null
    }
    isConnecting.value = false
    connected.value = false
    connecting.value = false
  }

  const send = (data: any) => {
    if (ws.value?.readyState === WebSocket.OPEN) {
      ws.value.send(JSON.stringify(data))
    } else {
      console.error('WebSocket is not connected')
    }
  }

  const clearMessages = () => {
    messages.value = []
    lastMessage.value = null
  }

  // Cleanup on component unmount
  onUnmounted(() => {
    console.log(`[useWebSocket #${instanceId}] Component unmounted, cleaning up`)
    // Don't close shared connection on unmount, only when explicitly disconnected
    // This prevents closing the connection when component re-renders
    ws.value = null
    connected.value = false
    connecting.value = false
  })

  return {
    connected,
    connecting,
    error,
    lastMessage,
    messages,
    connect,
    disconnect,
    send,
    clearMessages,
  }
}
