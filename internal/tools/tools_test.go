// SPDX-License-Identifier: MIT

package tools_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aspose-pdf-foss/aspose-pdf-foss-mcp/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// newTestSession starts a server with all tools registered and returns a
// connected client session.
func newTestSession(t *testing.T) *mcp.ClientSession {
	t.Helper()
	ctx := context.Background()
	clientT, serverT := mcp.NewInMemoryTransports()

	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "test"}, nil)
	tools.Register(server)
	serverSession, err := server.Connect(ctx, serverT, nil)
	if err != nil {
		t.Fatalf("server connect: %v", err)
	}
	t.Cleanup(func() { serverSession.Close() })

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client"}, nil)
	clientSession, err := client.Connect(ctx, clientT, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	t.Cleanup(func() { clientSession.Close() })
	return clientSession
}

// resultText concatenates all TextContent in a tool result.
func resultText(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	var b strings.Builder
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			b.WriteString(tc.Text)
		}
	}
	return b.String()
}

func TestPDFInfo(t *testing.T) {
	cs := newTestSession(t)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "pdf_info",
		Arguments: map[string]any{"input_path": "../../testdata/4pages.pdf"},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if res.IsError {
		t.Fatalf("tool error: %s", resultText(t, res))
	}
	text := resultText(t, res)
	if !strings.Contains(text, "\"page_count\": 4") {
		t.Fatalf("expected page_count 4 in output, got: %s", text)
	}
}

func TestPDFExtractText(t *testing.T) {
	cs := newTestSession(t)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "pdf_extract_text",
		Arguments: map[string]any{"input_path": "../../testdata/Hello world.pdf"},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if res.IsError {
		t.Fatalf("tool error: %s", resultText(t, res))
	}
	text := resultText(t, res)
	if !strings.Contains(strings.ToLower(text), "hello") {
		t.Fatalf("expected extracted text to contain 'hello', got: %s", text)
	}
}

func TestPDFExtractTextPageRange(t *testing.T) {
	cs := newTestSession(t)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "pdf_extract_text",
		Arguments: map[string]any{
			"input_path": "../../testdata/4pages.pdf",
			"pages":      "2",
		},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if res.IsError {
		t.Fatalf("tool error: %s", resultText(t, res))
	}
	text := resultText(t, res)
	if strings.Count(text, "--- Page") != 1 {
		t.Fatalf("expected exactly one page marker, got: %s", text)
	}
}

func TestPDFRenderPage(t *testing.T) {
	cs := newTestSession(t)
	out := filepath.Join(t.TempDir(), "page1.png")
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "pdf_render_page",
		Arguments: map[string]any{
			"input_path":  "../../testdata/4pages.pdf",
			"page":        1,
			"output_path": out,
			"dpi":         72,
		},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if res.IsError {
		t.Fatalf("tool error: %s", resultText(t, res))
	}
	fi, err := os.Stat(out)
	if err != nil {
		t.Fatalf("output not written: %v", err)
	}
	if fi.Size() == 0 {
		t.Fatal("output file is empty")
	}
}
