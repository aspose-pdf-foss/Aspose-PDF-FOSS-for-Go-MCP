// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"encoding/json"
	"fmt"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type searchParams struct {
	InputPath       string `json:"input_path" jsonschema:"path to the input PDF file"`
	Query           string `json:"query" jsonschema:"text to search for"`
	Password        string `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
	CaseInsensitive bool   `json:"case_insensitive,omitempty" jsonschema:"fold case when matching"`
	Regex           bool   `json:"regex,omitempty" jsonschema:"treat query as an RE2 regular expression"`
}

type searchMatchJSON struct {
	Text string     `json:"text"`
	Page int        `json:"page"`
	Rect [4]float64 `json:"rect"`
}

type searchResult struct {
	Count   int               `json:"count"`
	Matches []searchMatchJSON `json:"matches"`
}

func handleSearch(ctx context.Context, req *mcp.CallToolRequest, args searchParams) (*mcp.CallToolResult, any, error) {
	if args.Query == "" {
		return nil, nil, fmt.Errorf("query is required")
	}
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	matches, err := doc.SearchText(args.Query, pdf.SearchOptions{
		CaseInsensitive: args.CaseInsensitive,
		Regex:           args.Regex,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("search %q in %q: %w", args.Query, args.InputPath, err)
	}
	out := searchResult{Count: len(matches), Matches: []searchMatchJSON{}}
	for _, m := range matches {
		out.Matches = append(out.Matches, searchMatchJSON{
			Text: m.Text,
			Page: m.PageNumber,
			Rect: [4]float64{m.Rect.LLX, m.Rect.LLY, m.Rect.URX, m.Rect.URY},
		})
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("marshal matches: %w", err)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(data)}}}, nil, nil
}
