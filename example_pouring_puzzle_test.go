// Water pouring puzzle.
//
// Measure out 6 ounces of water using two glasses of 9 and 4 oz.
// You're allowed to pour water from one glass to another as well as emptying
// and refilling them.
//
// Use A* Search to solve it.
//
//
//                    =()=
//                   .-||--|
//                  .   ___|
//                  '==’
//                   ||
//                   ||
//       |     |     ||
//       |     |
//       |     |   |    |
//       |     |   |    |
//       |     |   |    |
//       |     |   |    |
//       +-----+   +----+
//        9 oz.     4 oz.
//
package astar_test

import (
	"os"
	"text/template"

	"github.com/pietv/astar"
)

// Glasses’ capacities and the goal.
type Puzzle struct {
	CapFirst  int
	CapSecond int
	Goal      int
}

// Glasses state.
type State struct {
	p      *Puzzle
	Action string
	First  int
	Second int
}

func (s State) Start() interface{} { return State{s.p, "Both Empty", 0, 0} }
func (s State) Finish() bool {
	// One of the glasses contains the goal amount.
	return s.p.Goal == s.First || s.p.Goal == s.Second
}
func (s *State) Move(x interface{})            { *s = x.(State) }
func (s State) Cost(x interface{}) float64     { return 1 }
func (s State) Estimate(x interface{}) float64 { return 1 }
func (s State) Successors() []interface{} {
	succ := []interface{}{}

	First, Second, CapFirst, CapSecond := s.First, s.Second, s.p.CapFirst, s.p.CapSecond

	// Fill first glass.
	succ = append(succ, State{s.p, "Fill First", CapFirst, Second})

	// Fill second glass.
	succ = append(succ, State{s.p, "Fill Second", First, CapSecond})

	// Empty first glass.
	succ = append(succ, State{s.p, "Empty First", 0, Second})

	// Empty second glass.
	succ = append(succ, State{s.p, "Empty Second", First, 0})

	// Pour from the first glass into the second.
	if First+Second > CapSecond {
		succ = append(succ, State{s.p, "First –> Second", First - (CapSecond - Second), CapSecond})
	} else {
		succ = append(succ, State{s.p, "First –> Second", 0, First + Second})
	}

	// Pour from the second glass into the first.
	if First+Second > CapFirst {
		succ = append(succ, State{s.p, "Second –> First", CapFirst, Second - (CapFirst - First)})
	} else {
		succ = append(succ, State{s.p, "Second –> First", First + Second, 0})
	}
	return succ
}

const PouringTmpl = `
{{range .}}  {{printf "%-16s (%v %v)\n" .Action .First .Second}}{{end}}
`

func ExampleSearch_pouringPuzzle() {
	if path, _, err := astar.Search(&State{p: &Puzzle{
		CapFirst:  9,
		CapSecond: 4,
		Goal:      6,
	}}); err == nil {
		template.Must(template.New("Pouring Puzzle").Parse(PouringTmpl)).Execute(os.Stdout, path)
	}
	// Output:
	//   Both Empty       (0 0)
	//   Fill First       (9 0)
	//   First –> Second  (5 4)
	//   Empty Second     (5 0)
	//   First –> Second  (1 4)
	//   Empty Second     (1 0)
	//   First –> Second  (0 1)
	//   Fill First       (9 1)
	//   First –> Second  (6 4)
}
