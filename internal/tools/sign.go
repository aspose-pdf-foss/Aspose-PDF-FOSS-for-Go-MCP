// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type signParams struct {
	InputPath   string    `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputPath  string    `json:"output_path" jsonschema:"path to write the signed PDF"`
	CertPath    string    `json:"cert_path" jsonschema:"path to the signer certificate in PEM format (extra CERTIFICATE blocks are embedded as the chain)"`
	KeyPath     string    `json:"key_path" jsonschema:"path to the private key in PEM format (PKCS#8, PKCS#1 RSA, or SEC1 EC)"`
	Password    string    `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
	Reason      string    `json:"reason,omitempty" jsonschema:"signing reason shown by viewers (optional)"`
	Location    string    `json:"location,omitempty" jsonschema:"signing location shown by viewers (optional)"`
	ContactInfo string    `json:"contact_info,omitempty" jsonschema:"signer contact info (optional)"`
	SignerName  string    `json:"signer_name,omitempty" jsonschema:"signer name shown by viewers (optional)"`
	Visible     bool      `json:"visible,omitempty" jsonschema:"draw a visible signature widget (requires rect)"`
	Page        int       `json:"page,omitempty" jsonschema:"1-based page for the visible widget; defaults to 1"`
	Rect        []float64 `json:"rect,omitempty" jsonschema:"visible widget rectangle [llx, lly, urx, ury] in points (required when visible)"`
}

func handleSign(ctx context.Context, req *mcp.CallToolRequest, args signParams) (*mcp.CallToolResult, any, error) {
	cert, chain, err := loadCertificatePEM(args.CertPath)
	if err != nil {
		return nil, nil, err
	}
	key, err := loadPrivateKeyPEM(args.KeyPath)
	if err != nil {
		return nil, nil, err
	}
	opts := pdf.SignOptions{
		Certificate: cert,
		PrivateKey:  key,
		Chain:       chain,
		Reason:      args.Reason,
		Location:    args.Location,
		ContactInfo: args.ContactInfo,
		Name:        args.SignerName,
		Visible:     args.Visible,
		Page:        args.Page,
	}
	if args.Visible {
		if len(args.Rect) != 4 {
			return nil, nil, fmt.Errorf("a visible signature requires rect [llx, lly, urx, ury], got %d values", len(args.Rect))
		}
		opts.Rect = pdf.Rectangle{LLX: args.Rect[0], LLY: args.Rect[1], URX: args.Rect[2], URY: args.Rect[3]}
	}
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	if err := doc.Sign(opts); err != nil {
		return nil, nil, fmt.Errorf("sign %q: %w", args.InputPath, err)
	}
	if err := doc.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	msg := fmt.Sprintf("Signed %s as %q: %s.", args.InputPath, cert.Subject.CommonName, args.OutputPath)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}

// loadCertificatePEM reads a PEM file and returns the first CERTIFICATE block
// as the signer certificate and any further CERTIFICATE blocks as the chain.
func loadCertificatePEM(path string) (*x509.Certificate, []*x509.Certificate, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("read certificate %q: %w", path, err)
	}
	var certs []*x509.Certificate
	for len(data) > 0 {
		var block *pem.Block
		block, data = pem.Decode(data)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" {
			continue
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, nil, fmt.Errorf("parse certificate %q: %w", path, err)
		}
		certs = append(certs, cert)
	}
	if len(certs) == 0 {
		return nil, nil, fmt.Errorf("no CERTIFICATE block found in %q", path)
	}
	return certs[0], certs[1:], nil
}

// loadPrivateKeyPEM reads a PEM private key (PKCS#8, PKCS#1 RSA, or SEC1 EC)
// and returns it as a crypto.Signer.
func loadPrivateKeyPEM(path string) (crypto.Signer, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read key %q: %w", path, err)
	}
	for len(data) > 0 {
		var block *pem.Block
		block, data = pem.Decode(data)
		if block == nil {
			break
		}
		var key any
		var parseErr error
		switch block.Type {
		case "PRIVATE KEY":
			key, parseErr = x509.ParsePKCS8PrivateKey(block.Bytes)
		case "RSA PRIVATE KEY":
			key, parseErr = x509.ParsePKCS1PrivateKey(block.Bytes)
		case "EC PRIVATE KEY":
			key, parseErr = x509.ParseECPrivateKey(block.Bytes)
		default:
			continue
		}
		if parseErr != nil {
			return nil, fmt.Errorf("parse key %q: %w", path, parseErr)
		}
		signer, ok := key.(crypto.Signer)
		if !ok {
			return nil, fmt.Errorf("key %q is not a signing key (%T)", path, key)
		}
		return signer, nil
	}
	return nil, fmt.Errorf("no private key block found in %q", path)
}
