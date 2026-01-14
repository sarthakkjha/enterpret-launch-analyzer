package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// getEnv returns the value of an environment variable or a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getPort returns the server port from environment variable or default
func getPort() int {
	portStr := getEnv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Invalid PORT value '%s', using default 8080", portStr)
		return 8080
	}
	return port
}

// CORSMiddleware adds CORS headers to responses
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request for debugging
		log.Printf("📥 Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Set simple permissive CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		// w.Header().Set("Access-Control-Allow-Credentials", "true") // Not needed for wildcards

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Server represents the HTTP server
type Server struct {
	handler *APIHandler
	port    int
}

// NewServer creates a new server instance
func NewServer(handler *APIHandler, port int) *Server {
	return &Server{
		handler: handler,
		port:    port,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Enterpret Data Analysis API is running 🚀"))
	})
	mux.HandleFunc("/api/health", s.handler.HandleHealth)
	mux.HandleFunc("/api/upload", s.handler.HandleUpload)
	mux.HandleFunc("/api/analyze", s.handler.HandleAnalyze)

	// Wrap with CORS middleware
	handler := CORSMiddleware(mux)

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("🚀 Server starting on http://localhost%s", addr)
	log.Printf("📊 Enterpret Pre/Post Launch Analysis Dashboard API")
	log.Printf("📁 Endpoints:")
	log.Printf("   GET  /api/health  - Health check")
	log.Printf("   POST /api/upload  - Upload CSV files")
	log.Printf("   POST /api/analyze - Run analysis")

	return http.ListenAndServe(addr, handler)
}

func main() {
	// Get API key from environment variable
	apiKey := getEnv("GROQ_API_KEY", "")
	if apiKey == "" {
		log.Fatal("GROQ_API_KEY environment variable is required")
	}

	// Initialize dependencies using dependency injection
	groqClient := NewGroqClient(apiKey)
	csvParser := NewCSVReviewParser()
	analysisService := NewAnalysisService(groqClient)
	apiHandler := NewAPIHandler(csvParser, analysisService)

	// Create and start server
	port := getPort()
	server := NewServer(apiHandler, port)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
