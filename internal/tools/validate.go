// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type validateParams struct {
	InputPath string `json:"input_path" jsonschema:"path to the PDF file to validate"`
	Profile   string `json:"profile,omitempty" jsonschema:"validation profile: structural (default), pdfa-1b, pdfa-2b, pdfa-3b, pdfa-1a, pdfa-2a, pdfa-3a, or pdfua"`
	Password  string `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional; PDF/A and PDF/UA profiles only)"`
}

type validateIssueJSON struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type validateResult struct {
	Valid   bool                `json:"valid"`
	Profile string              `json:"profile"`
	Issues  []validateIssueJSON `json:"issues"`
}

func handleValidate(ctx context.Context, req *mcp.CallToolRequest, args validateParams) (*mcp.CallToolResult, any, error) {
	profile := strings.ToLower(strings.TrimSpace(args.Profile))
	if profile == "" {
		profile = "structural"
	}

	out := validateResult{Profile: profile, Issues: []validateIssueJSON{}}
	switch {
	case profile == "structural":
		report, err := pdf.Validate(args.InputPath)
		if err != nil {
			return nil, nil, fmt.Errorf("validate %q: %w", args.InputPath, err)
		}
		out.Valid = report.Valid
		for _, iss := range report.Issues {
			out.Issues = append(out.Issues, validateIssueJSON{Code: iss.Code, Message: iss.Message})
		}

	case profile == "pdfua":
		doc, err := openDocument(args.InputPath, args.Password)
		if err != nil {
			return nil, nil, err
		}
		report := doc.ValidatePDFUA()
		out.Valid = report.Conformant
		for _, iss := range report.Issues {
			out.Issues = append(out.Issues, validateIssueJSON{Code: iss.Rule, Message: iss.Message})
		}

	case strings.HasPrefix(profile, "pdfa-"):
		level, err := parsePDFALevel(profile)
		if err != nil {
			return nil, nil, err
		}
		doc, err := openDocument(args.InputPath, args.Password)
		if err != nil {
			return nil, nil, err
		}
		report := doc.ValidatePDFA(level)
		out.Valid = report.Conformant
		for _, iss := range report.Issues {
			out.Issues = append(out.Issues, validateIssueJSON{Code: iss.Rule, Message: iss.Message})
		}

	default:
		return nil, nil, fmt.Errorf("unknown profile %q (use structural, pdfa-1b, pdfa-2b, pdfa-3b, pdfa-1a, pdfa-2a, pdfa-3a, or pdfua)", args.Profile)
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("marshal report: %w", err)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(data)}}}, nil, nil
}
