package main

import (
	"fmt"
	"math"
	"os/exec"
	"strings"
)

// Maze drawing sequences.
var (
	SpaceRune  = " "
	WallRune   = "*"
	StartRune  = "S"
	FinishRune = "F"
	StepRune   = "·"
	PathRune   = "•"

	StepColor = Tput("setaf", 61)  // Cursor color StaleBlue3
	PathColor = Tput("setaf", 128) // Cursor color DarkViolet
	Reset     = Tput("sgr0")       // Cursor highlight reset
)

// Maze as printed to the terminal or file.
var TerminalTmpl = `
{{define "Terminal"}}
 {{.Title}}
{{range .Maze}}  {{range .}}{{colorize .}}{{end}}{{println}}{{end}}
 {{colorize legend}}
 Run with “-help” for available options.
{{end}}

{{define "File"}}{{.Title}}

{{range .Maze}}{{range .}}{{.}}{{end}}{{println}}{{end}}
{{legend}}
{{end}}`

var (
	helpers = map[string]interface{}{
		"legend": func() string {
			return WallRune + " - wall  " +
				StartRune + " - start  " +
				strings.Repeat(StepRune, 3) + " - explored  " +
				strings.Repeat(PathRune, 3) + " - shortest path"
		},
		"colorize": func(in string) (out string) {
			for _, s := range strings.Split(in, "") {
				switch s {
				case StartRune, FinishRune, PathRune:
					out += PathColor + s + Reset
				case StepRune:
					out += StepColor + s + Reset
				default:
					out += s
				}
			}
			return out
		},
	}

	estimateFunc = genEuclidEstimate(*estimFlag)
)

// genManhattanEstimate generates a Manhattan distance Estimate() function
// with custom multiplier.
func genManhattanEstimate(multiplier float64) func(interface{}, interface{}) float64 {
	return func(finish, neighbor interface{}) float64 {
		var (
			aI = float64(neighbor.(Location).i)
			bI = float64(finish.(Location).i)
			aJ = float64(neighbor.(Location).j)
			bJ = float64(finish.(Location).j)
		)
		return math.Abs(aI-bI) + math.Abs(aJ-bJ)*multiplier
	}
}

// genEuclidEstimate generates a Euclidean distance Estimate() function with
// custom multiplier.
func genEuclidEstimate(multiplier float64) func(interface{}, interface{}) float64 {
	return func(finish, neighbor interface{}) float64 {
		var (
			aI = float64(neighbor.(Location).i)
			bI = float64(finish.(Location).i)
			aJ = float64(neighbor.(Location).j)
			bJ = float64(finish.(Location).j)
		)
		return math.Sqrt(math.Pow(aI-bI, 2) + math.Pow(aJ-bJ, 2)*multiplier)
	}
}

type Location struct {
	i, j int
}

type Maze struct {
	maze                [][]string
	start, finish, curr Location
}

func (m Maze) Start() interface{}                { return m.start }
func (m Maze) Finish() bool                      { return m.curr == m.finish }
func (m *Maze) Move(t interface{})               { m.curr = t.(Location) }
func (m Maze) Cost(neighbor interface{}) float64 { return *costFlag }

func (m Maze) Estimate(neighbor interface{}) float64 {
	return estimateFunc(m.finish, neighbor)
}

func (m Maze) Successors() []interface{} {
	successors := []interface{}{}

	checkLocation := func(i, j int) {
		// The matrix is not necessarily rectangular.
		if i < 0 || j < 0 || i >= len(m.maze) || j >= len(m.maze[i]) {
			return
		}

		switch m.maze[i][j] {
		case SpaceRune, FinishRune:
			successors = append(successors, Location{i, j})
		}
	}

	i, j := m.curr.i, m.curr.j

	// North.
	checkLocation(i-1, j)
	// South.
	checkLocation(i+1, j)
	// West.
	checkLocation(i, j-1)
	// East.
	checkLocation(i, j+1)

	return successors
}

// Tput controls the terminal using tput(1). Command is a capability name and
// params are either strings or integers. Does nothing if a capability
// or tput(1) is not available. See terminfo(5) and infocmp(1) for a list of
// capability names.
func Tput(command string, params ...interface{}) string {
	args := []string{command}

	for _, param := range params {
		switch param.(type) {
		case string:
			args = append(args, param.(string))
		case int:
			args = append(args, fmt.Sprintf("%d", param))
		}
	}

	out, err := exec.Command("tput", args...).Output()
	if err != nil {
		return ""
	}
	return string(out)
}

// New initialize a maze with given slice of strings.
func New(m []string) *Maze {
	maze := make([][]string, len(m))

	var start, finish Location

	for i, line := range m {
		maze[i] = []string{}

		for j, rune := range strings.Split(line, "") {
			maze[i] = append(maze[i], rune)
			switch rune {
			case StartRune:
				start = Location{i, j}
			case FinishRune:
				finish = Location{i, j}
			}
		}
	}

	return &Maze{
		maze:   maze,
		start:  start,
		finish: finish,
		curr:   start,
	}
}

// DrawMaze applies path and explored state lists to a maze
// and returns it as a matrix of strings.
func (m *Maze) DrawMaze(path, steps []interface{}) [][]string {
	states := map[Location]string{}

	// Apply explored states first.
	for _, state := range steps {
		states[state.(Location)] = StepRune
	}

	// Apply path.
	for _, state := range path {
		states[state.(Location)] = PathRune
	}

	maze := make([][]string, len(m.maze))
	for i := 0; i < len(m.maze); i++ {
		maze[i] = make([]string, len(m.maze[i]))
		for j := 0; j < len(maze[i]); j++ {
			state, ok := states[Location{i, j}]

			if ok && m.maze[i][j] == SpaceRune {
				maze[i][j] = state
			} else {
				maze[i][j] = m.maze[i][j]
			}
		}
	}

	return maze
}
