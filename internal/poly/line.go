// Copyright 2023 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package poly

import "github.com/fzipp/geom"

// A LineSeg represents a line segment between two points A and B.
type LineSeg struct {
	A, B geom.Vec2
}

// Len returns the length of a line segment.
func (l LineSeg) Len() float32 {
	return l.A.Dist(l.B)
}

// ClosestPt returns the point on the line segment l that is closest to point p.
// This is either the orthogonal projection of p onto l or one of l's end
// points if the projection is not within the line segment.
func (l LineSeg) ClosestPt(p geom.Vec2) geom.Vec2 {
	v := l.B.Sub(l.A)
	w := p.Sub(l.A)
	c1 := w.Dot(v)
	if c1 <= 0 {
		return l.A
	}
	c2 := v.Dot(v)
	if c2 <= c1 {
		return l.B
	}
	return l.A.Add(v.Mul(c1 / c2))
}

// Crosses returns true if line segments l and m cross each other,
// otherwise false.
func (l LineSeg) Crosses(m LineSeg) bool {
	u := l.A.Sub(l.B)
	v := m.A.Sub(m.B)
	D := u.CrossLen(v)
	if D == 0 {
		// The line segments are parallel.
		return false
	}
	w := l.B.Sub(m.B)
	n1 := u.CrossLen(w)
	n2 := v.CrossLen(w)
	if n1 == 0 || n2 == 0 {
		return false
	}
	r := n1 / D
	s := n2 / D
	return (0 < r && r < 1) && (0 < s && s < 1)
}

// Middle returns the middle of the line segment.
func (l LineSeg) Middle() geom.Vec2 {
	return l.A.Add(l.B).Div(2)
}

// NearEq determines whether two line segments are equal or not.
// The order of the end points is relevant.
func (l LineSeg) NearEq(m LineSeg) bool {
	return l.A.NearEq(m.A) && l.B.NearEq(m.B)
}

// String returns a string representation of l like "L(1.5, 1):(2, 3.4)".
func (l LineSeg) String() string {
	return "L" + l.A.String() + ":" + l.B.String()
}

// Line represents a straight line that goes through the two points of a
// line segment and continues to infinity in both directions.
type Line struct {
	Seg LineSeg
}

// Intersect returns the intersection point p of two lines l and m.
// Returns false if the lines are parallel and therefore no such
// intersection point exists.
func (l Line) Intersect(m Line) (p geom.Vec2, exists bool) {
	u := l.Seg.A.Sub(l.Seg.B)
	v := m.Seg.A.Sub(m.Seg.B)
	D := u.CrossLen(v)
	if D == 0 {
		// The lines are parallel.
		return geom.Vec2{}, false
	}
	r := l.Seg.A.CrossLen(l.Seg.B) / D
	s := m.Seg.A.CrossLen(m.Seg.B) / D
	return v.Mul(r).Sub(u.Mul(s)), true
}

// Side reports on which side of the line point p is.
// It is +1 on one side, -1 on the other side, and 0 on the line.
func (l Line) Side(p geom.Vec2) int {
	ap := p.Sub(l.Seg.A)
	ab := l.Seg.B.Sub(l.Seg.A)
	return sgn(ap.CrossLen(ab))
}

// sgn returns -1 if x is negative, +1 if x is positive, and 0 otherwise.
func sgn(x float32) int {
	switch {
	case x < 0:
		return -1
	case x > 0:
		return +1
	default:
		return 0
	}
}
