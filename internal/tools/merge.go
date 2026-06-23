// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type mergeParams struct {
	InputPaths []string `json:"input_paths" jsonschema:"paths of the PDFs to merge, in order (at least two)"`
	OutputPath string   `json:"output_path" jsonschema:"path to write the merged PDF"`
}

func handleMerge(ctx context.Context, req *mcp.CallToolRequest, args mergeParams) (*mcp.CallToolResult, any, error) {
	if len(args.InputPaths) < 2 {
		return nil, nil, fmt.Errorf("pdf_merge needs at least two input_paths, got %d", len(args.InputPaths))
	}
	base, err := openDocument(args.InputPaths[0], "")
	if err != nil {
		return nil, nil, err
	}
	rest := make([]*pdf.Document, 0, len(args.InputPaths)-1)
	for _, p := range args.InputPaths[1:] {
		doc, err := openDocument(p, "")
		if err != nil {
			return nil, nil, err
		}
		rest = append(rest, doc)
	}
	base.Append(rest...)
	if err := base.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	msg := fmt.Sprintf("Merged %d PDFs into %s (%d pages).", len(args.InputPaths), args.OutputPath, base.PageCount())
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}
