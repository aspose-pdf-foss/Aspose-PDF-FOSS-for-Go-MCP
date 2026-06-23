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
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_render_page",
		Description: "Render a single PDF page to a PNG or JPEG image file on disk.",
	}, handleRender)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_merge",
		Description: "Merge two or more PDFs into a single output file, in the order given.",
	}, handleMerge)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_split",
		Description: "Split a PDF into one file per page, written to output_dir as page_001.pdf, page_002.pdf, ...",
	}, handleSplit)
}
