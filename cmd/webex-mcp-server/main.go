package main

import (
	"flag"
	"log"

	"github.com/raja-aiml/webex-mcp-server-go/internal/app"
)

const (
	ServerName    = "webex-mcp-server"
	ServerVersion = "0.1.0"
)

func main() {
	var httpAddr string
	flag.StringVar(&httpAddr, "http", "", "if set, use streamable HTTP at this address, instead of stdin/stdout")
	flag.Parse()

	application := app.New(app.Config{
		Name:     ServerName,
		Version:  ServerVersion,
		HTTPAddr: httpAddr,
	})

	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
