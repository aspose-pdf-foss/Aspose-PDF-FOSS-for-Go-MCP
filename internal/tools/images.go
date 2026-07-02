// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type extractImagesParams struct {
	InputPath string `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputDir string `json:"output_dir" jsonschema:"directory to write the extracted images (image_p001_01.png, ...)"`
	Password  string `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
}

func handleExtractImages(ctx context.Context, req *mcp.CallToolRequest, args extractImagesParams) (*mcp.CallToolResult, any, error) {
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	pages, err := doc.ExtractImages()
	if err != nil {
		return nil, nil, fmt.Errorf("extract images from %q: %w", args.InputPath, err)
	}
	total := 0
	for _, imgs := range pages {
		total += len(imgs)
	}
	if total == 0 {
		msg := fmt.Sprintf("No images found in %s.", args.InputPath)
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
	}
	if err := os.MkdirAll(args.OutputDir, 0o755); err != nil {
		return nil, nil, fmt.Errorf("create output dir %q: %w", args.OutputDir, err)
	}
	written := make([]string, 0, total)
	for i, imgs := range pages {
		for j, img := range imgs {
			ext := ".png"
			if img.Format == pdf.ImageFormatJPEG {
				ext = ".jpg"
			}
			name := filepath.Join(args.OutputDir, fmt.Sprintf("image_p%03d_%02d%s", i+1, j+1, ext))
			if err := img.Save(name); err != nil {
				return nil, nil, fmt.Errorf("save %q: %w", name, err)
			}
			written = append(written, name)
		}
	}
	msg := fmt.Sprintf("Extracted %d images from %s:\n%s", len(written), args.InputPath, strings.Join(written, "\n"))
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}
