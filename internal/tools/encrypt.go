// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"
	"strings"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type encryptParams struct {
	InputPath     string `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputPath    string `json:"output_path" jsonschema:"path to write the encrypted PDF"`
	UserPassword  string `json:"user_password,omitempty" jsonschema:"password required to open the document (at least one of user_password/owner_password is required)"`
	OwnerPassword string `json:"owner_password,omitempty" jsonschema:"password that unlocks all permissions; defaults to user_password"`
	Algorithm     string `json:"algorithm,omitempty" jsonschema:"encryption algorithm: aes-128 (default), aes-256, or rc4-128"`
	Password      string `json:"password,omitempty" jsonschema:"password for an already-encrypted input PDF (optional)"`
}

func handleEncrypt(ctx context.Context, req *mcp.CallToolRequest, args encryptParams) (*mcp.CallToolResult, any, error) {
	if args.UserPassword == "" && args.OwnerPassword == "" {
		return nil, nil, fmt.Errorf("at least one of user_password or owner_password is required")
	}
	var alg pdf.EncryptionAlgorithm
	switch strings.ToLower(strings.TrimSpace(args.Algorithm)) {
	case "", "aes-128":
		alg = pdf.EncryptionAlgAES128
	case "aes-256":
		alg = pdf.EncryptionAlgAES256
	case "rc4-128":
		alg = pdf.EncryptionAlgRC4_128
	default:
		return nil, nil, fmt.Errorf("unknown algorithm %q (use aes-128, aes-256, or rc4-128)", args.Algorithm)
	}
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	doc.SetEncryption(pdf.EncryptionOptions{
		UserPassword:  args.UserPassword,
		OwnerPassword: args.OwnerPassword,
		Algorithm:     alg,
	})
	if err := doc.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	msg := fmt.Sprintf("Encrypted %s: %s.", args.InputPath, args.OutputPath)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}
