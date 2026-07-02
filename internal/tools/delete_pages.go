// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type deletePagesParams struct {
	InputPath  string `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputPath string `json:"output_path" jsonschema:"path to write the PDF without the deleted pages"`
	Pages      string `json:"pages" jsonschema:"pages to delete, e.g. \"1-3,5\" (1-based); at least one page must remain"`
	Password   string `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
}

func handleDeletePages(ctx context.Context, req *mcp.CallToolRequest, args deletePagesParams) (*mcp.CallToolResult, any, error) {
	if strings.TrimSpace(args.Pages) == "" {
		return nil, nil, fmt.Errorf("pages is required")
	}
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	pages, err := parsePageRange(args.Pages, doc.PageCount())
	if err != nil {
		return nil, nil, err
	}
	if err := doc.DeletePages(pages...); err != nil {
		return nil, nil, fmt.Errorf("delete pages from %q: %w", args.InputPath, err)
	}
	if err := doc.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	msg := fmt.Sprintf("Deleted %d pages from %s (%d pages remain): %s.", len(pages), args.InputPath, doc.PageCount(), args.OutputPath)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}
