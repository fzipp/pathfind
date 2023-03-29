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

func TestParsePolygon(t *testing.T) {
	tests := []struct {
		coords string
		want   poly.Polygon
	}{
		{"", poly.Polygon{}},
		{"1.2,3.4", poly.Polygon{
			geom.V2(1.2, 3.4),
		}},
		{"132.7,-234.3,11.34,982,112.2,932", poly.Polygon{
			geom.V2(132.7, -234.3),
			geom.V2(11.34, 982),
			geom.V2(112.2, 932),
		}},
		{"    -224.33 ,43.7  ,  37, -13,   -32.4,9 , 99,-1,  34", poly.Polygon{
			geom.V2(-224.33, 43.7),
			geom.V2(37, -13),
			geom.V2(-32.4, 9),
			geom.V2(99, -1),
		}},
		{"0,.3,,-1,,4,2", poly.Polygon{
			geom.V2(0, 0.3),
			geom.V2(0, -1),
			geom.V2(0, 4),
		}},
	}
	for _, tt := range tests {
		polygon := poly.ParsePolygon(tt.coords)
		if !reflect.DeepEqual(polygon, tt.want) {
			t.Errorf("ParsePolygon(%q)\n got: %v\nwant: %v", tt.coords, polygon, tt.want)
		}
	}
}

func TestPolygonEdge(t *testing.T) {
	tests := []struct {
		polygon    poly.Polygon
		edgeNumber int
		want       poly.LineSeg
	}{
		{
			poly.Polygon{
				geom.V2(2.5, 3),
				geom.V2(1, 3.2),
				geom.V2(4, 5.1),
			},
			0,
			poly.LineSeg{A: geom.V2(2.5, 3), B: geom.V2(1, 3.2)},
		},
		{
			poly.ParsePolygon("4.5,3,-2.7,5,9,5"), 1,
			poly.LineSeg{A: geom.V2(-2.7, 5), B: geom.V2(9, 5)},
		},
	}
	for _, tt := range tests {
		if edge := tt.polygon.Edge(tt.edgeNumber); edge != tt.want {
			t.Errorf("Polygon: %v\nEdge(%d) = %v, want %v", tt.polygon, tt.edgeNumber, edge, tt.want)
		}
	}
}

// A square-shaped polygon with the origin at the top left corner.
//
//	 0,0 >---+ 10,0
//	     |   |
//	0,10 +---+ 10,10
var polygonSquare = poly.Polygon{
	geom.V2(0, 0),
	geom.V2(10, 0),
	geom.V2(10, 10),
	geom.V2(0, 10),
}

// A diamond/rhombus-shaped polygon.
//
//	       > 5,0
//	     /   \
//	0,5 +     + 10,5
//	     \   /
//	       + 5,10
var polygonDiamond = poly.Polygon{
	geom.V2(5, 0),
	geom.V2(10, 5),
	geom.V2(5, 10),
	geom.V2(0, 5),
}

func TestPolygonContains(t *testing.T) {
	tests := []struct {
		polygon            poly.Polygon
		point              geom.Vec2
		toleranceOnOutside bool
		want               bool
	}{
		//   +---+
		//   + x |
		//   +---+
		{polygonSquare, geom.V2(5, 5), false, true},
		//   +---+
		//   +   | x
		//   +---+
		{polygonSquare, geom.V2(15, 5), false, false},
		// x +---+
		//   +   |
		//   +---+
		{polygonSquare, geom.V2(-5, 0), false, false},
		//   +-x-+
		//   +   |
		//   +---+
		{polygonSquare, geom.V2(5, 0), true, true},
		//   +-x-+
		//   +   |
		//   +---+
		{polygonSquare, geom.V2(5, 0), false, false},
		//   +---+
		//   +   x
		//   +---+
		{polygonSquare, geom.V2(10, 5), false, false},
		//   +---+
		//   +   x
		//   +---+
		{polygonSquare, geom.V2(10, 5), true, true},
	}
	for _, tt := range tests {
		got := tt.polygon.Contains(tt.point, tt.toleranceOnOutside)
		if got != tt.want {
			t.Errorf("Polygon: %v\nContains(pt: %v, toleranceOnOutside: %v) = %v, want: %v",
				tt.polygon, tt.point, tt.toleranceOnOutside, got, tt.want)
		}
	}
}

func TestPolygonIsCrossedBy(t *testing.T) {
	tests := []struct {
		name    string
		polygon poly.Polygon
		ls      poly.LineSeg
		want    bool
	}{
		{
			"square crossed by line through left side",
			//   +---+
			// --+-- |
			//   +---+
			polygonSquare,
			poly.LineSeg{A: geom.V2(-5, 5), B: geom.V2(5, 5)},
			true,
		},
		{
			"square crossed by line through right side",
			//   +---+
			//   + --+--
			//   +---+
			polygonSquare,
			poly.LineSeg{A: geom.V2(5, 5), B: geom.V2(15, 5)},
			true,
		},
		{
			"square crossed by line through both vertical sides",
			//   +---+
			// --+---+--
			//   +---+
			polygonSquare,
			poly.LineSeg{A: geom.V2(-5, 5), B: geom.V2(15, 5)},
			true,
		},
		{
			"square crossed by line through top side",
			//     |
			//   +-+-+
			//   | | |
			//   +---+
			polygonSquare,
			poly.LineSeg{A: geom.V2(5, -5), B: geom.V2(5, 5)},
			true,
		},
		{
			"square crossed by diagonal line through left and right sides",
			//
			//  \+---+
			//   | \ |
			//   +---+\
			polygonSquare,
			poly.LineSeg{A: geom.V2(-5, 0), B: geom.V2(15, 10)},
			true,
		},
		{
			"square crossed by diagonal line through corners",
			// \
			//  +---+
			//  | \ |
			//  +---+
			//       \
			polygonSquare,
			poly.LineSeg{A: geom.V2(-5, -5), B: geom.V2(15, 15)},
			true,
		},
		{
			"square not crossed by outside line parallel to side",
			//   +---+
			//   |   |
			//   +---+
			// ---------
			polygonSquare,
			poly.LineSeg{A: geom.V2(-5, 15), B: geom.V2(20, 15)},
			false,
		},
		{
			"square not crossed by line only touching on corner",
			//    /
			//   +---+
			//  /|   |
			//   +---+
			polygonSquare,
			poly.LineSeg{A: geom.V2(-5, 5), B: geom.V2(5, -5)},
			false,
		},
		{
			"square not crossed by line on top of side",
			//   +===+
			//   |   |
			//   +---+
			polygonSquare,
			poly.LineSeg{A: geom.V2(0, 0), B: geom.V2(10, 0)},
			false,
		},
		{
			"triangle not crossed by line touching on corner",
			//      |
			// +----+
			// |   /|
			// | /
			// +
			poly.Polygon{
				geom.V2(0, 0),
				geom.V2(10, 0),
				geom.V2(0, 10),
			},
			poly.LineSeg{A: geom.V2(10, -5), B: geom.V2(10, 5)},
			false,
		},
		{
			"polygon crossed by line through flat vertex",
			//       |
			// +-----+-----+
			// |     |     |
			// +-----------+
			poly.Polygon{
				geom.V2(0, 0),
				geom.V2(10, 0),
				geom.V2(20, 0),
				geom.V2(20, 10),
				geom.V2(0, 10),
			},
			poly.LineSeg{A: geom.V2(10, -5), B: geom.V2(10, 5)},
			true,
		},
		{
			"polygon crossed by line through convex vertex",
			//    |
			//    +
			//   /|\
			//  /   \
			// +-----+
			poly.Polygon{
				geom.V2(10, 0),
				geom.V2(20, 10),
				geom.V2(0, 10),
			},
			poly.LineSeg{A: geom.V2(10, -5), B: geom.V2(10, 5)},
			true,
		},
		{
			"polygon crossed by line through concave vertex",
			// +       +
			// | \ | / |
			// |   +   |
			// |   |   |
			// +-------+
			poly.Polygon{
				geom.V2(0, -10),
				geom.V2(10, 0),
				geom.V2(20, -10),
				geom.V2(20, 10),
				geom.V2(0, 10),
			},
			poly.LineSeg{A: geom.V2(10, -5), B: geom.V2(10, 5)},
			true,
		},
		{
			"polygon crossed by line through another concave vertex",
			//         +
			//     | / |
			// +---+   |
			// |   |   |
			// +-------+
			poly.Polygon{
				geom.V2(0, 0),
				geom.V2(10, 0),
				geom.V2(20, -10),
				geom.V2(20, 10),
				geom.V2(0, 10),
			},
			poly.LineSeg{A: geom.V2(10, -5), B: geom.V2(10, 5)},
			true,
		},
		{
			"touch polygon corner and cross edge",
			poly.Polygon{
				geom.V2(0, 0),
				geom.V2(30, 0),
				geom.V2(30, 10),
				geom.V2(20, 10),
				geom.V2(20, 25),
				geom.V2(0, 25),
			},
			poly.LineSeg{A: geom.V2(10, 20), B: geom.V2(40, 5)},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.polygon.IsCrossedBy(tt.ls)
			if got != tt.want {
				t.Errorf("Polygon: %v\nIsCrossedBy(%v) = %v, want: %v",
					tt.polygon, tt.ls, got, tt.want)
			}
		})
	}
}

func TestPolygonClosestPt(t *testing.T) {
	tests := []struct {
		polygon poly.Polygon
		pt      geom.Vec2
		want    geom.Vec2
	}{
		//   +---+
		// x |   |
		//   +---+
		{polygonSquare, geom.V2(-5, 5), geom.V2(0, 5)},
		//     x
		//   +---+
		//   |   |
		//   +---+
		{polygonSquare, geom.V2(5, -5), geom.V2(5, 0)},
		// x
		//   +---+
		//   |   |
		//   +---+
		{polygonSquare, geom.V2(-5, -5), geom.V2(0, 0)},
		//   +---+
		//   x   |
		//   +---+
		{polygonSquare, geom.V2(0, 5), geom.V2(0, 5)},
		//   +---+
		//   |x  |
		//   +---+
		{polygonSquare, geom.V2(2.5, 5), geom.V2(0, 5)},
		// x  +
		//  /   \
		// +     +
		//  \   /
		//    +
		{polygonDiamond, geom.V2(0, 0), geom.V2(2.5, 2.5)},
		//    +
		//  /   \
		// +     +  x
		//  \   /
		//    +
		{polygonDiamond, geom.V2(15, 5), geom.V2(10, 5)},
		//    +
		//  /   \
		// +     +
		//  \  x/
		//    +
		{polygonDiamond, geom.V2(6, 6), geom.V2(7.5, 7.5)},
	}
	for _, tt := range tests {
		got := tt.polygon.ClosestPt(tt.pt)
		if got != tt.want {
			t.Errorf("Polygon: %v\nClosestPt(%v) = %v, want: %v",
				tt.polygon, tt.pt, got, tt.want)
		}
	}
}

// A U-shaped polygon with a downward slope as the right edge.
// The origin is at the top left corner.
//
//	 0,0 >---+   +---+ 30,0
//	     |   |   |    \
//	     |   +---+     \
//	     |              \
//	0,20 +---------------+ 40,20
var polygonSlopedU = poly.Polygon{
	geom.V2(0, 0),
	geom.V2(10, 0),
	geom.V2(10, 10),
	geom.V2(20, 10),
	geom.V2(20, 0),
	geom.V2(30, 0),
	geom.V2(40, 20),
	geom.V2(0, 20),
}

// A polygon with a concave vertex on the right side.
// The origin is at the top left corner.
//
//	 0,0 >-------+ 20,0
//	     |     /
//	     |   + 10,10
//	     |     \
//	0,20 +-------+ 20,20
var polygonK = poly.Polygon{
	geom.V2(0, 0),
	geom.V2(20, 0),
	geom.V2(10, 10),
	geom.V2(20, 20),
	geom.V2(0, 20),
}

func TestPolygonIsConcaveAt(t *testing.T) {
	tests := []struct {
		polygon poly.Polygon
		i       int
		want    bool
	}{
		{polygonSquare, 0, false},
		{polygonSquare, 1, false},
		{polygonSquare, 2, false},
		{polygonSquare, 3, false},
		{polygonSlopedU, 0, false},
		{polygonSlopedU, 1, false},
		{polygonSlopedU, 2, true},
		{polygonSlopedU, 3, true},
		{polygonSlopedU, 4, false},
		{polygonSlopedU, 5, false},
		{polygonSlopedU, 6, false},
		{polygonSlopedU, 7, false},
		{polygonK, 0, false},
		{polygonK, 1, false},
		{polygonK, 2, true},
		{polygonK, 3, false},
		{polygonK, 4, false},
	}
	for _, tt := range tests {
		got := tt.polygon.IsConcaveAt(tt.i)
		if got != tt.want {
			t.Errorf("Polygon: %v\nIsConcaveAt(%v) = %v, want: %v",
				tt.polygon, tt.i, got, tt.want)
		}
	}
}

func TestPolygonWrapIndex(t *testing.T) {
	tests := []struct {
		n    int
		i    int
		want int
	}{
		{n: 3, i: -6, want: 0},
		{n: 3, i: -5, want: 1},
		{n: 3, i: -4, want: 2},
		{n: 3, i: -3, want: 0},
		{n: 3, i: -2, want: 1},
		{n: 3, i: -1, want: 2},
		{n: 3, i: 0, want: 0},
		{n: 3, i: 1, want: 1},
		{n: 3, i: 2, want: 2},
		{n: 3, i: 3, want: 0},
		{n: 3, i: 4, want: 1},
		{n: 3, i: 5, want: 2},
		{n: 3, i: 6, want: 0},
	}
	for _, tt := range tests {
		p := make(poly.Polygon, tt.n)
		got := p.WrapIndex(tt.i)
		if got != tt.want {
			t.Errorf("Polygon[len: %v].WrapIndex(%v) = %v, want: %v",
				tt.n, tt.i, got, tt.want)
		}
	}
}
