// Copyright 2023 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package poly provides types and functions for working with polygons
// in support of the pathfind package.
package poly

import "github.com/fzipp/geom"

// A Polygon is a polygon in 2-dimensional space, represented as a slice
// of its vertices.
type Polygon []geom.Vec2

// ParsePolygon parses a new polygon from a comma-separated coordinate
// string, for example "186.5,364.7,303.25,374,303.1,412". Should have an
// even number of coordinate values. Rest is ignored.
func ParsePolygon(coords string) Polygon {
	floats := ParseFloats(coords)
	n := len(floats)
	p := make(Polygon, 0, n/2)
	for i := 0; i < n-1; i += 2 {
		v := geom.V2(floats[i], floats[i+1])
		p = append(p, v)
	}
	return p
}

// Edge returns the edge with index i of polygon p.
func (p Polygon) Edge(i int) LineSeg {
	j := (i + 1) % len(p)
	return LineSeg{p[i], p[j]}
}

// Contains checks if point pt lies inside the boundary of polygon p.
func (p Polygon) Contains(pt geom.Vec2, toleranceOnOutside bool) bool {
	// Ray casting algorithm: if a ray from point pt in any direction
	// (in our case horizontally to the east) crosses an odd number
	// of polygon edges, then pt lies inside the polygon, otherwise
	// outside.
	in := false
	for i := range p {
		edge := p.Edge(i)
		if edge.ClosestPt(pt).NearEq(pt) {
			return toleranceOnOutside
		}
		if hRayIntersects(pt, edge) {
			in = !in
		}
	}
	return in
}

// IsCrossedBy checks if any side of polygon p is crossed by line segment ls.
func (p Polygon) IsCrossedBy(ls LineSeg) bool {
	for i, v := range p {
		if ls.A == v || ls.B == v {
			continue
		}
		if ls.Crosses(p.Edge(i)) {
			return true
		}
		if ls.ClosestPt(v) == v {
			prev := p[p.WrapIndex(i-1)]
			next := p[p.WrapIndex(i+1)]
			l := Line{ls}
			if l.Side(prev) != l.Side(next) {
				return true
			}
		}
	}
	return false
}

// hRayIntersects checks if a horizontal ray from point p to the right
// intersects a line segment.
func hRayIntersects(p geom.Vec2, ls LineSeg) bool {
	if !hLineIntersects(p, ls) {
		return false
	}
	hRay := Line{LineSeg{p, geom.V2(p.X+1, p.Y)}}
	q, _ := hRay.Intersect(Line{ls})

	// Checks whether p is on the left-hand side of the line segment
	// by comparing p.X to the x coordinate of the intersection point q.
	return p.X <= q.X
}

// hLineIntersects checks if a horizontal line through point p intersects a
// line segment.
func hLineIntersects(p geom.Vec2, ls LineSeg) bool {
	// True, if each end point of the line segment lies on a
	// different side of the horizontal line.
	return (ls.A.Y >= p.Y) != (ls.B.Y >= p.Y)
}

// match is a helper structure for closest point algorithms. Used to hold the
// current best match and its distance.
type match struct {
	pt   geom.Vec2
	dist float32
}

// ClosestPt returns the closest point to point pt on the outline of
// polygon p.
func (p Polygon) ClosestPt(pt geom.Vec2) geom.Vec2 {
	var best match
	best.pt = p[0]
	best.dist = best.pt.SqDist(pt)
	for i := range p {
		var current match
		current.pt = p.Edge(i).ClosestPt(pt)
		current.dist = current.pt.SqDist(pt)
		if current.dist < best.dist {
			best = current
		}
	}
	return best.pt
}

// IsConcaveAt checks, whether the vertex with index i of polygon p is
// concave or not.
func (p Polygon) IsConcaveAt(i int) bool {
	v := p[i]
	prev := p[p.WrapIndex(i-1)]
	next := p[p.WrapIndex(i+1)]
	left := v.Sub(prev)
	right := next.Sub(v)
	return left.CrossLen(right) < 0
}

// WrapIndex returns an index based on i that can be safely used to access an
// element of p. It wraps around if i < 0 or i >= len(p).
func (p Polygon) WrapIndex(i int) int {
	n := len(p)
	return ((i % n) + n) % n
}
