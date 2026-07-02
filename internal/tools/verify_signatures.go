// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type verifySignaturesParams struct {
	InputPath string `json:"input_path" jsonschema:"path to the PDF file to verify"`
	Password  string `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
}

type signatureJSON struct {
	FieldName           string `json:"field_name"`
	SignerName          string `json:"signer_name,omitempty"`
	Reason              string `json:"reason,omitempty"`
	Location            string `json:"location,omitempty"`
	SigningTime         string `json:"signing_time,omitempty"`
	Valid               bool   `json:"valid"`
	IntegrityOK         bool   `json:"integrity_ok"`
	CoversWholeDocument bool   `json:"covers_whole_document"`
	CertificateSubject  string `json:"certificate_subject,omitempty"`
	Error               string `json:"error,omitempty"`
}

type verifySignaturesResult struct {
	Count      int             `json:"count"`
	Signatures []signatureJSON `json:"signatures"`
}

func handleVerifySignatures(ctx context.Context, req *mcp.CallToolRequest, args verifySignaturesParams) (*mcp.CallToolResult, any, error) {
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	results, err := doc.VerifySignatures()
	if err != nil {
		return nil, nil, fmt.Errorf("verify signatures in %q: %w", args.InputPath, err)
	}
	out := verifySignaturesResult{Count: len(results), Signatures: []signatureJSON{}}
	for _, r := range results {
		s := signatureJSON{
			FieldName:           r.FieldName,
			SignerName:          r.SignerName,
			Reason:              r.Reason,
			Location:            r.Location,
			Valid:               r.Valid,
			IntegrityOK:         r.IntegrityOK,
			CoversWholeDocument: r.CoversWholeDocument,
		}
		if !r.SigningTime.IsZero() {
			s.SigningTime = r.SigningTime.Format(time.RFC3339)
		}
		if r.Certificate != nil {
			s.CertificateSubject = r.Certificate.Subject.String()
		}
		if r.Err != nil {
			s.Error = r.Err.Error()
		}
		out.Signatures = append(out.Signatures, s)
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("marshal report: %w", err)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(data)}}}, nil, nil
}
