// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type addImageParams struct {
	InputPath  string  `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputPath string  `json:"output_path" jsonschema:"path to write the modified PDF (may equal input_path to overwrite)"`
	Password   string  `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
	Page       int     `json:"page,omitempty" jsonschema:"1-based page number; defaults to 1"`
	ImagePath  string  `json:"image_path" jsonschema:"path to a PNG or JPEG image file"`
	LLX        float64 `json:"llx" jsonschema:"rectangle lower-left x in points (PDF coordinates: origin at the page's bottom-left, y grows upward); the image is stretched to fill the rectangle, so match its aspect ratio"`
	LLY        float64 `json:"lly" jsonschema:"rectangle lower-left y in points"`
	URX        float64 `json:"urx" jsonschema:"rectangle upper-right x in points"`
	URY        float64 `json:"ury" jsonschema:"rectangle upper-right y in points"`
}

func handleAddImage(ctx context.Context, req *mcp.CallToolRequest, args addImageParams) (*mcp.CallToolResult, any, error) {
	if args.ImagePath == "" {
		return nil, nil, fmt.Errorf("image_path is required")
	}
	if args.URX <= args.LLX || args.URY <= args.LLY {
		return nil, nil, fmt.Errorf("invalid rectangle [%g %g %g %g]: urx must exceed llx and ury must exceed lly", args.LLX, args.LLY, args.URX, args.URY)
	}

	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	pageNum := args.Page
	if pageNum == 0 {
		pageNum = 1
	}
	page, err := doc.Page(pageNum)
	if err != nil {
		return nil, nil, err
	}
	rect := pdf.Rectangle{LLX: args.LLX, LLY: args.LLY, URX: args.URX, URY: args.URY}
	if err := page.AddImage(args.ImagePath, rect); err != nil {
		return nil, nil, fmt.Errorf("add image %q: %w", args.ImagePath, err)
	}
	if err := doc.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	msg := fmt.Sprintf("Added image %s to page %d of %s: %s.", args.ImagePath, pageNum, args.InputPath, args.OutputPath)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}
