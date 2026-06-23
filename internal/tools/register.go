// SPDX-License-Identifier: MIT

package tools

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds every PDF tool to the server.
func Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_info",
		Description: "Inspect a PDF: page count, per-page sizes (points), and Info metadata (title, author, etc.).",
	}, handleInfo)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_extract_text",
		Description: "Extract text from a PDF, one section per page. Optional page range and separator.",
	}, handleText)
}
