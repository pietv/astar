// Finding the shortest path between Arad and Bucharest 
// on a Romanian road map fragment. The road map is represented 
// as an undirected graph; edge costs are distances between cities.
//
// An example from Stuart Russell and Peter Norvig's
// “Artificial Intelligence. A Modern Approach”, 3rd ed., 2009, p. 68.
package astar_test

import (
	"fmt"

	"github.com/pietv/astar"
)

type Graph struct {
	edges map[string]map[string]float64
	curr  string
}

func (g Graph) Start() interface{} { return "Arad" }
func (g Graph) Finish() bool       { return "Bucharest" == g.curr }

func (g *Graph) Move(state interface{}) { g.curr = state.(string) }
func (g Graph) Successors() []interface{} {
	successors := []interface{}{}
	for succ := range g.edges[g.curr] {
		successors = append(successors, succ)
	}

	return successors
}

func (g Graph) Cost(given interface{}) float64 {
	return g.edges[g.curr][given.(string)]
}

// Precalculated distances to Bucharest (Ibid., figure 3.22, p. 93.)
func (g Graph) Estimate(given interface{}) float64 {
	estimates := map[string]float64{
		"Arad":           366,
		"Bucharest":      0,
		"Craiova":        160,
		"Drobeta":        242,
		"Eforie":         161,
		"Făgăraș":        176,
		"Giurgiu":        77,
		"Hârșova":        151,
		"Iași":           226,
		"Lugoj":          244,
		"Mehadia":        241,
		"Neamt":          234,
		"Oradea":         380,
		"Pitești":        100,
		"Râmnicu Vâlcea": 193,
		"Sibiu":          253,
		"Timișoara":      329,
		"Urziceni":       80,
		"Vaslui":         199,
		"Zerind":         374,
	}

	return estimates[given.(string)]
}

func ExampleSearch_graphTraversal() {
	g := &Graph{
		edges: map[string]map[string]float64{
			"Arad":           {"Zerind": 75, "Timișoara": 118, "Sibiu": 140},
			"Bucharest":      {"Pitești": 101, "Făgăraș": 211, "Urziceni": 85, "Giurgiu": 90},
			"Craiova":        {"Drobeta": 120, "Râmnicu Vâlcea": 146, "Pitești": 138},
			"Drobeta":        {"Mehadia": 75, "Craiova": 120},
			"Eforie":         {"Hârșova": 86},
			"Făgăraș":        {"Sibiu": 99, "Bucharest": 211},
			"Giurgiu":        {"Bucharest": 90},
			"Hârșova":        {"Urziceni": 98, "Eforie": 86},
			"Iași":           {"Neamt": 87, "Vaslui": 92},
			"Lugoj":          {"Timișoara": 111, "Mehadia": 70},
			"Mehadia":        {"Lugoj": 70, "Drobeta": 75},
			"Oradea":         {"Zerind": 71, "Sibiu": 151},
			"Pitești":        {"Râmnicu Vâlcea": 97, "Craiova": 138, "Bucharest": 101},
			"Râmnicu Vâlcea": {"Sibiu": 80, "Pitești": 97, "Craiova": 146},
			"Sibiu":          {"Arad": 140, "Oradea": 151, "Făgăraș": 99, "Râmnicu Vâlcea": 80},
			"Timișoara":      {"Arad": 118, "Lugoj": 111},
			"Urziceni":       {"Bucharest": 85, "Vaslui": 142, "Hârșova": 98},
			"Vaslui":         {"Urziceni": 142, "Iași": 92},
			"Zerind":         {"Arad": 75, "Oradea": 71},
		}}
	if path, _, err := astar.Search(g); err == nil {
		fmt.Printf("%q\n", path)
	}
	// Output: ["Arad" "Sibiu" "Râmnicu Vâlcea" "Pitești" "Bucharest"]
}
