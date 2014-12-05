package main

import (
	"fmt"
	"math"
	"os/exec"
	"strings"
)

// Maze drawing sequences.
var (
	spaceRune  = " "
	wallRune   = "*"
	startRune  = "S"
	finishRune = "F"
	stepRune   = "·"
	pathRune   = "•"

	stepColor = tput("setaf", 61)  // Cursor color StaleBlue3
	pathColor = tput("setaf", 128) // Cursor color DarkViolet
	reset     = tput("sgr0")       // Cursor highlight reset
)

// Maze as printed to the terminal or file.
var terminalTmpl = `
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
			return wallRune + " - wall  " +
				startRune + " - start  " +
				strings.Repeat(stepRune, 3) + " - explored  " +
				strings.Repeat(pathRune, 3) + " - shortest path"
		},
		"colorize": func(in string) (out string) {
			for _, s := range strings.Split(in, "") {
				switch s {
				case startRune, finishRune, pathRune:
					out += pathColor + s + reset
				case stepRune:
					out += stepColor + s + reset
				default:
					out += s
				}
			}
			return out
		},
	}

	estimateFunc = genEuclidEstimate(*estimateFlag)
)

// genManhattanEstimate generates a Manhattan distance Estimate() function
// with a custom multiplier.
func genManhattanEstimate(multiplier float64) func(interface{}, interface{}) float64 {
	return func(finish, neighbor interface{}) float64 {
		var (
			aI = float64(neighbor.(location).i)
			bI = float64(finish.(location).i)
			aJ = float64(neighbor.(location).j)
			bJ = float64(finish.(location).j)
		)
		return (math.Abs(aI-bI) + math.Abs(aJ-bJ)) * multiplier
	}
}

// genEuclidEstimate generates a Euclidean distance Estimate() function with a
// custom multiplier.
func genEuclidEstimate(multiplier float64) func(interface{}, interface{}) float64 {
	return func(finish, neighbor interface{}) float64 {
		var (
			aI = float64(neighbor.(location).i)
			bI = float64(finish.(location).i)
			aJ = float64(neighbor.(location).j)
			bJ = float64(finish.(location).j)
		)
		return math.Sqrt(math.Pow(aI-bI, 2)+math.Pow(aJ-bJ, 2)) * multiplier
	}
}

type location struct {
	i, j int
}

type maze struct {
	maze                [][]string
	start, finish, curr location
}

func (m maze) Start() interface{}                { return m.start }
func (m maze) Finish() bool                      { return m.curr == m.finish }
func (m *maze) Move(t interface{})               { m.curr = t.(location) }
func (m maze) Cost(neighbor interface{}) float64 { return *costFlag }

func (m maze) Estimate(neighbor interface{}) float64 {
	return estimateFunc(m.finish, neighbor)
}

func (m maze) Successors() []interface{} {
	successors := []interface{}{}

	checkLocation := func(i, j int) {
		// The matrix is not necessarily rectangular.
		if i < 0 || j < 0 || i >= len(m.maze) || j >= len(m.maze[i]) {
			return
		}

		switch m.maze[i][j] {
		case spaceRune, finishRune:
			successors = append(successors, location{i, j})
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

// tput controls the terminal using tput(1). Command is a capability name and
// params are either strings or integers. Does nothing if a capability
// or tput(1) is not available. See terminfo(5) and infocmp(1) for a list of
// capability names.
func tput(command string, params ...interface{}) string {
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

// new initialize a maze with a given slice of strings.
func new(lines []string) *maze {
	m := make([][]string, len(lines))

	var start, finish location

	for i, line := range lines {
		m[i] = []string{}

		for j, rune := range strings.Split(line, "") {
			m[i] = append(m[i], rune)
			switch rune {
			case startRune:
				start = location{i, j}
			case finishRune:
				finish = location{i, j}
			}
		}
	}

	return &maze{
		maze:   m,
		start:  start,
		finish: finish,
		curr:   start,
	}
}

// drawMaze applies path and explored states to a maze
// and returns it as a matrix of strings.
func (m *maze) drawMaze(path, steps []interface{}) [][]string {
	states := map[location]string{}

	for _, state := range steps {
		states[state.(location)] = stepRune
	}

	for _, state := range path {
		states[state.(location)] = pathRune
	}

	maze := make([][]string, len(m.maze))
	for i := 0; i < len(m.maze); i++ {
		maze[i] = make([]string, len(m.maze[i]))
		for j := 0; j < len(maze[i]); j++ {
			state, ok := states[location{i, j}]
			if ok && m.maze[i][j] == spaceRune {
				maze[i][j] = state
			} else {
				maze[i][j] = m.maze[i][j]
			}
		}
	}

	return maze
}
