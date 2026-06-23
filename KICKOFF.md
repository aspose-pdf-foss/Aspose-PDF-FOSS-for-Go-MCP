# aspose-pdf-foss-mcp — Kickoff Brief

> Handoff from the library-build context. Goal: build a **separate** repo — an
> **MCP server** that exposes the pure-Go PDF library `aspose.pdf-for-go-foss`
> to AI agents. This file is the starting point; read it, then begin.

## 1. What we're building

`aspose-pdf-foss-mcp` — a Model Context Protocol server (stdio transport) that
wraps the public API of our pure-Go, zero-dependency PDF library so any MCP
client (Claude, IDEs, etc.) can extract/render/merge/validate/sign/convert PDFs.

It's a **thin wrapper**: each MCP tool maps to one or a few public library calls.
No PDF internals are reimplemented here. The big win is making the whole library
(rendering, PDF/A, signatures, forms, grayscale, …) agent-accessible in one
server — a strong showcase for the Aspose brand in the Go ecosystem.

## 2. The library it wraps

- **Repo (local):** `D:\aspose\claude\aspose.pdf-for-go-foss`
- **Module path:** `github.com/aspose-pdf-foss/aspose-pdf-foss-for-go`
- **Go:** 1.24 (toolchain 1.26 installed). Pure Go, **zero deps**, MIT.
- **Not published** — consume it locally with a `replace` directive.

This MCP repo lives at `D:\aspose\claude\aspose-pdf-foss-mcp` (sibling folder), so
the relative path back to the library is `../aspose.pdf-for-go-foss`.

`go.mod` to start with:

```
module github.com/aspose-pdf-foss/aspose-pdf-foss-mcp

go 1.24

require github.com/aspose-pdf-foss/aspose-pdf-foss-for-go v0.0.0-00010101000000-000000000000

replace github.com/aspose-pdf-foss/aspose-pdf-foss-for-go => ../aspose.pdf-for-go-foss
```

Import in code: `import pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"`.

> **Note on dependencies:** the *library* is zero-dependency by rule. This *MCP
> server* is a separate program and MAY take dependencies (the MCP SDK). The
> zero-dep rule does NOT apply here.

**Full public API:** read the library's `README.md` and the "Public API" section
of its `CLAUDE.md` (both in the library repo) for exact, current signatures. The
subset below is the high-value surface for tools.

## 3. Public API surface worth exposing (verify signatures against README)

Open / create / save:
- `pdf.Open(path string) (*pdf.Document, error)` — `ErrEncrypted` if password-protected
- `pdf.OpenWithPassword(path, password string) (*pdf.Document, error)`
- `pdf.OpenStream(io.Reader)`, `pdf.OpenStreamWithPassword`
- `(*Document).Save(path string) error`, `WriteTo(io.Writer)`, `SaveLinearized(path)`

Inspect / extract:
- `(*Document).PageCount() int`, `Page(n int) (*Page, error)`
- `(*Document).ExtractText() ([]string, error)` (one string/page)
- `(*Document).SearchText(query string, opts...) ([]TextMatch, error)`
- `(*Document).ExtractImages() ([][]pdf.Image, error)`; `(*pdf.Image).Save(path)`
- `(*Document).Info() (DocumentInfo, error)`, `XMP()`
- `pdf.PageSizes(path) ([]PageSize, error)`

Render (returns `image.Image` / writes encoded bytes):
- `(*Document).RenderImage(pageNum int, opts RenderOptions) (image.Image, error)`
- `(*Page).RenderPNG(w, opts) / RenderJPEG(w, opts, quality)`; `(*Document).RenderTIFF(w, opts, pages...)`
- `RenderOptions{ DPI float64; Background *Color }` (DPI 0 → 150)

Manipulate:
- `(*Document).Split() ([]*Document, error)`, `Extract(ranges ...PageRange) (*Document, error)`
- `(*Document).Append(others ...*Document)` (merge), `DeletePage(n)`, `DeletePages(...)`, `Reorder([]int)`
- `(*Document).Rotate(angle RotationAngle, pageNums ...int)`, `SetRotation(...)`
- `(*Document).RemoveUnusedObjects() int`, `OptimizeImages(opts) (int, error)`
- `(*Document).NUp(NUpOptions)`, `Booklet(BookletOptions)`
- `(*Document).Flatten() error` (bake form fields)

Compliance / conversion (this session's headline features):
- `pdf.Validate(path) (*ValidationReport, error)` — structural integrity
- `(*Document).ValidatePDFA(format PDFAFormat) *PDFAValidationReport` — `PDFA1B/2B/3B/1A/2A/3A`
- `(*Document).ConvertToPDFA(format PDFAFormat) (*PDFAValidationReport, error)`
- `(*Document).ValidatePDFUA() *PDFUAValidationReport`
- `(*Document).ConvertToGrayscale() error`

Security:
- `(*Document).SetEncryption(EncryptionOptions)`, `RemoveEncryption()`, `ChangePassword(new,newOwner)`
- `pdf.Encrypt(in, out, user, owner string) error`
- `(*Document).Sign(SignOptions) error` (needs cert + `crypto.Signer`), `VerifySignatures() ([]SignatureVerification, error)`

Forms:
- `(*Document).Form() *Form`; `(*Form).Fields() []Field`, `Field(name)`, `AddTextField/…`
- `(Field).Value()/SetValue(s)`, `FullName()`

Authoring (less likely as MCP tools, but available): `NewDocumentFromFormat`,
`(*Page).AddText/AddImage/AddTable/Draw*/AddSVG/AddStamp`, `TaggedContent()`,
watermarks, TOC, outlines, annotations.

## 4. Recommended stack & conventions

- **Transport:** stdio (the standard for local MCP servers). Add HTTP/SSE later if needed.
- **MCP SDK (decide in first step):** two solid Go options —
  1. **official** `github.com/modelcontextprotocol/go-sdk` — canonical, actively developed;
  2. `github.com/mark3labs/mcp-go` — mature, widely used, ergonomic.
  Recommendation: start with the **official go-sdk**; fall back to mark3labs if its
  ergonomics/stability fit better. Check current versions/READMEs before choosing.
- **I/O convention (important design choice):** PDFs are binary, so the cleanest
  for a local stdio server is **file paths**: tools take `input_path` (and
  `output_path` for producers) and operate on the local filesystem. Text tools
  return text directly. For rendered pages, write a PNG to `output_path` and/or
  return it as MCP image content (base64). Pick one convention and keep it
  consistent across tools. (Avoid shoving multi-MB PDFs through JSON.)
- **Error handling:** map library errors to MCP tool errors with clear messages;
  surface `ErrEncrypted` specially (suggest the password tool).
- **License:** MIT, `// SPDX-License-Identifier: MIT` header on every `.go` file
  (match the library's convention).
- **Tests:** table-driven; use small fixture PDFs (can copy a couple from the
  library's `testdata/`). `go test ./...` green before commits.

## 5. Proposed initial tool set (MVP → then expand)

MVP (build + verify these first):
- `pdf_info` — page count, size, title/author/producer, encrypted? (Open + Info + PageSizes)
- `pdf_extract_text` — `{input_path}` → text (per page or joined)
- `pdf_render_page` — `{input_path, page, dpi, output_path}` → PNG on disk (+ optional base64)
- `pdf_merge` — `{input_paths[], output_path}` (Open each + Append + Save)
- `pdf_split` — `{input_path, output_dir}` (Split → Save each) / or extract ranges
- `pdf_validate` — `{input_path, profile: structural|pdfa-1b|…|pdfua}` → report

Then expand:
- `pdf_convert_pdfa` — `{input_path, output_path, level}` → ConvertToPDFA + report
- `pdf_to_grayscale` — `{input_path, output_path}`
- `pdf_rotate` / `pdf_delete_pages` / `pdf_extract_pages`
- `pdf_encrypt` / `pdf_decrypt` / `pdf_change_password`
- `pdf_extract_images` — `{input_path, output_dir}`
- `pdf_form_fields` (read) / `pdf_fill_form` (write) / `pdf_flatten`
- `pdf_optimize` (OptimizeImages + RemoveUnusedObjects)
- `pdf_linearize`, `pdf_sign` / `pdf_verify_signatures`, `pdf_add_watermark`

(Signatures need a cert + key — design that tool's inputs carefully, e.g. paths
to a PEM cert + key; the library takes `*x509.Certificate` + `crypto.Signer`.)

## 6. Suggested project structure

```
aspose-pdf-foss-mcp/
  go.mod / go.sum
  main.go                 // wire up the server + register tools, serve stdio
  internal/tools/         // one file per tool group (text.go, render.go, convert.go, …)
  internal/tools/tools_test.go
  testdata/               // a few small fixture PDFs
  README.md               // what it is, install, tool list, client config example
  CLAUDE.md               // conventions for this repo (mirror the library's style)
  LICENSE                 // MIT
```

Include a client-config snippet in the README (how to register the server in an
MCP client, e.g. the `command`/`args` for stdio).

## 7. First concrete steps in the new context

1. `cd` into `D:\aspose\claude\aspose-pdf-foss-mcp`, `git init`.
2. Create `go.mod` (above) + the `replace` directive; `go get` the chosen MCP SDK;
   confirm the library imports and builds (`go build ./...`).
3. Smoke-test the dependency: a tiny `main` that `pdf.Open`s a fixture and prints
   `PageCount()` — proves the wiring before any MCP code.
4. Stand up a minimal MCP server over stdio with ONE tool (`pdf_info`), test it
   end-to-end with an MCP client / the SDK's test harness.
5. Add the rest of the MVP tools, then expand. Write the README + CLAUDE.md.

## 8. Open decisions to make up front

- MCP SDK: official go-sdk vs mark3labs/mcp-go (recommend trying official first).
- I/O: file-path-based (recommended) vs base64 in/out vs hybrid.
- Tool granularity & naming (`pdf_*` prefix recommended for discoverability).
- How render returns images (disk path vs MCP image content vs both).
- Whether to tag a release of the library (`v0.1.0`) so the MCP repo can require
  it normally instead of a `replace` — for now `replace` is simplest.

---
*Library state at handoff: rendering, text/image extraction, forms, annotations,
SVG, tables, encryption (RC4/AES-128/AES-256), digital signatures (PKCS#7/PAdES/
timestamp/multi/encrypted), PDF/A (validate+convert, 1b/2b/3b + 1a/2a/3a, auto
font embedding), PDF/UA (validate + tagged-content authoring), linearization,
grayscale conversion, N-up/booklet, optimize — all pure Go. See the library's
README/CLAUDE.md for the authoritative API.*
