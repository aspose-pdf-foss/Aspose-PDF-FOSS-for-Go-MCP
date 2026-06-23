# aspose-pdf-foss-mcp — Design Spec (MVP)

**Date:** 2026-06-23
**Status:** Approved for implementation planning

## 1. Purpose

A Model Context Protocol (MCP) server, stdio transport, written in Go, that wraps
the **public API** of the pure-Go PDF library `aspose.pdf-for-go-foss`
(`github.com/aspose-pdf-foss/aspose-pdf-foss-for-go`) so any MCP client (Claude,
IDEs, etc.) can inspect, extract from, render, merge, split, and validate PDFs.

It is a **thin wrapper**: each MCP tool maps to one or a few public library calls.
No PDF internals are reimplemented. The MCP server is a separate program from the
library and MAY take dependencies (the MCP SDK) — the library's zero-dependency
rule does NOT apply here.

## 2. Decisions (locked)

| Decision | Choice |
|---|---|
| MCP SDK | **official** `github.com/modelcontextprotocol/go-sdk` **v1.6.1** |
| Library version | **v0.3.0** (released tag), consumed via `go get` (no `replace`) |
| I/O convention | **File paths** — `input_path` / `output_path`; text returned directly |
| Scope | **MVP** — 6 tools |
| Repo location | `c:\Users\ANDREY-PC\Aspose\aspose-pdf-foss-for-go-mcp` (current folder), `git init` here |
| Transport | stdio |
| License | MIT (`// SPDX-License-Identifier: MIT` on every `.go` file) |

### Note on library scope vs. the kickoff brief

The kickoff brief listed PDF/A, digital signatures, grayscale conversion,
linearization, and N-up/booklet as headline features. These exist in the
library's source but are **NOT in the v0.3.0 release tag**. This MCP server
targets v0.3.0, so those tools are deferred until a library release ships them.
The README documents this explicitly.

## 3. Library public API actually available in v0.3.0 (verified)

Confirmed against the v0.3.0 source in the module cache. Relevant surface:

- Open/save: `pdf.Open(path) (*Document, error)`,
  `pdf.OpenWithPassword(path, pw) (*Document, error)`,
  `(*Document).Save(path) error`, `(*Document).Append(others ...*Document)`,
  `(*Document).Split() ([]*Document, error)`.
- Sentinel: `pdf.ErrEncrypted` — checked with `errors.Is(err, pdf.ErrEncrypted)`.
- Inspect: `(*Document).PageCount() int`, `(*Document).Page(n int) (*Page, error)`
  (1-based), `(*Document).Info() (DocumentInfo, error)`,
  `pdf.PageSizes(path) ([]PageSize, error)`.
- `DocumentInfo`: Title, Author, Subject, Keywords, Creator, Producer,
  CreationDate, ModDate, Custom map[string]string.
- `PageSize{ Width, Height float64 }`.
- Text: `(*Document).ExtractText() ([]string, error)` (one string per page).
- Render: `(*Page).RenderPNG(w io.Writer, opts RenderOptions) error`,
  `(*Page).RenderJPEG(w io.Writer, opts RenderOptions, quality int) error`,
  `RenderOptions{ DPI float64; Background *Color }` (DPI 0 → 150).
- Validate: `pdf.Validate(path) (*ValidationReport, error)`,
  `ValidationReport{ Valid bool; Issues []ValidationIssue }`,
  `ValidationIssue{ Code, Message string }`.

## 4. MCP SDK usage pattern (v1.6.1)

```go
server := mcp.NewServer(&mcp.Implementation{Name: "...", Version: "..."}, nil)

mcp.AddTool(server, &mcp.Tool{Name: "pdf_info", Description: "..."},
    func(ctx context.Context, req *mcp.CallToolRequest, args InfoParams) (*mcp.CallToolResult, any, error) {
        // ...
    })

server.Run(ctx, &mcp.StdioTransport{})
```

`mcp.AddTool` is generic: it auto-generates the JSON input schema from the
`*Params` struct's `json` and `jsonschema` struct tags. Handler signature is
`func(ctx, *mcp.CallToolRequest, ArgsStruct) (*mcp.CallToolResult, any, error)`.
Tests use `mcp.NewInMemoryTransports()` to run server + client in-process.

## 5. Architecture & project structure

```
aspose-pdf-foss-for-go-mcp/
  go.mod / go.sum            // module github.com/aspose-pdf-foss/aspose-pdf-foss-mcp
  main.go                    // build server, register tools, server.Run(stdio)
  internal/server/
    server.go                // New() *mcp.Server + RegisterAll(server): one place that AddTools every tool
  internal/tools/
    info.go                  // pdf_info
    text.go                  // pdf_extract_text
    render.go                // pdf_render_page
    merge.go                 // pdf_merge
    split.go                 // pdf_split
    validate.go              // pdf_validate
    errors.go                // map library errors -> MCP errors (ErrEncrypted special-cased)
    pages.go                 // shared page-range parsing helper ("1-3,5")
    tools_test.go            // end-to-end tests via mcp.NewInMemoryTransports()
  testdata/                  // small fixture PDFs copied from the library's testdata/
  README.md
  CLAUDE.md
  LICENSE
```

- Dependencies: `pdf` v0.3.0 + `go-sdk` v1.6.1, both via `go get`.
- One tool = one file, one `*Params` struct (json/jsonschema tags), one handler.
  Each handler is independent and testable in isolation.
- Data flow: client JSON args → typed struct → handler opens PDF by `input_path`
  → 1–2 public `pdf` calls → result returned as text, or written to
  `output_path`/`output_dir` with a short text confirmation.

## 6. Tool contracts

All page numbers are 1-based. Reading tools accept an optional `password`. If a
file is encrypted and no password is supplied, the tool errors with a message
suggesting the `password` parameter.

### 6.1 `pdf_info`
- **Params:** `input_path` (req), `password` (opt).
- **Calls:** `Open`/`OpenWithPassword` + `Info()` + `PageCount()` + `pdf.PageSizes()`.
- **Returns:** JSON text — `page_count`, `page_sizes` (W×H pt per page),
  title/author/subject/keywords/creator/producer/creation_date/mod_date,
  `encrypted` (bool).

### 6.2 `pdf_extract_text`
- **Params:** `input_path` (req), `password` (opt), `pages` (opt range string
  e.g. "1-3,5"; default all), `separator` (opt; default `\n\n--- Page N ---\n`).
- **Calls:** `ExtractText()` → `[]string` per page; select requested pages.
- **Returns:** extracted text, pages joined by the separator.

### 6.3 `pdf_render_page`
- **Params:** `input_path` (req), `page` (req, 1-based), `output_path` (req),
  `dpi` (opt; default 150), `format` (opt: png|jpeg; default inferred from
  output_path extension, fallback png), `jpeg_quality` (opt; default 90).
- **Calls:** `Page(page)` + `RenderPNG` or `RenderJPEG` to the output file.
- **Returns:** text confirming `output_path` + rendered pixel dimensions.

### 6.4 `pdf_merge`
- **Params:** `input_paths` (req, ≥2), `output_path` (req).
- **Calls:** open first, `Append(rest...)`, `Save(output_path)`.
- **Returns:** text confirming `output_path` + total `page_count`.

### 6.5 `pdf_split`
- **Params:** `input_path` (req), `output_dir` (req), `password` (opt).
- **Calls:** `Split()` + `Save` each as `page_001.pdf`, `page_002.pdf`, …
- **Returns:** list of created file paths.

### 6.6 `pdf_validate`
- **Params:** `input_path` (req).
- **Calls:** `pdf.Validate(input_path)`.
- **Returns:** `valid` (bool) + issues (code + message). Description states this
  is a **structural-integrity** check, NOT PDF/A or PDF/UA conformance.

## 7. Error handling

A shared helper in `errors.go` maps library errors to MCP tool errors:
- `errors.Is(err, pdf.ErrEncrypted)` → "PDF is encrypted — supply the `password` parameter."
- Missing file / not a PDF → message including the path.
- Other library errors → wrapped with context (which tool, which path).

Tool errors are returned as `error` from the handler; the SDK surfaces them to
the client so the agent can self-correct.

## 8. Testing

- **End-to-end** via `mcp.NewInMemoryTransports()`: spin up the server in-process,
  connect a client, `CallTool` each tool, assert on the result. Exercises the
  full chain (args schema → handler → library) without an external MCP client.
- **Fixtures** copied from the library `testdata/`: `4pages.pdf` (multi-page —
  info/split/render), `Hello world.pdf` (text), `PdfWithImages.pdf`; merge uses
  two of them.
- Table-driven tests. `go build ./...`, `go vet ./...`, `go test ./...` green
  before commits.

## 9. Documentation deliverables

- **README.md** — what it is; install (`go get` / `go build`); the 6 tools with
  their parameters; an MCP client config snippet (`command`/`args` for stdio); a
  note that PDF/A, signatures, grayscale, etc. arrive with a future library release.
- **CLAUDE.md** — repo conventions mirroring the library's style: SPDX MIT header
  on every `.go` file; one tool = one file; file-path I/O convention; 1-based
  pages; how to add a new tool; how to run tests.
- **LICENSE** — MIT.

## 10. Out of scope (MVP)

- PDF/A validate/convert, PDF/UA, grayscale, digital signatures, linearization,
  N-up/booklet, encryption/decryption tools, image extraction, form
  read/fill/flatten, optimize, rotate/delete/extract-pages, watermark, search,
  metadata write. All are candidates for the post-MVP expansion once needed (and,
  for the unreleased library features, once a library release ships them).
- HTTP/SSE transport (stdio only for now).
- base64 image return (file-path convention chosen).
```
