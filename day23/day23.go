package main

import (
	"advent-of-code/shared"
	"fmt"
	"strconv"
)

type CaveContent string
type MovementFunction func(int, int) (int, int)

const (
	Covered CaveContent = "#"
	Empty               = "."
)

type Elf struct {
	x        int
	y        int
	isMoving bool
	nextX    int
	nextY    int
}

var maxRounds int = 10
var yParserPos int = 0
var cave [][]CaveContent = make([][]CaveContent, 0)
var elves []*Elf = make([]*Elf, 0)
var movementFunctions []MovementFunction = []MovementFunction{isMovingNorth, isMovingSouth, isMovingWest, isMovingEast}
var roundCount int = 0

func main() {
	filePuzzle := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle, parseElvesPositions)
	extendCave()

	solvePuzzle1()
	solvePuzzle2()
}

func solvePuzzle1() {
	simulateElvesMovement()
	emptyTiles := countEmptyGroundTiles()
	printCave()

	fmt.Println("Empty ground tiles at round " + strconv.Itoa(roundCount) + ": " + strconv.Itoa(emptyTiles))
}

func solvePuzzle2() {
	maxRounds = 999999
	simulateElvesMovement()
	printCave()

	fmt.Println("First round without elves moving: " + strconv.Itoa(roundCount))
}

func parseElvesPositions(line string) {
	cave = append(cave, make([]CaveContent, len(line)))

	for x := 0; x < len(line); x++ {
		content := CaveContent(line[x : x+1])
		cave[yParserPos][x] = content

		if content == Covered {
			elves = append(elves, &Elf{x: x, y: yParserPos})
		}
	}

	yParserPos++
}

func extendCave() {
	extendTop()
	extendBottom()
	extendLeft()
	extendRight()
}

func extendTop() {
	isTopCovered := false
	for x := 0; x < len(cave[0]); x++ {
		isTopCovered = isTopCovered || cave[0][x] == Covered
	}

	if isTopCovered {
		newTop := make([][]CaveContent, 1)
		newTop[0] = createEmptyRow()
		cave = append(newTop, cave...)

		for i := 0; i < len(elves); i++ {
			elves[i].y++
		}
	}
}

func extendBottom() {
	isBottomCovered := false
	for x := 0; x < len(cave[len(cave)-1]); x++ {
		isBottomCovered = isBottomCovered || cave[len(cave)-1][x] == Covered
	}

	if isBottomCovered {
		cave = append(cave, createEmptyRow())
	}
}

func extendLeft() {
	isLeftCovered := false
	for y := 0; y < len(cave); y++ {
		isLeftCovered = isLeftCovered || cave[y][0] == Covered
	}

	if isLeftCovered {
		for i := 0; i < len(cave); i++ {
			newLeft := make([]CaveContent, 1)
			newLeft[0] = Empty
			cave[i] = append(newLeft, cave[i]...)
		}
		for i := 0; i < len(elves); i++ {
			elves[i].x++
		}
	}
}

func extendRight() {
	isRightCovered := false
	for y := 0; y < len(cave); y++ {
		isRightCovered = isRightCovered || cave[y][len(cave[0])-1] == Covered
	}

	if isRightCovered {
		for i := 0; i < len(cave); i++ {
			cave[i] = append(cave[i], Empty)
		}
	}
}

func createEmptyRow() []CaveContent {
	var emptyRow []CaveContent

	for i := 0; i < len(cave[0]); i++ {
		emptyRow = append(emptyRow, Empty)
	}

	return emptyRow
}

func simulateElvesMovement() {
	isMovePossible := true

	for isMovePossible && roundCount < maxRounds {
		checkIfElvesAreMoving()
		nextPositionToElfCount := calculateNextElfMoves()
		isMovePossible = len(nextPositionToElfCount) > 0
		moveToEmptyTiles(nextPositionToElfCount)
		reorderMovementFunctions()
		extendCave()
		roundCount++
	}
}

func checkIfElvesAreMoving() {
	for i := 0; i < len(elves); i++ {
		emptyCount := 0
		elf := elves[i]

		for x := elf.x - 1; x <= elf.x+1; x++ {
			for y := elf.y - 1; y <= elf.y+1; y++ {
				if (x != elf.x || y != elf.y) && cave[y][x] == Empty {
					emptyCount++
				}
			}
		}
		elf.isMoving = emptyCount != 8
	}
}

func calculateNextElfMoves() map[int]map[int]int {
	nextPositionToElfCount := make(map[int]map[int]int)

	for i := 0; i < len(elves); i++ {
		if elves[i].isMoving {
			newX, newY := proposeNewPosition(elves[i])
			elves[i].nextX = newX
			elves[i].nextY = newY

			if nextPositionToElfCount[newX] == nil {
				nextPositionToElfCount[newX] = make(map[int]int)
			}
			nextPositionToElfCount[newX][newY]++
		}
	}

	return nextPositionToElfCount
}

func proposeNewPosition(elf *Elf) (int, int) {
	for i := 0; i < len(movementFunctions); i++ {
		newX, newY := movementFunctions[i](elf.x, elf.y)
		if newX != elf.x || newY != elf.y {
			return newX, newY
		}
	}
	return elf.x, elf.y
}

func isMovingNorth(x int, y int) (int, int) {
	return checkXPositions(x, y, y-1)
}

func isMovingSouth(x int, y int) (int, int) {
	return checkXPositions(x, y, y+1)
}

func isMovingWest(x int, y int) (int, int) {
	return checkYPositions(x, y, x-1)
}

func isMovingEast(x int, y int) (int, int) {
	return checkYPositions(x, y, x+1)
}

func checkXPositions(currentX int, currentY int, targetY int) (int, int) {
	emptyCount := 0
	for newX := currentX - 1; newX <= currentX+1; newX++ {
		if cave[targetY][newX] == Empty {
			emptyCount++
		}
	}

	if emptyCount == 3 {
		return currentX, targetY
	}
	return currentX, currentY
}

func checkYPositions(currentX int, currentY int, targetX int) (int, int) {
	emptyCount := 0
	for newY := currentY - 1; newY <= currentY+1; newY++ {
		if cave[newY][targetX] == Empty {
			emptyCount++
		}
	}

	if emptyCount == 3 {
		return targetX, currentY
	}
	return currentX, currentY
}

func moveToEmptyTiles(nextPositionToElfCount map[int]map[int]int) {
	for i := 0; i < len(elves); i++ {
		elf := elves[i]
		if elf.isMoving && nextPositionToElfCount[elf.nextX][elf.nextY] == 1 {
			cave[elf.y][elf.x] = Empty
			elf.x = elf.nextX
			elf.y = elf.nextY
			cave[elf.y][elf.x] = Covered
		}
	}
}

func reorderMovementFunctions() {
	firstFunction := movementFunctions[0]
	movementFunctions = movementFunctions[1:]
	movementFunctions = append(movementFunctions, firstFunction)
}

func countEmptyGroundTiles() int {
	emptyCount := 0
	minX, maxX, minY, maxY := getSmallestRectangleContainingAllElves()

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if cave[y][x] == Empty {
				emptyCount++
			}
		}
	}

	return emptyCount
}

func getSmallestRectangleContainingAllElves() (int, int, int, int) {
	var maxY, maxX int = 0, 0
	var minY, minX int = len(cave), len(cave)

	for i := 0; i < len(elves); i++ {
		elf := elves[i]

		if elf.x > maxX {
			maxX = elf.x
		}
		if elf.x < minX {
			minX = elf.x
		}
		if elf.y > maxY {
			maxY = elf.y
		}
		if elf.y < minY {
			minY = elf.y
		}
	}

	return minX, maxX, minY, maxY
}

func printCave() {
	shared.PrintXHeader(0, len(cave[0]))

	for y := 0; y < len(cave); y++ {
		var printedLine string = strconv.Itoa((y)) + shared.GetYPadding(y)
		for x := 0; x < len(cave[y]); x++ {
			printedLine += string(cave[y][x])
		}
		fmt.Println(printedLine)
	}
}
