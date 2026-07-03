# CLAUDE.md — Repo Conventions

This file documents the conventions that govern this repository. Read it before making any change.

## Purpose

This is a **thin MCP wrapper** over the `aspose-pdf-foss-for-go` PDF library. No PDF internals are reimplemented here. Every tool is a thin adapter: validate inputs, call a library function, format the result as MCP tool output.

## Module

```
github.com/aspose-pdf-foss/aspose-pdf-foss-for-go-mcp
```

Go directive: `go 1.25.0` (required by the MCP SDK v1.6.1).

## License header

Every `.go` file **must** begin with:

```go
// SPDX-License-Identifier: MIT
```

No exceptions. `.md` and other non-Go files do not need this header.

## Tool file layout

One tool = one file under `internal/tools/`. The file is named after the tool without the `pdf_` prefix:

| File | Tool |
|------|------|
| `internal/tools/info.go` | `pdf_info` |
| `internal/tools/text.go` | `pdf_extract_text` |
| `internal/tools/render.go` | `pdf_render_page` |
| `internal/tools/merge.go` | `pdf_merge` |
| `internal/tools/split.go` | `pdf_split` |
| `internal/tools/validate.go` | `pdf_validate` |
| `internal/tools/convert_pdfa.go` | `pdf_convert_pdfa` |
| `internal/tools/sign.go` | `pdf_sign` |
| `internal/tools/verify_signatures.go` | `pdf_verify_signatures` |
| `internal/tools/grayscale.go` | `pdf_to_grayscale` |
| `internal/tools/linearize.go` | `pdf_linearize` |
| `internal/tools/nup.go` | `pdf_nup` |
| `internal/tools/booklet.go` | `pdf_booklet` |
| `internal/tools/encrypt.go` | `pdf_encrypt` |
| `internal/tools/decrypt.go` | `pdf_decrypt` |
| `internal/tools/images.go` | `pdf_extract_images` |
| `internal/tools/form_fields.go` | `pdf_form_fields` |
| `internal/tools/fill_form.go` | `pdf_fill_form` |
| `internal/tools/flatten.go` | `pdf_flatten` |
| `internal/tools/optimize.go` | `pdf_optimize` |
| `internal/tools/rotate.go` | `pdf_rotate` |
| `internal/tools/delete_pages.go` | `pdf_delete_pages` |
| `internal/tools/extract_pages.go` | `pdf_extract_pages` |
| `internal/tools/watermark.go` | `pdf_watermark` |
| `internal/tools/search.go` | `pdf_search` |
| `internal/tools/create.go` | `pdf_create` |
| `internal/tools/add_text.go` | `pdf_add_text` |
| `internal/tools/add_table.go` | `pdf_add_table` |
| `internal/tools/add_image.go` | `pdf_add_image` |
| `internal/tools/draw.go` | `pdf_draw` |

Every tool is registered in `internal/tools/register.go` via `mcp.AddTool`. Do not register tools anywhere else.

## Adding a new tool

1. Create `internal/tools/<name>.go` with the SPDX header.
2. Define a `*Params` struct with `json` and `jsonschema` struct tags. Use `omitempty` on optional fields.
3. Write a `handleX` function with the signature:
   ```go
   func handleX(ctx context.Context, req *mcp.CallToolRequest, args xParams) (*mcp.CallToolResult, any, error)
   ```
4. Add `mcp.AddTool` in `internal/tools/register.go` with `Name` and `Description`.
5. Add an end-to-end test in `internal/tools/tools_test.go` using the `newTestSession` helper.

## I/O conventions

- **File paths** are passed as `input_path`, `output_path`, or `output_dir` string parameters.
- **Text results** are returned directly in the MCP tool result as `TextContent`.
- **Page numbers** are **1-based** throughout (matches the library).
- Errors are returned as Go `error` values; the MCP SDK converts them to tool-error responses automatically.

## Shared helpers

| File | Exported symbol | Purpose |
|------|----------------|---------|
| `internal/tools/errors.go` | `openDocument(path, password string) (*pdf.Document, error)` | Opens a PDF, handling the encrypted-but-no-password case with an actionable error message |
| `internal/tools/pages.go` | `parsePageRange(spec string, total int) ([]int, error)` | Parses a page range string such as `"1-3,5"` into a slice of 1-based page numbers; empty spec means all pages |
| `internal/tools/pdfa.go` | `parsePDFALevel(level string) (pdf.PDFAFormat, error)` | Maps a level string such as `"pdfa-1b"` to the library's `PDFAFormat` constant |
| `internal/tools/style.go` | `parseHexColor`, `resolveFont`, `boldVariant`, `parseHAlign`, `parseVAlign`, `parsePageFormat`, `contentRect` | Shared parsing/defaults for the create & compose tools: `#RRGGBB[AA]` colors, standard-14 fonts, alignment and page-format names, and the default content rectangle (page minus 36 pt margin) |

Do not duplicate this logic in individual tool files.

## Library

The PDF library is pinned at:

```
github.com/aspose-pdf-foss/aspose-pdf-foss-for-go v0.4.0
```

Consumed via `go get` (no `replace` directive in `go.mod`). Import alias: `pdf`:

```go
import pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
```

v0.4.0 ships the features behind every item on the README roadmap (PDF/A validation/conversion, digital signatures, grayscale, linearization, imposition, forms data interchange, and more), so tools for them may now be added following the conventions above. Do not add tools for library features that are not present in the pinned release tag.

## Running tests

Before every commit, all three commands must be green:

```bash
go build ./...
go vet ./...
go test ./...
```

The test suite is in `internal/tools/tools_test.go` and uses `newTestSession` to start an in-process MCP server and call tools end-to-end. Tests require real PDF fixtures; see the `testdata/` directory.
