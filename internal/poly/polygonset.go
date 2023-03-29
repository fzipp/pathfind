// Copyright 2023 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package poly

import "github.com/fzipp/geom"

// A PolygonSet represents multiple polygons.
type PolygonSet []Polygon

// ParsePolygons parses polygons from n comma-separated coordinate
// strings and returns them as a PolygonSet. See ParsePolygon for
// details on the format.
func ParsePolygons(coords []string) PolygonSet {
	ps := make([]Polygon, len(coords))
	for i, cs := range coords {
		ps[i] = ParsePolygon(cs)
	}
	return ps
}

// Contains checks if point pt lies inside the boundaries of a polygon set.
// Overlapping polygons can form holes and islands.
func (ps PolygonSet) Contains(pt geom.Vec2) bool {
	in := false
	for _, p := range ps {
		if p.Contains(pt, !in) {
			in = !in
		}
	}
	return in
}

// ClosestPt returns the closest point to point pt on any of the outlines of
// polygon set ps.
func (ps PolygonSet) ClosestPt(pt geom.Vec2) geom.Vec2 {
	var best match
	best.pt = ps[0].ClosestPt(pt)
	best.dist = best.pt.SqDist(pt)
	for _, p := range ps {
		var current match
		current.pt = p.ClosestPt(pt)
		current.dist = current.pt.SqDist(pt)
		if current.dist < best.dist {
			best = current
		}
	}
	return best.pt
}
