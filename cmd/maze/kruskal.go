package main

import (
	"container/heap"
	"math/rand"

	"github.com/pietv/unionfind"
)

type edge struct {
	v1, v2 location
	w      int
}

type priorityQueue []edge

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].w < pq[j].w }
func (pq *priorityQueue) Push(x interface{}) {
	edge := x.(edge)
	*pq = append(*pq, edge)
}
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	edge := old[n-1]
	*pq = old[0 : n-1]
	return edge

}

// newRandomKruskal returns a random rectangular maze of rows by cols size.
func newRandomKruskal(rows, cols int) *maze {
	pq := priorityQueue{}
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
				heap.Push(&pq, edge{location{i, j}, location{i, j + 1}, rand.Int()})
			}

			// Vertical edge.
			if i < rows-1 {
				heap.Push(&pq, edge{location{i, j}, location{i + 1, j}, rand.Int()})
			}

			// Make a set for each vertex.
			u.MakeSet(u, location{i, j})
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
	m := make([][]string, rows*2+1)
	for i := 0; i <= rows*2; i++ {
		m[i] = make([]string, cols*2+1)
		for j := 0; j <= cols*2; j++ {
			// Gaps (spaces) in odd-numbered columns and rows (zero-based).
			if i%2 == 1 && j%2 == 1 {
				m[i][j] = spaceRune
			} else {
				m[i][j] = wallRune
			}
		}
	}

	for {
		if pq.Len() == 0 {
			break
		}

		// Pick an edge belonging to a random MST.
		// Break a wall at that edge location.
		e := heap.Pop(&pq).(edge)
		if !u.Connected(e.v1, e.v2) {
			u.Union(e.v1, e.v2)

			if e.v1.i == e.v2.i {
				// Break vertical wall.
				m[e.v1.i*2+1][e.v1.j*2+2] = spaceRune
			} else {
				// Break horizontal wall.
				m[e.v1.i*2+2][e.v1.j*2+1] = spaceRune
			}
		}
	}

	// Fixed Start and Finish locations.
	start := location{rows*2 - 1, 1}
	m[start.i][start.j] = startRune

	finish := location{1, cols*2 - 1}
	m[finish.i][finish.j] = finishRune

	return &maze{
		maze:   m,
		start:  start,
		finish: finish,
		curr:   start,
	}
}
