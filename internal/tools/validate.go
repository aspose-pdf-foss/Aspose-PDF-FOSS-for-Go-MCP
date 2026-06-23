// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"encoding/json"
	"fmt"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type validateParams struct {
	InputPath string `json:"input_path" jsonschema:"path to the PDF file to validate"`
}

type validateIssueJSON struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type validateResult struct {
	Valid  bool                `json:"valid"`
	Issues []validateIssueJSON `json:"issues"`
}

func handleValidate(ctx context.Context, req *mcp.CallToolRequest, args validateParams) (*mcp.CallToolResult, any, error) {
	report, err := pdf.Validate(args.InputPath)
	if err != nil {
		return nil, nil, fmt.Errorf("validate %q: %w", args.InputPath, err)
	}
	out := validateResult{Valid: report.Valid, Issues: []validateIssueJSON{}}
	for _, iss := range report.Issues {
		out.Issues = append(out.Issues, validateIssueJSON{Code: iss.Code, Message: iss.Message})
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("marshal report: %w", err)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(data)}}}, nil, nil
}
