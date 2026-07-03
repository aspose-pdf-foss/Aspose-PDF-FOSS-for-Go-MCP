# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- 5 create & compose tools (30 tools total), so an AI client can build a PDF
  from scratch instead of only transforming existing files — the need surfaced
  by the AI-editor pillar of the ecosystem:
  - `pdf_create` — new blank PDF: a4/a3/letter/legal or a custom size in
    points, portrait or landscape, N pages.
  - `pdf_add_text` — word-wrapped text in a rectangle: standard-14 font, size,
    `#RRGGBB[AA]` color, horizontal/vertical alignment, rotation.
  - `pdf_add_table` — table of text cells: column widths, borders, bold
    repeating header rows; long tables continue onto auto-appended pages.
  - `pdf_add_image` — place a PNG/JPEG into a rectangle on a page.
  - `pdf_draw` — vector shapes (line, rectangle, circle, ellipse) with stroke
    color/width, dash pattern, and optional fill.

## [0.2.0] — 2026-07-02

First public release. The server now wraps library v0.4.0 and exposes 25 tools.

### Added

- 19 new tools covering the library v0.4.0 feature set:
  - **Compliance & conversion** — `pdf_convert_pdfa` (convert toward PDF/A-1b…3a with a remaining-violations report), `pdf_to_grayscale`, `pdf_linearize` (fast web view), `pdf_optimize` (image recompression/downscaling + unused-object removal).
  - **Digital signatures** — `pdf_sign` (PKCS#7 detached, SHA-256; PEM certificate + PKCS#8/PKCS#1/SEC1 key; optional visible signature widget) and `pdf_verify_signatures` (validity, integrity, whole-document coverage, signer certificate).
  - **Security** — `pdf_encrypt` (AES-128 default, AES-256, RC4-128) and `pdf_decrypt`.
  - **Forms** — `pdf_form_fields` (list name/type/value/page/flags), `pdf_fill_form` (fill by full field name, optional flatten), `pdf_flatten`.
  - **Pages** — `pdf_rotate`, `pdf_delete_pages`, `pdf_extract_pages`, `pdf_nup` (grid imposition), `pdf_booklet` (saddle-stitch imposition).
  - **Content** — `pdf_extract_images`, `pdf_watermark` (text or image stamp), `pdf_search` (literal or RE2 regex, optional case folding; returns pages and bounding boxes).
- `pdf_validate` gained a `profile` parameter: `structural` (default), `pdfa-1b`…`pdfa-3a`, and `pdfua`.
- End-to-end tests for every tool; test fixtures (AcroForm PDF, AES-encrypted PDF, self-signed ECDSA certificate) are generated on the fly.
- GitHub Actions CI (build + vet + test on Linux and Windows).

### Changed

- Upgraded `aspose-pdf-foss-for-go` from v0.3.0 to v0.4.0.
- Renamed the module to `github.com/aspose-pdf-foss/aspose-pdf-foss-for-go-mcp` to match the GitHub repository (lowercased).

### Fixed

- `pdf_info` failed on encrypted PDFs even with the correct password: page sizes are now read from the opened document instead of the password-less `pdf.PageSizes`, and the `encrypted` flag is probed from the file rather than inferred from whether a password was supplied.

## [0.1.0] — 2026-06-23

Internal MVP (not tagged). A stdio MCP server wrapping library v0.3.0 with six tools: `pdf_info`, `pdf_extract_text`, `pdf_render_page`, `pdf_merge`, `pdf_split`, and `pdf_validate` (structural only).
