package broadcast

import (
	"log"
	"sync"
)

// InboundMessage represents a message received from a client
type InboundMessage struct {
	client  *Client
	message []byte
}

// MessageHandler is a function that processes inbound messages from clients
type MessageHandler func(client *Client, message []byte)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from clients
	inbound chan *InboundMessage

	// Outbound messages to broadcast to clients
	broadcast chan []byte

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Shutdown signal
	shutdown chan struct{}

	// Wait group for graceful shutdown
	wg sync.WaitGroup

	// Handler for inbound messages
	messageHandler MessageHandler
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		inbound:    make(chan *InboundMessage, 256),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		shutdown:   make(chan struct{}),
	}
}

// SetMessageHandler sets the handler function for inbound messages
func (h *Hub) SetMessageHandler(handler MessageHandler) {
	h.messageHandler = handler
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	h.wg.Add(1)
	defer h.wg.Done()

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client registered. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client unregistered. Total clients: %d", len(h.clients))
			}

		case message := <-h.broadcast:
			// Broadcast to all clients
			h.broadcastToClients(message)

		case msg := <-h.inbound:
			// Process inbound message from client
			if h.messageHandler != nil {
				go h.messageHandler(msg.client, msg.message)
			}

		case <-h.shutdown:
			// Close all client connections
			for client := range h.clients {
				close(client.send)
				delete(h.clients, client)
			}
			log.Println("Hub shutdown complete")
			return
		}
	}
}

// broadcastToClients sends a message to all registered clients
func (h *Hub) broadcastToClients(message []byte) {
	for client := range h.clients {
		select {
		case client.send <- message:
			// Message sent successfully
		default:
			// Client's send channel is full, close and remove the client
			log.Printf("Client send buffer full, removing slow client")
			close(client.send)
			delete(h.clients, client)
		}
	}
}

// Broadcast queues a message to be broadcast to all clients
func (h *Hub) Broadcast(message []byte) {
	select {
	case h.broadcast <- message:
		// Message queued successfully
	default:
		// Broadcast buffer full, drop message
		log.Printf("Broadcast buffer full, dropping message")
	}
}

// Register queues a client for registration
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister queues a client for removal
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// Shutdown gracefully shuts down the hub
func (h *Hub) Shutdown() {
	close(h.shutdown)
	h.wg.Wait()
}

// ClientCount returns the number of connected clients
func (h *Hub) ClientCount() int {
	return len(h.clients)
}
