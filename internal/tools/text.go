// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type textParams struct {
	InputPath string `json:"input_path" jsonschema:"path to the input PDF file"`
	Password  string `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
	Pages     string `json:"pages,omitempty" jsonschema:"page range to extract, e.g. \"1-3,5\" (1-based); empty means all pages"`
	Separator string `json:"separator,omitempty" jsonschema:"separator inserted between pages; defaults to a page header"`
}

func handleText(ctx context.Context, req *mcp.CallToolRequest, args textParams) (*mcp.CallToolResult, any, error) {
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	perPage, err := doc.ExtractText()
	if err != nil {
		return nil, nil, fmt.Errorf("extract text %q: %w", args.InputPath, err)
	}
	selected, err := parsePageRange(args.Pages, len(perPage))
	if err != nil {
		return nil, nil, err
	}

	var b strings.Builder
	for _, p := range selected {
		if args.Separator != "" {
			b.WriteString(args.Separator)
		} else {
			fmt.Fprintf(&b, "\n\n--- Page %d ---\n", p)
		}
		b.WriteString(perPage[p-1])
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: b.String()}}}, nil, nil
}
