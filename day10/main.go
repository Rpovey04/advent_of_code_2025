package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func transition(current string, button []int) string {
	currentByte := []byte(current)
	for _, idx := range button {
		if currentByte[idx] == '.' {
			currentByte[idx] = '#'
		} else {
			currentByte[idx] = '.'
		}
	}
	return string(currentByte)
}

func initialDistances(length int) map[string]int {
	states := make(map[string]int)
	for i := 0; i < int(math.Pow(2.0, float64(length))); i++ {
		newState := ""
		for j := 0; j < length; j++ {
			if int(math.Floor(float64(i)/math.Pow(2.0, float64(j))))%2 == 0 {
				newState += "."
			} else {
				newState += "#"
			}
		}
		states[newState] = math.MaxInt32
	}
	// fmt.Printf("%d states found\n", len(states))
	return states
}

func calculateShortest(desired string, buttons [][]int) int {
	initial := ""
	for i := 0; i < len(desired); i++ {
		initial += "."
	}
	dist := initialDistances(len(desired))
	dist[initial] = 0
	queue := []string{initial}
	// Dijkstras
	for len(queue) > 0 {
		current := queue[0]
		if len(queue) == 1 {
			queue = []string{}
		} else {
			queue = queue[1:]
		}
		for _, button := range buttons {
			next := transition(current, button)
			// fmt.Printf("%s->%s|\tdist[next]=%d,dist[current]+1=%d\n", current, next, dist[next], dist[current]+1)
			if dist[next] > dist[current]+1 {
				dist[next] = dist[current] + 1
				queue = append(queue, next)
			}
		}
	}
	return dist[desired]
}

// PART ONE
// Idea: Treat every state (configuration of lights) as a node and the buttons as actions / edges
// Then use Dijkstra's algorithm to find shortest path from initial state (all off) to desired state
// Possible since the number of states is 2^n (where n is the number of light).
// Although this is exponential, desired states do not have more than 10 characters (max states = 1024)
// This is tractable with Dijkstra's algorithm, which has complexity 0(n^2), however would explode with largest states
func main() {
	data, _ := os.ReadFile(os.Args[1])
	lines := strings.Split(string(data), "\n")
	lines = lines[:len(lines)-1]

	total := 0
	initialDistances(10)
	for _, l := range lines {
		/*
			PARSING. No need to save the parsed data since the problem can be solved one line at a time
		*/
		objects := strings.Split(l, " ")
		desiredStr := objects[0]
		// joltageStr := objects[len(objects)-1]
		buttonsStr := objects[1 : len(objects)-1]

		desired := desiredStr[1 : len(desiredStr)-1]
		// joltage here
		buttons := make([][]int, 0)
		for _, button := range buttonsStr {
			lightsStr := strings.Split(button[1:len(button)-1], ",")
			lights := make([]int, 0)
			for _, light := range lightsStr {
				v, _ := strconv.ParseInt(light, 10, 32)
				lights = append(lights, int(v))
			}
			buttons = append(buttons, lights)
		}
		// fmt.Println(desired)
		// fmt.Printf("[%d] %v\n", i, buttons)
		total += calculateShortest(desired, buttons)
	}
	fmt.Printf("Answer: %d\n", total)
}
