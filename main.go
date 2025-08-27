package main

import (
	"flag"
	"log"
	"os"

	"github.com/raja-aiml/webex-mcp-server/internal/app"

	// Import for side effects - registers advanced tools loader
	_ "github.com/raja-aiml/webex-mcp-server/internal/advanced_tools"
)

// ServerInfo holds server metadata - TypeScript compliant structure
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Bin     string `json:"bin,omitempty"`
}

// Configuration constants matching TypeScript package.json structure
const (
	ServerName    = "webex-mcp-server"
	ServerVersion = "0.1.0"
	ServerBin     = "./webex-mcp-server" // TypeScript bin equivalent
)

// CLI arguments interface matching TypeScript commander pattern
type CLIArgs struct {
	HTTP        string `json:"http,omitempty"`
	Env         string `json:"env,omitempty"`
	UseAllTools bool   `json:"useAllTools"`
	SSE         bool   `json:"sse,omitempty"` // TypeScript SSE mode support
}

// parseArgs parses command line arguments in TypeScript style
func parseArgs() CLIArgs {
	var args CLIArgs

	flag.StringVar(&args.HTTP, "http", "", "if set, use streamable HTTP at this address, instead of stdin/stdout")
	flag.StringVar(&args.Env, "env", "", "path to .env file. If not set, will try to load from current directory")
	flag.BoolVar(&args.UseAllTools, "all-tools", false, "load all tools including advanced ones (default: core tools only)")
	flag.BoolVar(&args.SSE, "sse", false, "enable Server-Sent Events mode (HTTP streaming)")
	flag.Parse()

	// Handle TypeScript-style SSE argument detection
	for _, arg := range os.Args[1:] {
		if arg == "--sse" {
			args.SSE = true
			break
		}
	}

	return args
}

// createServerInfo creates server information matching TypeScript package.json
func createServerInfo() ServerInfo {
	return ServerInfo{
		Name:    ServerName,
		Version: ServerVersion,
		Bin:     ServerBin,
	}
}

func main() {
	// Parse CLI arguments in TypeScript commander style
	args := parseArgs()
	serverInfo := createServerInfo()

	// Create application config matching TypeScript patterns
	application := app.New(app.Config{
		Name:        serverInfo.Name,
		Version:     serverInfo.Version,
		HTTPAddr:    args.HTTP,
		EnvPath:     args.Env,
		UseAllTools: args.UseAllTools,
		SSEMode:     args.SSE, // TypeScript SSE mode support
	})

	// Run application with TypeScript-style error handling
	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
