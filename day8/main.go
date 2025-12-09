package main

import (
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

// Each line is a 3d coordinate
// Connect pairs of the closest coordinates and keep track of circuits (connected coordinates)
// Result is the product of the three largest circuits (test should be 40)

type coord struct {
	x int64
	y int64
	z int64
}

type pair struct {
	a coord
	b coord
}

func coordKey(c coord) string {
	return strconv.FormatInt(c.x, 10) + "." + strconv.FormatInt(c.y, 10) + "." + strconv.FormatInt(c.z, 10)
}

func pairKey(p pair) string {
	return coordKey(p.a) + "-" + coordKey(p.b)
}

func eqCoord(c1 coord, c2 coord) bool {
	return c1.x == c2.x && c1.y == c2.y && c1.z == c2.z
}

func diff(c1 coord, c2 coord) float64 {
	return math.Sqrt(math.Pow(float64(c1.x-c2.x), 2) + math.Pow(float64(c1.y-c2.y), 2) + math.Pow(float64(c1.z-c2.z), 2))
}

func mergeCircuits(circuits [][]coord, i1 int, i2 int) [][]coord {
	if i1 == i2 {
		return circuits
	}
	// add all to i1 and delete i2
	circuits[i1] = append(circuits[i1], circuits[i2]...)
	circuits = append(circuits[:i2], circuits[i2+1:]...)
	return circuits
}

func findCircuit(circuits [][]coord, c1 coord) int {
	for i := 0; i < len(circuits); i++ {
		for _, c2 := range circuits[i] {
			if eqCoord(c1, c2) {
				return i
			}
		}
	}
	return -1
}

func main() {
	data, _ := os.ReadFile(os.Args[1])
	lines := strings.Split(string(data), "\n")
	/*
		n, err := strconv.ParseInt(os.Args[2], 10, 32)
		if err != nil {
			fmt.Printf("Need to provide the number of pairs to process as the second argument")
		}
	*/
	// first index is the circuit index, then an array (initially of length 1) of each box
	circuits := make([][]coord, 0)
	allCoords := make([]coord, 0)
	var coordString []string
	for i, l := range lines[:len(lines)-1] {
		circuits = append(circuits, make([]coord, 1))
		coordString = strings.Split(l, ",")
		v1, _ := strconv.ParseInt(coordString[0], 10, 32)
		v2, _ := strconv.ParseInt(coordString[1], 10, 32)
		v3, _ := strconv.ParseInt(coordString[2], 10, 32)
		newCoord := coord{x: int64(v1), y: int64(v2), z: int64(v3)}
		circuits[i][0] = newCoord
		allCoords = append(allCoords, newCoord)
	}
	// make a map of all possible pairs with their distance as a key
	distanceMap := make(map[float64]pair)
	distances := make([]float64, 0)
	var res float64
	for i, c1 := range allCoords {
		for _, c2 := range allCoords[i+1:] {
			res = diff(c1, c2)
			distances = append(distances, res)
			distanceMap[res] = pair{a: c1, b: c2}
			// fmt.Printf("[%d] Calculate pair %d.%d.%d - %d.%d.%d\t\tDistance: %f\n", i, c1.x, c1.y, c1.z, c2.x, c2.y, c2.z, res)
		}
	}
	// sort the list of all results and use to query map to find closest pairs
	fmt.Printf("\n\n")
	slices.Sort(distances)
	var i1 int // circuit index of the first and second coord
	var i2 int
	i := 0
	// for i := 0; i < int(n); i++ {
	for i < len(distances) && len(circuits) > 1 {
		toMerge := distanceMap[distances[i]]
		i1 = findCircuit(circuits, toMerge.a)
		i2 = findCircuit(circuits, toMerge.b)
		fmt.Printf("\n\nMerging circuit %d and %d\n", i1, i2)
		circuits = mergeCircuits(circuits, i1, i2)
		/*
			for _, c := range circuits {
				fmt.Println(c)
			}
		*/
		if len(circuits) == 1 {
			fmt.Printf("Last merge was %v and %v\n", toMerge.a, toMerge.b)
			fmt.Printf("Answer: %d\n", toMerge.a.x*toMerge.b.x)
		}
		i += 1
	}
	// now find the size of each circuit and order
	circuitSizes := make([]int, 0)
	for _, c := range circuits {
		circuitSizes = append(circuitSizes, len(c))
	}
	slices.Sort(circuitSizes)
	slices.Reverse(circuitSizes)
	fmt.Printf("Sizes: %v\n", circuitSizes)
	// fmt.Printf("Answer: %d\n", circuitSizes[0]*circuitSizes[1]*circuitSizes[2])
}
