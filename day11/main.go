package main

import (
	"fmt"
	"os"
	"strings"
)

var cache map[string][][]string = make(map[string][][]string)

func findPaths(start string, end string, nodes map[string][]string) [][]string {
	// base case
	if start == end {
		return [][]string{[]string{end}}
	}
	var paths [][]string
	// check cache
	paths, exists := cache[start]
	if exists {
		fmt.Printf("|")
		// fmt.Println("Cached!")
	} else { // standard cases
		paths = make([][]string, 0)
		nextNodes := nodes[start]
		for _, nextNode := range nextNodes {
			nextPaths := findPaths(nextNode, end, nodes)
			for _, nextPath := range nextPaths {
				paths = append(paths, append([]string{start}, nextPath...))
			}
		}
		cache[start] = paths
	}
	return paths
}

// Saw the trick of coupling fft and dac together to use as a key on the subreddit
// But I had the idea that paths that did or didn't satisfy the fft and dac requirement would have to be cached seperately
// So I don't consider this cheating
type position struct {
	start string
	fft   bool
	dac   bool
}

// not finding the paths, only the number of them
var cacheCount map[position]int = make(map[position]int)

func findNumPaths(pos position, end string, nodes map[string][]string) int {
	// base case
	if pos.start == end {
		if pos.fft && pos.dac {
			return 1
		} else {
			return 0
		}
	}
	// updating found fft and dac
	if pos.start == "fft" {
		pos.fft = true
	}
	if pos.start == "dac" {
		pos.dac = true
	}

	// checking cache
	sum, exists := cacheCount[pos]
	if exists {
		fmt.Printf("|")
	} else { // standard case
		sum = 0
		nextNodes := nodes[pos.start]
		for _, nextNode := range nextNodes {
			nextPos := position{start: nextNode, fft: pos.fft, dac: pos.dac}
			sum += findNumPaths(nextPos, end, nodes)
		}
		cacheCount[pos] = sum
		fmt.Printf("\nCached %s\n", pos.start)
	}
	return sum
}

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

	/*
		paths := findPaths(start, "out", nodes)
		for i, p := range paths {
			fmt.Printf("%d) %v\n", i, p)
		}
	*/
	numPaths := findNumPaths(position{start, false, false}, "out", nodes)
	fmt.Printf("Found %d paths\n", numPaths)
}
