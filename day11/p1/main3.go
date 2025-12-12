package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

func reverseNodes(nodes map[string][]string) map[string][]string {
	nodesRev := make(map[string][]string)
	for src, dst := range nodes {
		for _, v := range dst {
			_, exists := nodesRev[v]
			if !exists {
				nodesRev[v] = []string{src}
			} else {
				nodesRev[v] = append(nodesRev[v], src)
			}
		}
	}
	return nodesRev
}

// assumes we are going backwards (start working backwards from end)
func findPaths(nodesBackwards map[string][]string, nodesForwards map[string][]string, end string, start string) [][]string {
	foundPaths := make(map[string][][]string)
	foundPaths[end] = [][]string{[]string{"end"}} // special case, won't try and find a path to itself
	queue := []string{end}
	for len(queue) != 0 {
		// BFS
		current := queue[0]
		if len(queue) == 1 {
			queue = []string{}
		} else {
			queue = queue[1:]
		}
		for _, next := range nodesBackwards[current] {
			_, exists := foundPaths[next]
			if !exists {
				queue = append(queue, next)
			}
		}
		// UPDATING PATHS (since every value of current is being evaluated)
		if current != end {
			foundPaths[current] = make([][]string, 0)
			for _, nextForward := range nodesForwards[current] {
				pathsToAdd, exists := foundPaths[nextForward]
				if exists { // if it doesn't exist, there is not a path to end here
					for _, pathToAdd := range pathsToAdd {
						foundPaths[current] = append(foundPaths[current], append(pathToAdd, current))
					}
				}
			}
		}
	}
	return foundPaths[start]
}

// This works but is way slower than the previous two, I obviously misunderstood what time complexity it would have
// Last thing I will try is doing all of this recursively since that will make caching the results of node traversal
// (and catching failures as well as found paths) much easier. In a seperate file again of course...
func main() {
	data, _ := os.ReadFile(os.Args[1])
	lines := strings.Split(string(data), "\n")
	lines = lines[0 : len(lines)-1]
	start := os.Args[2]

	nodes := make(map[string][]string)
	for _, line := range lines {
		srcSplit := strings.Split(line, ":")
		destinations := strings.Split(srcSplit[1], " ")
		destinations = destinations[1:len(destinations)]
		nodes[srcSplit[0]] = destinations
	}
	nodesRev := reverseNodes(nodes)
	paths := findPaths(nodesRev, nodes, "out", start)
	for i, p := range paths {
		slices.Reverse(p)
		fmt.Printf("%d) %v\n", i, p)
	}
	fmt.Printf("Found %d paths:\n", len(paths))
}
