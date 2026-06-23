// SPDX-License-Identifier: MIT

package tools

import (
	"strings"
	"testing"
)

func TestOpenDocumentPlain(t *testing.T) {
	doc, err := openDocument("../../testdata/4pages.pdf", "")
	if err != nil {
		t.Fatalf("openDocument: %v", err)
	}
	if doc.PageCount() != 4 {
		t.Fatalf("PageCount = %d, want 4", doc.PageCount())
	}
}

func TestOpenDocumentMissingFile(t *testing.T) {
	_, err := openDocument("../../testdata/does-not-exist.pdf", "")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !strings.Contains(err.Error(), "does-not-exist.pdf") {
		t.Fatalf("error should mention the path, got: %v", err)
	}
}
