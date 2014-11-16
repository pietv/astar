package main

import (
	"container/heap"
	"math/rand"

	"github.com/pietv/unionfind"
)

type Edge struct {
	v1, v2 Location
	w      int
}

type PriorityQueue []Edge

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].w < pq[j].w }
func (pq *PriorityQueue) Push(x interface{}) {
	edge := x.(Edge)
	*pq = append(*pq, edge)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	edge := old[n-1]
	*pq = old[0 : n-1]
	return edge

}

func NewRandomKruskal(rows, cols int) *Maze {
	pq := PriorityQueue{}
	heap.Init(&pq)

	// Kruskal's minimum spanning tree (MST) algorithm
	// uses the disjoint set (union-find) data structure.
	u := unionfind.New()

	// Initialize a complete (fully connected) graph with random edge weights.
	// Vertices are cell locations surrounded by walls.
	//
	// (1,1)--(1,2)--(1,3)--...  (M = cols)
	//   |      |      |
	// (2,1)--(2,2)--(2,3)--...
	//   |      |      |
	// (3,1)--(3,2)--(3,3)--...
	//   |      |      |
	//   .      .      .
	//
	// (N = rows)
	//
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			// Horizontal edge.
			if j < cols-1 {
				heap.Push(&pq, Edge{Location{i, j}, Location{i, j + 1}, rand.Int()})
			}

			// Vertical edge.
			if i < rows-1 {
				heap.Push(&pq, Edge{Location{i, j}, Location{i + 1, j}, rand.Int()})
			}

			// Make a set for each vertex.
			u.MakeSet(u, Location{i, j})
		}
	}

	// Maze template.
	//
	//  *******
	//  * * * *
	//  *******
	//  * * * *
	//  *******
	//
	maze := make([][]string, rows*2+1)
	for i := 0; i <= rows*2; i++ {
		maze[i] = make([]string, cols*2+1)
		for j := 0; j <= cols*2; j++ {
			// Gaps (spaces) in odd-numbered columns and rows (zero-based).
			if i%2 == 1 && j%2 == 1 {
				maze[i][j] = SpaceRune
			} else {
				maze[i][j] = WallRune
			}
		}
	}

	for {
		if pq.Len() == 0 {
			break
		}

		// Pick an edge belonging to a random MST.
		// Break a wall at that edge location.
		e := heap.Pop(&pq).(Edge)
		if !u.Connected(e.v1, e.v2) {
			u.Union(e.v1, e.v2)

			if e.v1.i == e.v2.i {
				// Break vertical wall.
				maze[e.v1.i*2+1][e.v1.j*2+2] = SpaceRune
			} else {
				// Break horizontal wall.
				maze[e.v1.i*2+2][e.v1.j*2+1] = SpaceRune
			}
		}
	}

	// Fixed Start and Finish locations.
	start := Location{rows*2 - 1, 1}
	maze[start.i][start.j] = StartRune

	finish := Location{1, cols*2 - 1}
	maze[finish.i][finish.j] = FinishRune

	return &Maze{
		maze:   maze,
		start:  start,
		finish: finish,
		curr:   start,
	}
}
