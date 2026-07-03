// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type createParams struct {
	OutputPath string  `json:"output_path" jsonschema:"path to write the new PDF"`
	PageFormat string  `json:"page_format,omitempty" jsonschema:"page size: a4 (default), a3, letter, or legal"`
	Landscape  bool    `json:"landscape,omitempty" jsonschema:"landscape orientation (width and height swapped)"`
	PageCount  int     `json:"page_count,omitempty" jsonschema:"number of blank pages; defaults to 1"`
	Width      float64 `json:"width,omitempty" jsonschema:"custom page width in points (1/72 inch); overrides page_format when both width and height are given"`
	Height     float64 `json:"height,omitempty" jsonschema:"custom page height in points; overrides page_format when both width and height are given"`
}

func handleCreate(ctx context.Context, req *mcp.CallToolRequest, args createParams) (*mcp.CallToolResult, any, error) {
	count := args.PageCount
	if count == 0 {
		count = 1
	}
	if count < 1 || count > 10000 {
		return nil, nil, fmt.Errorf("page_count must be between 1 and 10000, got %d", args.PageCount)
	}

	var format pdf.PageFormat
	if args.Width > 0 && args.Height > 0 {
		format = pdf.PageFormat{Width: args.Width, Height: args.Height}
	} else if args.Width != 0 || args.Height != 0 {
		return nil, nil, fmt.Errorf("width and height must be given together and be positive")
	} else {
		f, err := parsePageFormat(args.PageFormat)
		if err != nil {
			return nil, nil, err
		}
		format = f
	}
	if args.Landscape {
		format = format.Landscape()
	}

	doc := pdf.NewDocumentFromFormat(format)
	for i := 1; i < count; i++ {
		if err := doc.AddBlankPageFromFormat(format); err != nil {
			return nil, nil, fmt.Errorf("add blank page: %w", err)
		}
	}
	if err := doc.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	msg := fmt.Sprintf("Created %s: %d blank page(s), %g x %g pt.", args.OutputPath, count, format.Width, format.Height)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}
