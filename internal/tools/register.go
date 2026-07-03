// SPDX-License-Identifier: MIT

package tools

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds every PDF tool to the server.
func Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_info",
		Description: "Inspect a PDF: page count, per-page sizes (points), and Info metadata (title, author, etc.).",
	}, handleInfo)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_extract_text",
		Description: "Extract text from a PDF, one section per page. Optional page range and separator.",
	}, handleText)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_render_page",
		Description: "Render a single PDF page to a PNG or JPEG image file on disk.",
	}, handleRender)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_merge",
		Description: "Merge two or more PDFs into a single output file, in the order given.",
	}, handleMerge)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_split",
		Description: "Split a PDF into one file per page, written to output_dir as page_001.pdf, page_002.pdf, ...",
	}, handleSplit)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_validate",
		Description: "Validate a PDF: structural integrity (default), PDF/A conformance (pdfa-1b ... pdfa-3a), or PDF/UA accessibility (pdfua).",
	}, handleValidate)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_convert_pdfa",
		Description: "Convert a PDF toward a PDF/A conformance level (pdfa-1b ... pdfa-3a) and report any remaining violations.",
	}, handleConvertPDFA)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_sign",
		Description: "Digitally sign a PDF (PKCS#7 detached, SHA-256) with a PEM certificate and private key; optionally draw a visible signature.",
	}, handleSign)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_verify_signatures",
		Description: "Verify every digital signature in a PDF: integrity, whole-document coverage, signer certificate.",
	}, handleVerifySignatures)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_to_grayscale",
		Description: "Convert a PDF to grayscale: text, vector graphics, images, shadings, and annotations.",
	}, handleGrayscale)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_linearize",
		Description: "Save a PDF linearized for fast web view (page 1 renders before the whole file downloads).",
	}, handleLinearize)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_nup",
		Description: "Impose multiple pages per sheet in a rows x cols grid (e.g. 2x2 for 4-up handouts).",
	}, handleNUp)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_booklet",
		Description: "Impose a PDF two-up reordered for saddle-stitch booklet printing (print double-sided, fold, read in order).",
	}, handleBooklet)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_encrypt",
		Description: "Encrypt a PDF with user/owner passwords (AES-128 default, AES-256, or RC4-128).",
	}, handleEncrypt)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_decrypt",
		Description: "Remove encryption from a PDF (requires the password that opens it).",
	}, handleDecrypt)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_extract_images",
		Description: "Extract all raster images from a PDF into a directory as PNG/JPEG files.",
	}, handleExtractImages)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_form_fields",
		Description: "List a PDF's AcroForm fields: name, type, value, page, read-only and required flags.",
	}, handleFormFields)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_fill_form",
		Description: "Fill AcroForm fields by full field name and save; optionally flatten the result.",
	}, handleFillForm)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_flatten",
		Description: "Bake all form fields into static page content and remove the AcroForm (renders the same, no longer fillable).",
	}, handleFlatten)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_optimize",
		Description: "Reduce PDF file size: recompress/downscale images and remove unreachable objects.",
	}, handleOptimize)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_rotate",
		Description: "Rotate all or selected pages clockwise by 90, 180, or 270 degrees.",
	}, handleRotate)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_delete_pages",
		Description: "Delete the given pages from a PDF (at least one page must remain).",
	}, handleDeletePages)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_extract_pages",
		Description: "Extract the given pages into a new PDF, in ascending order.",
	}, handleExtractPages)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_watermark",
		Description: "Add a translucent text or image watermark to all or selected pages.",
	}, handleWatermark)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_search",
		Description: "Find text in a PDF (literal or RE2 regex, optionally case-insensitive); returns page numbers and bounding boxes.",
	}, handleSearch)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_create",
		Description: "Create a new blank PDF: a4/a3/letter/legal or a custom size in points, portrait or landscape, N pages.",
	}, handleCreate)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_add_text",
		Description: "Draw text into a rectangle on a page: standard-14 font, size, color, alignment, rotation; word-wraps inside the rectangle.",
	}, handleAddText)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_add_table",
		Description: "Draw a table of text cells on a page: column widths, borders, optional bold header rows; long tables continue onto appended pages.",
	}, handleAddTable)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_add_image",
		Description: "Place a PNG or JPEG image into a rectangle on a page (stretched to fill the rectangle).",
	}, handleAddImage)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pdf_draw",
		Description: "Draw a vector shape on a page: line, rectangle, circle, or ellipse with stroke and optional fill.",
	}, handleDraw)
}
