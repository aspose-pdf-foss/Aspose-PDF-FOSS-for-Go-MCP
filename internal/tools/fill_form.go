// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type fillFormParams struct {
	InputPath  string            `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputPath string            `json:"output_path" jsonschema:"path to write the filled PDF"`
	Values     map[string]string `json:"values" jsonschema:"field values keyed by full field name; checkbox values are On/Off"`
	Flatten    bool              `json:"flatten,omitempty" jsonschema:"bake the filled fields into static page content (no longer editable)"`
	Password   string            `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
}

func handleFillForm(ctx context.Context, req *mcp.CallToolRequest, args fillFormParams) (*mcp.CallToolResult, any, error) {
	if len(args.Values) == 0 {
		return nil, nil, fmt.Errorf("values must contain at least one field")
	}
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	form := doc.Form()

	names := make([]string, 0, len(args.Values))
	for name := range args.Values {
		names = append(names, name)
	}
	sort.Strings(names)

	var unknown []string
	for _, name := range names {
		if form.Field(name) == nil {
			unknown = append(unknown, name)
		}
	}
	if len(unknown) > 0 {
		return nil, nil, fmt.Errorf("no such form fields in %q: %s (use pdf_form_fields to list them)", args.InputPath, strings.Join(unknown, ", "))
	}

	for _, name := range names {
		if err := form.Field(name).SetValue(args.Values[name]); err != nil {
			return nil, nil, fmt.Errorf("set field %q: %w", name, err)
		}
	}
	if args.Flatten {
		if err := doc.Flatten(); err != nil {
			return nil, nil, fmt.Errorf("flatten %q: %w", args.InputPath, err)
		}
	}
	if err := doc.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	msg := fmt.Sprintf("Filled %d fields in %s: %s.", len(names), args.InputPath, args.OutputPath)
	if args.Flatten {
		msg = fmt.Sprintf("Filled %d fields in %s and flattened the form: %s.", len(names), args.InputPath, args.OutputPath)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}
