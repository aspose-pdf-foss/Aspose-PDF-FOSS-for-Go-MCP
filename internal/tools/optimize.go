// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type optimizeParams struct {
	InputPath        string  `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputPath       string  `json:"output_path" jsonschema:"path to write the optimized PDF"`
	Password         string  `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
	MaxDPI           float64 `json:"max_dpi,omitempty" jsonschema:"downscale images above this effective DPI; 0 = no downscaling"`
	JPEGQuality      int     `json:"jpeg_quality,omitempty" jsonschema:"JPEG re-encode quality 1-100; defaults to 75"`
	ConvertPNGToJPEG bool    `json:"convert_png_to_jpeg,omitempty" jsonschema:"convert opaque PNG images (no alpha) to JPEG"`
}

type optimizeResult struct {
	ImagesOptimized int    `json:"images_optimized"`
	ObjectsRemoved  int    `json:"objects_removed"`
	BytesBefore     int64  `json:"bytes_before"`
	BytesAfter      int64  `json:"bytes_after"`
	OutputPath      string `json:"output_path"`
}

func handleOptimize(ctx context.Context, req *mcp.CallToolRequest, args optimizeParams) (*mcp.CallToolResult, any, error) {
	before, err := os.Stat(args.InputPath)
	if err != nil {
		return nil, nil, fmt.Errorf("stat %q: %w", args.InputPath, err)
	}
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	optimized, err := doc.OptimizeImages(pdf.OptimizeImageOptions{
		MaxDPI:           args.MaxDPI,
		JPEGQuality:      args.JPEGQuality,
		ConvertPNGToJPEG: args.ConvertPNGToJPEG,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("optimize images in %q: %w", args.InputPath, err)
	}
	removed := doc.RemoveUnusedObjects()
	if err := doc.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	after, err := os.Stat(args.OutputPath)
	if err != nil {
		return nil, nil, fmt.Errorf("stat %q: %w", args.OutputPath, err)
	}
	out := optimizeResult{
		ImagesOptimized: optimized,
		ObjectsRemoved:  removed,
		BytesBefore:     before.Size(),
		BytesAfter:      after.Size(),
		OutputPath:      args.OutputPath,
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("marshal result: %w", err)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(data)}}}, nil, nil
}
