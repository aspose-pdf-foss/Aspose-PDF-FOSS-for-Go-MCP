// SPDX-License-Identifier: MIT

package tools_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"image"
	"image/png"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/aspose-pdf-foss/aspose-pdf-foss-for-go-mcp/internal/tools"
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

func TestPDFRenderPageUnsupportedFormat(t *testing.T) {
	cs := newTestSession(t)
	out := filepath.Join(t.TempDir(), "page1.gif")
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "pdf_render_page",
		Arguments: map[string]any{
			"input_path":  "../../testdata/4pages.pdf",
			"page":        1,
			"output_path": out,
			"format":      "gif",
		},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected an error for an unsupported format")
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Fatalf("unsupported format must not leave an output file, stat err = %v", err)
	}
}

func TestPDFMerge(t *testing.T) {
	cs := newTestSession(t)
	out := filepath.Join(t.TempDir(), "merged.pdf")
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "pdf_merge",
		Arguments: map[string]any{
			"input_paths": []any{"../../testdata/4pages.pdf", "../../testdata/Hello world.pdf"},
			"output_path": out,
		},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if res.IsError {
		t.Fatalf("tool error: %s", resultText(t, res))
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("merged file not written: %v", err)
	}
	if !strings.Contains(resultText(t, res), "page") {
		t.Fatalf("expected page count in result, got: %s", resultText(t, res))
	}
}

func TestPDFMergeTooFew(t *testing.T) {
	cs := newTestSession(t)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "pdf_merge",
		Arguments: map[string]any{
			"input_paths": []any{"../../testdata/4pages.pdf"},
			"output_path": filepath.Join(t.TempDir(), "x.pdf"),
		},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected an error for fewer than 2 inputs")
	}
}

func TestPDFSplit(t *testing.T) {
	cs := newTestSession(t)
	dir := t.TempDir()
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "pdf_split",
		Arguments: map[string]any{
			"input_path": "../../testdata/4pages.pdf",
			"output_dir": dir,
		},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if res.IsError {
		t.Fatalf("tool error: %s", resultText(t, res))
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 4 {
		t.Fatalf("expected 4 split files, got %d", len(entries))
	}
	if _, err := os.Stat(filepath.Join(dir, "page_001.pdf")); err != nil {
		t.Fatalf("page_001.pdf missing: %v", err)
	}
}

func TestPDFValidate(t *testing.T) {
	cs := newTestSession(t)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "pdf_validate",
		Arguments: map[string]any{"input_path": "../../testdata/4pages.pdf"},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if res.IsError {
		t.Fatalf("tool error: %s", resultText(t, res))
	}
	if !strings.Contains(resultText(t, res), "\"valid\": true") {
		t.Fatalf("expected valid true, got: %s", resultText(t, res))
	}
}

// callOK calls a tool and fails the test on a transport or tool error.
func callOK(t *testing.T, cs *mcp.ClientSession, name string, args map[string]any) string {
	t.Helper()
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatalf("CallTool %s: %v", name, err)
	}
	if res.IsError {
		t.Fatalf("%s tool error: %s", name, resultText(t, res))
	}
	return resultText(t, res)
}

// callErr calls a tool and fails the test unless the tool reports an error.
func callErr(t *testing.T, cs *mcp.ClientSession, name string, args map[string]any) string {
	t.Helper()
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatalf("CallTool %s: %v", name, err)
	}
	if !res.IsError {
		t.Fatalf("%s: expected a tool error, got: %s", name, resultText(t, res))
	}
	return resultText(t, res)
}

// makeFormPDF writes a one-page PDF with a text field named "name" and a
// checkbox named "agree", and returns its path.
func makeFormPDF(t *testing.T) string {
	t.Helper()
	doc := pdf.NewDocument(612, 792)
	if _, err := doc.Form().AddTextField(1, pdf.Rectangle{LLX: 100, LLY: 700, URX: 300, URY: 720}, "name"); err != nil {
		t.Fatalf("AddTextField: %v", err)
	}
	if _, err := doc.Form().AddCheckbox(1, pdf.Rectangle{LLX: 100, LLY: 660, URX: 120, URY: 680}, "agree"); err != nil {
		t.Fatalf("AddCheckbox: %v", err)
	}
	path := filepath.Join(t.TempDir(), "form.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save form fixture: %v", err)
	}
	return path
}

// makeEncryptedPDF writes an AES-128-encrypted copy of Hello world.pdf with
// user password "secret" and returns its path.
func makeEncryptedPDF(t *testing.T) string {
	t.Helper()
	doc, err := pdf.Open("../../testdata/Hello world.pdf")
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	doc.SetEncryption(pdf.EncryptionOptions{UserPassword: "secret"})
	path := filepath.Join(t.TempDir(), "encrypted.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save encrypted fixture: %v", err)
	}
	return path
}

// makeCertAndKey writes a self-signed ECDSA P-256 certificate and its private
// key as PEM files and returns their paths.
func makeCertAndKey(t *testing.T) (certPath, keyPath string) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "MCP Test Signer"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("marshal key: %v", err)
	}
	dir := t.TempDir()
	certPath = filepath.Join(dir, "cert.pem")
	keyPath = filepath.Join(dir, "key.pem")
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})
	if err := os.WriteFile(certPath, certPEM, 0o644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyPath, keyPEM, 0o600); err != nil {
		t.Fatalf("write key: %v", err)
	}
	return certPath, keyPath
}

func TestPDFValidatePDFA(t *testing.T) {
	cs := newTestSession(t)
	text := callOK(t, cs, "pdf_validate", map[string]any{
		"input_path": "../../testdata/4pages.pdf",
		"profile":    "pdfa-1b",
	})
	if !strings.Contains(text, "\"profile\": \"pdfa-1b\"") || !strings.Contains(text, "\"issues\"") {
		t.Fatalf("expected a pdfa-1b report, got: %s", text)
	}
}

func TestPDFValidatePDFUA(t *testing.T) {
	cs := newTestSession(t)
	text := callOK(t, cs, "pdf_validate", map[string]any{
		"input_path": "../../testdata/4pages.pdf",
		"profile":    "pdfua",
	})
	if !strings.Contains(text, "\"profile\": \"pdfua\"") {
		t.Fatalf("expected a pdfua report, got: %s", text)
	}
}

func TestPDFValidateUnknownProfile(t *testing.T) {
	cs := newTestSession(t)
	callErr(t, cs, "pdf_validate", map[string]any{
		"input_path": "../../testdata/4pages.pdf",
		"profile":    "pdfa-9z",
	})
}

func TestPDFConvertPDFA(t *testing.T) {
	cs := newTestSession(t)
	out := filepath.Join(t.TempDir(), "pdfa.pdf")
	text := callOK(t, cs, "pdf_convert_pdfa", map[string]any{
		"input_path":  "../../testdata/Hello world.pdf",
		"output_path": out,
		"level":       "pdfa-1b",
	})
	if !strings.Contains(text, "\"level\": \"PDF/A-1B\"") {
		t.Fatalf("expected a PDF/A-1B report, got: %s", text)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("output not written: %v", err)
	}
}

func TestPDFToGrayscale(t *testing.T) {
	cs := newTestSession(t)
	out := filepath.Join(t.TempDir(), "gray.pdf")
	callOK(t, cs, "pdf_to_grayscale", map[string]any{
		"input_path":  "../../testdata/PdfWithImages.pdf",
		"output_path": out,
	})
	text := callOK(t, cs, "pdf_validate", map[string]any{"input_path": out})
	if !strings.Contains(text, "\"valid\": true") {
		t.Fatalf("grayscale output not structurally valid: %s", text)
	}
}

func TestPDFLinearize(t *testing.T) {
	cs := newTestSession(t)
	out := filepath.Join(t.TempDir(), "linear.pdf")
	callOK(t, cs, "pdf_linearize", map[string]any{
		"input_path":  "../../testdata/4pages.pdf",
		"output_path": out,
	})
	text := callOK(t, cs, "pdf_validate", map[string]any{"input_path": out})
	if !strings.Contains(text, "\"valid\": true") {
		t.Fatalf("linearized output not structurally valid: %s", text)
	}
}

func TestPDFNUp(t *testing.T) {
	cs := newTestSession(t)
	out := filepath.Join(t.TempDir(), "nup.pdf")
	callOK(t, cs, "pdf_nup", map[string]any{
		"input_path":  "../../testdata/4pages.pdf",
		"output_path": out,
		"rows":        2,
		"cols":        2,
	})
	text := callOK(t, cs, "pdf_info", map[string]any{"input_path": out})
	if !strings.Contains(text, "\"page_count\": 1") {
		t.Fatalf("expected 4 pages imposed onto 1 sheet, got: %s", text)
	}
}

func TestPDFBooklet(t *testing.T) {
	cs := newTestSession(t)
	out := filepath.Join(t.TempDir(), "booklet.pdf")
	callOK(t, cs, "pdf_booklet", map[string]any{
		"input_path":  "../../testdata/4pages.pdf",
		"output_path": out,
	})
	text := callOK(t, cs, "pdf_info", map[string]any{"input_path": out})
	if !strings.Contains(text, "\"page_count\": 2") {
		t.Fatalf("expected 4 pages imposed as 2 spreads, got: %s", text)
	}
}

func TestPDFEncrypt(t *testing.T) {
	cs := newTestSession(t)
	out := filepath.Join(t.TempDir(), "locked.pdf")
	callOK(t, cs, "pdf_encrypt", map[string]any{
		"input_path":    "../../testdata/Hello world.pdf",
		"output_path":   out,
		"user_password": "secret",
	})
	callErr(t, cs, "pdf_info", map[string]any{"input_path": out})
	text := callOK(t, cs, "pdf_info", map[string]any{"input_path": out, "password": "secret"})
	if !strings.Contains(text, "\"page_count\": 1") {
		t.Fatalf("expected the encrypted copy to open with the password, got: %s", text)
	}
}

func TestPDFEncryptNoPasswords(t *testing.T) {
	cs := newTestSession(t)
	callErr(t, cs, "pdf_encrypt", map[string]any{
		"input_path":  "../../testdata/Hello world.pdf",
		"output_path": filepath.Join(t.TempDir(), "x.pdf"),
	})
}

func TestPDFDecrypt(t *testing.T) {
	cs := newTestSession(t)
	encrypted := makeEncryptedPDF(t)
	out := filepath.Join(t.TempDir(), "plain.pdf")
	callOK(t, cs, "pdf_decrypt", map[string]any{
		"input_path":  encrypted,
		"output_path": out,
		"password":    "secret",
	})
	text := callOK(t, cs, "pdf_info", map[string]any{"input_path": out})
	if !strings.Contains(text, "\"encrypted\": false") {
		t.Fatalf("expected decrypted output to open without a password, got: %s", text)
	}
}

func TestPDFSignAndVerify(t *testing.T) {
	cs := newTestSession(t)
	certPath, keyPath := makeCertAndKey(t)
	out := filepath.Join(t.TempDir(), "signed.pdf")
	callOK(t, cs, "pdf_sign", map[string]any{
		"input_path":  "../../testdata/4pages.pdf",
		"output_path": out,
		"cert_path":   certPath,
		"key_path":    keyPath,
		"reason":      "e2e test",
	})
	text := callOK(t, cs, "pdf_verify_signatures", map[string]any{"input_path": out})
	if !strings.Contains(text, "\"count\": 1") || !strings.Contains(text, "\"valid\": true") {
		t.Fatalf("expected one valid signature, got: %s", text)
	}
	if !strings.Contains(text, "MCP Test Signer") {
		t.Fatalf("expected the signer certificate subject, got: %s", text)
	}
}

func TestPDFExtractImages(t *testing.T) {
	cs := newTestSession(t)
	dir := filepath.Join(t.TempDir(), "images")
	text := callOK(t, cs, "pdf_extract_images", map[string]any{
		"input_path": "../../testdata/PdfWithImages.pdf",
		"output_dir": dir,
	})
	if !strings.Contains(text, "Extracted") {
		t.Fatalf("expected extracted images, got: %s", text)
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) == 0 {
		t.Fatalf("expected image files in %s (err=%v, n=%d)", dir, err, len(entries))
	}
}

func TestPDFFormFieldsAndFill(t *testing.T) {
	cs := newTestSession(t)
	form := makeFormPDF(t)
	text := callOK(t, cs, "pdf_form_fields", map[string]any{"input_path": form})
	if !strings.Contains(text, "\"count\": 2") || !strings.Contains(text, "\"name\": \"name\"") || !strings.Contains(text, "\"type\": \"checkbox\"") {
		t.Fatalf("expected a text and a checkbox field, got: %s", text)
	}

	out := filepath.Join(t.TempDir(), "filled.pdf")
	callOK(t, cs, "pdf_fill_form", map[string]any{
		"input_path":  form,
		"output_path": out,
		"values":      map[string]any{"name": "Alice"},
	})
	text = callOK(t, cs, "pdf_form_fields", map[string]any{"input_path": out})
	if !strings.Contains(text, "\"value\": \"Alice\"") {
		t.Fatalf("expected the filled value to round-trip, got: %s", text)
	}
}

func TestPDFFillFormUnknownField(t *testing.T) {
	cs := newTestSession(t)
	form := makeFormPDF(t)
	text := callErr(t, cs, "pdf_fill_form", map[string]any{
		"input_path":  form,
		"output_path": filepath.Join(t.TempDir(), "x.pdf"),
		"values":      map[string]any{"no_such_field": "x"},
	})
	if !strings.Contains(text, "no_such_field") {
		t.Fatalf("expected the unknown field name in the error, got: %s", text)
	}
}

func TestPDFFlatten(t *testing.T) {
	cs := newTestSession(t)
	form := makeFormPDF(t)
	out := filepath.Join(t.TempDir(), "flat.pdf")
	callOK(t, cs, "pdf_flatten", map[string]any{
		"input_path":  form,
		"output_path": out,
	})
	text := callOK(t, cs, "pdf_form_fields", map[string]any{"input_path": out})
	if !strings.Contains(text, "\"count\": 0") {
		t.Fatalf("expected no fields after flattening, got: %s", text)
	}
}

func TestPDFOptimize(t *testing.T) {
	cs := newTestSession(t)
	out := filepath.Join(t.TempDir(), "optimized.pdf")
	text := callOK(t, cs, "pdf_optimize", map[string]any{
		"input_path":  "../../testdata/PdfWithImages.pdf",
		"output_path": out,
		"max_dpi":     72,
	})
	if !strings.Contains(text, "\"images_optimized\"") {
		t.Fatalf("expected an optimization report, got: %s", text)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("output not written: %v", err)
	}
}

func TestPDFRotate(t *testing.T) {
	cs := newTestSession(t)
	out := filepath.Join(t.TempDir(), "rotated.pdf")
	callOK(t, cs, "pdf_rotate", map[string]any{
		"input_path":  "../../testdata/4pages.pdf",
		"output_path": out,
		"angle":       90,
		"pages":       "1,3",
	})
	text := callOK(t, cs, "pdf_info", map[string]any{"input_path": out})
	if !strings.Contains(text, "\"page_count\": 4") {
		t.Fatalf("expected 4 pages after rotation, got: %s", text)
	}
}

func TestPDFRotateBadAngle(t *testing.T) {
	cs := newTestSession(t)
	callErr(t, cs, "pdf_rotate", map[string]any{
		"input_path":  "../../testdata/4pages.pdf",
		"output_path": filepath.Join(t.TempDir(), "x.pdf"),
		"angle":       45,
	})
}

func TestPDFDeletePages(t *testing.T) {
	cs := newTestSession(t)
	out := filepath.Join(t.TempDir(), "deleted.pdf")
	callOK(t, cs, "pdf_delete_pages", map[string]any{
		"input_path":  "../../testdata/4pages.pdf",
		"output_path": out,
		"pages":       "1,4",
	})
	text := callOK(t, cs, "pdf_info", map[string]any{"input_path": out})
	if !strings.Contains(text, "\"page_count\": 2") {
		t.Fatalf("expected 2 remaining pages, got: %s", text)
	}
}

func TestPDFExtractPages(t *testing.T) {
	cs := newTestSession(t)
	out := filepath.Join(t.TempDir(), "extracted.pdf")
	callOK(t, cs, "pdf_extract_pages", map[string]any{
		"input_path":  "../../testdata/4pages.pdf",
		"output_path": out,
		"pages":       "2-3",
	})
	text := callOK(t, cs, "pdf_info", map[string]any{"input_path": out})
	if !strings.Contains(text, "\"page_count\": 2") {
		t.Fatalf("expected 2 extracted pages, got: %s", text)
	}
}

func TestPDFWatermarkText(t *testing.T) {
	cs := newTestSession(t)
	out := filepath.Join(t.TempDir(), "watermarked.pdf")
	callOK(t, cs, "pdf_watermark", map[string]any{
		"input_path":  "../../testdata/4pages.pdf",
		"output_path": out,
		"text":        "CONFIDENTIAL",
	})
	text := callOK(t, cs, "pdf_search", map[string]any{
		"input_path": out,
		"query":      "CONFIDENTIAL",
	})
	if strings.Contains(text, "\"count\": 0") {
		t.Fatalf("expected the watermark text to be searchable, got: %s", text)
	}
}

func TestPDFWatermarkBothSources(t *testing.T) {
	cs := newTestSession(t)
	callErr(t, cs, "pdf_watermark", map[string]any{
		"input_path":  "../../testdata/4pages.pdf",
		"output_path": filepath.Join(t.TempDir(), "x.pdf"),
		"text":        "X",
		"image_path":  "y.png",
	})
}

func TestPDFSearch(t *testing.T) {
	cs := newTestSession(t)
	text := callOK(t, cs, "pdf_search", map[string]any{
		"input_path":       "../../testdata/Hello world.pdf",
		"query":            "hello",
		"case_insensitive": true,
	})
	if strings.Contains(text, "\"count\": 0") || !strings.Contains(text, "\"page\": 1") {
		t.Fatalf("expected a match on page 1, got: %s", text)
	}
}

func TestPDFCreate(t *testing.T) {
	cs := newTestSession(t)
	out := filepath.Join(t.TempDir(), "new.pdf")
	callOK(t, cs, "pdf_create", map[string]any{
		"output_path": out,
		"page_format": "letter",
		"page_count":  3,
	})
	text := callOK(t, cs, "pdf_info", map[string]any{"input_path": out})
	if !strings.Contains(text, "\"page_count\": 3") {
		t.Fatalf("expected page_count 3, got: %s", text)
	}
}

func TestPDFCreateCustomLandscape(t *testing.T) {
	cs := newTestSession(t)
	out := filepath.Join(t.TempDir(), "wide.pdf")
	callOK(t, cs, "pdf_create", map[string]any{
		"output_path": out,
		"width":       200,
		"height":      400,
		"landscape":   true,
	})
	text := callOK(t, cs, "pdf_info", map[string]any{"input_path": out})
	if !strings.Contains(text, "\"width\": 400") || !strings.Contains(text, "\"height\": 200") {
		t.Fatalf("expected a 400x200 page, got: %s", text)
	}
}

func TestPDFCreateBadFormat(t *testing.T) {
	cs := newTestSession(t)
	callErr(t, cs, "pdf_create", map[string]any{
		"output_path": filepath.Join(t.TempDir(), "x.pdf"),
		"page_format": "tabloid",
	})
}

func TestPDFAddText(t *testing.T) {
	cs := newTestSession(t)
	dir := t.TempDir()
	src := filepath.Join(dir, "blank.pdf")
	out := filepath.Join(dir, "text.pdf")
	callOK(t, cs, "pdf_create", map[string]any{"output_path": src})
	callOK(t, cs, "pdf_add_text", map[string]any{
		"input_path":  src,
		"output_path": out,
		"text":        "Invoice INV-001",
		"font":        "helvetica-bold",
		"font_size":   24,
		"color":       "#003366",
		"halign":      "center",
	})
	text := callOK(t, cs, "pdf_extract_text", map[string]any{"input_path": out})
	if !strings.Contains(text, "Invoice INV-001") {
		t.Fatalf("expected added text in extraction, got: %s", text)
	}
}

func TestPDFAddTextUnknownFont(t *testing.T) {
	cs := newTestSession(t)
	dir := t.TempDir()
	src := filepath.Join(dir, "blank.pdf")
	callOK(t, cs, "pdf_create", map[string]any{"output_path": src})
	callErr(t, cs, "pdf_add_text", map[string]any{
		"input_path":  src,
		"output_path": filepath.Join(dir, "x.pdf"),
		"text":        "X",
		"font":        "comic-sans",
	})
}

func TestPDFAddTable(t *testing.T) {
	cs := newTestSession(t)
	dir := t.TempDir()
	src := filepath.Join(dir, "blank.pdf")
	out := filepath.Join(dir, "table.pdf")
	callOK(t, cs, "pdf_create", map[string]any{"output_path": src})
	callOK(t, cs, "pdf_add_table", map[string]any{
		"input_path":  src,
		"output_path": out,
		"rows": [][]string{
			{"Item", "Qty", "Price"},
			{"Widget", "2", "10.00"},
			{"Gadget", "1", "25.50"},
		},
		"header_rows": 1,
	})
	text := callOK(t, cs, "pdf_extract_text", map[string]any{"input_path": out})
	if !strings.Contains(text, "Widget") || !strings.Contains(text, "25.50") {
		t.Fatalf("expected table cells in extraction, got: %s", text)
	}
}

func TestPDFAddTableRaggedRows(t *testing.T) {
	cs := newTestSession(t)
	dir := t.TempDir()
	src := filepath.Join(dir, "blank.pdf")
	callOK(t, cs, "pdf_create", map[string]any{"output_path": src})
	callErr(t, cs, "pdf_add_table", map[string]any{
		"input_path":  src,
		"output_path": filepath.Join(dir, "x.pdf"),
		"rows":        [][]string{{"a", "b"}, {"c"}},
	})
}

// makeTestPNG writes a small opaque PNG and returns its path.
func makeTestPNG(t *testing.T, dir string) string {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for i := range img.Pix {
		img.Pix[i] = 0xff
	}
	path := filepath.Join(dir, "pixel.png")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create png: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	return path
}

func TestPDFAddImage(t *testing.T) {
	cs := newTestSession(t)
	dir := t.TempDir()
	src := filepath.Join(dir, "blank.pdf")
	out := filepath.Join(dir, "image.pdf")
	callOK(t, cs, "pdf_create", map[string]any{"output_path": src})
	callOK(t, cs, "pdf_add_image", map[string]any{
		"input_path":  src,
		"output_path": out,
		"image_path":  makeTestPNG(t, dir),
		"llx":         100,
		"lly":         600,
		"urx":         200,
		"ury":         700,
	})
	imgDir := filepath.Join(dir, "extracted")
	text := callOK(t, cs, "pdf_extract_images", map[string]any{
		"input_path": out,
		"output_dir": imgDir,
	})
	if strings.Contains(text, "\"count\": 0") {
		t.Fatalf("expected the placed image to be extracted, got: %s", text)
	}
}

func TestPDFDraw(t *testing.T) {
	cs := newTestSession(t)
	dir := t.TempDir()
	src := filepath.Join(dir, "blank.pdf")
	out := filepath.Join(dir, "drawn.pdf")
	callOK(t, cs, "pdf_create", map[string]any{"output_path": src})
	callOK(t, cs, "pdf_draw", map[string]any{
		"input_path":  src,
		"output_path": out,
		"shape":       "rectangle",
		"llx":         50,
		"lly":         50,
		"urx":         250,
		"ury":         150,
		"fill_color":  "#ffcc00",
	})
	callOK(t, cs, "pdf_draw", map[string]any{
		"input_path":   out,
		"output_path":  out,
		"shape":        "line",
		"x1":           50,
		"y1":           40,
		"x2":           250,
		"y2":           40,
		"stroke_width": 2,
		"dash_pattern": []float64{3, 2},
	})
	text := callOK(t, cs, "pdf_validate", map[string]any{"input_path": out})
	if !strings.Contains(text, "\"valid\": true") {
		t.Fatalf("expected the drawn PDF to validate, got: %s", text)
	}
}

func TestPDFDrawBadShape(t *testing.T) {
	cs := newTestSession(t)
	dir := t.TempDir()
	src := filepath.Join(dir, "blank.pdf")
	callOK(t, cs, "pdf_create", map[string]any{"output_path": src})
	callErr(t, cs, "pdf_draw", map[string]any{
		"input_path":  src,
		"output_path": filepath.Join(dir, "x.pdf"),
		"shape":       "triangle",
	})
}
