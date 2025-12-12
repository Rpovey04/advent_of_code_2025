package main

import (
	"fmt"
	"os"
	"strings"
)

type pathNode struct {
	prev *pathNode
	curr string
}

func displayPath(end pathNode) {
	current := &end
	str := "out"
	for current.curr != "you" {
		str = str + "->" + current.prev.curr
		current = current.prev
	}
	fmt.Println(str)
}

func findPaths(nodes map[string][]string) []pathNode {
	stack := make([]pathNode, 0)
	paths := make([]pathNode, 0) // will be path nodes where curr='end'
	stack = append(stack, pathNode{prev: nil, curr: "you"})
	var current pathNode
	for len(stack) != 0 {
		current = stack[len(stack)-1]
		stack = stack[0 : len(stack)-1]
		// fmt.Printf("%s: %v\n", current.curr, nodes[current.curr])
		for _, next := range nodes[current.curr] {
			currentClone := pathNode{prev: current.prev, curr: current.curr}
			nextNode := pathNode{prev: &currentClone, curr: next}
			// fmt.Printf("%s has prev %s\n", nextNode.curr, nextNode.prev.curr)
			if nextNode.curr == "out" {
				paths = append(paths, nextNode)
				displayPath(nextNode)
			} else {
				stack = append(stack, nextNode)
			}
		}
	}
	return paths
}

// PART ONE
// Idea:
// Simply use a DFS to exhaustively search all paths, adding complete paths to a slice instead of terminating
// (length of this slice is then the answer. Can print for validity)
// Verify that this works for test, then optimise to allow for this to work on input
// [Didn't have to do any optimizations, literally just a DFS]
//
//	When a path is complete and added to the slice, it will have all nodes from 'you' to 'end' added to an array
//	When a new node is evaluated, check the complete paths slice to see if the node appears anywhere
//	If the path appears anywhere in this slice, we can add:
//		The path from our 'you' to our current node to
//		The path from our current node in the slice to 'end'
//	May want to try and use maps for accessing the path from a current node to the end to allow for fast access of 'current node'
//	If a node is fully evaluated and doesn't have any paths that lead to 'end', note this so dead ends are not continually evaluated
func main() {
	data, _ := os.ReadFile(os.Args[1])
	lines := strings.Split(string(data), "\n")
	lines = lines[0 : len(lines)-1]

	nodes := make(map[string][]string)
	for _, line := range lines {
		srcSplit := strings.Split(line, ":")
		destinations := strings.Split(srcSplit[1], " ")
		destinations = destinations[1:len(destinations)]
		nodes[srcSplit[0]] = destinations
	}

	paths := findPaths(nodes)
	fmt.Printf("Found %d paths", len(paths))
}
