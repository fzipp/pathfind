// Copyright 2023 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pathfind

// graph is represented by an adjacency list.
type graph[Node comparable] map[Node][]Node

// link creates a directed edge from node a to node b.
func (g graph[Node]) link(a, b Node) graph[Node] {
	g[a] = append(g[a], b)
	return g
}

// Neighbours returns the neighbour nodes of node n in the graph.
// This method makes graph[Node] implement the astar.Graph[Node] interface.
func (g graph[Node]) Neighbours(n Node) []Node {
	return g[n]
}
