// Copyright 2023 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package poly_test

import (
	"reflect"
	"testing"

	"github.com/fzipp/geom"
	"github.com/fzipp/pathfind/internal/poly"
)

func TestParsePolygons(t *testing.T) {
	tests := []struct {
		coords []string
		want   poly.PolygonSet
	}{
		{nil, poly.PolygonSet{}},
		{
			[]string{
				"0,0,10,0,10,10,0,10",
				"0.5,-1.2,3.7,5.4",
			}, poly.PolygonSet{
				poly.Polygon{
					geom.V2(0, 0),
					geom.V2(10, 0),
					geom.V2(10, 10),
					geom.V2(0, 10),
				},
				poly.Polygon{
					geom.V2(0.5, -1.2),
					geom.V2(3.7, 5.4),
				},
			},
		},
	}
	for _, tt := range tests {
		polygons := poly.ParsePolygons(tt.coords)
		if !reflect.DeepEqual(polygons, tt.want) {
			t.Errorf("ParsePolygons(%#v)\n got: %v\nwant: %v", tt.coords, polygons, tt.want)
		}
	}
}

// Two square-shaped polygons (20x20, 10x10) nested within each other.
// The origin of the coordinate system is in the center.
//
//	-20,-20 >-------+ 20,-20
//	        | >---+ |
//	        | |   | |
//	        | +---+ |
//	 -20,20 +-------+ 20,20
var twoSquaresNested = poly.PolygonSet{
	// Outer square
	poly.Polygon{
		geom.V2(-20, -20),
		geom.V2(20, -20),
		geom.V2(20, 20),
		geom.V2(-20, 20),
	},
	// Inner square
	poly.Polygon{
		geom.V2(-10, -10),
		geom.V2(10, -10),
		geom.V2(10, 10),
		geom.V2(-10, 10),
	},
}

// Three square-shaped polygons (30x30, 20x20, 10x10) nested within each other.
// The origin of the coordinate system is in the center.
//
//	-30,-30 >-----------+ 30,-30
//	        | >-------+ |
//	        | | >---+ | |
//	        | | |   | | |
//	        | | +---+ | |
//	        | +-------+ |
//	 -30,30 +-----------+ 30,30
var threeSquaresNested = poly.PolygonSet{
	twoSquaresNested[0],
	twoSquaresNested[1],
	poly.Polygon{
		geom.V2(-30, -30),
		geom.V2(30, -30),
		geom.V2(30, 30),
		geom.V2(-30, 30),
	},
}

// Two equally sized (10x10), disjoint, square-shaped polygons next to each
// other with the origin at the top left corner.
//
//	 0,0 >---+   >---+ 30,0
//	     |   |   |   |
//	0,10 +---+   +---+ 30,10
var twoDisjointSquares = poly.PolygonSet{
	poly.Polygon{
		geom.V2(0, 0),
		geom.V2(10, 0),
		geom.V2(10, 10),
		geom.V2(0, 10),
	},
	poly.Polygon{
		geom.V2(20, 0),
		geom.V2(30, 0),
		geom.V2(30, 10),
		geom.V2(20, 10),
	},
}

func TestPolygonSetContains(t *testing.T) {
	tests := []struct {
		polygonSet poly.PolygonSet
		pt         geom.Vec2
		want       bool
	}{
		{nil, geom.V2(0, 0), false},
		{twoSquaresNested, geom.V2(5, 5), false},
		{twoSquaresNested, geom.V2(15, 15), true},
		{twoSquaresNested, geom.V2(25, 25), false},
		{threeSquaresNested, geom.V2(5, 5), true},
		{threeSquaresNested, geom.V2(15, 15), false},
		{threeSquaresNested, geom.V2(25, 25), true},
		{threeSquaresNested, geom.V2(35, 35), false},
		{twoDisjointSquares, geom.V2(5, 5), true},
		{twoDisjointSquares, geom.V2(25, 5), true},
		{twoDisjointSquares, geom.V2(15, 5), false},
	}
	for _, tt := range tests {
		got := tt.polygonSet.Contains(tt.pt)
		if got != tt.want {
			t.Errorf("PolygonSet: %v\nContains(%v) = %v, want: %v",
				tt.polygonSet, tt.pt, got, tt.want)
		}
	}
}

func TestPolygonSetClosestPt(t *testing.T) {
	tests := []struct {
		polygonSet poly.PolygonSet
		pt         geom.Vec2
		want       geom.Vec2
	}{
		{twoSquaresNested, geom.V2(5, 0), geom.V2(10, 0)},
		{twoSquaresNested, geom.V2(14, 0), geom.V2(10, 0)},
		{twoSquaresNested, geom.V2(16, 0), geom.V2(20, 0)},
		{twoSquaresNested, geom.V2(25, 25), geom.V2(20, 20)},
		{twoSquaresNested, geom.V2(10, 25), geom.V2(10, 20)},
	}
	for _, tt := range tests {
		got := tt.polygonSet.ClosestPt(tt.pt)
		if got != tt.want {
			t.Errorf("PolygonSet: %v\nClosestPt(%v) = %v, want: %v",
				tt.polygonSet, tt.pt, got, tt.want)
		}
	}
}
