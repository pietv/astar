// Package astar implements the A* (“A Star”) search algorithm as described in
// the paper by Peter E. Hart et al, “A Formal Basis for the Heuristic Determination
// of Minimum Cost Paths” http://ai.stanford.edu/~nilsson/OnlinePubs-Nils/PublishedPapers/astar.pdf
//
// The “open” and “closed” sets in this implementation are named “priority queue”
// and “explored” respectively.
//
// Time complexity of the algorithm depends on the quality of the heuristic function Estimate().
//
// If Estimate() is constant, the complexity is the same as for the uniform cost search (UCS)
// algorithm – O(b^m), where b is the branching factor (how many successor states on average)
// and m is the maximum depth of the decision tree.
//
// If Estimate() is optimal, the complexity is O(n).
//
// The algorithm is implemented as a Search() function which takes astar.Interface as a parameter.
//
//
// Basic usage (counting to 10):
//
//   type Number int
//
//   func (n Number) Start() interface{}             { return Number(1) }
//   func (n Number) Finish() bool                   { return n == Number(10) }
//   func (n *Number) Move(x interface{})            { *n = x.(Number) }
//   func (n Number) Successors() []interface{}      { return []interface{}{n - 1, n + 1} }
//   func (n Number) Cost(x interface{}) float64     { return 1 }
//   func (n Number) Estimate(x interface{}) float64 {
//     return math.Abs(10 - float64(x.(Number)))
//   }
//
//   var n Number = 10
//   if path, walk, err := astar.Search(&n); err == nil {
//     fmt.Println(path)
//     fmt.Println(walk)
//   }
//   // Output: [1 2 3 4 5 6 7 8 9 10] —— the solution.
//   // [1 2 3 4 5 6 7 8 9 10]         —— states explored by the algorithm
//                                        before it could find the best solution.
//
// You could allow only “subtract 7” and “add 5” operations to get to 10:
//
//   func (n Number) Successors() []interface{} { return []interface{}{n - 7, n + 5} }
//
//   // Output: [1 6 11 4 9 14 7 12 5 10]
//   // [1 6 11 4 9 16 14 7 12 -1 2 5 10]
//
// Or subtract 3, 7, and multiply by 9:
//
//   func (n Number) Successors() []interface{} { return []interface{}{n - 3, n - 7, n * 9} }
//
//   // Output: [1 9 6 3 27 20 13 10]
//   // [1 9 6 2 3 18 11 8 15 12 4 5 -2 -1 0 -5 -6 -4 -3 27 20 13 10]
//
// Etc.
//
package astar

import (
	"container/heap"
	"errors"
)

// ErrNotFound means that the final state cannot be reached from the given start state.
var ErrNotFound = errors.New("final state is not reachable")

// Any type is suitable for A* search as long as it can change its current state and tell
// legal moves from it.  Knowing costs and estimates helps, but not necessary.
type Interface interface {
	// Initial state.
	Start() interface{}

	// Is this state final?
	Finish() bool

	// Move to a new state.
	Move(interface{})

	// Available moves from the current state.
	Successors() []interface{}

	// Path cost between the current and the given state.
	Cost(interface{}) float64

	// Heuristic estimate of “how far to go?” between the given
	// and the final state. Smaller values mean closer.
	Estimate(interface{}) float64
}

type state struct {
	state          interface{}
	cost, estimate float64
	index          int
}

type states []*state

func (pq states) Len() int           { return len(pq) }
func (pq states) Empty() bool        { return len(pq) == 0 }
func (pq states) Less(n, j int) bool { return pq[n].cost+pq[n].estimate < pq[j].cost+pq[j].estimate }
func (pq states) Swap(n, j int) {
	pq[n], pq[j] = pq[j], pq[n]

	// Index is maintained for heap.Fix().
	pq[n].index = n
	pq[j].index = j
}

func (pq *states) Push(x interface{}) {
	n := len(*pq)
	item := x.(*state)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *states) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

// Search finds the p.Finish() state from a given p.Start() state by
// invoking p.Successors() and p.Move() at each step. Search returns two slices:
// 1) the shortest path to the final state, and a 2) sequence of explored states.
// If the shortest path cannot be found, ErrNotFound error is returned.
func Search(p Interface) ([]interface{}, []interface{}, error) {
	// Priority queue of states on the frontier.
	// Initialized with the start state.
	pq := states{{state: p.Start(), estimate: p.Estimate(p.Start())}}
	heap.Init(&pq)

	// States currently on the frontier.
	queuedLinks := map[interface{}]*state{}

	// States explored so far.
	explored := map[interface{}]bool{}

	// State transitions from start to finish (to reconstruct
	// the shortest path at the end of the search).
	transitions := map[interface{}]interface{}{}

	// Sequence of states in the order they have been explored.
	steps := []interface{}{}

	p.Move(p.Start())

	// Exhaust all successor states.
	for !pq.Empty() {
		// Pick a state with a minimum Cost() + Estimate() value.
		current := heap.Pop(&pq).(*state)
		delete(queuedLinks, current.state)
		explored[current.state] = true

		// Move to the new state.
		p.Move(current.state)

		steps = append(steps, current.state)

		// If the state is final, terminate.
		if p.Finish() {
			// Reconstruct the path from finish to start.
			return func() []interface{} {
				path := []interface{}{current.state}
				for {
					if _, ok := transitions[current.state]; !ok {
						break
					}
					current.state = transitions[current.state]

					// Reverse.
					path = append([]interface{}{current.state}, path...)

				}
				return path
			}(), steps, nil
		}

		for _, succ := range p.Successors() {
			// Don't re-explore.
			if explored[succ] {
				continue
			}

			// Path cost so far.
			cost := current.cost + p.Cost(succ)

			// Add a successor to the frontier.
			if queuedState, ok := queuedLinks[succ]; ok {
				// If the successor is already on the frontier,
				// update its path cost.
				if cost < queuedState.cost {
					queuedState.cost = cost
					heap.Fix(&pq, queuedState.index)
					transitions[succ] = current.state
				}
			} else {
				state := state{
					state:    succ,
					cost:     cost,
					estimate: p.Estimate(succ),
				}
				heap.Push(&pq, &state)
				queuedLinks[succ] = &state
				transitions[succ] = current.state
			}
		}
	}

	return nil, steps, ErrNotFound
}
