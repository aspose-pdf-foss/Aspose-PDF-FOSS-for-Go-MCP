// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type addTextParams struct {
	InputPath  string  `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputPath string  `json:"output_path" jsonschema:"path to write the modified PDF (may equal input_path to overwrite)"`
	Password   string  `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
	Page       int     `json:"page,omitempty" jsonschema:"1-based page number; defaults to 1"`
	Text       string  `json:"text" jsonschema:"text to draw; word-wraps inside the rectangle"`
	LLX        float64 `json:"llx,omitempty" jsonschema:"rectangle lower-left x in points (PDF coordinates: origin at the page's bottom-left, y grows upward); omit all four coordinates to use the page minus a 36 pt margin"`
	LLY        float64 `json:"lly,omitempty" jsonschema:"rectangle lower-left y in points"`
	URX        float64 `json:"urx,omitempty" jsonschema:"rectangle upper-right x in points"`
	URY        float64 `json:"ury,omitempty" jsonschema:"rectangle upper-right y in points"`
	Font       string  `json:"font,omitempty" jsonschema:"standard-14 font name: helvetica (default), helvetica-bold, helvetica-oblique, times-roman, times-bold, times-italic, courier, courier-bold, symbol, zapfdingbats, etc."`
	FontSize   float64 `json:"font_size,omitempty" jsonschema:"font size in points; defaults to 12"`
	Color      string  `json:"color,omitempty" jsonschema:"text color as #RRGGBB or #RRGGBBAA; defaults to black"`
	HAlign     string  `json:"halign,omitempty" jsonschema:"horizontal alignment inside the rectangle: left (default), center, or right"`
	VAlign     string  `json:"valign,omitempty" jsonschema:"vertical alignment inside the rectangle: top (default), middle, or bottom"`
	Rotation   float64 `json:"rotation,omitempty" jsonschema:"rotation in degrees counter-clockwise about the rectangle's lower-left corner"`
}

func handleAddText(ctx context.Context, req *mcp.CallToolRequest, args addTextParams) (*mcp.CallToolResult, any, error) {
	if args.Text == "" {
		return nil, nil, fmt.Errorf("text is required")
	}
	font, err := resolveFont(args.Font)
	if err != nil {
		return nil, nil, err
	}
	color, err := parseHexColor(args.Color)
	if err != nil {
		return nil, nil, err
	}
	halign, err := parseHAlign(args.HAlign)
	if err != nil {
		return nil, nil, err
	}
	valign, err := parseVAlign(args.VAlign)
	if err != nil {
		return nil, nil, err
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
	rect, err := contentRect(page, args.LLX, args.LLY, args.URX, args.URY)
	if err != nil {
		return nil, nil, err
	}

	style := pdf.TextStyle{
		Font:     font,
		Size:     args.FontSize,
		Color:    color,
		HAlign:   halign,
		VAlign:   valign,
		Rotation: args.Rotation,
	}
	if err := page.AddText(args.Text, style, rect); err != nil {
		return nil, nil, fmt.Errorf("add text: %w", err)
	}
	if err := doc.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	msg := fmt.Sprintf("Added text (%d chars) to page %d of %s: %s.", len(args.Text), pageNum, args.InputPath, args.OutputPath)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}
