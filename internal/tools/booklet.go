// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type bookletParams struct {
	InputPath  string  `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputPath string  `json:"output_path" jsonschema:"path to write the booklet PDF"`
	Password   string  `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
	PageWidth  float64 `json:"page_width,omitempty" jsonschema:"output spread width in points; defaults to twice the source page width (both page_width and page_height must be given together)"`
	PageHeight float64 `json:"page_height,omitempty" jsonschema:"output spread height in points; defaults to the source page height"`
}

func handleBooklet(ctx context.Context, req *mcp.CallToolRequest, args bookletParams) (*mcp.CallToolResult, any, error) {
	if (args.PageWidth == 0) != (args.PageHeight == 0) {
		return nil, nil, fmt.Errorf("page_width and page_height must be given together")
	}
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	opts := pdf.BookletOptions{
		PageSize: pdf.PageFormat{Width: args.PageWidth, Height: args.PageHeight},
	}
	out, err := doc.Booklet(opts)
	if err != nil {
		return nil, nil, fmt.Errorf("booklet %q: %w", args.InputPath, err)
	}
	if err := out.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	msg := fmt.Sprintf("Imposed %s as a saddle-stitch booklet of %d spreads: %s.", args.InputPath, out.PageCount(), args.OutputPath)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}
