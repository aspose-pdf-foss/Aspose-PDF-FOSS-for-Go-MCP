// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type convertPDFAParams struct {
	InputPath  string `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputPath string `json:"output_path" jsonschema:"path to write the converted PDF"`
	Level      string `json:"level" jsonschema:"target PDF/A level: pdfa-1b, pdfa-2b, pdfa-3b, pdfa-1a, pdfa-2a, or pdfa-3a"`
	Password   string `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
}

type convertPDFAResult struct {
	Conformant bool                `json:"conformant"`
	Level      string              `json:"level"`
	OutputPath string              `json:"output_path"`
	Issues     []validateIssueJSON `json:"remaining_issues"`
}

func handleConvertPDFA(ctx context.Context, req *mcp.CallToolRequest, args convertPDFAParams) (*mcp.CallToolResult, any, error) {
	level, err := parsePDFALevel(args.Level)
	if err != nil {
		return nil, nil, err
	}
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	report, err := doc.ConvertToPDFA(level)
	if err != nil {
		return nil, nil, fmt.Errorf("convert %q to %s: %w", args.InputPath, level, err)
	}
	if err := doc.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	out := convertPDFAResult{
		Conformant: report.Conformant,
		Level:      report.Format.String(),
		OutputPath: args.OutputPath,
		Issues:     []validateIssueJSON{},
	}
	for _, iss := range report.Issues {
		out.Issues = append(out.Issues, validateIssueJSON{Code: iss.Rule, Message: iss.Message})
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("marshal report: %w", err)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(data)}}}, nil, nil
}
