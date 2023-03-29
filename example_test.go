// Copyright 2023 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pathfind_test

import (
	"fmt"
	"image"

	"github.com/fzipp/pathfind"
)

func ExamplePathfinder_Path() {
	//  (0,0) >---+   +-----------+ (50,0)
	//        | s |   |   >---+   |
	//        |   +---+   |   | d |
	//        |           +---+   |
	// (0,20) +-------------------+ (50,20)
	//
	// s = start, d = destination
	polygons := [][]image.Point{
		// Outer shape
		{
			image.Pt(0, 0),
			image.Pt(10, 0),
			image.Pt(10, 10),
			image.Pt(20, 10),
			image.Pt(20, 0),
			image.Pt(50, 0),
			image.Pt(50, 20),
			image.Pt(0, 20),
		},
		// Inner rectangle ("hole")
		{
			image.Pt(30, 5),
			image.Pt(40, 5),
			image.Pt(40, 15),
			image.Pt(30, 15),
		},
	}
	start := image.Pt(5, 5)
	destination := image.Pt(45, 10)

	pathfinder := pathfind.NewPathfinder(polygons)
	path := pathfinder.Path(start, destination)
	fmt.Println(path)
	// Output:
	// [(5,5) (10,10) (30,15) (40,15) (45,10)]
}
