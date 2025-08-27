package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// HandlerConfig represents handler configuration - TypeScript interface pattern
type HandlerConfig struct {
	ServiceName string            `json:"serviceName"`
	Version     string            `json:"version"`
	Headers     map[string]string `json:"headers,omitempty"`
	Timeout     time.Duration     `json:"timeout,omitempty"`
}

// HealthResponse represents health check response - TypeScript interface pattern
type HealthResponse struct {
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime,omitempty"`
}

// ErrorResponse represents error response - TypeScript error interface pattern
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// HandlerOptions represents handler options - TypeScript options pattern
type HandlerOptions struct {
	EnableCORS   bool              `json:"enableCors"`
	AllowOrigins []string          `json:"allowOrigins,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
}

// HealthHandler returns an HTTP handler for health checks - TypeScript factory function pattern
func HealthHandler(serviceName, version string) http.HandlerFunc {
	startTime := time.Now() // TypeScript closure pattern

	return func(w http.ResponseWriter, r *http.Request) {
		// TypeScript-style method validation
		if r.Method != http.MethodGet {
			writeErrorResponse(w, ErrorResponse{
				Error:   "method_not_allowed",
				Message: "Only GET method is allowed",
				Code:    http.StatusMethodNotAllowed,
			}, http.StatusMethodNotAllowed)
			return
		}

		// Set response headers - TypeScript header management pattern
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Calculate uptime - TypeScript duration calculation pattern
		uptime := time.Since(startTime).String()

		// Create response object - TypeScript object creation pattern
		response := HealthResponse{
			Status:    "healthy",
			Service:   serviceName,
			Version:   version,
			Timestamp: time.Now().UTC(),
			Uptime:    uptime,
		}

		// TypeScript-style error handling
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("[Health Handler] Failed to encode response: %v", err)
		}
	}
}

// CORSMiddleware provides CORS support - TypeScript middleware pattern
func CORSMiddleware(options HandlerOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if options.EnableCORS {
				// TypeScript-style array join equivalent
				origins := "*"
				if len(options.AllowOrigins) > 0 {
					origins = options.AllowOrigins[0] // Simplified for demo
				}

				// Set CORS headers - TypeScript header setting pattern
				w.Header().Set("Access-Control-Allow-Origin", origins)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}

			// Handle preflight requests - TypeScript OPTIONS handling pattern
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Set custom headers - TypeScript object iteration pattern
			for key, value := range options.Headers {
				w.Header().Set(key, value)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// LoggingMiddleware provides request logging - TypeScript middleware pattern
func LoggingMiddleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// TypeScript-style logging with template literals equivalent
			log.Printf("[%s] %s %s - Started", serviceName, r.Method, r.URL.Path)

			next.ServeHTTP(w, r)

			// Calculate duration - TypeScript duration pattern
			duration := time.Since(start)
			log.Printf("[%s] %s %s - Completed in %v", serviceName, r.Method, r.URL.Path, duration)
		})
	}
}

// writeErrorResponse writes error response - TypeScript utility function pattern
func writeErrorResponse(w http.ResponseWriter, errorResp ErrorResponse, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(errorResp); err != nil {
		log.Printf("[Error Handler] Failed to encode error response: %v", err)
	}
}

// SetupHTTPHandlers configures HTTP handlers for the server - TypeScript setup function pattern
func SetupHTTPHandlers(server *mcp.Server, serviceName, version string) *http.ServeMux {
	// Create mux - TypeScript router pattern
	mux := http.NewServeMux()

	// Health check endpoint - TypeScript route registration pattern
	mux.HandleFunc("/health", HealthHandler(serviceName, version))

	// Server info endpoint - TypeScript additional endpoint pattern
	mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeErrorResponse(w, ErrorResponse{
				Error:   "method_not_allowed",
				Message: "Only GET method is allowed",
				Code:    http.StatusMethodNotAllowed,
			}, http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		info := map[string]interface{}{
			"name":    serviceName,
			"version": version,
			"type":    "mcp-server",
			"capabilities": map[string]interface{}{
				"tools": true,
			},
		}
		json.NewEncoder(w).Encode(info)
	})

	// MCP handler - TypeScript main handler pattern
	mcpHandler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)
	mux.Handle("/", mcpHandler)

	return mux
}

// ValidateRequest validates incoming requests - TypeScript validation pattern
func ValidateRequest(r *http.Request, allowedMethods []string) error {
	// TypeScript-style array includes equivalent
	methodAllowed := false
	for _, method := range allowedMethods {
		if r.Method == method {
			methodAllowed = true
			break
		}
	}

	if !methodAllowed {
		return fmt.Errorf("method %s not allowed", r.Method)
	}

	return nil
}
