// SPDX-License-Identifier: MIT

package tools

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// parsePageRange parses a 1-based page spec such as "1-3,5" against a document
// of pageCount pages. An empty spec selects all pages. The result is sorted and
// de-duplicated. Any page outside [1, pageCount] or any malformed token errors.
func parsePageRange(spec string, pageCount int) ([]int, error) {
	if strings.TrimSpace(spec) == "" {
		all := make([]int, pageCount)
		for i := range all {
			all[i] = i + 1
		}
		return all, nil
	}
	seen := map[int]bool{}
	for _, part := range strings.Split(spec, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			from, err1 := strconv.Atoi(strings.TrimSpace(bounds[0]))
			to, err2 := strconv.Atoi(strings.TrimSpace(bounds[1]))
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid range %q", part)
			}
			if from < 1 || to > pageCount || from > to {
				return nil, fmt.Errorf("range %q out of bounds (document has %d pages)", part, pageCount)
			}
			for p := from; p <= to; p++ {
				seen[p] = true
			}
			continue
		}
		p, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid page %q", part)
		}
		if p < 1 || p > pageCount {
			return nil, fmt.Errorf("page %d out of bounds (document has %d pages)", p, pageCount)
		}
		seen[p] = true
	}
	out := make([]int, 0, len(seen))
	for p := range seen {
		out = append(out, p)
	}
	sort.Ints(out)
	return out, nil
}
