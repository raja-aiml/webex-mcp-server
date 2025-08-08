package main

import (
	"flag"
	"log"

	"github.com/raja-aiml/webex-mcp-server/internal/app"

	// Import for side effects - registers advanced tools loader
	_ "github.com/raja-aiml/webex-mcp-server/internal/advanced_tools"
)

const (
	ServerName    = "webex-mcp-server"
	ServerVersion = "0.1.0"
)

func main() {
	var (
		httpAddr    string
		envPath     string
		useAllTools bool
	)
	flag.StringVar(&httpAddr, "http", "", "if set, use streamable HTTP at this address, instead of stdin/stdout")
	flag.StringVar(&envPath, "env", "", "path to .env file. If not set, will try to load from current directory")
	flag.BoolVar(&useAllTools, "all-tools", false, "load all tools including advanced ones (default: core tools only)")
	flag.Parse()

	application := app.New(app.Config{
		Name:        ServerName,
		Version:     ServerVersion,
		HTTPAddr:    httpAddr,
		EnvPath:     envPath,
		UseAllTools: useAllTools,
	})

	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
