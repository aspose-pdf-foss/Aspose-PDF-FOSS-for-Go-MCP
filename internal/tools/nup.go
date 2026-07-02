// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type nupParams struct {
	InputPath  string  `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputPath string  `json:"output_path" jsonschema:"path to write the imposed PDF"`
	Rows       int     `json:"rows" jsonschema:"grid rows per sheet (>= 1)"`
	Cols       int     `json:"cols" jsonschema:"grid columns per sheet (>= 1)"`
	Password   string  `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
	PageWidth  float64 `json:"page_width,omitempty" jsonschema:"output sheet width in points; defaults to A4 (both page_width and page_height must be given together)"`
	PageHeight float64 `json:"page_height,omitempty" jsonschema:"output sheet height in points; defaults to A4"`
	Margin     float64 `json:"margin,omitempty" jsonschema:"blank border around the grid, in points"`
	Gutter     float64 `json:"gutter,omitempty" jsonschema:"spacing between adjacent cells, in points"`
	DrawBorder bool    `json:"draw_border,omitempty" jsonschema:"draw a thin frame around each placed page"`
}

func handleNUp(ctx context.Context, req *mcp.CallToolRequest, args nupParams) (*mcp.CallToolResult, any, error) {
	if args.Rows < 1 || args.Cols < 1 {
		return nil, nil, fmt.Errorf("rows and cols must be >= 1, got %dx%d", args.Rows, args.Cols)
	}
	if (args.PageWidth == 0) != (args.PageHeight == 0) {
		return nil, nil, fmt.Errorf("page_width and page_height must be given together")
	}
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	opts := pdf.NUpOptions{
		Rows:       args.Rows,
		Cols:       args.Cols,
		PageSize:   pdf.PageFormat{Width: args.PageWidth, Height: args.PageHeight},
		Margin:     args.Margin,
		Gutter:     args.Gutter,
		DrawBorder: args.DrawBorder,
	}
	out, err := doc.NUp(opts)
	if err != nil {
		return nil, nil, fmt.Errorf("n-up %q: %w", args.InputPath, err)
	}
	if err := out.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	msg := fmt.Sprintf("Imposed %s %dx%d onto %d sheets: %s.", args.InputPath, args.Rows, args.Cols, out.PageCount(), args.OutputPath)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}
