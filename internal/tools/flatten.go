// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type flattenParams struct {
	InputPath  string `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputPath string `json:"output_path" jsonschema:"path to write the flattened PDF"`
	Password   string `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
}

func handleFlatten(ctx context.Context, req *mcp.CallToolRequest, args flattenParams) (*mcp.CallToolResult, any, error) {
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	if err := doc.Flatten(); err != nil {
		return nil, nil, fmt.Errorf("flatten %q: %w", args.InputPath, err)
	}
	if err := doc.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	msg := fmt.Sprintf("Flattened form fields of %s into static content: %s.", args.InputPath, args.OutputPath)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}
