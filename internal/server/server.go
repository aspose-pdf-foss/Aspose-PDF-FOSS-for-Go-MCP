// SPDX-License-Identifier: MIT

// Package server wires up the MCP server and registers all PDF tools.
package server

import (
	"github.com/aspose-pdf-foss/aspose-pdf-foss-mcp/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// New builds the MCP server with all tools registered.
func New() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "aspose-pdf-foss-mcp",
		Version: "0.1.0",
	}, nil)
	tools.Register(server)
	return server
}
