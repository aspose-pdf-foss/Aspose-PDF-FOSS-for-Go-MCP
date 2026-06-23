// SPDX-License-Identifier: MIT

package tools

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds every PDF tool to the server.
func Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_info",
		Description: "Inspect a PDF: page count, per-page sizes (points), and Info metadata (title, author, etc.).",
	}, handleInfo)
}
