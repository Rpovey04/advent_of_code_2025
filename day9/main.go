package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type coord struct {
	x int64
	y int64
}

func rect(c1 coord, c2 coord) int64 {
	return (abs(c1.x-c2.x) + 1) * (abs(c1.y-c2.y) + 1)
}

func abs(i int64) int64 {
	if i >= 0 {
		return i
	} else {
		return -i
	}
}

// Input is the position of tiles
// Find the biggest rectangle bound by two tiles
// Can be bound on the same column / row ((7,3) and (2, 3) form a rectangle of area 5)
func main() {
	data, _ := os.ReadFile(os.Args[1])
	lines := strings.Split(string(data), "\n")
	lines = lines[:len(lines)-1]

	// format coordinates
	tiles := make([]coord, 0)
	var coordString []string
	for _, l := range lines {
		coordString = strings.Split(l, ",")
		v1, _ := strconv.ParseInt(coordString[0], 10, 32)
		v2, _ := strconv.ParseInt(coordString[1], 10, 32)
		tiles = append(tiles, coord{x: v1, y: v2})
	}
	// find size of each tile combo and select the highest
	largest := int64(0)
	var currentRect int64
	for i, c1 := range tiles {
		for _, c2 := range tiles[i+1:] {
			currentRect = rect(c1, c2)
			if currentRect > largest {
				largest = currentRect
			}
		}
	}

	fmt.Printf("Answer: %d", largest)
}
