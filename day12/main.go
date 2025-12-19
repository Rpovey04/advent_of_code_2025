package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

/*
P1
*/
type weightedRegion struct {
	r region
	w int
}

func cmpWeighted(a, b weightedRegion) int {
	return a.w - b.w
}

func generateNextState(current region, shape [][]bool, nShapes []int, i int, j int, lowestH int) (weightedRegion, bool) {
	if !(shape[0][0] && current.grid[i-1][j-1]) &&
		!(shape[0][1] && current.grid[i-1][j]) &&
		!(shape[0][2] && current.grid[i-1][j+1]) &&
		!(shape[1][0] && current.grid[i][j-1]) &&
		!(shape[1][1] && current.grid[i][j]) &&
		!(shape[1][2] && current.grid[i][j+1]) &&
		!(shape[2][0] && current.grid[i+1][j-1]) &&
		!(shape[2][1] && current.grid[i+1][j]) &&
		!(shape[2][2] && current.grid[i+1][j+1]) {
		newGrid := make([][]bool, len(current.grid))
		for g := 0; g < len(current.grid); g++ {
			newGrid[g] = make([]bool, len(current.grid[g]))
			copy(newGrid[g], current.grid[g])
		}
		newState := region{numShapes: nShapes, x: current.x, y: current.y, grid: newGrid}
		newState.grid[i-1][j-1] = shape[0][0] || newState.grid[i-1][j-1]
		newState.grid[i-1][j] = shape[0][1] || newState.grid[i-1][j]
		newState.grid[i-1][j+1] = shape[0][2] || newState.grid[i-1][j+1]
		newState.grid[i][j-1] = shape[1][0] || newState.grid[i][j-1]
		newState.grid[i][j] = shape[1][1] || newState.grid[i][j]
		newState.grid[i][j+1] = shape[1][2] || newState.grid[i][j+1]
		newState.grid[i+1][j-1] = shape[2][0] || newState.grid[i+1][j-1]
		newState.grid[i+1][j] = shape[2][1] || newState.grid[i+1][j]
		newState.grid[i+1][j+1] = shape[2][2] || newState.grid[i+1][j+1]
		newStateHeuristic := bbHeuristic(newState)
		if newStateHeuristic < lowestH {
			return weightedRegion{r: newState, w: newStateHeuristic}, true
		}
	}
	return weightedRegion{}, false
}

// Assuming that all shapes have dimensions 3x3 (because they all do)
// nShapes []int passed here to initialise the newState variable to avoid having to do an extra loop
// Changed to return only the next state with the lowest valued heuristic
// Avoids having to sort and process large slices
// Idea for optimization: try only searching for new states around the edges of the bounding box
// Since optimal packings should always have shapes that are densly packed
// Search with a width of 6 along both outer lines of the box
// Can do this with just the br corner
func generateNextStates(current region, shape [][]bool, nShapes []int) (weightedRegion, bool) {
	b, r := getBRCorner(current)
	lowestHeuristic := math.MaxInt32
	bestRegion := weightedRegion{}
	found := false
	// "bottom" row
	for i := max(1, b-3); i < min(len(current.grid)-1, b+3); i++ {
		for j := 1; j < min(len(current.grid[0])-1, r+3); j++ {
			newState, exists := generateNextState(current, shape, nShapes, i, j, lowestHeuristic)
			if exists {
				found = true
				bestRegion = newState
				lowestHeuristic = newState.w
			}
		}
	}
	// "right" column
	for i := 1; i < min(len(current.grid)-1, b+3); i++ {
		for j := max(1, r-3); j < min(len(current.grid[0])-1, r+3); j++ {
			newState, exists := generateNextState(current, shape, nShapes, i, j, lowestHeuristic)
			if exists {
				found = true
				bestRegion = newState
				lowestHeuristic = newState.w
			}
		}
	}

	return bestRegion, found
}

func getBRCorner(current region) (int, int) {
	b := 0 // bottom
	r := 0 // right
	for i := 0; i < len(current.grid); i++ {
		for j := 0; j < len(current.grid[0]); j++ {
			if current.grid[i][j] {
				b = max(b, i)
				r = max(r, j)
			}
		}
	}
	return b, r
}

// area of the smallest bounding box of all current shapes
// find the furthest point in each orthogonal direction and construct (tl, bl) corners
func bbHeuristic(current region) int {
	b, r := getBRCorner(current)
	if b == 0 && r == 0 {
		return 0
	} else {
		return (b + 1) * (r + 1)
	}
}

func closedSpaceHeuristic(current region) int {
	total := 0
	for i := 0; i < len(current.grid); i++ {
		for j := 0; j < len(current.grid[0]); j++ {
			if !current.grid[i][j] {
				count := 0
				if i == 0 {
					count += 1
				} else if current.grid[i-1][j] {
					count += 1
				}
				if i == len(current.grid)-1 {
					count += 1
				} else if current.grid[i+1][j] {
					count += 1
				}
				if j == 0 {
					count += 1
				} else if current.grid[i][j-1] {
					count += 1
				}
				if j == len(current.grid[0])-1 {
					count += 1
				} else if current.grid[i][j+1] {
					count += 1
				}
				if count == 4 {
					total += 1
				}
			}
		}
	}
	return total
}

var runCurr int = 0
var runNum int = 0
var runSucc int = 0

// region is the state being searched (combination of grid and required shapes at any given time form the state)
func fitShape(current region, shapes []shape, depth int) (bool, region) {
	// comment out, horribly inefficient but quite satisfying to watch
	b, r := getBRCorner(current)
	if false {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		fmt.Printf("Run %d/%d (%d successful)\tDepth: %d\tCorner: (%d,%d)\t", runCurr, runNum, runSucc, depth, b, r)
		displayRegion(current)
	}
	frontier := make([]weightedRegion, 0)
	added := false
	for i, s := range shapes {
		if current.numShapes[i] > 0 {
			added = true
			nextNumShapes := make([]int, len(current.numShapes))
			copy(nextNumShapes, current.numShapes)
			nextNumShapes[i] -= 1
			for _, v := range s.s {
				newState, found := generateNextStates(current, v, nextNumShapes)
				if found {
					frontier = append(frontier, newState)
				}
			}
		}
	}
	slices.SortFunc(frontier, cmpWeighted)

	// numShapes[i] = 0 for all i, success
	if !added {
		return true, current
	} else if len(frontier) == 0 {
		return false, region{}
	} else {
		return fitShape(frontier[0].r, shapes, depth+1)
	}
	/*
				for _, r := range frontier {
					// fmt.Printf("Heuristic: %d\n", r.w)
					// displayRegion(r.r)
					success, successRegion := fitShape(r.r, shapes, depth+1)
					if success {
						return true, successRegion
					}
				}
		}
		return false, region{}
	*/
}

/*
Main function and loading data
*/
type region struct {
	numShapes []int
	grid      [][]bool
	x         int
	y         int
}

type shape struct {
	s map[string][][]bool
}

func displayShape(shape [][]bool) {
	fmt.Printf("%dx%d\n", len(shape), len(shape[0]))
	for i := 0; i < len(shape); i++ {
		for j := 0; j < len(shape); j++ {
			if shape[i][j] {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func displayAllShapes(s shape) {
	for k, v := range s.s {
		fmt.Printf("%s: ", k)
		displayShape(v)
	}
}

func displayRegion(r region) {
	fmt.Printf("Dimensions: %d, %d\t Requires shapes: %v\n", r.x, r.y, r.numShapes)
	for i := 0; i < r.y; i++ {
		for j := 0; j < r.x; j++ {
			if r.grid[j][i] {
				fmt.Printf("@")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func initRegionGrid(x int, y int) [][]bool {
	grid := make([][]bool, 0)
	for i := 0; i < x; i++ {
		grid = append(grid, make([]bool, 0))
		for j := 0; j < y; j++ {
			grid[i] = append(grid[i], false)
		}
	}
	return grid
}

func rotateShape(grid [][]bool) [][]bool {
	newShape := make([][]bool, len(grid))
	for i := 0; i < len(grid); i++ {
		newShape[i] = make([]bool, len(grid[0]))
		for j := 0; j < len(grid[0]); j++ {
			newShape[i][j] = grid[j][(len(grid)-i)-1]
		}
	}
	return newShape
}

func flipShape(grid [][]bool) [][]bool {
	newShape := make([][]bool, len(grid))
	for i := 0; i < len(newShape); i++ {
		newShape[i] = make([]bool, len(grid[len(grid)-i-1]))
		copy(newShape[i], grid[len(grid)-i-1])
	}
	return newShape
}

// generates a hash map with format map["direciton/flip"] = bool grid
// keys follow the format u (standard), r (->), d (\/), l (<-)
// followed by an f if the shape should be flipped after rotation
func fromShapeGrid(grid [][]bool) shape {
	newShape := shape{s: make(map[string][][]bool)}
	// rotation
	newShape.s["u"] = grid
	newShape.s["r"] = rotateShape(newShape.s["u"])
	newShape.s["d"] = rotateShape(newShape.s["r"])
	newShape.s["l"] = rotateShape(newShape.s["d"])
	// flips
	newShape.s["uf"] = flipShape(newShape.s["u"])
	newShape.s["rf"] = flipShape(newShape.s["r"])
	newShape.s["df"] = flipShape(newShape.s["d"])
	newShape.s["lf"] = flipShape(newShape.s["l"])

	return newShape
}

// PART 1: Cheated a little bit by looking at the subreddit for ideas
// I had the idea of how it would be some version of a DFS with some heuristic and that's basically all I saw online
// An idea for the heuristic I had was just the amount of empty space after the next step forward (since this would encourage tight packing)
// The state is created when each line begins being processed (All false seperated by , at the start of a new line)
// State is changed by adding true to all tiles that would be covered by a shape
//
// Main thing I'm unsure about is how to realise that a line cannot be satisfied
// If the heuristic is suitable, then this would be known when the state space cannot be expanded for the first time
// So I think I'll try implementing a search and trying a few heuristics as I can't think of anything else
// Heuristic ideas: (Assume that the minimum heuristic value is selected)
//
//  1. "Lost / Bound tiles", a tile is lost / bound if the 4 tiles adjacent tiles are occupied
//     A lost / bound tile can never be filled since doing so would cause overlap
//     Should also count the edges as being "occupied tiles" when calculating this
//     0(n) time complexity (where n is the total number of tiles)
//     1 doesn't seem to work but it's hard to know if that's because it's a bad heuristic or bad implementation
//     It's not great, for example for the "n" shape it doesn't detect an enclosed area since there's always two empty tiles
//
//  2. I saw someone on the subreddit talk about a "smallest bounding rectangle" heuristic
//     While I'm happy to try implementing that, I don't see how it would be a complete heuristic
//     For example when placing the first shape, there is no difference in which rotation it is placed since all
//     shapes have area 3x3 and so will have the same smallest bounding rectangle
//     however the optimal placement of the first shape needs to have any empty space exposed
//
//  3. Perimeter of the filled in area could be a good one, but still doesn't quite feel complete
func main() {
	data, _ := os.ReadFile(os.Args[1])
	lines := strings.Split(string(data), "\n")
	lines = lines[0 : len(lines)-1]

	lineIdx := 0
	// Read shapes
	var newShape [][]bool
	// foundShapes := make([][][]bool, 0)
	foundShapes := make([]shape, 0)
	for lineIdx < len(lines) && len(lines[lineIdx]) < 10 {
		if len(lines[lineIdx]) == 0 {
			foundShapes = append(foundShapes, fromShapeGrid(newShape))
		} else {
			firstChar := byte(lines[lineIdx][0])
			if firstChar == byte('#') || firstChar == byte('.') {
				// Add to detected shape
				newShape = append(newShape, make([]bool, 0))
				for _, c := range lines[lineIdx][0:len(lines[lineIdx])] {
					newShape[len(newShape)-1] = append(newShape[len(newShape)-1], c == rune('#'))
				}
			} else {
				// Start reading new detected shape
				newShape = make([][]bool, 0)
			}
		}
		lineIdx += 1
	}
	// Read constraints
	regions := make([]region, 0)
	for lineIdx < len(lines) {
		splitStr := strings.Split(lines[lineIdx], " ")
		dimStr := strings.Split(splitStr[0], "x")
		xDim, _ := strconv.ParseInt(dimStr[0], 10, 32)
		yDim, _ := strconv.ParseInt(dimStr[1][0:len(dimStr[1])-1], 10, 32)
		numShapes := make([]int, len(splitStr)-1)
		for idx, num := range splitStr[1:len(splitStr)] {
			tmp, _ := strconv.ParseInt(num, 10, 16)
			numShapes[idx] = int(tmp)
		}
		newRegion := region{numShapes: numShapes, x: int(xDim), y: int(yDim)}
		newRegion.grid = initRegionGrid(int(xDim), int(yDim))
		displayRegion(newRegion)
		regions = append(regions, newRegion)
		// fmt.Println(regions[len(regions)-1])
		lineIdx += 1
	}

	/*
		for _, shape := range foundShapes {
			displayAllShapes(shape)
		}
	*/
	total := 0
	runNum = len(regions)
	for i, r := range regions {
		runCurr = i
		success, successRegion := fitShape(r, foundShapes, 0)
		if success {
			total += 1
			runSucc = total
			fmt.Printf("Region %d successful! (total=%d)\n", i, total)
			displayRegion(successRegion)
		} else {
			fmt.Printf("Region %d unsuccessful! (total=%d)\n", i, total)
		}
	}
	fmt.Printf("Answer: %d\n", total)
}
