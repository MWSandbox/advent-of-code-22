package main

import (
	"advent-of-code/shared"
	"fmt"
	"strconv"
	"strings"
)

var moveToMoveToMatchScore map[string]map[string]int
var moveToMoveScore map[string]int
var strategyToMatchResult map[string]int
var totalScorePuzzle1 int = 0
var totalScorePuzzle2 int = 0

func main() {
	intMaps()
	filePuzzle1 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle1, calculatePuzzle1)
	filePuzzle2 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle2, calculatePuzzle2)

	fmt.Println("Total Score Puzzle 1: " + strconv.Itoa(totalScorePuzzle1))
	fmt.Println("Total Score Puzzle 2: " + strconv.Itoa(totalScorePuzzle2))
}

func intMaps() {
	moveToMoveToMatchScore = make(map[string]map[string]int)
	initMatchResults("A", 3, 6, 0)
	initMatchResults("B", 0, 3, 6)
	initMatchResults("C", 6, 0, 3)

	moveToMoveScore = make(map[string]int)
	moveToMoveScore["X"] = 1
	moveToMoveScore["Y"] = 2
	moveToMoveScore["Z"] = 3

	strategyToMatchResult = make(map[string]int)
	strategyToMatchResult["X"] = 0
	strategyToMatchResult["Y"] = 3
	strategyToMatchResult["Z"] = 6
}

func initMatchResults(opponentMove string, valueX int, valueY int, valueZ int) {
	moveToMoveToMatchScore[opponentMove] = make(map[string]int)
	moveToMoveToMatchScore[opponentMove]["X"] = valueX
	moveToMoveToMatchScore[opponentMove]["Y"] = valueY
	moveToMoveToMatchScore[opponentMove]["Z"] = valueZ
}

func calculatePuzzle1(line string) {
	moves := strings.Fields(line)
	opponentMove := moves[0]
	ownMove := moves[1]
	roundScore := getMoveScore(ownMove) + getMatchScore(opponentMove, ownMove)
	totalScorePuzzle1 += roundScore
}

func calculatePuzzle2(line string) {
	fields := strings.Fields(line)
	opponentMove := fields[0]
	expectedResult := fields[1]
	expectedMatchScore := strategyToMatchResult[expectedResult]
	var ownMove string

	if moveToMoveToMatchScore[opponentMove]["X"] == expectedMatchScore {
		ownMove = "X"
	} else if moveToMoveToMatchScore[opponentMove]["Y"] == expectedMatchScore {
		ownMove = "Y"
	} else {
		ownMove = "Z"
	}
	totalScorePuzzle2 += getMoveScore(ownMove) + getMatchScore(opponentMove, ownMove)
}

func getMoveScore(ownMove string) int {
	return moveToMoveScore[ownMove]
}

func getMatchScore(opponentMove string, ownMove string) int {
	return moveToMoveToMatchScore[opponentMove][ownMove]
}
