// Copyright 2023 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image"
	"math"

	"github.com/fzipp/canvas"
)

func drawPolygons(ctx *canvas.Context, polygons [][]image.Point) {
	for _, p := range polygons {
		drawPolygon(ctx, p)
	}
}

func drawPolygon(ctx *canvas.Context, polygon []image.Point) {
	n := len(polygon)
	if n < 2 {
		return
	}
	drawPath(ctx, polygon)
	drawLine(ctx, polygon[n-1], polygon[0])
}

func drawLine(ctx *canvas.Context, a, b image.Point) {
	ctx.BeginPath()
	ctx.MoveTo(float64(a.X), float64(a.Y))
	ctx.LineTo(float64(b.X), float64(b.Y))
	ctx.ClosePath()
	ctx.Stroke()
}

func drawPoint(ctx *canvas.Context, pt image.Point) {
	drawCircle(ctx, pt, 5.0)
	ctx.Fill()
}

func drawGraph(ctx *canvas.Context, g map[image.Point][]image.Point) {
	for start, ends := range g {
		for _, end := range ends {
			drawLine(ctx, start, end)
		}
	}
}

func drawPath(ctx *canvas.Context, path []image.Point) {
	for i := 0; i < len(path)-1; i++ {
		drawLine(ctx, path[i], path[i+1])
	}
}

func drawCircle(ctx *canvas.Context, pos image.Point, radius float64) {
	ctx.BeginPath()
	ctx.Arc(float64(pos.X), float64(pos.Y), radius, 0, 2*math.Pi, false)
	ctx.ClosePath()
	ctx.Stroke()
}
