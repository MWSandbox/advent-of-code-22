package main

import (
	"advent-of-code/shared"
	"fmt"
	"strconv"
	"strings"
)

var fullyContainedRanges int = 0
var partlyContainedRanges int = 0

func main() {
	filePuzzle1 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle1, calculateFullOverlaps)

	filePuzzle2 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle2, calculatePartlyOverlaps)

	fmt.Println("Fully contained ranges: " + strconv.Itoa(fullyContainedRanges))
	fmt.Println("Partly contained ranges: " + strconv.Itoa(partlyContainedRanges))
}

func calculateFullOverlaps(line string) {
	firstRange, secondRange := extractRanges(line)

	if isRangeFullyContained(firstRange, secondRange) || isRangeFullyContained(secondRange, firstRange) {
		fullyContainedRanges++
	}
}

func calculatePartlyOverlaps(line string) {
	firstRange, secondRange := extractRanges(line)

	if isRangePartlyContained(firstRange, secondRange) || isRangePartlyContained(secondRange, firstRange) {
		partlyContainedRanges++
	}
}

func extractRanges(text string) ([2]int, [2]int) {
	assignments := strings.Split(text, ",")
	firstRangeAsString := strings.Split(assignments[0], "-")
	secondRangeAsString := strings.Split(assignments[1], "-")
	firstRange := convertStringsToInts(firstRangeAsString)
	secondRange := convertStringsToInts(secondRangeAsString)
	return firstRange, secondRange
}

func convertStringsToInts(input []string) [2]int {
	var output [2]int
	
	for i := 0; i < len(input); i ++ {
		value, err := strconv.Atoi(input[i])

		if err != nil {
			fmt.Println("Couldn't convert string to int")
		}
		output[i] = value		
	}

	return output
}

func isRangeFullyContained(containerRange [2]int, enclosedRange [2]int) bool {
	return containerRange[0] <= enclosedRange[0] && containerRange[1] >= enclosedRange[1]
}

func isRangePartlyContained(containerRange [2]int, enclosedRange [2]int) bool {
	return containerRange[0] <= enclosedRange[0] && containerRange[1] >= enclosedRange[0]
}