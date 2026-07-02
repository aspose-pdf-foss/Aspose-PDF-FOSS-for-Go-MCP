// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"
	"strings"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type watermarkParams struct {
	InputPath   string  `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputPath  string  `json:"output_path" jsonschema:"path to write the watermarked PDF"`
	Text        string  `json:"text,omitempty" jsonschema:"watermark text (exactly one of text or image_path is required)"`
	ImagePath   string  `json:"image_path,omitempty" jsonschema:"path to a PNG or JPEG watermark image (exactly one of text or image_path is required)"`
	Password    string  `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
	Opacity     float64 `json:"opacity,omitempty" jsonschema:"watermark opacity in (0, 1]; defaults to 0.3"`
	RotateAngle float64 `json:"rotate_angle,omitempty" jsonschema:"rotation in degrees counter-clockwise about the page centre; e.g. 45 for a diagonal watermark"`
	FontSize    float64 `json:"font_size,omitempty" jsonschema:"text watermark font size in points; defaults to 48"`
	Behind      bool    `json:"behind,omitempty" jsonschema:"draw the watermark behind the page content instead of on top"`
	Pages       string  `json:"pages,omitempty" jsonschema:"page range, e.g. \"1-3,5\" (1-based); omit for all pages"`
}

func handleWatermark(ctx context.Context, req *mcp.CallToolRequest, args watermarkParams) (*mcp.CallToolResult, any, error) {
	if (args.Text == "") == (args.ImagePath == "") {
		return nil, nil, fmt.Errorf("exactly one of text or image_path is required")
	}
	opacity := args.Opacity
	if opacity == 0 {
		opacity = 0.3
	}
	if opacity < 0 || opacity > 1 {
		return nil, nil, fmt.Errorf("opacity must be in (0, 1], got %g", args.Opacity)
	}

	var stamp pdf.Stamp
	var kind string
	if args.Text != "" {
		kind = fmt.Sprintf("text %q", args.Text)
		fontSize := args.FontSize
		if fontSize == 0 {
			fontSize = 48
		}
		ts := pdf.NewTextStamp(args.Text, pdf.TextStyle{Size: fontSize})
		ts.HAlign = pdf.HAlignCenter
		ts.VAlign = pdf.VAlignMiddle
		ts.Opacity = opacity
		ts.RotateAngle = args.RotateAngle
		ts.Background = args.Behind
		stamp = ts
	} else {
		kind = fmt.Sprintf("image %s", args.ImagePath)
		is, err := pdf.NewImageStamp(args.ImagePath)
		if err != nil {
			return nil, nil, fmt.Errorf("load watermark image %q: %w", args.ImagePath, err)
		}
		is.HAlign = pdf.HAlignCenter
		is.VAlign = pdf.VAlignMiddle
		is.Opacity = opacity
		is.RotateAngle = args.RotateAngle
		is.Background = args.Behind
		stamp = is
	}

	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	scope := "all pages"
	if strings.TrimSpace(args.Pages) == "" {
		err = doc.AddStamp(stamp)
	} else {
		pages, perr := parsePageRange(args.Pages, doc.PageCount())
		if perr != nil {
			return nil, nil, perr
		}
		err = doc.AddStamp(stamp, pages...)
		scope = fmt.Sprintf("%d pages", len(pages))
	}
	if err != nil {
		return nil, nil, fmt.Errorf("watermark %q: %w", args.InputPath, err)
	}
	if err := doc.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	msg := fmt.Sprintf("Added %s watermark to %s of %s: %s.", kind, scope, args.InputPath, args.OutputPath)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}
