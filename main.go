package main

import (
	"log"

	"aayushsiwa/spendwise-mcp/backend"
	"aayushsiwa/spendwise-mcp/config"
	"aayushsiwa/spendwise-mcp/handlers"
	"aayushsiwa/spendwise-mcp/routes"
	"aayushsiwa/spendwise-mcp/services"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	cfg := config.Load()
	apiClient := backend.NewHTTPClient(cfg.BackendBaseURL)
	svc := services.NewService(apiClient)
	handler := handlers.NewHandler(svc)

	mcpServer := server.NewMCPServer(
		cfg.ServerName,
		cfg.ServerVersion,
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	apiRoutes := routes.NewRoutes(handler)
	routes.AttachRoutes(mcpServer, apiRoutes)

	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("mcp server error: %v", err)
	}
}
