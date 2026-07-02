// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"
	"strings"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type rotateParams struct {
	InputPath  string `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputPath string `json:"output_path" jsonschema:"path to write the rotated PDF"`
	Angle      int    `json:"angle" jsonschema:"clockwise rotation in degrees: 90, 180, or 270 (added to any existing rotation)"`
	Pages      string `json:"pages,omitempty" jsonschema:"page range, e.g. \"1-3,5\" (1-based); omit for all pages"`
	Password   string `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
}

func handleRotate(ctx context.Context, req *mcp.CallToolRequest, args rotateParams) (*mcp.CallToolResult, any, error) {
	switch args.Angle {
	case 90, 180, 270:
	default:
		return nil, nil, fmt.Errorf("angle must be 90, 180, or 270, got %d", args.Angle)
	}
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	scope := "all pages"
	if strings.TrimSpace(args.Pages) == "" {
		err = doc.Rotate(pdf.RotationAngle(args.Angle))
	} else {
		pages, perr := parsePageRange(args.Pages, doc.PageCount())
		if perr != nil {
			return nil, nil, perr
		}
		err = doc.Rotate(pdf.RotationAngle(args.Angle), pages...)
		scope = fmt.Sprintf("%d pages", len(pages))
	}
	if err != nil {
		return nil, nil, fmt.Errorf("rotate %q: %w", args.InputPath, err)
	}
	if err := doc.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	msg := fmt.Sprintf("Rotated %s of %s by %d° clockwise: %s.", scope, args.InputPath, args.Angle, args.OutputPath)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}
