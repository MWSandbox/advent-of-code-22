package main

import (
	"advent-of-code/shared"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const stackCount int = 9
var stackVsInstructionSeparator string = ""
var isInstructionParsed bool = false
var stacksOfCrates [stackCount][]string
var isFirstPuzzle bool

func main() {
	fmt.Println("Puzzle 1")
	filePuzzle1 := shared.OpenFile("./input.txt")
	isFirstPuzzle = true
	shared.ReadFileLineByLine(filePuzzle1, calculateTopCrates)
	printTopCrates()

	fmt.Println("Puzzle 2")
	isInstructionParsed = false
	stacksOfCrates = [stackCount][]string{}
	filePuzzle2 := shared.OpenFile("./input.txt")
	isFirstPuzzle = false
	shared.ReadFileLineByLine(filePuzzle2, calculateTopCrates)
	printTopCrates()
}

func calculateTopCrates(line string) {

	if line == stackVsInstructionSeparator {
		isInstructionParsed = true
		reverseStacks()
	} else if isInstructionParsed {
		if isFirstPuzzle {
			executeInstructionFirstPuzzle(line)
		} else {
			executeInstructionSecondPuzzle(line)
		}
	} else {
		parseStackInput(line)
	}

}

func parseStackInput(line string) {
	var separatorSpaces int = 0

	for i := 0; i < stackCount; i++ {
		currentStack := line[i*3+1+separatorSpaces : i*3+2+separatorSpaces]
		separatorSpaces++
		isNumeric := regexp.MustCompile(`\d`).MatchString(currentStack)

		if currentStack != " " && !isNumeric {
			stacksOfCrates[i] = append(stacksOfCrates[i], currentStack)
		}
	}
}

func reverseStacks() {
	for stackIndex := 0; stackIndex < stackCount; stackIndex++ {
		reverseStack(stacksOfCrates[stackIndex])
	}
}

func reverseStack(stack []string) {
	for i, j := 0, len(stack)-1; i < j; i, j = i+1, j-1 {
		stack[i], stack[j] = stack[j], stack[i]
	}
}

func executeInstructionFirstPuzzle(line string) {
	crateCount, startStackIndex, endStackIndex := parseInstruction(line)

	for i := 0; i < crateCount; i++ {
		topCrate := removeCrate(startStackIndex)
		stacksOfCrates[endStackIndex-1] = append(stacksOfCrates[endStackIndex-1], topCrate)
	}
}

func executeInstructionSecondPuzzle(line string) {
	var cratesToMove []string
	crateCount, startStackIndex, endStackIndex := parseInstruction(line)

	for i := 0; i < crateCount; i++ {
		topCrate := removeCrate(startStackIndex)
		cratesToMove = append(cratesToMove, topCrate)
	}

	reverseStack(cratesToMove)

	for i := 0; i < len(cratesToMove); i++ {
		topCrate := cratesToMove[i]
		stacksOfCrates[endStackIndex-1] = append(stacksOfCrates[endStackIndex-1], topCrate)
	}
}

func parseInstruction(line string) (crateCount int, startStackIndex int, endStackIndex int) {
	split := strings.Split(line, " ")
	crateCount, err1 := strconv.Atoi(split[1])
	startStackIndex, err2 := strconv.Atoi(split[3])
	endStackIndex, err3 := strconv.Atoi(split[5])

	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Println("Couldn't parse instruction: " + line)
	}

	return crateCount, startStackIndex, endStackIndex
}

func removeCrate(startStackIndex int) string {
	topCrate, topCrateIndex := getTopElementAndIndex(stacksOfCrates[startStackIndex-1])
	stacksOfCrates[startStackIndex-1] = stacksOfCrates[startStackIndex-1][:topCrateIndex]
	return topCrate
}

func printTopCrates() {
	var topCrates string = ""
	var crateCounts string = ""

	for i := 0; i < stackCount; i++ {
		topCrate, topCrateIndex := getTopElementAndIndex(stacksOfCrates[i])

		topCrates += topCrate
		crateCounts += strconv.Itoa(topCrateIndex + 1)
	}

	fmt.Println("Top crates: " + topCrates)
	fmt.Println("Crate counts: " + crateCounts)
}

func getTopElementAndIndex(stack []string) (crate string, index int) {
	if len(stack) > 0 {
		topCrateIndex := len(stack) - 1
		topCrate := stack[topCrateIndex]
		return topCrate, topCrateIndex
	}
	return " ", 0
}
