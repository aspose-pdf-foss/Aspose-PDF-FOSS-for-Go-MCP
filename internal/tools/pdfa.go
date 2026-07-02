// SPDX-License-Identifier: MIT

package tools

import (
	"fmt"
	"strings"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
)

// parsePDFALevel maps a level string such as "pdfa-1b" to the library's
// PDFAFormat constant.
func parsePDFALevel(level string) (pdf.PDFAFormat, error) {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "pdfa-1b":
		return pdf.PDFA1B, nil
	case "pdfa-2b":
		return pdf.PDFA2B, nil
	case "pdfa-3b":
		return pdf.PDFA3B, nil
	case "pdfa-1a":
		return pdf.PDFA1A, nil
	case "pdfa-2a":
		return pdf.PDFA2A, nil
	case "pdfa-3a":
		return pdf.PDFA3A, nil
	default:
		return 0, fmt.Errorf("unknown PDF/A level %q (use pdfa-1b, pdfa-2b, pdfa-3b, pdfa-1a, pdfa-2a, or pdfa-3a)", level)
	}
}
