// Copyright 2023 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Pathfinddemo demonstrates an algorithm for finding the shortest path
// between two points constrained by a set of polygons.
// The demo starts an HTTP server that provides a graphical interface
// accessible via a web browser.
// It will try to open the demo in the system's default browser if possible.
//
// Usage:
//
//	pathfinddemo [-http address]
//
// Flags:
//
//	-http  HTTP service address (e.g., '127.0.0.1:8080' or just ':8080').
//	       The default is ':8080'.
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os/exec"
	"runtime"
	"time"

	"github.com/fzipp/canvas"
	"github.com/fzipp/pathfind"
)

var (
	backgroundColor      = color.RGBA{R: 0x26, G: 0x46, B: 0x53, A: 0xFF}
	polygonColor         = color.RGBA{R: 0xf4, G: 0xa2, B: 0x61, A: 0xFF}
	concaveVertexColor   = color.RGBA{R: 0xe9, G: 0xc4, B: 0x6a, A: 0xFF}
	visibilityGraphColor = color.White
	pathColor            = color.RGBA{R: 0x2a, G: 0x9d, B: 0x8f, A: 0xFF}
	startPointColor      = color.RGBA{R: 0x90, G: 0xee, B: 0x90, A: 0xff}
	destPointColor       = color.RGBA{R: 0xe7, G: 0x6f, B: 0x51, A: 0xFF}
)

// Created with Polygon Constructor: https://alaricus.github.io/PolygonConstructor/
const floorPlan = `{"canvas":{"w":800,"h":600},"polygons":[[{"x":73.4375,"y":55},{"x":288.4375,"y":54},{"x":287.4375,"y":172},{"x":231.4375,"y":171},{"x":233.4375,"y":203},{"x":665.4375,"y":203},{"x":665.4375,"y":299},{"x":630.4375,"y":299},{"x":630.4375,"y":331},{"x":711.4375,"y":331},{"x":710.4375,"y":166},{"x":508.4375,"y":167},{"x":508.4375,"y":113},{"x":471.4375,"y":112},{"x":473.4375,"y":149},{"x":381.4375,"y":149},{"x":379.4375,"y":39},{"x":470.4375,"y":37},{"x":470.4375,"y":79},{"x":510.4375,"y":81},{"x":509.4375,"y":37},{"x":600.4375,"y":37},{"x":600.4375,"y":129},{"x":613.4375,"y":128},{"x":613.4375,"y":39},{"x":783.4375,"y":38},{"x":782.4375,"y":165},{"x":736.4375,"y":166},{"x":735.4375,"y":361},{"x":602.4375,"y":361},{"x":602.4375,"y":296},{"x":551.4375,"y":296},{"x":551.4375,"y":242},{"x":530.4375,"y":241},{"x":531.4375,"y":291},{"x":411.4375,"y":291},{"x":410.4375,"y":241},{"x":388.4375,"y":240},{"x":393.4375,"y":290},{"x":205.4375,"y":288},{"x":205.4375,"y":238},{"x":179.4375,"y":238},{"x":164.4375,"y":265},{"x":157.4375,"y":316},{"x":195.4375,"y":337},{"x":257.4375,"y":335},{"x":362.4375,"y":408},{"x":404.4375,"y":409},{"x":465.4375,"y":336},{"x":522.4375,"y":336},{"x":596.4375,"y":407},{"x":595.4375,"y":476},{"x":564.4375,"y":511},{"x":595.4375,"y":526},{"x":644.4375,"y":526},{"x":645.4375,"y":548},{"x":666.4375,"y":548},{"x":663.4375,"y":525},{"x":691.4375,"y":525},{"x":690.4375,"y":487},{"x":642.4375,"y":486},{"x":639.4375,"y":396},{"x":764.4375,"y":395},{"x":762.4375,"y":485},{"x":719.4375,"y":485},{"x":716.4375,"y":521},{"x":769.4375,"y":520},{"x":770.4375,"y":582},{"x":601.4375,"y":583},{"x":598.4375,"y":551},{"x":573.4375,"y":551},{"x":576.4375,"y":582},{"x":522.4375,"y":581},{"x":519.4375,"y":510},{"x":472.4375,"y":512},{"x":408.4375,"y":439},{"x":362.4375,"y":438},{"x":277.4375,"y":507},{"x":198.4375,"y":506},{"x":140.4375,"y":447},{"x":139.4375,"y":392},{"x":172.4375,"y":363},{"x":139.4375,"y":340},{"x":118.4375,"y":311},{"x":129.4375,"y":261},{"x":123.4375,"y":235},{"x":78.4375,"y":233},{"x":77.4375,"y":328},{"x":23.4375,"y":328},{"x":26.4375,"y":199},{"x":197.4375,"y":202},{"x":198.4375,"y":172},{"x":76.4375,"y":168}],[{"x":512.4375,"y":386},{"x":539.4375,"y":430},{"x":492.4375,"y":464},{"x":462.4375,"y":419},{"x":496.4375,"y":385}],[{"x":245.4375,"y":383},{"x":298.4375,"y":414},{"x":264.4375,"y":461},{"x":220.4375,"y":462},{"x":192.4375,"y":433},{"x":188.4375,"y":398}],[{"x":683.4375,"y":80},{"x":686.4375,"y":141},{"x":652.4375,"y":143},{"x":652.4375,"y":82}]]}`

func main() {
	http := flag.String("http", ":8080", "HTTP service address (e.g., '127.0.0.1:8080' or just ':8080')")
	flag.Parse()

	url := httpLink(*http)
	startBrowser(url)
	fmt.Println("Serving demo at " + url)
	err := canvas.ListenAndServe(*http, run,
		canvas.Size(800, 600),
		canvas.ScaleFullPage(false, true),
		canvas.Title("Path finding demo"),
		canvas.EnableEvents(
			canvas.MouseDownEvent{},
			canvas.KeyDownEvent{},
			canvas.MouseMoveEvent{},
		),
		canvas.Reconnect(time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx *canvas.Context) {
	polygons, size, err := polygonsFromJSON([]byte(floorPlan))
	if err != nil {
		log.Fatal(err)
	}
	d := &demo{
		polygons:            polygons,
		size:                size,
		start:               image.Pt(245, 78),
		dest:                image.Pt(420, 125),
		pathfinder:          pathfind.NewPathfinder(polygons),
		showVisibilityGraph: false,
	}
	ctx.SetFont("14px sans-serif")

	for !d.quit {
		d.update()
		d.draw(ctx)
		ctx.Flush()
		event := <-ctx.Events()
		d.handle(event)
	}
}

type demo struct {
	quit                bool
	polygons            [][]image.Point
	size                image.Point
	start               image.Point
	dest                image.Point
	showVisibilityGraph bool
	pathfinder          *pathfind.Pathfinder
	path                []image.Point
}

func (d *demo) handle(ev canvas.Event) {
	switch e := ev.(type) {
	case canvas.CloseEvent:
		d.quit = true
	case canvas.MouseDownEvent:
		if len(d.path) == 0 {
			d.start = image.Pt(e.X, e.Y)
			break
		}
		d.start = d.path[len(d.path)-1]
	case canvas.MouseMoveEvent:
		d.dest = image.Pt(e.X, e.Y)
	case canvas.KeyDownEvent:
		if e.Key == " " {
			d.showVisibilityGraph = !d.showVisibilityGraph
		}
	}
}

func (d *demo) update() {
	d.path = d.pathfinder.Path(d.start, d.dest)
}

func (d *demo) draw(ctx *canvas.Context) {
	ctx.SetFillStyle(backgroundColor)
	ctx.FillRect(0, 0, float64(d.size.X), float64(d.size.Y))

	ctx.SetStrokeStyle(polygonColor)
	ctx.SetLineWidth(3)
	drawPolygons(ctx, d.polygons)

	if d.showVisibilityGraph {
		vis := d.pathfinder.VisibilityGraph()
		ctx.SetStrokeStyle(visibilityGraphColor)
		ctx.SetLineWidth(0.3)
		drawGraph(ctx, vis)

		ctx.SetStrokeStyle(concaveVertexColor)
		ctx.SetLineWidth(3)
		for v := range vis {
			drawPoint(ctx, v)
		}
	}

	ctx.SetStrokeStyle(pathColor)
	ctx.SetLineWidth(3)
	drawPath(ctx, d.path)

	ctx.SetStrokeStyle(startPointColor)
	ctx.SetFillStyle(startPointColor)
	drawPoint(ctx, d.start)
	if len(d.path) > 0 {
		ctx.SetStrokeStyle(destPointColor)
		ctx.SetFillStyle(destPointColor)
		drawPoint(ctx, d.path[len(d.path)-1])
	}

	ctx.SetFillStyle(color.White)
	ctx.FillText("Press [Space] to toggle the visibility graph. Click to change the position of the start node.", 10, 20)
}

func httpLink(addr string) string {
	if addr[0] == ':' {
		addr = "localhost" + addr
	}
	return "http://" + addr
}

// startBrowser tries to open the URL in a browser
// and reports whether it succeeds.
func startBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
