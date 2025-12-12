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

func pathNodeToArr(path pathNode, start string) []string {
	current := &path
	res := make([]string, 0)
	for current.curr != start {
		res = append([]string{current.curr}, res...)
		current = current.prev
	}
	res = append([]string{start}, res...)
	return res
}

var useCache bool = false

func findExistingPaths(foundPaths [][]string, currentNode string) [][]string {
	res := make([][]string, 0)
	for _, foundPath := range foundPaths {
		for i, foundNode := range foundPath {
			if foundNode == currentNode {
				res = append(res, foundPath[i+1:])
			}
		}
	}

	if useCache {
		return res
	} else {
		return make([][]string, 0)
	}
}

func findPaths(nodes map[string][]string, start string) [][]string {
	stack := make([]pathNode, 0)
	paths := make([][]string, 0)
	stack = append(stack, pathNode{prev: nil, curr: start})
	var current pathNode
	for len(stack) != 0 {
		current = stack[len(stack)-1]
		stack = stack[0 : len(stack)-1]
		// fmt.Printf("%s: %v\n", current.curr, nodes[current.curr])
		for _, next := range nodes[current.curr] {
			currentClone := pathNode{prev: current.prev, curr: current.curr}
			existing := findExistingPaths(paths, next)
			if len(existing) == 0 || next == "out" {
				nextNode := pathNode{prev: &currentClone, curr: next}
				if nextNode.curr == "out" {
					paths = append(paths, pathNodeToArr(nextNode, start))
				} else {
					stack = append(stack, nextNode)
				}
			} else {
				// new logic here
				currentPath := pathNodeToArr(currentClone, start)
				// fmt.Printf("Checking segment | %s -> %s\n", current.curr, next)
				for _, existingPath := range existing {
					// fmt.Printf("Create new path by appending | %v to %v\n", currentPath, existingPath)
					fmt.Printf("Cached!\n")
					paths = append(paths, append(currentPath, existingPath...))
				}
			}
		}
	}
	return paths
}

func removeDuplicates(paths [][]string) [][]string {
	foundStr := make([]string, 0)
	foundIdx := make([]int, 0)
	dupe := false
	for i, p := range paths {
		checkStr := ""
		for _, c := range p {
			checkStr += c
		}
		dupe = false
		for _, str := range foundStr {
			if str == checkStr {
				dupe = true
			}
		}
		if !dupe {
			foundStr = append(foundStr, checkStr)
			foundIdx = append(foundIdx, i)
		}
	}
	res := make([][]string, 0)
	for _, idx := range foundIdx {
		res = append(res, paths[idx])
	}
	return res
}

// (see p1/main.go for more on the idea behind this solution)
// Changing a few things to make this easier to work with:
//  1. Store paths as []string instead of a linked list
//     Will still use a linked list for calculating the path initially, but convert to []string when storing
//     Makes working with them a lot easier: Can calculate new paths by just concatenating paths
//     Just generally more readable
//
// [Implement first change and see if that is enough, the second idea may not be necassary since
//
//	I can achieve what I was going for in p1/main.go with this idea alone]
//	2) For each node, store all known paths to 'out' in a map
//		Significantly more storage intensive but much faster: leads to a lot of duplicate data being stored
//		Will just update this when a new path is found for every node in the path
//
// [Implementing the first idea doesn't make as much sense as I thought it would
//
//	It misses paths since there is no guarantee that nodes are fully evaluated before being cached
//	It doesn't do anything to cache dead ends
//
// New idea:
// Start at the end node and work backwards: use a BFS to do this
// For every found node, run a DFS which goes back towards end
// In a seperate map of type map[string][][]string, save all of the found paths which lead to end
// Can then reference this map during the next evaluation
// Because this is a BFS from end -> start, evaluation should only require looking forward one node forward
// using the new map we are building
// when start is reached, evaluating this node will give us our answer
// radically different so will do in yet another file
// VERY MEMORY INTENSIVE
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

	paths := findPaths(nodes, start)
	paths = removeDuplicates(paths)
	useCache = true
	cachePath := findPaths(nodes, start)
	cachePath = removeDuplicates(cachePath)
	/*
		for _, p := range paths {
			fmt.Println(p)
		}
	*/
	fmt.Printf("No caching:%v\n\n-----\nWith caching:%v\n", paths, cachePath)
}
