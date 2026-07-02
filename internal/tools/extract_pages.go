// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"
	"strings"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type extractPagesParams struct {
	InputPath  string `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputPath string `json:"output_path" jsonschema:"path to write the PDF with only the selected pages"`
	Pages      string `json:"pages" jsonschema:"pages to extract, e.g. \"1-3,5\" (1-based); extracted in ascending order"`
	Password   string `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
}

func handleExtractPages(ctx context.Context, req *mcp.CallToolRequest, args extractPagesParams) (*mcp.CallToolResult, any, error) {
	if strings.TrimSpace(args.Pages) == "" {
		return nil, nil, fmt.Errorf("pages is required")
	}
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	pages, err := parsePageRange(args.Pages, doc.PageCount())
	if err != nil {
		return nil, nil, err
	}
	out, err := doc.Extract(toPageRanges(pages)...)
	if err != nil {
		return nil, nil, fmt.Errorf("extract pages from %q: %w", args.InputPath, err)
	}
	if err := out.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	msg := fmt.Sprintf("Extracted %d pages from %s: %s.", out.PageCount(), args.InputPath, args.OutputPath)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}

// toPageRanges compresses a sorted, de-duplicated page list into inclusive
// PageRange runs, e.g. [1 2 3 5] -> [{1,3} {5,5}].
func toPageRanges(pages []int) []pdf.PageRange {
	var ranges []pdf.PageRange
	for _, p := range pages {
		if n := len(ranges); n > 0 && ranges[n-1].To == p-1 {
			ranges[n-1].To = p
			continue
		}
		ranges = append(ranges, pdf.PageRange{From: p, To: p})
	}
	return ranges
}
