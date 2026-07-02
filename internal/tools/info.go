// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"encoding/json"
	"errors"
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
	pageSizes := make([]pageSizeJSON, 0, doc.PageCount())
	for i, p := range doc.Pages() {
		size, err := p.Size()
		if err != nil {
			return nil, nil, fmt.Errorf("size of page %d in %q: %w", i+1, args.InputPath, err)
		}
		pageSizes = append(pageSizes, pageSizeJSON{Width: size.Width, Height: size.Height})
	}
	// The document is already open (possibly via a password), so probe the
	// file itself for the encrypted flag.
	encrypted := false
	if args.Password != "" {
		if _, err := pdf.Open(args.InputPath); errors.Is(err, pdf.ErrEncrypted) {
			encrypted = true
		}
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
		Encrypted:    encrypted,
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("marshal info: %w", err)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(data)}}}, nil, nil
}
