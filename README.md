A* [![GoDoc](https://godoc.org/github.com/pietv/astar?status.png)](https://godoc.org/github.com/pietv/astar) [![Build Status](https://drone.io/github.com/pietv/astar/status.png)](https://drone.io/github.com/pietv/astar/latest)
==

Implementation of the A* Search algorithm.

Install
=======

```shell
$ go get github.com/pietv/astar
```

Overview
========

A* Search is one of the intelligent exhaustive search algorithms which gets from
Start to Finish by exploring successor states.

![A* Steps](/screenshots/astar.png)
It's intelligent because it uses a special guidance in selecting the states that
are going to be explored next. The algorithm uses a minimum value of a sum of
next successor cost and a heuristic estimate (distance, for example) of how close
that next successor to Finish. They are named Cost and Estimate in this implementation.

Without any guidance (that is when both Cost and Estimate values are constant),
it explores all successor states, making it essentially the Breadth First Search
algorithm (Go's `container/heap` implementation behaves like a queue if keys are equal).

Depending on whether Cost and Estimate are constant or not, A* Search behaves
like other well-known algorithms:

Cost    |Estimate  |Behavior
--------|----------|-----------------------------------------------
const   |const     |[Breadth First Search](http://en.wikipedia.org/wiki/Breadth-first_search)
variable|const     |[Dijkstra's Shortest Path] (http://en.wikipedia.org/wiki/Dijkstra%27s_algorithm) / [Uniform Cost Search](http://en.wikipedia.org/wiki/Uniform-cost_search)
const   |variable  |[Greedy Best First Search](http://en.wikipedia.org/wiki/Best-first_search)
variable|variable  |[A*](http://en.wikipedia.org/wiki/A*_search_algorithm)


It is not necessary to use a graph data structure. Crawling the internet
and feeding harvested links as successors would do or, as another example,
the provided [maze demo](/cmd/maze/maze.go) uses a rectangular matrix
and uses surrounding cells that are not walls as successors.

![Maze Demo](/screenshots/maze.png)
Install the demo:

```shell
$ go install github.com/pietv/astar/cmd/maze
```




