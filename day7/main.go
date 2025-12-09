package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	data, _ := os.ReadFile(os.Args[1])
	lines := strings.Split(string(data), "\n")
	byteLines := make([][]byte, 0)
	for _, l := range lines {
		byteLines = append(byteLines, []byte(l))
	}

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
}
