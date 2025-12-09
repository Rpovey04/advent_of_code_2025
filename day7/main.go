package main

import (
	"fmt"
	"os"
	"strings"
)

var totalCalls int64 = 0
var cached map[string]int64 = make(map[string]int64)

func cloneState(state [][]byte) [][]byte {
	newState := make([][]byte, len(state))
	for i := 0; i < len(newState); i++ {
		newState[i] = make([]byte, len(state[i]))
		copy(newState[i], state[i])
	}
	return newState
}

func evalState(state [][]byte) int64 {
	// convert to string to use as key
	key := ""
	for _, l := range state {
		key += string(l)
	}
	res, exists := cached[key]
	if exists {
		fmt.Printf("Cached!\n")
		return res
	}
	newRes := propogate(state)
	cached[key] = newRes
	return newRes
}

// only pass the everything the line where the split happens and everything after
func propogate(state [][]byte) int64 {
	/*
		for _, l := range state {
			fmt.Printf("%s\n", string(l))
		}

		fmt.Printf("Started: %d with %d lines\n", totalCalls, len(state))
	*/
	totalCalls += 1
	if totalCalls%1000000 == 0 {
		fmt.Printf("|")
	}
	for l := 0; l < len(state)-2; l++ {
		for i := 0; i < len(state[0])-1; i++ {
			// handle first line
			if state[l][i] == 'S' {
				state[l+1][i] = '|'
			}
			// splitting
			if state[l][i] == '|' && state[l+1][i] == '^' {
				left := cloneState(state[l+1:])
				left[0][i-1] = '|'
				right := cloneState(state[l+1:])
				right[0][i+1] = '|'
				return evalState(left) + evalState(right)
			}
			// beam goes down
			if state[l][i] == '|' && state[l+1][i] == '.' {
				state[l+1][i] = '|'
			}
		}
	}
	return 1
}

func main() {
	data, _ := os.ReadFile(os.Args[1])
	lines := strings.Split(string(data), "\n")
	byteLines := make([][]byte, 0)
	for _, l := range lines {
		byteLines = append(byteLines, []byte(l))
	}

	total := propogate(byteLines)
	fmt.Printf("Total: %d", total)
	/*
		total := 0
		for l := 0; l < len(byteLines)-2; l++ {
			for i := 0; i < len(byteLines[0])-1; i++ {
				// handle first line
				if byteLines[l][i] == 'S' {
					byteLines[l+1][i] = '|'
				}
				// splitting
				if byteLines[l][i] == '|' && byteLines[l+1][i] == '^' {
					total += 1
					byteLines[l+1][i-1] = '|'
					byteLines[l+1][i+1] = '|'
				}
				// beam goes down
				if byteLines[l][i] == '|' && byteLines[l+1][i] == '.' {
					byteLines[l+1][i] = '|'
				}
			}
		}

		for _, l := range byteLines {
			fmt.Printf("%s\n", string(l))
		}
		fmt.Printf("Total: %d\n", total)
	*/
}
