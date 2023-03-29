// Copyright 2023 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package poly_test

import (
	"testing"

	"github.com/fzipp/geom"
	"github.com/fzipp/pathfind/internal/poly"
)

func TestLineSegLen(t *testing.T) {
	tests := []struct {
		lineSeg poly.LineSeg
		want    float32
	}{
		{poly.LineSeg{A: geom.V2Zero, B: geom.V2Zero}, 0},
		{poly.LineSeg{A: geom.V2Zero, B: geom.V2UnitX}, 1},
		{poly.LineSeg{A: geom.V2Zero, B: geom.V2UnitY}, 1},
		{poly.LineSeg{A: geom.V2Zero, B: geom.V2Unit}, 1.4142135},
		{poly.LineSeg{A: geom.V2UnitY, B: geom.V2UnitX}, 1.4142135},
		{poly.LineSeg{A: geom.V2(1, 3), B: geom.V2(5, 2)}, 4.1231055},
		{poly.LineSeg{A: geom.V2(2.5, -1.25), B: geom.V2(4, -3.2)}, 2.460183},
	}
	for _, tt := range tests {
		if length := tt.lineSeg.Len(); length != tt.want {
			t.Errorf("length of line segment %v was %g, want %g", tt.lineSeg, length, tt.want)
		}
	}
}

func TestLineSegClosestPt(t *testing.T) {
	tests := []struct {
		lineSeg poly.LineSeg
		p       geom.Vec2
		want    geom.Vec2
	}{
		{poly.LineSeg{A: geom.V2(-3, 2), B: geom.V2(3, 2)}, geom.V2(1.5, 4), geom.V2(1.5, 2)},
		{poly.LineSeg{A: geom.V2(-3, 2), B: geom.V2(3, 2)}, geom.V2(1.5, 3), geom.V2(1.5, 2)},
		{poly.LineSeg{A: geom.V2(-3, 2), B: geom.V2(3, 2)}, geom.V2(2.8, -2), geom.V2(2.8, 2)},
		{poly.LineSeg{A: geom.V2(-3, 2), B: geom.V2(3, 2)}, geom.V2(6, 3), geom.V2(3, 2)},
		{poly.LineSeg{A: geom.V2(-3, 2), B: geom.V2(3, 2)}, geom.V2(-4, -4), geom.V2(-3, 2)},
		{poly.LineSeg{A: geom.V2(-4, -4), B: geom.V2(4, 4)}, geom.V2(7, 6), geom.V2(4, 4)},
		{poly.LineSeg{A: geom.V2(-4, -4), B: geom.V2(4, 4)}, geom.V2Zero, geom.V2Zero},
		{poly.LineSeg{A: geom.V2(-4, -4), B: geom.V2(4, 4)}, geom.V2(4, 0), geom.V2(2, 2)},
		{poly.LineSeg{A: geom.V2(-4, -4), B: geom.V2(4, 4)}, geom.V2(0, 4), geom.V2(2, 2)},
	}
	for _, tt := range tests {
		if closestPt := tt.lineSeg.ClosestPt(tt.p); !closestPt.NearEq(tt.want) {
			t.Errorf("closest point on line segment %v to %v was %v, want %v", tt.lineSeg, tt.p, closestPt, tt.want)
		}
	}
}

func TestLineSegCrosses(t *testing.T) {
	tests := []struct {
		name string
		l1   poly.LineSeg
		l2   poly.LineSeg
		want bool
	}{
		{
			"line segments on top of each other don't cross",
			poly.LineSeg{A: geom.V2(-3, 2), B: geom.V2(3, 2)},
			poly.LineSeg{A: geom.V2(-3, 2), B: geom.V2(3, 2)},
			false,
		},
		{
			"parallel line segments don't cross",
			poly.LineSeg{A: geom.V2(1, 3), B: geom.V2(5, 7)},
			poly.LineSeg{A: geom.V2(1, 4), B: geom.V2(5, 8)},
			false,
		},
		{
			"perpendicular line segments with gap don't cross",
			//      |
			// ---- |
			//      |
			poly.LineSeg{A: geom.V2(0, 0), B: geom.V2(4, 0)},
			poly.LineSeg{A: geom.V2(5, 2), B: geom.V2(5, -2)},
			false,
		},
		{
			"perpendicular line segments touching each other don't cross",
			//     |
			// ----|
			//     |
			poly.LineSeg{A: geom.V2(0, 0), B: geom.V2(4, 0)},
			poly.LineSeg{A: geom.V2(4, 2), B: geom.V2(4, -2)},
			false,
		},
		{
			"line segments with common end point don't cross",
			poly.LineSeg{A: geom.V2(3, 2), B: geom.V2(5, 3)},
			poly.LineSeg{A: geom.V2(3, 2), B: geom.V2(5, 7)},
			false,
		},
		{
			"X-shaped line segments do cross",
			poly.LineSeg{A: geom.V2(-2, -1), B: geom.V2(2, 1)},
			poly.LineSeg{A: geom.V2(-2, 1), B: geom.V2(2, -1)},
			true,
		},
		{
			"perpendicular line segments with small protrusion do cross",
			//     |
			// ----x
			//     |
			poly.LineSeg{A: geom.V2(0, 0), B: geom.V2(4, 0)},
			poly.LineSeg{A: geom.V2(3.999999, 2), B: geom.V2(3.999999, -2)},
			true,
		},
		{
			"perpendicular line segments crossing each other",
			//    |
			// ---x-
			//    |
			poly.LineSeg{A: geom.V2(0, 0), B: geom.V2(4, 0)},
			poly.LineSeg{A: geom.V2(3, 2), B: geom.V2(3, -2)},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if b := tt.l1.Crosses(tt.l2); b != tt.want {
				t.Errorf("line segment (%v).Crosses(%v) = %v, want %v", tt.l1, tt.l2, b, tt.want)
			}
		})
	}
}

func TestLineSegMiddle(t *testing.T) {
	tests := []struct {
		lineSeg poly.LineSeg
		want    geom.Vec2
	}{
		{poly.LineSeg{A: geom.V2Zero, B: geom.V2Zero}, geom.V2Zero},
		{poly.LineSeg{A: geom.V2(-2, 0), B: geom.V2(2, 0)}, geom.V2Zero},
		{poly.LineSeg{A: geom.V2(1, 3), B: geom.V2(4, 3)}, geom.V2(2.5, 3)},
		{poly.LineSeg{A: geom.V2(0, -3), B: geom.V2(0, 3)}, geom.V2Zero},
		{poly.LineSeg{A: geom.V2(8, 0), B: geom.V2(8, 10)}, geom.V2(8, 5)},
		{poly.LineSeg{A: geom.V2(-2.5, -2.5), B: geom.V2(2.5, 2.5)}, geom.V2Zero},
		{poly.LineSeg{A: geom.V2(1, 1), B: geom.V2(5, 2)}, geom.V2(3, 1.5)},
		{poly.LineSeg{A: geom.V2(0, 1), B: geom.V2(1, 0)}, geom.V2(0.5, 0.5)},
	}
	for _, tt := range tests {
		if s := tt.lineSeg.Middle(); s != tt.want {
			t.Errorf("middle of line segment %v was %v, want %v", tt.lineSeg, s, tt.want)
		}
	}
}

func TestLineSegNearEq(t *testing.T) {
	tests := []struct {
		lineSegA poly.LineSeg
		lineSegB poly.LineSeg
		want     bool
	}{
		{
			poly.LineSeg{A: geom.V2Zero, B: geom.V2Zero},
			poly.LineSeg{A: geom.V2Zero, B: geom.V2Zero},
			true,
		},
		{
			poly.LineSeg{A: geom.V2(1.2345678, 0.9876543), B: geom.V2(42.4626734, -165.452349)},
			poly.LineSeg{A: geom.V2(1.2345678, 0.9876543), B: geom.V2(42.4626734, -165.452349)},
			true,
		},
		{
			poly.LineSeg{A: geom.V2(1.234559, 0.987651), B: geom.V2(42.462669, -165.45235)},
			poly.LineSeg{A: geom.V2(1.23456, 0.98765), B: geom.V2(42.4626734, -165.452349)},
			true,
		},
		{
			poly.LineSeg{A: geom.V2(1.23449, 0.987651), B: geom.V2(42.462669, -165.45235)},
			poly.LineSeg{A: geom.V2(1.23456, 0.98765), B: geom.V2(42.4626734, -165.452349)},
			false,
		},
		{
			poly.LineSeg{A: geom.V2(1.234559, 0.987651), B: geom.V2(42.462669, -165.45235)},
			poly.LineSeg{A: geom.V2(1.23456, 0.98765), B: geom.V2(42.4626734, -165.452339)},
			false,
		},
	}
	for _, tt := range tests {
		if b := tt.lineSegA.NearEq(tt.lineSegB); b != tt.want {
			t.Errorf("line segment (%v).NearEq(%v) was %v, want %v", tt.lineSegA, tt.lineSegB, b, tt.want)
		}
	}
}

func TestLineSegString(t *testing.T) {
	tests := []struct {
		lineSeg poly.LineSeg
		want    string
	}{
		{poly.LineSeg{A: geom.V2Zero, B: geom.V2Zero}, "L(0, 0):(0, 0)"},
		{poly.LineSeg{A: geom.V2(3, 4), B: geom.V2(0, 5)}, "L(3, 4):(0, 5)"},
		{poly.LineSeg{A: geom.V2(1.5, 1), B: geom.V2(2, 3.4)}, "L(1.5, 1):(2, 3.4)"},
		{poly.LineSeg{A: geom.V2(-4.54, 2.0), B: geom.V2(23.5, -2.643)}, "L(-4.54, 2):(23.5, -2.643)"},
		{poly.LineSeg{A: geom.V2(42.5, 12.78), B: geom.V2(0.003, -0.004)}, "L(42.5, 12.78):(0.003, -0.004)"},
	}
	for _, tt := range tests {
		if s := tt.lineSeg.String(); s != tt.want {
			t.Errorf("string representation of line segment was %q, want %q", s, tt.want)
		}
	}
}

func TestLineIntersect(t *testing.T) {
	tests := []struct {
		l1               poly.Line
		l2               poly.Line
		wantIntersection geom.Vec2
		wantExists       bool
	}{
		{
			poly.Line{Seg: poly.LineSeg{A: geom.V2(-1, -1), B: geom.V2(1, 1)}},
			poly.Line{Seg: poly.LineSeg{A: geom.V2(-1, 1), B: geom.V2(1, -1)}},
			geom.V2Zero, true,
		},
		{
			poly.Line{Seg: poly.LineSeg{A: geom.V2(-2, 1), B: geom.V2(2, 1)}},
			poly.Line{Seg: poly.LineSeg{A: geom.V2(-2, 2), B: geom.V2(2, 2)}},
			geom.V2Zero, false,
		},
		{
			poly.Line{Seg: poly.LineSeg{A: geom.V2(0, 1), B: geom.V2(0, 2)}},
			poly.Line{Seg: poly.LineSeg{A: geom.V2(3, 2), B: geom.V2(3, 4)}},
			geom.V2Zero, false,
		},
		{
			poly.Line{Seg: poly.LineSeg{A: geom.V2(0, 0), B: geom.V2(2, 1)}},
			poly.Line{Seg: poly.LineSeg{A: geom.V2(1, 0), B: geom.V2(3, 1)}},
			geom.V2Zero, false,
		},
		{
			poly.Line{Seg: poly.LineSeg{A: geom.V2(2, 3), B: geom.V2(5, 6)}},
			poly.Line{Seg: poly.LineSeg{A: geom.V2(2, 3), B: geom.V2(5, 6)}},
			geom.V2Zero, false,
		},
	}
	for _, tt := range tests {
		p, exists := tt.l1.Intersect(tt.l2)
		if exists != tt.wantExists {
			t.Errorf("existence of intersection of lines %v and %v was %v, want %v", tt.l1, tt.l2, exists, tt.wantExists)
		}
		if p != tt.wantIntersection {
			t.Errorf("intersection of lines %v and %v was %v, want %v", tt.l1, tt.l2, p, tt.wantIntersection)
		}
	}
}

func TestLineSide(t *testing.T) {
	tests := []struct {
		name string
		l    poly.Line
		p    geom.Vec2
		want int
	}{
		{
			name: "point on left side",
			l: poly.Line{Seg: poly.LineSeg{
				A: geom.V2(0, 0),
				B: geom.V2(1, 1),
			}},
			p:    geom.V2(0, 2),
			want: -1,
		},
		{
			name: "point on right side",
			l: poly.Line{Seg: poly.LineSeg{
				A: geom.V2(0, 0),
				B: geom.V2(1, 1),
			}},
			p:    geom.V2(2, 0),
			want: +1,
		},
		{
			name: "point on line",
			l: poly.Line{Seg: poly.LineSeg{
				A: geom.V2(0, 0),
				B: geom.V2(1, 1),
			}},
			p:    geom.V2(2, 2),
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.l.Side(tt.p)
			if got != tt.want {
				t.Errorf("%v.Side(%v) = %d, want %d", tt.l, tt.p, got, tt.want)
			}
		})
	}
}
