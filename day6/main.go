package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func evaluate(digits []byte) int64 {
	byteString := ""
	for i := 0; i < len(digits)-1; i++ {
		if digits[i] != byte(' ') {
			byteString += string(digits[i])
		}
	}
	val, _ := strconv.ParseInt(byteString, 10, 32)
	fmt.Printf("Evaluate: %s byteString as %d\n", byteString, val)
	return int64(val)
}

func main() {
	data, _ := os.ReadFile(os.Args[1])
	lines := strings.Split(string(data), "\n")
	numberString := lines[:len(lines)-1]
	operatorString := lines[len(lines)-2]
	// get operators
	operators := make([]string, 0)
	for k := len(operatorString) - 1; k >= 0; k-- {
		if operatorString[k] != byte(' ') {
			operators = append(operators, string(operatorString[k]))
		}
	}
	// get numbers
	i := len(numberString[0]) - 1 // i is the "horizontal" position
	found := false
	// formattedNumbers := make([]string, 0)
	formatted := make([]byte, len(numberString))
	operatorIndex := 0
	total := int64(0)
	var current int64
	if operators[0] == "*" {
		current = 1
	} else {
		current = 0
	}
	for i >= 0 {
		found = false
		for j := 0; j < len(numberString)-1; j++ {
			found = found || numberString[j][i] != ' '
		}
		if found {
			for j := 0; j < len(numberString)-1; j++ {
				formatted[j] = numberString[j][i]
			}
			if operators[operatorIndex] == "*" {
				current *= evaluate(formatted)
			} else {
				current += evaluate(formatted)
			}
		}
		if !found || i == 0 {
			operatorIndex += 1
			total += current
			if operatorIndex < len(operators) {
				if operators[operatorIndex] == "*" {
					current = 1
				} else {
					current = 0
				}
			}
		}
		i -= 1
	}
	fmt.Printf("Total: %d", total)
}
