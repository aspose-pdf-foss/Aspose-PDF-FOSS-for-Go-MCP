// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"
	"os"
	"strings"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type renderParams struct {
	InputPath   string  `json:"input_path" jsonschema:"path to the input PDF file"`
	Page        int     `json:"page" jsonschema:"1-based page number to render"`
	OutputPath  string  `json:"output_path" jsonschema:"path to write the rendered image"`
	Password    string  `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
	DPI         float64 `json:"dpi,omitempty" jsonschema:"render resolution in DPI; defaults to 150"`
	Format      string  `json:"format,omitempty" jsonschema:"image format: png or jpeg; defaults to the output_path extension, else png"`
	JPEGQuality int     `json:"jpeg_quality,omitempty" jsonschema:"JPEG quality 1-100; defaults to 90"`
}

func handleRender(ctx context.Context, req *mcp.CallToolRequest, args renderParams) (*mcp.CallToolResult, any, error) {
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	page, err := doc.Page(args.Page)
	if err != nil {
		return nil, nil, fmt.Errorf("page %d of %q: %w", args.Page, args.InputPath, err)
	}

	format := strings.ToLower(args.Format)
	if format == "" {
		if strings.HasSuffix(strings.ToLower(args.OutputPath), ".jpg") ||
			strings.HasSuffix(strings.ToLower(args.OutputPath), ".jpeg") {
			format = "jpeg"
		} else {
			format = "png"
		}
	}

	f, err := os.Create(args.OutputPath)
	if err != nil {
		return nil, nil, fmt.Errorf("create %q: %w", args.OutputPath, err)
	}
	defer f.Close()

	opts := pdf.RenderOptions{DPI: args.DPI}
	switch format {
	case "jpeg", "jpg":
		quality := args.JPEGQuality
		if quality == 0 {
			quality = 90
		}
		err = page.RenderJPEG(f, opts, quality)
	case "png":
		err = page.RenderPNG(f, opts)
	default:
		return nil, nil, fmt.Errorf("unsupported format %q (use png or jpeg)", args.Format)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("render page %d: %w", args.Page, err)
	}

	msg := fmt.Sprintf("Rendered page %d of %s to %s (%s, DPI %.0f).",
		args.Page, args.InputPath, args.OutputPath, format, effectiveDPI(args.DPI))
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}

func effectiveDPI(dpi float64) float64 {
	if dpi == 0 {
		return pdf.DefaultDPI
	}
	return dpi
}
