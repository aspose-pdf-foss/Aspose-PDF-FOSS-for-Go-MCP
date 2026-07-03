// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"fmt"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type addTableParams struct {
	InputPath    string     `json:"input_path" jsonschema:"path to the input PDF file"`
	OutputPath   string     `json:"output_path" jsonschema:"path to write the modified PDF (may equal input_path to overwrite)"`
	Password     string     `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
	Page         int        `json:"page,omitempty" jsonschema:"1-based page number; defaults to 1"`
	Rows         [][]string `json:"rows" jsonschema:"table cells: array of rows, each an array of cell strings; every row must have the same number of cells"`
	HeaderRows   int        `json:"header_rows,omitempty" jsonschema:"number of leading rows styled bold and repeated at the top of continuation pages"`
	ColumnWidths []float64  `json:"column_widths,omitempty" jsonschema:"column widths in points, one per column; defaults to equal widths across the rectangle"`
	LLX          float64    `json:"llx,omitempty" jsonschema:"rectangle lower-left x in points (PDF coordinates: origin at the page's bottom-left, y grows upward); omit all four coordinates to use the page minus a 36 pt margin"`
	LLY          float64    `json:"lly,omitempty" jsonschema:"rectangle lower-left y in points"`
	URX          float64    `json:"urx,omitempty" jsonschema:"rectangle upper-right x in points"`
	URY          float64    `json:"ury,omitempty" jsonschema:"rectangle upper-right y in points"`
	Font         string     `json:"font,omitempty" jsonschema:"standard-14 font name for cell text; defaults to helvetica"`
	FontSize     float64    `json:"font_size,omitempty" jsonschema:"cell font size in points; defaults to 10"`
	BorderWidth  float64    `json:"border_width,omitempty" jsonschema:"table and cell border width in points; defaults to 0.5"`
	NoBorders    bool       `json:"no_borders,omitempty" jsonschema:"draw no table or cell borders at all"`
}

func handleAddTable(ctx context.Context, req *mcp.CallToolRequest, args addTableParams) (*mcp.CallToolResult, any, error) {
	if len(args.Rows) == 0 {
		return nil, nil, fmt.Errorf("rows is required and must not be empty")
	}
	cols := len(args.Rows[0])
	if cols == 0 {
		return nil, nil, fmt.Errorf("rows[0] has no cells")
	}
	for i, r := range args.Rows {
		if len(r) != cols {
			return nil, nil, fmt.Errorf("row %d has %d cells, want %d (all rows must have the same number of cells)", i+1, len(r), cols)
		}
	}
	if args.HeaderRows < 0 || args.HeaderRows > len(args.Rows) {
		return nil, nil, fmt.Errorf("header_rows %d out of range (0..%d)", args.HeaderRows, len(args.Rows))
	}
	font, err := resolveFont(args.Font)
	if err != nil {
		return nil, nil, err
	}
	fontSize := args.FontSize
	if fontSize == 0 {
		fontSize = 10
	}

	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	pageNum := args.Page
	if pageNum == 0 {
		pageNum = 1
	}
	page, err := doc.Page(pageNum)
	if err != nil {
		return nil, nil, err
	}
	rect, err := contentRect(page, args.LLX, args.LLY, args.URX, args.URY)
	if err != nil {
		return nil, nil, err
	}

	widths := args.ColumnWidths
	if len(widths) == 0 {
		w := (rect.URX - rect.LLX) / float64(cols)
		widths = make([]float64, cols)
		for i := range widths {
			widths[i] = w
		}
	} else if len(widths) != cols {
		return nil, nil, fmt.Errorf("column_widths has %d entries, want %d (one per column)", len(widths), cols)
	}

	table := pdf.NewTable().
		SetColumnWidths(widths).
		SetDefaultCellMargin(pdf.MarginInfo{Top: 4, Right: 6, Bottom: 4, Left: 6}).
		SetDefaultCellStyle(pdf.TextStyle{Font: font, Size: fontSize})
	if !args.NoBorders {
		borderWidth := args.BorderWidth
		if borderWidth == 0 {
			borderWidth = 0.5
		}
		border := pdf.BorderInfo{Sides: pdf.BorderSideAll, Width: borderWidth}
		table.SetBorder(border).SetDefaultCellBorder(border)
	}
	rows := table.AddRows(args.Rows)
	if args.HeaderRows > 0 {
		table.SetRepeatingRowsCount(args.HeaderRows)
		headerStyle := pdf.TextStyle{Font: boldVariant(font), Size: fontSize}
		for _, r := range rows[:args.HeaderRows] {
			r.SetTextStyle(headerStyle)
		}
	}

	pagesAdded, err := page.AddTable(table, rect)
	if err != nil {
		return nil, nil, fmt.Errorf("add table: %w", err)
	}
	if err := doc.Save(args.OutputPath); err != nil {
		return nil, nil, fmt.Errorf("save %q: %w", args.OutputPath, err)
	}
	msg := fmt.Sprintf("Added a %dx%d table to page %d of %s: %s.", len(args.Rows), cols, pageNum, args.InputPath, args.OutputPath)
	if pagesAdded > 0 {
		msg += fmt.Sprintf(" The table continued onto %d appended page(s).", pagesAdded)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}
