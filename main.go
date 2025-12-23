package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/danhigham/ergometer.live/server"
)

func main() {
	log.Println("Starting Ergometer.Live WebSocket Server...")

	// Create server
	srv := server.NewServer(":8080")

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("\nReceived shutdown signal...")

	// Shutdown server
	srv.Shutdown()

	log.Println("Server stopped")
}
