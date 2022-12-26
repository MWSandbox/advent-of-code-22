package main

import (
	"advent-of-code/shared"
	"fmt"
	"strconv"
	"strings"
)

type SimulationTerminationFunction func(int, int) bool

type CaveMaterial string

const (
	Air         CaveMaterial = "."
	FallingSand              = "+"
	StillSand                = "o"
	Rock                     = "#"
)

const minX int = 200
const maxX int = 800
const minY int = 0
const maxY int = 200
const sandStartPosX int = 500
const barrierYPos int = 160

var cave [][]CaveMaterial = make([][]CaveMaterial, maxX-minX)

func main() {
	fillCaveWithAir()
	filePuzzle1 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle1, drawRocks)
	stillSandCount1 := simulateSand(terminateSimulationOnceAbyssIsReached)
	printCave()

	cave = make([][]CaveMaterial, maxX-minX)
	fillCaveWithAir()
	filePuzzle2 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle2, drawRocks)
	drawBarrier()
	stillSandCount2 := simulateSand(terminateSimulationOnceSourceIsStill)
	printCave()

	fmt.Println("Count of still sand for puzzle 1: " + strconv.Itoa(stillSandCount1))
	fmt.Println("Count of still sand for puzzle 2: " + strconv.Itoa(stillSandCount2))
}

func fillCaveWithAir() {
	for x := 0; x < maxX-minX; x++ {
		cave[x] = make([]CaveMaterial, maxY-minY)
		for y := 0; y < maxY-minY; y++ {
			cave[x][y] = Air
		}
	}
}

func drawRocks(line string) {
	coordinates := strings.Split(line, " -> ")

	for i := 0; i < len(coordinates)-1; i++ {
		startX, startY := extractPosition(coordinates[i])
		endX, endY := extractPosition(coordinates[i+1])

		if startX == endX {
			drawVerticalRockLine(startX, startY, endY)
		} else {
			drawHorizontalRockLine(startY, startX, endX)
		}
	}
}

func extractPosition(coordinates string) (int, int) {
	positions := strings.Split(coordinates, ",")

	positionX, errX := strconv.Atoi(positions[0])
	positionY, errY := strconv.Atoi(positions[1])

	if errX != nil || errY != nil {
		fmt.Println("Couldn't convert coordinates to int")
	}

	return positionX, positionY
}

func drawVerticalRockLine(posX int, startY int, endY int) {
	if endY < startY {
		startY, endY = endY, startY
	}

	posX = posX - minX
	startY = startY - minY
	endY = endY - minY

	for y := startY; y <= endY; y++ {
		cave[posX][y] = Rock
	}
}

func drawHorizontalRockLine(posY int, startX int, endX int) {
	if endX < startX {
		startX, endX = endX, startX
	}

	posY = posY - minY
	startX = startX - minX
	endX = endX - minX

	for x := startX; x <= endX; x++ {
		cave[x][posY] = Rock
	}
}

func simulateSand(isTerminationReached SimulationTerminationFunction) int {
	stillSandCount := 0
	currentSandXPos := sandStartPosX - minX
	currentSandYPos := 0
	cave[currentSandXPos][currentSandYPos] = FallingSand

	for !isTerminationReached(currentSandXPos, currentSandYPos) {
		if cave[currentSandXPos][currentSandYPos+1] == Air {
			moveFallinSandTo(currentSandXPos, currentSandYPos, currentSandXPos, currentSandYPos+1)
			currentSandYPos++
		} else if cave[currentSandXPos-1][currentSandYPos+1] == Air {
			moveFallinSandTo(currentSandXPos, currentSandYPos, currentSandXPos-1, currentSandYPos+1)
			currentSandXPos--
			currentSandYPos++
		} else if cave[currentSandXPos+1][currentSandYPos+1] == Air {
			moveFallinSandTo(currentSandXPos, currentSandYPos, currentSandXPos+1, currentSandYPos+1)
			currentSandXPos++
			currentSandYPos++
		} else {
			cave[currentSandXPos][currentSandYPos] = StillSand
			currentSandXPos = sandStartPosX - minX
			currentSandYPos = 0
			if cave[currentSandXPos][currentSandYPos] == Air {
				cave[currentSandXPos][currentSandYPos] = FallingSand
			}
			stillSandCount++
		}
	}
	return stillSandCount
}

func moveFallinSandTo(startPosX int, startPosY int, endPosX int, endPosY int) {
	cave[startPosX][startPosY] = Air
	cave[endPosX][endPosY] = FallingSand
}

func printCave() {
	xHeader1 := "    "
	xHeader2 := "    "
	xHeader3 := "    "

	for x := minX; x < maxX; x++ {
		firstDigit := x / 100
		secondDigit := (x % 100) / 10
		thirdDigit := x % 10

		xHeader1 += strconv.Itoa(firstDigit)
		xHeader2 += strconv.Itoa(secondDigit)
		xHeader3 += strconv.Itoa(thirdDigit)
	}

	fmt.Println(xHeader1)
	fmt.Println(xHeader2)
	fmt.Println(xHeader3)

	for y := 0; y < maxY-minY; y++ {
		padding := " "

		if y < 10 {
			padding += " "
		}

		if y < 100 {
			padding += " "
		}

		var printedLine string = strconv.Itoa(y) + padding
		for x := 0; x < maxX-minX; x++ {
			printedLine += string(cave[x][y])
		}
		fmt.Println(printedLine)
	}
}

func drawBarrier() {
	for x := 0; x < maxX-minX; x++ {
		cave[x][barrierYPos] = Rock
	}
}

func terminateSimulationOnceAbyssIsReached(currentSandXPos int, currentSandYPos int) bool {
	return currentSandYPos == maxY-1
}

func terminateSimulationOnceSourceIsStill(currentSandXPos int, currentSandYPos int) bool {
	return cave[sandStartPosX-minX][0] == StillSand
}
