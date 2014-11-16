// Use A* search algorithm to traverse a maze.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"text/template"
	"time"

	"code.google.com/p/go.crypto/ssh/terminal"
	"github.com/pietv/astar"
)

var (
	// Default random maze size.
	defaultSize = "3x18"

	// Command line flags.
	euclidFlag    = flag.Bool("euclid", false, "use Euclid distance")
	manhattanFlag = flag.Bool("manhattan", true, "use Manhattan distance")
	estimFlag     = flag.Float64("estim", 1.5, "estimate multiplier")
	costFlag      = flag.Float64("cost", 1.0, "cost multiplier")
	demoFlag      = flag.Int("demo", 0, "run demo #")
	randomFlag    = flag.Bool("random", false, "generate a random maze")
	sizeFlag      = flag.String("size", defaultSize, "generate a random maze of size NxM")
	helpFlag      = flag.Bool("help", false, "show help")
)

var usage = `maze: demonstrate A* search algorithm traversing a maze.
Usage: maze [FILE] [-demo N] [-random] [-size NxM] [-help]
            [-euclid | -manhattan] [-mult MULTIPLIER]

With no FILE, use a demo or a random maze.

Flags:
 -demo N                 show a specific demo. #1..` + fmt.Sprintf("#%d", len(demos)) + `

 -manhattan              use Manhattan distance as a heuristic estimate (default).
 -euclid                 use Euclidean distance.
 -estim MULTIPLIER       multiply estimate value by MULTIPLIER.
 -cost MULTIPLIER        multiply cost value by MULTIPLIER.


 -random                 show a random maze.
 -size NxM               show a random maze of size NxM.

 -help                   show this help.`

func init() {
	rand.Seed(time.Now().UnixNano())
}

func readFile(filename string) []string {
	in, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read a maze from %q: %s\n", filename, err)
		os.Exit(1)
	}
	defer in.Close()

	maze := []string{}
	r := bufio.NewReader(in)
	for {
		line, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Cannot read a maze from %q: %s\n",
				filename, err)
			os.Exit(1)
		}
		if err == io.EOF {
			break
		}
		maze = append(maze, strings.Trim(line, "\n"))
	}
	return maze
}

func main() {
	flag.Parse()
	if *helpFlag {
		fmt.Println(usage)
		os.Exit(2)
	}

	var (
		demo  int
		title string
		maze  *Maze
	)

	// Random or specified demo screen.
	if *demoFlag == 0 {
		demo = rand.Intn(len(demos))

		// Randomly choose between showing a demo or a generated maze.
		if rand.Intn(2) == 0 {
			*randomFlag = true
		}
	} else {
		demo = *demoFlag - 1
		if demo >= len(demos) {
			fmt.Fprintf(os.Stderr, "Available demos are from #1 upto #%d.\n", len(demos))
			os.Exit(1)
		}
	}

	if flag.NArg() > 0 {
		// From FILE.
		maze = New(readFile(flag.Args()[0]))
		title = "Charming maze"
	} else if *randomFlag || *sizeFlag != defaultSize {
		// Random.
		var n, m int
		fmt.Sscanf(*sizeFlag, "%dx%d", &n, &m)

		if n == 0 || m == 0 {
			fmt.Fprintf(os.Stderr, "You should provide the size in the form of “-size=NxM”")
			os.Exit(1)
		}

		maze = NewRandomKruskal(n, m)
		title = "Randomly generated maze"
	} else {
		// Demo.
		maze = New(demos[demo].maze)
		title = fmt.Sprintf("Demo #%d. %s", demo+1, demos[demo].title)
	}

	// By default use Manhattan distance.
	if *euclidFlag {
		estimateFunc = genEuclidEstimate(*estimFlag)
	} else {
		estimateFunc = genManhattanEstimate(*estimFlag)
	}
	// Don't use fancy colorings if output is redirected.
	var medium string
	if terminal.IsTerminal(1) {
		medium = "Terminal"
	} else {
		medium = "File"
	}

	path, steps, err := astar.Search(maze)
	if err != nil {
		title = "Yikes! Could not find the path for this one"
	}

	template.Must(template.New("Maze").Funcs(helpers).Parse(TerminalTmpl)).ExecuteTemplate(os.Stdout, medium, struct {
		Title string
		Maze  [][]string
	}{
		Title: title,
		Maze:  maze.DrawMaze(path, steps),
	})
}
