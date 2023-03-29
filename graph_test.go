// Copyright 2023 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pathfind

import (
	"reflect"
	"testing"
)

func TestGraphNeighbours(t *testing.T) {
	g := make(graph[string])
	g.link("a", "b").link("a", "c").link("a", "d")
	g.link("b", "a").link("b", "d")
	g.link("c", "d").link("c", "b")

	tests := []struct {
		node string
		want []string
	}{
		{"a", []string{"b", "c", "d"}},
		{"b", []string{"a", "d"}},
		{"c", []string{"d", "b"}},
		{"d", nil},
	}
	for _, tt := range tests {
		got := g.Neighbours(tt.node)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Neighbors of node %q: got %v, want %v", tt.node, got, tt.want)
		}
	}
}
