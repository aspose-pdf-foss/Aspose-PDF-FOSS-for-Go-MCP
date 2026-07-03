// SPDX-License-Identifier: MIT

package tools

import (
	"fmt"
	"strconv"
	"strings"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
)

// parseHexColor converts "#RRGGBB" or "#RRGGBBAA" (leading '#' optional) into
// a library Color. An empty string returns nil so callers can apply their own
// default.
func parseHexColor(s string) (*pdf.Color, error) {
	if s == "" {
		return nil, nil
	}
	h := strings.TrimPrefix(s, "#")
	if len(h) != 6 && len(h) != 8 {
		return nil, fmt.Errorf("color %q: want #RRGGBB or #RRGGBBAA", s)
	}
	v, err := strconv.ParseUint(h, 16, 64)
	if err != nil {
		return nil, fmt.Errorf("color %q: %w", s, err)
	}
	c := &pdf.Color{A: 1}
	if len(h) == 8 {
		c.A = float64(v&0xff) / 255
		v >>= 8
	}
	c.B = float64(v&0xff) / 255
	c.G = float64((v>>8)&0xff) / 255
	c.R = float64((v>>16)&0xff) / 255
	return c, nil
}

// resolveFont maps a standard-14 PostScript name to a Font; empty → Helvetica.
func resolveFont(name string) (pdf.Font, error) {
	if name == "" {
		return pdf.FontHelvetica, nil
	}
	return pdf.FindFont(name)
}

// boldVariant returns the bold counterpart of a standard-14 font, or the font
// itself when there is none (Symbol, ZapfDingbats) or it is already bold.
func boldVariant(f pdf.Font) pdf.Font {
	switch f {
	case pdf.FontHelvetica:
		return pdf.FontHelveticaBold
	case pdf.FontHelveticaOblique:
		return pdf.FontHelveticaBoldOblique
	case pdf.FontTimesRoman:
		return pdf.FontTimesBold
	case pdf.FontTimesItalic:
		return pdf.FontTimesBoldItalic
	case pdf.FontCourier:
		return pdf.FontCourierBold
	case pdf.FontCourierOblique:
		return pdf.FontCourierBoldOblique
	}
	return f
}

func parseHAlign(s string) (pdf.HAlign, error) {
	switch strings.ToLower(s) {
	case "", "left":
		return pdf.HAlignLeft, nil
	case "center":
		return pdf.HAlignCenter, nil
	case "right":
		return pdf.HAlignRight, nil
	}
	return 0, fmt.Errorf("halign %q: want left, center, or right", s)
}

func parseVAlign(s string) (pdf.VAlign, error) {
	switch strings.ToLower(s) {
	case "", "top":
		return pdf.VAlignTop, nil
	case "middle":
		return pdf.VAlignMiddle, nil
	case "bottom":
		return pdf.VAlignBottom, nil
	}
	return 0, fmt.Errorf("valign %q: want top, middle, or bottom", s)
}

// parsePageFormat maps a format name to the library PageFormat; empty → A4.
func parsePageFormat(name string) (pdf.PageFormat, error) {
	switch strings.ToLower(name) {
	case "", "a4":
		return pdf.PageFormatA4, nil
	case "a3":
		return pdf.PageFormatA3, nil
	case "letter":
		return pdf.PageFormatLetter, nil
	case "legal":
		return pdf.PageFormatLegal, nil
	}
	return pdf.PageFormat{}, fmt.Errorf("page_format %q: want a4, a3, letter, or legal", name)
}

// contentRect returns the rectangle given by the four coordinates, or, when
// all four are zero, the page rectangle inset by a 36 pt margin.
func contentRect(page *pdf.Page, llx, lly, urx, ury float64) (pdf.Rectangle, error) {
	if llx == 0 && lly == 0 && urx == 0 && ury == 0 {
		size, err := page.Size()
		if err != nil {
			return pdf.Rectangle{}, err
		}
		const margin = 36
		return pdf.Rectangle{LLX: margin, LLY: margin, URX: size.Width - margin, URY: size.Height - margin}, nil
	}
	if urx <= llx || ury <= lly {
		return pdf.Rectangle{}, fmt.Errorf("invalid rectangle [%g %g %g %g]: urx must exceed llx and ury must exceed lly", llx, lly, urx, ury)
	}
	return pdf.Rectangle{LLX: llx, LLY: lly, URX: urx, URY: ury}, nil
}
