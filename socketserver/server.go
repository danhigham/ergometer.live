package socketserver

import (
	"log"
	"net/http"

	"github.com/danhigham/ergometer.live/broadcast"
	"github.com/danhigham/ergometer.live/pm5"
)

// Server represents the HTTP/WebSocket server
type Server struct {
	hub     *broadcast.Hub
	manager *pm5.Manager
	addr    string
}

// NewServer creates a new Server instance
func NewServer(addr string) *Server {
	hub := broadcast.NewHub()
	manager := pm5.GetManager(hub)

	// Set message handler for inbound client messages
	hub.SetMessageHandler(func(client *broadcast.Client, message []byte) {
		handleClientMessage(manager, client, message)
	})

	return &Server{
		hub:     hub,
		manager: manager,
		addr:    addr,
	}
}

// Start starts the HTTP server and hub
func (s *Server) Start() error {
	// Start the hub
	go s.hub.Run()

	// Setup HTTP routes
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(s.hub, w, r)
	})
	http.HandleFunc("/", serveHome)

	log.Printf("Server starting on %s", s.addr)
	return http.ListenAndServe(s.addr, nil)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() {
	log.Println("Shutting down server...")
	s.manager.Shutdown()
	s.hub.Shutdown()
}
