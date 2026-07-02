# Contributing

Thanks for your interest in **aspose-pdf-foss-for-go-mcp** — a stdio MCP server
that exposes the pure-Go [`aspose-pdf-foss-for-go`](https://github.com/aspose-pdf-foss/aspose-pdf-foss-for-go)
PDF library to MCP-compatible AI clients.

Right now the most valuable contribution is a good **bug report** — especially a
real-world PDF that a tool handles incorrectly. See
[Reporting bugs](#reporting-bugs) below.

The rest of this document records the project's conventions and ground rules.
They apply to any code change and are a useful reference if you are reading the
source or proposing a fix.

## Ground rules

These are hard requirements for any change to the server:

- **Thin wrapper only.** No PDF internals are implemented here. Every tool is a
  thin adapter: validate inputs, call a library function, format the result as
  MCP tool output. If a feature needs PDF logic the library doesn't provide,
  it belongs in the [library](https://github.com/aspose-pdf-foss/aspose-pdf-foss-for-go),
  not in this repo.
- **Released library features only.** Tools may use only APIs present in the
  library version pinned in `go.mod`. Do not build against unreleased library
  code or add a `replace` directive.
- **Minimal dependencies.** The server depends on the MCP SDK and the PDF
  library. Do not add further `require` entries to `go.mod` without a very
  good reason.
- **SPDX license header.** Every `.go` file starts with the line
  `// SPDX-License-Identifier: MIT` on its own line, a blank line, then the
  `package` declaration.
- **One tool = one file.** Each tool lives in its own file under
  `internal/tools/`, named after the tool without the `pdf_` prefix, and is
  registered only in `internal/tools/register.go`. Follow the step-by-step
  recipe in [CLAUDE.md](CLAUDE.md) ("Adding a new tool").
- **Consistent I/O conventions.** File paths are passed as `input_path` /
  `output_path` / `output_dir` parameters; text results are returned as MCP
  `TextContent`; page numbers are 1-based; page ranges use the `"1-3,5"`
  syntax via the shared `parsePageRange` helper.
- **Keep docs in sync.** If you add or change a tool, update the tool tables in
  `README.md`, the file-layout table in `CLAUDE.md`, and `CHANGELOG.md`.

## Development

```bash
go build ./...      # server and packages build
go test ./...       # full test suite must pass
go vet ./...        # must be clean
go build -o aspose-pdf-foss-for-go-mcp .   # build the server binary
```

The server speaks MCP over stdio. To try it end-to-end, register the binary in
an MCP client (see "MCP Client Configuration" in the README), or drive it
manually by writing JSON-RPC messages to its stdin.

## Tests

- Every tool gets an end-to-end test in `internal/tools/tools_test.go`, using
  the `newTestSession` helper (an in-process MCP server + client session) and
  the `callOK` / `callErr` helpers.
- Tests use the small fixture PDFs in `testdata/`; fixtures that would bloat
  the repo (an AcroForm PDF, an encrypted PDF, a signing certificate) are
  generated on the fly in test helpers instead. If you need a new committed
  test PDF, open an issue first so we can agree on the fixture.
- A bug fix should come with a regression test that fails before the fix and
  passes after.

## Reporting bugs

Open a GitHub issue with the PDF that reproduces the problem (or a minimal
synthetic one), the tool call you made (tool name + arguments), what you
expected, and what happened.

If the problem is in PDF processing itself — parsing, rendering, text
extraction, conversion output — it almost certainly belongs to the underlying
library: please report it at
[aspose-pdf-foss-for-go](https://github.com/aspose-pdf-foss/aspose-pdf-foss-for-go/issues)
instead, ideally with a direct library-API reproduction. Issues in this repo
are for the MCP layer: tool schemas, parameter validation, error messages,
result formatting, and server behavior.

## Spec

The protocol reference is the [Model Context Protocol specification](https://modelcontextprotocol.io/specification/latest);
this server targets the official [Go SDK](https://github.com/modelcontextprotocol/go-sdk).
For PDF behavior questions, the normative reference is ISO 32000-1 (PDF 1.7) —
see the library's CONTRIBUTING.md.
