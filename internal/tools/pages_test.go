// SPDX-License-Identifier: MIT

package tools

import (
	"reflect"
	"testing"
)

func TestParsePageRange(t *testing.T) {
	tests := []struct {
		name      string
		spec      string
		pageCount int
		want      []int
		wantErr   bool
	}{
		{"empty returns all", "", 3, []int{1, 2, 3}, false},
		{"single", "2", 4, []int{2}, false},
		{"range", "1-3", 4, []int{1, 2, 3}, false},
		{"mixed dedup sorted", "3,1-2,2", 5, []int{1, 2, 3}, false},
		{"out of range", "5", 4, nil, true},
		{"zero invalid", "0", 4, nil, true},
		{"reversed range", "3-1", 4, nil, true},
		{"garbage", "a-b", 4, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePageRange(tt.spec, tt.pageCount)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
