// SPDX-License-Identifier: MIT

// Command aspose-pdf-foss-for-go-mcp is a stdio MCP server wrapping the
// aspose-pdf-foss-for-go PDF library.
package main

import (
	"context"
	"log"

	"github.com/aspose-pdf-foss/aspose-pdf-foss-for-go-mcp/internal/server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	srv := server.New()
	if err := srv.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
