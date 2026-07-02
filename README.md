# aspose-pdf-foss-for-go-mcp

A stdio MCP server that wraps the [`aspose-pdf-foss-for-go`](https://github.com/aspose-pdf-foss/aspose-pdf-foss-for-go) v0.4.0 PDF library and exposes 25 PDF tools to any MCP-compatible AI client. It speaks the Model Context Protocol over standard input/output and has no HTTP server, no daemon, and no dependencies beyond the Go standard library and the two declared modules.

## Install / Build

```bash
go install github.com/aspose-pdf-foss/aspose-pdf-foss-for-go-mcp@latest
```

or build from a checkout:

```bash
git clone https://github.com/aspose-pdf-foss/Aspose-PDF-FOSS-for-Go-MCP.git
cd Aspose-PDF-FOSS-for-Go-MCP
go build -o aspose-pdf-foss-for-go-mcp .
```

On Windows the binary is `aspose-pdf-foss-for-go-mcp.exe`. There are no additional install steps; the server binary is self-contained.

> **Note:** the GitHub repository is named `Aspose-PDF-FOSS-for-Go-MCP`, while the Go module path is all-lowercase `github.com/aspose-pdf-foss/aspose-pdf-foss-for-go-mcp`. GitHub URLs are case-insensitive, so `go install` with the lowercase path works as shown above — always use the lowercase form in Go imports and `go get`/`go install` commands.

## Tools

### Inspect & extract

| Tool | Required params | Optional params | Returns |
|------|----------------|-----------------|---------|
| `pdf_info` | `input_path` | `password` | JSON: page count, per-page sizes (points), and Info metadata (title, author, subject, keywords, creator, producer, creation/mod dates, encrypted flag) |
| `pdf_extract_text` | `input_path` | `password`, `pages`, `separator` | Plain text extracted from the requested pages |
| `pdf_extract_images` | `input_path`, `output_dir` | `password` | Writes every raster image as PNG/JPEG to `output_dir`; returns the list of written files |
| `pdf_search` | `input_path`, `query` | `password`, `case_insensitive`, `regex` | JSON: match count and matches with text, page, and bounding box |
| `pdf_render_page` | `input_path`, `page`, `output_path` | `password`, `dpi`, `format`, `jpeg_quality` | Writes an image file to `output_path`; returns a confirmation message |

### Organize pages

| Tool | Required params | Optional params | Returns |
|------|----------------|-----------------|---------|
| `pdf_merge` | `input_paths` (array, ≥2), `output_path` | — | Writes the merged PDF |
| `pdf_split` | `input_path`, `output_dir` | `password` | Writes one PDF per page (`page_001.pdf`, …); returns the list of written files |
| `pdf_extract_pages` | `input_path`, `output_path`, `pages` | `password` | Writes a PDF with only the selected pages (ascending order) |
| `pdf_delete_pages` | `input_path`, `output_path`, `pages` | `password` | Writes a PDF without the selected pages (at least one page must remain) |
| `pdf_rotate` | `input_path`, `output_path`, `angle` (90/180/270) | `password`, `pages` | Writes a PDF with the selected pages rotated clockwise |
| `pdf_nup` | `input_path`, `output_path`, `rows`, `cols` | `password`, `page_width`, `page_height`, `margin`, `gutter`, `draw_border` | Writes an N-up imposed PDF (grid of pages per sheet) |
| `pdf_booklet` | `input_path`, `output_path` | `password`, `page_width`, `page_height` | Writes a two-up saddle-stitch booklet PDF |

### Compliance & conversion

| Tool | Required params | Optional params | Returns |
|------|----------------|-----------------|---------|
| `pdf_validate` | `input_path` | `profile`, `password` | JSON: `valid`, `profile`, and `issues` array; profiles: `structural` (default), `pdfa-1b`…`pdfa-3a`, `pdfua` |
| `pdf_convert_pdfa` | `input_path`, `output_path`, `level` | `password` | Writes the converted PDF; returns a JSON report of remaining violations |
| `pdf_to_grayscale` | `input_path`, `output_path` | `password` | Writes a grayscale PDF (text, graphics, images, shadings, annotations) |
| `pdf_linearize` | `input_path`, `output_path` | `password` | Writes a linearized (fast web view) PDF |
| `pdf_optimize` | `input_path`, `output_path` | `password`, `max_dpi`, `jpeg_quality`, `convert_png_to_jpeg` | Writes a size-optimized PDF; returns JSON with image/object counts and byte sizes |

### Security & signatures

| Tool | Required params | Optional params | Returns |
|------|----------------|-----------------|---------|
| `pdf_encrypt` | `input_path`, `output_path`, `user_password` and/or `owner_password` | `algorithm`, `password` | Writes an encrypted PDF (AES-128 default, AES-256, RC4-128) |
| `pdf_decrypt` | `input_path`, `output_path`, `password` | — | Writes a decrypted (plaintext) PDF |
| `pdf_sign` | `input_path`, `output_path`, `cert_path`, `key_path` | `password`, `reason`, `location`, `contact_info`, `signer_name`, `visible`, `page`, `rect` | Writes a digitally signed PDF (PKCS#7 detached, SHA-256) |
| `pdf_verify_signatures` | `input_path` | `password` | JSON: one entry per signature with validity, integrity, coverage, and signer certificate subject |

### Forms

| Tool | Required params | Optional params | Returns |
|------|----------------|-----------------|---------|
| `pdf_form_fields` | `input_path` | `password` | JSON: every AcroForm field with name, type, value, page, read-only/required flags |
| `pdf_fill_form` | `input_path`, `output_path`, `values` (object) | `flatten`, `password` | Writes the filled PDF; errors on unknown field names |
| `pdf_flatten` | `input_path`, `output_path` | `password` | Writes a PDF with all form fields baked into static content |

### Watermarks

| Tool | Required params | Optional params | Returns |
|------|----------------|-----------------|---------|
| `pdf_watermark` | `input_path`, `output_path`, `text` or `image_path` | `password`, `opacity`, `rotate_angle`, `font_size`, `behind`, `pages` | Writes a watermarked PDF |

### Parameter notes

- **Page numbers are 1-based** everywhere. Page ranges use the syntax `"1-3,5"`.
- `password` opens an encrypted input PDF; tools that write output produce unencrypted files unless stated otherwise (`pdf_encrypt` encrypts, `pdf_sign` preserves input encryption).
- `pdf_validate` `profile` values: `structural` (parseable and internally consistent; default), `pdfa-1b`, `pdfa-2b`, `pdfa-3b`, `pdfa-1a`, `pdfa-2a`, `pdfa-3a` (ISO 19005 conformance), `pdfua` (PDF/UA-1 accessibility prerequisites).
- `pdf_convert_pdfa` `level` takes the same `pdfa-*` values. Conversion removes encryption/JavaScript/Launch actions, embeds non-embedded simple fonts, adds an sRGB OutputIntent, and writes the `pdfaid` XMP packet; the JSON report lists anything it could not fix.
- `pdf_sign` reads the certificate and private key from PEM files. `cert_path` may contain extra `CERTIFICATE` blocks (embedded as the chain); `key_path` accepts PKCS#8 (`PRIVATE KEY`), PKCS#1 (`RSA PRIVATE KEY`), or SEC1 (`EC PRIVATE KEY`). A visible signature needs `visible: true` and `rect: [llx, lly, urx, ury]` (points, PDF user space).
- `pdf_fill_form` `values` is an object keyed by full (dotted) field name; use `pdf_form_fields` to discover names. Checkbox values are `On`/`Off`.
- `pdf_watermark` requires exactly one of `text` or `image_path`. `opacity` defaults to 0.3, `font_size` to 48; `rotate_angle` is degrees counter-clockwise (e.g. 45 for a diagonal watermark); `behind: true` draws under the page content.
- `pdf_nup`/`pdf_booklet` sheet size defaults: A4 for N-up, twice the source page width for booklet; override with `page_width`/`page_height` (points).

## MCP Client Configuration

Add the following to your MCP client's configuration (e.g. Claude Desktop `claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "aspose-pdf": {
      "command": "/absolute/path/to/aspose-pdf-foss-for-go-mcp"
    }
  }
}
```

On Windows use the `.exe` path and forward slashes or escaped back slashes:

```json
{
  "mcpServers": {
    "aspose-pdf": {
      "command": "C:/Users/you/bin/aspose-pdf-foss-for-go-mcp.exe"
    }
  }
}
```

## Roadmap

Possible future tools on top of library features that are already released: change password, reorder pages, add/replace images, outlines and table-of-contents generation, form-data interchange (JSON/FDF/XFDF), annotations, page labels, and tagged-PDF authoring.

## License

MIT. See [LICENSE](LICENSE).
