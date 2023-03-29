// Copyright 2023 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package poly_test

import (
	"reflect"
	"testing"

	"github.com/fzipp/pathfind/internal/poly"
)

func TestParseFloats(t *testing.T) {
	tests := []struct {
		s    string
		want []float32
	}{
		{"", []float32{}},
		{"1.2", []float32{1.2}},
		{"186.5,364.7,303.25,374,303.1,412", []float32{186.5, 364.7, 303.25, 374, 303.1, 412}},
		{"   -435.23 ,56.9  ,  867, -123,   32.4,12 ", []float32{-435.23, 56.9, 867, -123, 32.4, 12}},
		{"3.14,,5,.3", []float32{3.14, 0, 5, 0.3}},
	}
	for _, tt := range tests {
		floats := poly.ParseFloats(tt.s)
		if !reflect.DeepEqual(floats, tt.want) {
			t.Errorf("parsed floats from %q was %v, want %v", tt.s, floats, tt.want)
		}
	}
}
