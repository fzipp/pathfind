# pathfind

[![PkgGoDev](https://pkg.go.dev/badge/github.com/fzipp/pathfind)](https://pkg.go.dev/github.com/fzipp/pathfind)
![Build Status](https://github.com/fzipp/pathfind/workflows/build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fzipp/pathfind)](https://goreportcard.com/report/github.com/fzipp/pathfind)

Package pathfind finds the shortest path between two points in a polygon set.

The algorithm works as follows:
- determine all concave polygon vertices
- add start and end points
- build a visibility graph
- use the A* search algorithm (package [astar](https://github.com/fzipp/astar))
  on the visibility graph to find the shortest path

## Demo

```
go run github.com/fzipp/pathfind/cmd/pathfinddemo@latest
```

https://user-images.githubusercontent.com/598327/228644017-a9800747-a096-46ae-befe-d725da18205a.mp4

## Example Code

```go
package main

import (
	"fmt"
	"image"
	
	"github.com/fzipp/pathfind"
)

func main() {
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
}
```

Output:

```
[(5,5) (10,10) (30,15) (40,15) (45,10)]
```

## License

This project is free and open source software licensed under the
[BSD 3-Clause License](LICENSE).
