package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// HealthHandler returns an HTTP handler for health checks
func HealthHandler(serviceName, version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]any{
			"status":  "healthy",
			"service": serviceName,
			"version": version,
			"time":    time.Now().UTC().Format(time.RFC3339),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode health response: %v", err)
		}
	}
}

// SetupHTTPHandlers configures HTTP handlers for the server
func SetupHTTPHandlers(server *mcp.Server, serviceName, version string) *http.ServeMux {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", HealthHandler(serviceName, version))

	// MCP handler
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)
	mux.Handle("/", handler)

	return mux
}
