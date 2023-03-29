// Copyright 2023 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pathfind

import (
	"image"
	"math"

	"github.com/fzipp/geom"
)

// ps2vs converts a []image.Point to a []geom.Vec2.
func ps2vs(ps []image.Point) []geom.Vec2 {
	return convert(ps, p2v)
}

// p2v converts an image.Point to a geom.Vec2.
func p2v(p image.Point) geom.Vec2 {
	return geom.Vec2{X: float32(p.X), Y: float32(p.Y)}
}

// v2p converts a geom.Vec2 to an image.Point. X and Y coordinates are rounded.
func v2p(v geom.Vec2) image.Point {
	return image.Point{
		X: int(math.Round(float64(v.X))),
		Y: int(math.Round(float64(v.Y))),
	}
}

// convert maps a slice s to a new slice of elements with target type To by
// applying the conversion function f to each element.
func convert[From, To any](s []From, f func(From) To) []To {
	res := make([]To, 0, len(s))
	for _, e := range s {
		res = append(res, f(e))
	}
	return res
}
