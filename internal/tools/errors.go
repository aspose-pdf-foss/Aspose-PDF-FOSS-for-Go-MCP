// SPDX-License-Identifier: MIT

// Package tools implements the PDF MCP tool handlers.
package tools

import (
	"errors"
	"fmt"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
)

// openDocument opens a PDF by path. If password is non-empty it uses
// OpenWithPassword (which also works on unencrypted files); otherwise it uses
// plain Open and turns ErrEncrypted into an actionable message.
func openDocument(path, password string) (*pdf.Document, error) {
	if password != "" {
		doc, err := pdf.OpenWithPassword(path, password)
		if err != nil {
			return nil, fmt.Errorf("open %q: %w", path, err)
		}
		return doc, nil
	}
	doc, err := pdf.Open(path)
	if err != nil {
		if errors.Is(err, pdf.ErrEncrypted) {
			return nil, fmt.Errorf("open %q: PDF is encrypted — supply the \"password\" parameter", path)
		}
		return nil, fmt.Errorf("open %q: %w", path, err)
	}
	return doc, nil
}
