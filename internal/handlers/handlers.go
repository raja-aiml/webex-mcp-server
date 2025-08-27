package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ErrorResponse represents error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// HandlerOptions represents handler options
type HandlerOptions struct {
	EnableCORS   bool
	AllowOrigins []string
	Headers      map[string]string
}

// HealthHandler returns an HTTP handler for health checks
func HealthHandler(serviceName, version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]interface{}{
			"status":  "healthy",
			"service": serviceName,
			"version": version,
			"time":    time.Now().UTC().Format(time.RFC3339),
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("[Health Handler] Failed to encode response: %v", err)
		}
	}
}

// CORSMiddleware provides CORS support
func CORSMiddleware(options HandlerOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if options.EnableCORS {
				origins := "*"
				if len(options.AllowOrigins) > 0 {
					origins = options.AllowOrigins[0]
				}

				w.Header().Set("Access-Control-Allow-Origin", origins)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			for key, value := range options.Headers {
				w.Header().Set(key, value)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// LoggingMiddleware provides request logging
func LoggingMiddleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			log.Printf("[%s] %s %s - Started", serviceName, r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			log.Printf("[%s] %s %s - Completed in %v", serviceName, r.Method, r.URL.Path, duration)
		})
	}
}

// writeErrorResponse writes error response
func writeErrorResponse(w http.ResponseWriter, errorResp ErrorResponse, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(errorResp); err != nil {
		log.Printf("[Error Handler] Failed to encode error response: %v", err)
	}
}

// SetupHTTPHandlers configures HTTP handlers for the server
func SetupHTTPHandlers(server *mcp.Server, serviceName, version string) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", HealthHandler(serviceName, version))

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

	mcpHandler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)
	mux.Handle("/", mcpHandler)

	return mux
}

// ValidateRequest validates incoming requests
func ValidateRequest(r *http.Request, allowedMethods []string) error {
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
