package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type coord struct {
	x int64
	y int64
}

func eq(c1 coord, c2 coord) bool {
	return c1.x == c2.x && c1.y == c2.y
}

func rect(c1 coord, c2 coord) int64 {
	return (abs(c1.x-c2.x) + 1) * (abs(c1.y-c2.y) + 1)
}

func direction(c1 coord, c2 coord) byte {
	if c1.x == c2.x { // vertical
		if c1.y < c2.y {
			return byte('d')
		} else {
			return byte('u')
		}
	} else if c1.y == c2.y { // horizontal
		if c1.x < c2.x {
			return byte('r')
		} else {
			return byte('l')
		}
	} else { // no connection (also use this function for checking connectivity)
		return byte('n')
	}
}

// getting really ugly
func intersects(a1 coord, a2 coord, d1 byte, b1 coord, b2 coord, d2 byte) (bool, coord) {
	switch d1 {
	case 'u':
		// assume d2 == r
		if (a1.x >= b1.x && a1.x <= b2.x) && (b1.y < a1.y && b1.y > a2.y) {
			return true, coord{x: a1.x, y: b1.y}
		}
	case 'd':
		// assume d2 == l
		if (a1.x <= b1.x && a1.x >= b2.x) && (b1.y > a1.y && b1.y < a2.y) {
			return true, coord{x: a1.x, y: b1.y}
		}
	case 'l':
		// assume d2 == u
		if (a1.y <= b1.y && a1.y >= b2.y) && (b1.x < a1.x && b1.x > a2.x) {
			return true, coord{x: b1.x, y: a1.y}
		}
	case 'r':
		// assume d2 == d
		if (a1.y >= b1.y && a1.y <= b2.y) && (b1.x > a1.x && b1.x < a2.x) {
			return true, coord{x: b1.x, y: a1.y}
		}
	}
	return false, coord{x: -1, y: -1}
}

func abs(i int64) int64 {
	if i >= 0 {
		return i
	} else {
		return -i
	}
}

var blockedBy map[byte]byte = map[byte]byte{'r': 'd', 'd': 'l', 'l': 'u', 'u': 'r'}

// Order matters in the connections (need to assume every coord in cs is connected to the one after it)
// Returns if the connection is valid and at what point it was intersected (for casting) if false
func checkConnection(c1 coord, c2 coord, cs []coord) (bool, coord) {
	// check if two coords can connect
	dir := direction(c1, c2)
	if dir == 'n' {
		return false, coord{x: -1, y: -1}
	}
	// check if two coords are blocked by anything (to ensure connections remain internal)
	for i := 0; i < len(cs)-2; i++ {
		checkDir := direction(cs[i], cs[i+1])
		if !(eq(c1, cs[i]) || eq(c1, cs[i+1]) || eq(c2, cs[i]) || eq(c2, cs[i+1])) && blockedBy[dir] == checkDir {
			/*
				CHECK IF THE LINES FORMED BY c1->c2 and cs[i]->cs[i+1] INTERSECT
				Should intersect if c1 is on a 'blocking' line but not c2
					This might already be solved by only considering lines in a certain direction
			*/
			res, intersectCoord := intersects(c1, c2, dir, cs[i], cs[i+1], checkDir)
			if res {
				return false, intersectCoord
			}
		}
	}
	return true, coord{x: -1, y: -1}
}

// Input is the position of tiles
// Find the biggest rectangle bound by two tiles
// Can be bound on the same column / row ((7,3) and (2, 3) form a rectangle of area 50)

// P2: Adjacent red tiles are connected by green tiles
// Connect the next and previous red tile to eachother with green tiles (list wraps) and fill in any area
// Must find the largest rectangle that fits within this area ((9,5), (2,3) form a rectangle of area 24)
// This is getting hard now

// Solution idea: Can think of the perimeter as moving in a "clockwise direction" (think of each edge as an arrow rather than a line)
// Can use this idea to check whether a line is blocked by another line [0(n)]
//
//	A line moving right will be blocked by lines moving down (r, d)
//	A line moving down will be blocked by lines moving left	 (d, l)
//	A line moving left will be blocked by lines moving up	 (l, u)
//	A line moving up will be blocked by lines moving right	 (u, r)
//
// Then, "cast" from each corner to get 4 points that are candidate corners for the rectangle [0(n^2)]}
//
//	Don't need to do this, can generate "candidate corners" based just on pairs of potential red corners
//
// Check through each candidate to see if it can connect to other candidates and save connections [0(n^2)]
//
//	Save connections in a map where the candidate is the key
//	Check using the previously described method to check if the connecting line is blocked by another
//
// Check through all candidates, use the connection info to form all possible rectangles [0(n)]
// Calculate their area and select the highest, excluding solutions where opposite corners are not original coordinates [Part of above routine, 0(1)]
//
//	Notes the opposing corners for clarity
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
	tiles = append(tiles, tiles[0])
	tiles = append(append(make([]coord, 0), tiles[len(tiles)-2]), tiles...)

	// find "candidate" rectangles and check the relevant connectivity
	largest := int64(0)
	xc := []int64{0, 0}
	yc := []int64{0, 0}
	for i := 0; i < len(tiles)-1; i++ {
		for j := i; j < len(tiles)-1; j++ {
			// construct tl, tr, dl and br corners
			xc[0] = tiles[i].x
			xc[1] = tiles[j].x
			yc[0] = tiles[i].y
			yc[1] = tiles[j].y
			slices.Sort(xc)
			slices.Sort(yc)
			tl := coord{x: xc[0], y: yc[0]}
			tr := coord{x: xc[1], y: yc[0]}
			bl := coord{x: xc[0], y: yc[1]}
			br := coord{x: xc[1], y: yc[1]}
			// need to order these connections to be the same "direction" as the rest of the shape
			c1Con1, c11 := checkConnection(tl, tr, tiles)
			c1Con2, c12 := checkConnection(tr, br, tiles)
			c2Con1, c21 := checkConnection(br, bl, tiles)
			c2Con2, c22 := checkConnection(bl, tl, tiles)
			if c1Con1 && c1Con2 && c2Con1 && c2Con2 && rect(tiles[i], tiles[j]) > largest {
				largest = rect(tiles[i], tiles[j])
				fmt.Printf("New largest rect found of size [%d] with corners %v and %v\n", largest, tiles[i], tiles[j])
			} else if false {
				if !c1Con1 {
					fmt.Printf("Candidate with corners %v, %v, %v, %v blocked at %v\n", tl, tr, bl, br, c11)
				}
				if !c1Con2 {
					fmt.Printf("Candidate with corners %v, %v, %v, %v blocked at %v\n", tl, tr, bl, br, c12)
				}
				if !c2Con1 {
					fmt.Printf("Candidate with corners %v, %v, %v, %v blocked at %v\n", tl, tr, bl, br, c21)
				}
				if !c2Con2 {
					fmt.Printf("Candidate with corners %v, %v, %v, %v blocked at %v\n", tl, tr, bl, br, c22)
				}
			}
		}
	}
	fmt.Printf("Answer: %d", largest)

	/* PART ONE
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
	*/
}
