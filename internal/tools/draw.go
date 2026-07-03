// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"
	"strings"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type drawParams struct {
	InputPath   string    `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputPath  string    `json:"output_path" jsonschema:"path to write the modified PDF (may equal input_path to overwrite)"`
	Password    string    `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
	Page        int       `json:"page,omitempty" jsonschema:"1-based page number; defaults to 1"`
	Shape       string    `json:"shape" jsonschema:"one of: line, rectangle, circle, ellipse. All coordinates are points in PDF coordinates (origin at the page's bottom-left, y grows upward)"`
	X1          float64   `json:"x1,omitempty" jsonschema:"line start x (shape=line)"`
	Y1          float64   `json:"y1,omitempty" jsonschema:"line start y (shape=line)"`
	X2          float64   `json:"x2,omitempty" jsonschema:"line end x (shape=line)"`
	Y2          float64   `json:"y2,omitempty" jsonschema:"line end y (shape=line)"`
	LLX         float64   `json:"llx,omitempty" jsonschema:"rectangle lower-left x (shape=rectangle)"`
	LLY         float64   `json:"lly,omitempty" jsonschema:"rectangle lower-left y (shape=rectangle)"`
	URX         float64   `json:"urx,omitempty" jsonschema:"rectangle upper-right x (shape=rectangle)"`
	URY         float64   `json:"ury,omitempty" jsonschema:"rectangle upper-right y (shape=rectangle)"`
	CX          float64   `json:"cx,omitempty" jsonschema:"center x (shape=circle or ellipse)"`
	CY          float64   `json:"cy,omitempty" jsonschema:"center y (shape=circle or ellipse)"`
	R           float64   `json:"r,omitempty" jsonschema:"radius (shape=circle)"`
	RX          float64   `json:"rx,omitempty" jsonschema:"horizontal semi-axis (shape=ellipse)"`
	RY          float64   `json:"ry,omitempty" jsonschema:"vertical semi-axis (shape=ellipse)"`
	StrokeColor string    `json:"stroke_color,omitempty" jsonschema:"stroke color as #RRGGBB or #RRGGBBAA; defaults to black"`
	StrokeWidth float64   `json:"stroke_width,omitempty" jsonschema:"stroke width in points; defaults to 1"`
	NoStroke    bool      `json:"no_stroke,omitempty" jsonschema:"draw no outline (fill-only; not valid for shape=line)"`
	FillColor   string    `json:"fill_color,omitempty" jsonschema:"fill color as #RRGGBB or #RRGGBBAA; omit for no fill (not valid for shape=line)"`
	DashPattern []float64 `json:"dash_pattern,omitempty" jsonschema:"dash pattern as alternating on/off lengths in points, e.g. [3, 2]; omit for a solid line"`
}

func handleDraw(ctx context.Context, req *mcp.CallToolRequest, args drawParams) (*mcp.CallToolResult, any, error) {
	shape := strings.ToLower(args.Shape)
	strokeColor, err := parseHexColor(args.StrokeColor)
	if err != nil {
		return nil, nil, err
	}
	fillColor, err := parseHexColor(args.FillColor)
	if err != nil {
		return nil, nil, err
	}
	strokeWidth := args.StrokeWidth
	if strokeWidth == 0 {
		strokeWidth = 1
	}
	if strokeWidth < 0 {
		return nil, nil, fmt.Errorf("stroke_width must be positive, got %g", args.StrokeWidth)
	}
	if args.NoStroke {
		if shape == "line" {
			return nil, nil, fmt.Errorf("no_stroke is not valid for shape=line")
		}
		if fillColor == nil {
			return nil, nil, fmt.Errorf("no_stroke without fill_color would draw nothing")
		}
		strokeWidth = 0
	}
	line := pdf.LineStyle{Color: strokeColor, Width: strokeWidth, DashPattern: args.DashPattern}
	style := pdf.ShapeStyle{LineStyle: line, FillColor: fillColor}

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

	switch shape {
	case "line":
		if fillColor != nil {
			return nil, nil, fmt.Errorf("fill_color is not valid for shape=line")
		}
		if args.X1 == args.X2 && args.Y1 == args.Y2 {
			return nil, nil, fmt.Errorf("line endpoints coincide: (%g, %g)", args.X1, args.Y1)
		}
		err = page.DrawLine(pdf.Point{X: args.X1, Y: args.Y1}, pdf.Point{X: args.X2, Y: args.Y2}, line)
	case "rectangle":
		if args.URX <= args.LLX || args.URY <= args.LLY {
			return nil, nil, fmt.Errorf("invalid rectangle [%g %g %g %g]: urx must exceed llx and ury must exceed lly", args.LLX, args.LLY, args.URX, args.URY)
		}
		err = page.DrawRectangle(pdf.Rectangle{LLX: args.LLX, LLY: args.LLY, URX: args.URX, URY: args.URY}, style)
	case "circle":
		if args.R <= 0 {
			return nil, nil, fmt.Errorf("r must be positive for shape=circle, got %g", args.R)
		}
		err = page.DrawCircle(pdf.Point{X: args.CX, Y: args.CY}, args.R, style)
	case "ellipse":
		if args.RX <= 0 || args.RY <= 0 {
			return nil, nil, fmt.Errorf("rx and ry must be positive for shape=ellipse, got rx=%g ry=%g", args.RX, args.RY)
		}
		err = page.DrawEllipse(pdf.Point{X: args.CX, Y: args.CY}, args.RX, args.RY, style)
	default:
		return nil, nil, fmt.Errorf("shape %q: want line, rectangle, circle, or ellipse", args.Shape)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("draw %s: %w", shape, err)
	}

	if err := doc.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	msg := fmt.Sprintf("Drew a %s on page %d of %s: %s.", shape, pageNum, args.InputPath, args.OutputPath)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}
