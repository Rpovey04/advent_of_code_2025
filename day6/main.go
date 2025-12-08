package main

import (
	"fmt"
	"os"
	"strings"
)

func removeIndexFromSlice(slice []string, i int) []string {
	return append(slice[:i], slice[i+1:]...)
}

func removeFromSlice(v []string, elem string) []string {
	i := 0
	for i < len(v) {
		if v[i] == elem {
			v = removeIndexFromSlice(v, i)
			i = 0
		} else {
			i += 1
		}
	}
	return v
}

// For every calculation <calcs>
// Need <longest> formats
// Each of length <params>
func getFormattedNumbers(slice [][]string, numCalcs int, numParams int, longest int) [][]string {
	formatted := make(string[][], numCalcs)
	for c := 0; c < numCalcs; c++ {
		formatted[c] = make(string[], longest)
		for p := 0; p < numParams; p++ {
			printf()
		}
		// fmt.Printf("%s|", slice[i][j])
	}
	return slice
}

func main() {
	data, _ := os.ReadFile(os.Args[1])
	lines := strings.Split(string(data), "\n")
	lines = removeIndexFromSlice(lines, len(lines)-1)
	longest := 0
	allNumbers := make([][]string, 0)
	for _, line := range lines {
		numbers := removeFromSlice(strings.Split(line, " "), "")
		numbers = removeFromSlice(numbers, "\n")
		for _, k := range numbers {
			if len(k) > longest {
				longest = len(k)
			}
		}
		allNumbers = append(allNumbers, numbers)
	}
	numCalcs := len(allNumbers[0])
	numParams := len(allNumbers) - 1
	fmt.Printf("Longest: %d\tParameters: %d\tCalculations: %d\n\n", longest, numParams, numCalcs)

	getFormattedNumbers(allNumbers[:len(allNumbers)-1], numCalcs, numParams, longest)
}
