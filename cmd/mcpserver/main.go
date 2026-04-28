package main

import (
	"datapulse/internal/mcpserver"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	addr := os.Getenv("MCP_ADDR")
	if addr == "" {
		addr = ":8091"
	}

	log.Printf("mcp server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mcpserver.NewServer()))
}
