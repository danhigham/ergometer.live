package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/danhigham/ergometer.live/api/config"
	"github.com/danhigham/ergometer.live/api/middleware"
	"github.com/danhigham/ergometer.live/api/services"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting Ergometer.Live REST API Server...")

	// Load configuration
	cfg := config.Load()

	// Initialize Firebase service
	firebaseService, err := services.NewFirebaseService(
		cfg.FirebaseCredentialsPath,
		cfg.FirebaseProjectID,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase service: %v", err)
	}

	// Create router
	router := mux.NewRouter()

	// Health check endpoint (no auth required)
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")

	// API v1 routes
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// Auth verification endpoint (requires auth)
	authRouter := apiV1.PathPrefix("/auth").Subrouter()
	authRouter.Use(middleware.Auth(firebaseService))
	authRouter.HandleFunc("/verify", verifyHandler).Methods("POST")

	// Apply middleware
	corsMiddleware := middleware.NewCORS(cfg.AllowedOrigins)
	handler := corsMiddleware.Handler(router)
	handler = middleware.Logger(handler)

	// Start server
	addr := ":" + cfg.Port
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-sigChan
	log.Println("\nReceived shutdown signal...")
	log.Println("Server stopped")
}

// healthCheckHandler handles health check requests
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

// verifyHandler handles auth verification requests
func verifyHandler(w http.ResponseWriter, r *http.Request) {
	// If we reach here, the token was valid (auth middleware passed)
	uid, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"uid":"` + uid + `","verified":true}`))
}
