// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type linearizeParams struct {
	InputPath  string `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputPath string `json:"output_path" jsonschema:"path to write the linearized (fast web view) PDF"`
	Password   string `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional; the output is not encrypted)"`
}

func handleLinearize(ctx context.Context, req *mcp.CallToolRequest, args linearizeParams) (*mcp.CallToolResult, any, error) {
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	if err := doc.SaveLinearized(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("linearize %q: %w", args.InputPath, err)
	}
	msg := fmt.Sprintf("Linearized %s for fast web view: %s (%d pages).", args.InputPath, args.OutputPath, doc.PageCount())
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}
