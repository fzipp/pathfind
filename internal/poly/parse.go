// Copyright 2023 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package poly

import (
	"strconv"
	"strings"
)

// ParseFloats parses a slice of float32s from a comma-separated
// string of numbers, for example "186.5,364.7,303.25,374,303.1,412".
// Spaces are ignored.
func ParseFloats(s string) []float32 {
	tokens := strings.Split(s, ",")
	floats := make([]float32, 0, len(tokens))
	for i, tok := range tokens {
		tok = strings.TrimSpace(tok)
		if i == 0 && len(tok) == 0 {
			break
		}
		f, _ := strconv.ParseFloat(tok, 32)
		floats = append(floats, float32(f))
	}
	return floats
}
