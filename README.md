# aspose-pdf-foss-mcp

A stdio MCP server that wraps the [`aspose-pdf-foss-for-go`](https://github.com/aspose-pdf-foss/aspose-pdf-foss-for-go) v0.4.0 PDF library and exposes six PDF tools to any MCP-compatible AI client. It speaks the Model Context Protocol over standard input/output and has no HTTP server, no daemon, and no dependencies beyond the Go standard library and the two declared modules.

## Install / Build

```bash
go install github.com/aspose-pdf-foss/aspose-pdf-foss-mcp@latest
```

or build from a checkout:

```bash
go build -o aspose-pdf-foss-mcp .
```

On Windows the binary is `aspose-pdf-foss-mcp.exe`. There are no additional install steps; the server binary is self-contained.

## Tools

| Tool | Required params | Optional params | Returns |
|------|----------------|-----------------|---------|
| `pdf_info` | `input_path` | `password` | JSON: page count, per-page sizes (points), and Info metadata (title, author, subject, keywords, creator, producer, creation/mod dates, encrypted flag) |
| `pdf_extract_text` | `input_path` | `password`, `pages`, `separator` | Plain text extracted from the requested pages |
| `pdf_render_page` | `input_path`, `page`, `output_path` | `password`, `dpi`, `format`, `jpeg_quality` | Writes an image file to `output_path`; returns a confirmation message |
| `pdf_merge` | `input_paths` (array, ≥2), `output_path` | — | Writes the merged PDF; returns a confirmation message |
| `pdf_split` | `input_path`, `output_dir` | `password` | Writes one PDF per page (`page_001.pdf`, …) to `output_dir`; returns the list of written files |
| `pdf_validate` | `input_path` | — | JSON: `valid` (bool) and `issues` array with `code`/`message` fields |

### Parameter reference

#### `pdf_info`

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `input_path` | string | yes | Path to the input PDF file |
| `password` | string | no | Password for an encrypted PDF |

#### `pdf_extract_text`

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `input_path` | string | yes | Path to the input PDF file |
| `password` | string | no | Password for an encrypted PDF |
| `pages` | string | no | Page range, e.g. `"1-3,5"` (1-based); omit for all pages |
| `separator` | string | no | String inserted between pages; defaults to `"\n\n--- Page N ---\n"` |

#### `pdf_render_page`

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `input_path` | string | yes | Path to the input PDF file |
| `page` | integer | yes | 1-based page number to render |
| `output_path` | string | yes | Path to write the rendered image |
| `password` | string | no | Password for an encrypted PDF |
| `dpi` | number | no | Resolution in DPI; defaults to 150 |
| `format` | string | no | `"png"` or `"jpeg"`; defaults to the `output_path` extension, else `png` |
| `jpeg_quality` | integer | no | JPEG quality 1–100; defaults to 90 |

#### `pdf_merge`

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `input_paths` | string array | yes | Paths of the PDFs to merge, in order (at least two) |
| `output_path` | string | yes | Path to write the merged PDF |

#### `pdf_split`

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `input_path` | string | yes | Path to the input PDF file |
| `output_dir` | string | yes | Directory to write one PDF per page |
| `password` | string | no | Password for an encrypted PDF |

#### `pdf_validate`

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `input_path` | string | yes | Path to the PDF file to validate |

Note: `pdf_validate` checks structural integrity (parseable and internally consistent). It is **not** a PDF/A or PDF/UA conformance check.

## MCP Client Configuration

Add the following to your MCP client's configuration (e.g. Claude Desktop `claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "aspose-pdf": {
      "command": "/absolute/path/to/aspose-pdf-foss-mcp"
    }
  }
}
```

On Windows use the `.exe` path and forward slashes or escaped back slashes:

```json
{
  "mcpServers": {
    "aspose-pdf": {
      "command": "C:/Users/you/bin/aspose-pdf-foss-mcp.exe"
    }
  }
}
```

## Roadmap

The following tools are planned. The underlying features all ship in library v0.4.0; the corresponding MCP tools have not been added to this server yet.

- PDF/A validation and conversion
- Digital signatures
- Grayscale conversion
- Linearization (fast-web-view)
- N-up and booklet imposition
- Encryption and decryption
- Image extraction
- Forms: read, fill, and flatten
- File size optimization
- Rotate, delete, and extract pages
- Watermark (text and image)
- Text search

## License

MIT. See [LICENSE](LICENSE).
