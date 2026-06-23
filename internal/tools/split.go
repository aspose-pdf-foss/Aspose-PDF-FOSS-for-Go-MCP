// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type splitParams struct {
	InputPath string `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputDir string `json:"output_dir" jsonschema:"directory to write one PDF per page (page_001.pdf, ...)"`
	Password  string `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
}

func handleSplit(ctx context.Context, req *mcp.CallToolRequest, args splitParams) (*mcp.CallToolResult, any, error) {
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	pages, err := doc.Split()
	if err != nil {
		return nil, nil, fmt.Errorf("split %q: %w", args.InputPath, err)
	}
	if err := os.MkdirAll(args.OutputDir, 0o755); err != nil {
		return nil, nil, fmt.Errorf("create output dir %q: %w", args.OutputDir, err)
	}
	written := make([]string, 0, len(pages))
	for i, p := range pages {
		name := filepath.Join(args.OutputDir, fmt.Sprintf("page_%03d.pdf", i+1))
		if err := p.Save(name); err != nil {
			return nil, nil, fmt.Errorf("save %q: %w", name, err)
		}
		written = append(written, name)
	}
	msg := fmt.Sprintf("Split %s into %d files:\n%s", args.InputPath, len(written), strings.Join(written, "\n"))
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}
