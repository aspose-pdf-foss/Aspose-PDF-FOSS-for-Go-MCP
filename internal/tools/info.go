// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"encoding/json"
	"fmt"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type infoParams struct {
	InputPath string `json:"input_path" jsonschema:"path to the input PDF file"`
	Password  string `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
}

type pageSizeJSON struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type infoResult struct {
	PageCount    int            `json:"page_count"`
	PageSizes    []pageSizeJSON `json:"page_sizes"`
	Title        string         `json:"title"`
	Author       string         `json:"author"`
	Subject      string         `json:"subject"`
	Keywords     string         `json:"keywords"`
	Creator      string         `json:"creator"`
	Producer     string         `json:"producer"`
	CreationDate string         `json:"creation_date"`
	ModDate      string         `json:"mod_date"`
	Encrypted    bool           `json:"encrypted"`
}

func handleInfo(ctx context.Context, req *mcp.CallToolRequest, args infoParams) (*mcp.CallToolResult, any, error) {
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	info, err := doc.Info()
	if err != nil {
		return nil, nil, fmt.Errorf("read info %q: %w", args.InputPath, err)
	}
	sizes, err := pdf.PageSizes(args.InputPath)
	if err != nil {
		return nil, nil, fmt.Errorf("page sizes %q: %w", args.InputPath, err)
	}
	pageSizes := make([]pageSizeJSON, len(sizes))
	for i, s := range sizes {
		pageSizes[i] = pageSizeJSON{Width: s.Width, Height: s.Height}
	}
	out := infoResult{
		PageCount:    doc.PageCount(),
		PageSizes:    pageSizes,
		Title:        info.Title,
		Author:       info.Author,
		Subject:      info.Subject,
		Keywords:     info.Keywords,
		Creator:      info.Creator,
		Producer:     info.Producer,
		CreationDate: info.CreationDate,
		ModDate:      info.ModDate,
		Encrypted:    args.Password != "",
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("marshal info: %w", err)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(data)}}}, nil, nil
}
