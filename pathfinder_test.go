// Copyright 2023 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pathfind_test

import (
	"image"
	"reflect"
	"testing"

	"github.com/fzipp/pathfind"
)

// A U-shaped polygon. Origin is at the top-left corner.
//
//	 0,0 >---+   +---+ 30,0
//	     |   |   |   |
//	     |   +---+   |
//	     |           |
//	0,20 +-----------+ 30,20
var polygonU = [][]image.Point{
	{
		image.Pt(0, 0),
		image.Pt(10, 0),
		image.Pt(10, 10),
		image.Pt(20, 10),
		image.Pt(20, 0),
		image.Pt(30, 0),
		image.Pt(30, 20),
		image.Pt(0, 20),
	},
}

// A square with a diamond shaped hole inside. Origin is at the top-left corner.
//
//	 0,0 >-----------+ 40,0
//	     |     >     |
//	     |    / \    |
//	     |   +   +   |
//	     |    \ /    |
//	     |     +     |
//	0,40 +-----------+ 40,40
var polygonO = [][]image.Point{
	{
		// Outer rectangle
		image.Pt(0, 0),
		image.Pt(40, 0),
		image.Pt(40, 40),
		image.Pt(0, 40),
	},
	{
		// Inner diamond
		image.Pt(20, 10),
		image.Pt(30, 20),
		image.Pt(20, 30),
		image.Pt(10, 20),
	},
}

func TestPathfinderPath(t *testing.T) {
	tests := []struct {
		name     string
		polygons [][]image.Point
		start    image.Point
		dest     image.Point
		want     []image.Point
	}{
		{
			// +---+   +---+
			// | s |   |   |
			// |   +---+   |
			// | d         |
			// +-----------+
			name:     "Direct connection",
			polygons: polygonU,
			start:    image.Pt(5, 5),
			dest:     image.Pt(5, 15),
			want: []image.Point{
				image.Pt(5, 5),
				image.Pt(5, 15),
			},
		},
		{
			// +---+   +---+
			// | s |   |   |
			// |   +---+   |
			// |         d |
			// +-----------+
			name:     "One corner",
			polygons: polygonU,
			start:    image.Pt(5, 5),
			dest:     image.Pt(25, 15),
			want: []image.Point{
				image.Pt(5, 5),
				image.Pt(10, 10),
				image.Pt(25, 15),
			},
		},
		{
			// >---+   +---+
			// | s |   | d |
			// |   +---+   |
			// |           |
			// +-----------+
			name:     "Two corners",
			polygons: polygonU,
			start:    image.Pt(5, 5),
			dest:     image.Pt(25, 5),
			want: []image.Point{
				image.Pt(5, 5),
				image.Pt(10, 10),
				image.Pt(20, 10),
				image.Pt(25, 5),
			},
		},
		{
			// +---+   +---+
			// | s | d |   |
			// |   +---+   |
			// |           |
			// +-----------+
			name:     "No path through wall: dest clamped to polygons",
			polygons: polygonU,
			start:    image.Pt(5, 5),
			dest:     image.Pt(15, 5),
			want: []image.Point{
				image.Pt(5, 5),
				image.Pt(10, 5),
			},
		},
		{
			// +---+ s +---+
			// |   | d |   |
			// |   +---+   |
			// |           |
			// +-----------+
			name:     "No path outside polygon",
			polygons: polygonU,
			start:    image.Pt(15, 0),
			dest:     image.Pt(15, 5),
			want:     nil,
		},
		{
			// >-----------+
			// | s   >     |
			// |    / \    |
			// |   +   +   |
			// |    \ /    |
			// |     + d   |
			// +-----------+
			name:     "Path around inner polygon",
			polygons: polygonO,
			start:    image.Pt(15, 10),
			dest:     image.Pt(30, 30),
			want: []image.Point{
				image.Pt(15, 10),
				image.Pt(20, 10),
				image.Pt(30, 20),
				image.Pt(30, 30),
			},
		},
		{
			// >
			// | \
			// | s \ d
			// |     +-----+
			// |           |
			// +-----+     |
			//         \   |
			//           \ |
			//             +
			name: "No path out of thunderbolt shape: dest clamped to polygons",
			polygons: [][]image.Point{
				{
					image.Pt(0, 0),
					image.Pt(100, 100),
					image.Pt(200, 100),
					image.Pt(200, 300),
					image.Pt(100, 200),
					image.Pt(0, 200),
				},
			},
			start: image.Pt(30, 70),
			dest:  image.Pt(100, 70),
			want: []image.Point{
				image.Pt(30, 70),
				image.Pt(85, 85),
			},
		},
		{
			name: "ensure clamped dest inside 1",
			polygons: [][]image.Point{
				{
					image.Pt(70, 55),
					image.Pt(250, 54),
					image.Pt(300, 100),
				},
			},
			start: image.Pt(180, 60),
			dest:  image.Pt(181, 54),
			want: []image.Point{
				image.Pt(180, 60),
				image.Pt(180, 55),
			},
		},
		{
			name: "ensure clamped dest inside 2",
			polygons: [][]image.Point{
				{
					image.Pt(73, 55),
					image.Pt(100, 100),
					image.Pt(76, 168),
				},
			},
			start: image.Pt(90, 100),
			dest:  image.Pt(74, 98),
			want: []image.Point{
				image.Pt(90, 100),
				image.Pt(75, 97),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathfinder := pathfind.NewPathfinder(tt.polygons)
			got := pathfinder.Path(tt.start, tt.dest)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`%s
polygons: %v
Path(%v, %v)
 got: %v
want: %v`,
					tt.name, tt.polygons, tt.start, tt.dest, got, tt.want)
			}
		})
	}
}

func TestPathfinderVisibilityGraph(t *testing.T) {
	tests := []struct {
		name     string
		polygons [][]image.Point
		start    image.Point
		dest     image.Point
		want     map[image.Point][]image.Point
	}{
		{
			// >---+   +---+
			// | s |   | d |
			// |   +---+   |
			// |           |
			// +-----------+
			name:     "Two corners",
			polygons: polygonU,
			start:    image.Pt(5, 5),
			dest:     image.Pt(25, 5),
			want: map[image.Point][]image.Point{
				image.Pt(5, 5):   {image.Pt(10, 10)},
				image.Pt(10, 10): {image.Pt(20, 10), image.Pt(5, 5)},
				image.Pt(20, 10): {image.Pt(10, 10), image.Pt(25, 5)},
				image.Pt(25, 5):  {image.Pt(20, 10)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathfinder := pathfind.NewPathfinder(tt.polygons)
			pathfinder.Path(tt.start, tt.dest)
			got := pathfinder.VisibilityGraph()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`%s
polygons: %v
Path(%v, %v)
VisibilityGraph()
 got: %v
want: %v`,
					tt.name, tt.polygons, tt.start, tt.dest, got, tt.want)
			}
		})
	}
}
