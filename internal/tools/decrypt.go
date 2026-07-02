// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type decryptParams struct {
	InputPath  string `json:"input_path" jsonschema:"path to the encrypted input PDF file"`
	OutputPath string `json:"output_path" jsonschema:"path to write the decrypted PDF"`
	Password   string `json:"password" jsonschema:"password that opens the encrypted PDF"`
}

func handleDecrypt(ctx context.Context, req *mcp.CallToolRequest, args decryptParams) (*mcp.CallToolResult, any, error) {
	if args.Password == "" {
		return nil, nil, fmt.Errorf("password is required")
	}
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	doc.RemoveEncryption()
	if err := doc.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	msg := fmt.Sprintf("Decrypted %s: %s.", args.InputPath, args.OutputPath)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}
