package main

import (
	"advent-of-code/shared"
	"fmt"
	"strconv"
	"strings"
)

type CaveContent string

const (
	Air         CaveContent = "."
	FallingRock             = "@"
	StillRock               = "#"
)

const rockCount int = 5
const width int = 7
const maxRockHeight int = 4
const airRowsOnSpawn int = 3

var rock1 [][]CaveContent = [][]CaveContent{
	{Air, Air, FallingRock, FallingRock, FallingRock, FallingRock, Air}}

var rock2 [][]CaveContent = [][]CaveContent{
	{Air, Air, Air, FallingRock, Air, Air, Air},
	{Air, Air, FallingRock, FallingRock, FallingRock, Air, Air},
	{Air, Air, Air, FallingRock, Air, Air, Air}}

var rock3 [][]CaveContent = [][]CaveContent{
	{Air, Air, FallingRock, FallingRock, FallingRock, Air, Air},
	{Air, Air, Air, Air, FallingRock, Air, Air},
	{Air, Air, Air, Air, FallingRock, Air, Air}}

var rock4 [][]CaveContent = [][]CaveContent{
	{Air, Air, FallingRock, Air, Air, Air, Air},
	{Air, Air, FallingRock, Air, Air, Air, Air},
	{Air, Air, FallingRock, Air, Air, Air, Air},
	{Air, Air, FallingRock, Air, Air, Air, Air}}

var rock5 [][]CaveContent = [][]CaveContent{
	{Air, Air, FallingRock, FallingRock, Air, Air, Air},
	{Air, Air, FallingRock, FallingRock, Air, Air, Air}}

var jets []string
var cave [][]CaveContent
var rocks [][][]CaveContent
var stillRockCount int = 0
var currentRockIndex int = 0
var currentJetIndex int = 0
var currentRockTopPosition int
var currentTurn int = 1
var heightAtFirstRepetition int = 0
var stillRockCountAtFirstRepetition int = 0
var stillRockCountPerRepetition int = 0
var heightPerRepetition int = 0
var unrepeatedMovesInTheEnd int = -1
var unrepeatedHeightInTheEnd int = -1
var stillRockTarget int = 2022
var rock1PosToJetIndex map[int]map[int]bool = make(map[int]map[int]bool)

func main() {
	filePuzzle1 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle1, parseJetPattern)

	initialize()
	heightPuzzle1 := playTetris()
	fmt.Println("Rock tower height after " + strconv.Itoa(stillRockTarget) + " rocks have fallen: " + strconv.Itoa(heightPuzzle1))

	initialize()
	stillRockTarget = 1000000000000
	heightPuzzle2 := playTetris()
	fmt.Println("Rock tower height after " + strconv.Itoa(stillRockTarget) + " rocks have fallen: " + strconv.Itoa(heightPuzzle2))
}

func parseJetPattern(line string) {
	jets = strings.Split(line, "")
}

func initialize() {
	resetRepetitionMarkers()

	stillRockCount = 0
	currentRockIndex = 0
	currentJetIndex = 0
	currentTurn = 1
	heightAtFirstRepetition = 0
	stillRockCountAtFirstRepetition = 0
	stillRockCountPerRepetition = 0
	heightPerRepetition = 0
	unrepeatedMovesInTheEnd = -1
	unrepeatedHeightInTheEnd = -1

	cave = make([][]CaveContent, 0)

	rocks = make([][][]CaveContent, rockCount)
	initSingleRock(0, rock1)
	initSingleRock(1, rock2)
	initSingleRock(2, rock3)
	initSingleRock(3, rock4)
	initSingleRock(4, rock5)
}

func resetRepetitionMarkers() {
	for i := 0; i < width; i++ {
		rock1PosToJetIndex[i] = make(map[int]bool)
	}

}

func initSingleRock(index int, rock [][]CaveContent) {
	rocks[index] = make([][]CaveContent, len(rock))

	for i := 0; i < len(rock); i++ {
		rocks[index][i] = make([]CaveContent, width)

		for j := 0; j < width; j++ {
			rocks[index][i][j] = rock[i][j]
		}
	}
}

func playTetris() int {
	needToSpawnNewRock := true

	for stillRockCount < stillRockTarget {
		if needToSpawnNewRock {
			spawnNewRock()
			needToSpawnNewRock = false
		}

		if isJetMovementPossible() {
			moveRockByJet()
		}

		if isDownwardMovementPossible() {
			moveRockDownwards()
		} else {
			needToSpawnNewRock = true
			rockPos := turnRockStill()
			checkRepetitionOfFirstRock(rockPos)
			if isEndOfSimulationReached() {
				return predictHeightOfRockTower()
			}
		}

		currentJetIndex++
		currentJetIndex = currentJetIndex % len(jets)
		currentTurn++
	}

	//printCave()
	height := calculateHeightOfRockTower()
	return height
}

func spawnNewRock() {
	removeTopLayersOfAir()
	appendNewLayersOfAir()

	cave = append(cave, createNewRock()...)
	currentRockTopPosition = len(cave) - 1
}

func removeTopLayersOfAir() {
	firstAirRowIndex := -1

	for i := 0; i < len(cave) && firstAirRowIndex == -1; i++ {
		airCount := 0
		for j := 0; j < width; j++ {
			if cave[i][j] == Air {
				airCount++
			}
		}
		if airCount == width {
			firstAirRowIndex = i
		}
	}

	if firstAirRowIndex != -1 {
		cave = cave[:firstAirRowIndex]
	}
}

func appendNewLayersOfAir() {
	for i := 0; i < airRowsOnSpawn; i++ {
		cave = append(cave, []CaveContent{Air, Air, Air, Air, Air, Air, Air})
	}
}

func createNewRock() [][]CaveContent {
	var newRock [][]CaveContent = make([][]CaveContent, len(rocks[currentRockIndex]))

	for i := 0; i < len(newRock); i++ {
		newRock[i] = make([]CaveContent, width)
		for j := 0; j < width; j++ {
			newRock[i][j] = rocks[currentRockIndex][i][j]
		}
	}

	return newRock
}

func moveRockByJet() {
	currentJet := jets[currentJetIndex]
	for i := currentRockTopPosition; i >= currentRockTopPosition-maxRockHeight && i >= 0; i-- {
		if currentJet == ">" {
			moveRockToRight(i)
		} else {
			moveRockToLeft(i)
		}
	}
}

func isJetMovementPossible() bool {
	currentJet := jets[currentJetIndex]
	return (currentJet == ">" && isMoveToRightPossible()) || (currentJet == "<" && isMoveToLeftPossible())
}

func isMoveToRightPossible() bool {
	for i := currentRockTopPosition; i >= currentRockTopPosition-maxRockHeight && i >= 0; i-- {
		for j := len(cave[i]) - 1; j >= 0; j-- {
			if cave[i][j] == FallingRock && (j == width-1 || cave[i][j+1] == StillRock) {
				return false
			}
		}
	}
	return true
}

func isMoveToLeftPossible() bool {
	for i := currentRockTopPosition; i >= currentRockTopPosition-maxRockHeight && i >= 0; i-- {
		for j := 0; j < len(cave[i]); j++ {
			if cave[i][j] == FallingRock && (j == 0 || cave[i][j-1] == StillRock) {
				return false
			}
		}
	}
	return true
}

func moveRockToRight(row int) {
	for j := len(cave[row]) - 2; j >= 0; j-- {
		if cave[row][j] == FallingRock && cave[row][j+1] == Air {
			cave[row][j+1] = FallingRock
			cave[row][j] = Air
		}
	}
}

func moveRockToLeft(row int) {
	for j := 1; j < len(cave[row]); j++ {
		if cave[row][j] == FallingRock && cave[row][j-1] == Air {
			cave[row][j-1] = FallingRock
			cave[row][j] = Air
		}
	}
}

func isDownwardMovementPossible() bool {
	for i := currentRockTopPosition; i >= currentRockTopPosition-maxRockHeight && i >= 0; i-- {
		for j := 0; j < width; j++ {
			if cave[i][j] == FallingRock && (i == 0 || cave[i-1][j] == StillRock) {
				return false
			}
		}
	}
	return true
}

func moveRockDownwards() {
	startPosition := currentRockTopPosition - maxRockHeight

	if startPosition <= 0 {
		startPosition = 1
	}

	for i := startPosition; i <= currentRockTopPosition; i++ {
		for j := 0; j < width; j++ {
			if cave[i][j] == FallingRock {
				cave[i][j] = Air
				cave[i-1][j] = FallingRock
			}
		}
	}
	currentRockTopPosition--
}

func turnRockStill() int {
	rockPos := -1

	for i := currentRockTopPosition; i >= currentRockTopPosition-maxRockHeight && i >= 0; i-- {
		for j := 0; j < width; j++ {
			if cave[i][j] == FallingRock {
				cave[i][j] = StillRock
				rockPos = j
			}
		}
	}

	currentRockIndex++
	currentRockIndex = currentRockIndex % len(rocks)
	stillRockCount++
	return rockPos
}

func checkRepetitionOfFirstRock(rockPos int) {
	if stillRockCountPerRepetition == 0 && currentRockIndex == 0 {
		if rock1PosToJetIndex[rockPos][currentJetIndex] {
			if stillRockCountAtFirstRepetition == 0 {
				stillRockCountAtFirstRepetition = stillRockCount
				heightAtFirstRepetition = calculateHeightOfRockTower()
			} else {
				stillRockCountPerRepetition = stillRockCount - stillRockCountAtFirstRepetition
				heightPerRepetition = calculateHeightOfRockTower() - heightAtFirstRepetition
				unrepeatedMovesInTheEnd = (stillRockTarget - stillRockCountAtFirstRepetition) % stillRockCountPerRepetition
			}
			resetRepetitionMarkers()
		}
		rock1PosToJetIndex[rockPos][currentJetIndex] = true
	}
}

func isEndOfSimulationReached() bool {
	if unrepeatedMovesInTheEnd == 0 {
		unrepeatedHeightInTheEnd = calculateHeightOfRockTower() - heightAtFirstRepetition - heightPerRepetition
		return true
	} else if unrepeatedMovesInTheEnd > 0 {
		unrepeatedMovesInTheEnd--
	}
	return false
}

func calculateHeightOfRockTower() int {
	for i := len(cave) - 1; i >= 0; i-- {
		for j := 0; j < width; j++ {
			if cave[i][j] == StillRock {
				return i + 1
			}
		}
	}
	return 0
}

func predictHeightOfRockTower() int {
	totalRepetitions := (stillRockTarget - stillRockCountAtFirstRepetition) / stillRockCountPerRepetition
	totalHeight := heightAtFirstRepetition + unrepeatedHeightInTheEnd + (totalRepetitions * heightPerRepetition)
	return totalHeight
}

func printCave() {
	for i := len(cave) - 1; i >= 0; i-- {
		line := ""
		for j := 0; j < width; j++ {
			line += string(cave[i][j])
		}
		fmt.Println(line)
	}
}
