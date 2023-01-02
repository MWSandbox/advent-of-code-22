package main

import (
	"advent-of-code/shared"
	"fmt"
	"strconv"
	"strings"
)

type TransitionFunction func(int, int, int)
type IsEndOfBoardFunction func() bool
type StepFunction func(bool, int, int)
type CubeTransitionFunction func(int, int, int)
type FlatTransitionFunction func(int)

type Position struct {
	x int
	y int
}

type WalkResult bool

const (
	Completed   WalkResult = true
	Interrupted            = false
)

type Direction string

const (
	Up    Direction = "^"
	Right           = ">"
	Down            = "v"
	Left            = "<"
)

type BoardContent string

const (
	NoContent   BoardContent = " "
	Empty                    = "."
	Wall                     = "#"
	RightFacing              = ">"
	LeftFacing               = "<"
	UpFacing                 = "^"
	DownFacing               = "v"
)

type PathElement struct {
	direction   Direction
	stepsToMove int
}

type CubeSide struct {
	boardCoords *Position
	board       [][]BoardContent
	left        *CubeSideTransition
	up          *CubeSideTransition
	right       *CubeSideTransition
	down        *CubeSideTransition
}

type CubeSideTransition struct {
	nextSide       *CubeSide
	rightRotations int
}

const cubeSideLength int = 50

var flatBoard [][]BoardContent
var currentPos *Position = &Position{x: -1, y: -1}
var currentDirIndex int = 3
var yParserPos int = 0
var maxY int = 0
var maxX int = 0
var isParsingMap bool = true
var directions [4]Direction = [4]Direction{Right, Down, Left, Up}
var path []*PathElement = make([]*PathElement, 0)
var cube []*CubeSide
var isUsingCube bool = false
var currentCubeSide *CubeSide
var stepFunctions [2]StepFunction = [2]StepFunction{walkOneStepOnX, walkOneStepOnY}
var isEndOfBoardFunctions [4]IsEndOfBoardFunction = [4]IsEndOfBoardFunction{isRightEndOfBoard, isBottomEndOfBoard, isLeftEndOfBoard, isTopEndOfBoard}
var cubeTransitionFunctions [4]CubeTransitionFunction = [4]CubeTransitionFunction{rightCubeTransition, downCubeTransition, leftCubeTransition, upCubeTransition}
var flatTransitionFunctions [2]FlatTransitionFunction = [2]FlatTransitionFunction{flatTransitionOnXEnd, flatTransitionOnYEnd}

func main() {
	solvePuzzle1()
	solvePuzzle2()
}

func solvePuzzle1() {
	parse()
	walkPath()
	printFlatBoard()
	calculateAndPrintResult(currentPos.x, currentPos.y)
}

func solvePuzzle2() {
	initStateForPuzzle2()
	parse()
	walkPath()
	printCube()
	absX, absY := localToWorldCoords(currentPos.x, currentPos.y, currentCubeSide)
	calculateAndPrintResult(absX, absY)
}

func parse() {
	filePuzzle := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle, parseMonkeyNotes)
	fillEmptyPositions()
	createCube()
}

func parseMonkeyNotes(line string) {
	if line == "" {
		isParsingMap = false
	} else if isParsingMap {
		parseMap(line)
	} else {
		parsePath(line, Right)
	}
}

func parseMap(line string) {
	maxY = yParserPos + 1
	flatBoard = append(flatBoard, make([]BoardContent, 0))

	for x := 0; x < len(line); x++ {
		if x > maxX {
			maxX = x + 1
		}
		currentChar := line[x : x+1]
		flatBoard[yParserPos] = append(flatBoard[yParserPos], BoardContent(currentChar))

		if currentPos.x == -1 && currentChar == string(Empty) {
			currentPos.x = x
			currentPos.y = yParserPos
		}
	}
	yParserPos++
}

func parsePath(line string, currentDir Direction) {
	var nextDir Direction
	nextDirChange := strings.IndexAny(line, "LR")

	if nextDirChange == -1 {
		stepsToMove := shared.ConvertStringToInt(line)
		path = append(path, &PathElement{direction: currentDir, stepsToMove: stepsToMove})
		return
	}

	if line[nextDirChange:nextDirChange+1] == "R" {
		nextDir = Right
	} else {
		nextDir = Left
	}
	stepsToMove := shared.ConvertStringToInt(line[:nextDirChange])
	path = append(path, &PathElement{direction: currentDir, stepsToMove: stepsToMove})

	parsePath(line[nextDirChange+1:], nextDir)
}

func initStateForPuzzle2() {
	isParsingMap = true
	isUsingCube = true
	currentPos.x = 0
	currentPos.y = 0
	currentDirIndex = 3
}

func createCube() {
	for y := 0; y < maxY; y += 50 {
		for x := 0; x < maxX; x += 50 {
			if flatBoard[y][x] != NoContent {
				cube = append(cube, createCubeSide(x, x+49, y, y+49))
			}
		}
	}

	addNeighbors()
	currentCubeSide = cube[0]
}

func createCubeSide(xStart int, xEnd int, yStart int, yEnd int) *CubeSide {
	var cubeBoard [][]BoardContent = make([][]BoardContent, cubeSideLength)

	for y := yStart; y <= yEnd; y++ {
		cubeBoard[y-yStart] = make([]BoardContent, cubeSideLength)
		for x := xStart; x <= xEnd; x++ {
			cubeBoard[y-yStart][x-xStart] = flatBoard[y][x]
		}
	}

	return &CubeSide{board: cubeBoard, boardCoords: &Position{x: xStart / 50, y: yStart / 50}}
}

func addNeighbors() {
	cube[0].right = &CubeSideTransition{nextSide: cube[1], rightRotations: 0}
	cube[0].down = &CubeSideTransition{nextSide: cube[2], rightRotations: 0}
	cube[0].left = &CubeSideTransition{nextSide: cube[3], rightRotations: 2}
	cube[0].up = &CubeSideTransition{nextSide: cube[5], rightRotations: 1}

	cube[1].right = &CubeSideTransition{nextSide: cube[4], rightRotations: 2}
	cube[1].down = &CubeSideTransition{nextSide: cube[2], rightRotations: 1}
	cube[1].left = &CubeSideTransition{nextSide: cube[0], rightRotations: 0}
	cube[1].up = &CubeSideTransition{nextSide: cube[5], rightRotations: 0}

	cube[2].right = &CubeSideTransition{nextSide: cube[1], rightRotations: 3}
	cube[2].down = &CubeSideTransition{nextSide: cube[4], rightRotations: 0}
	cube[2].left = &CubeSideTransition{nextSide: cube[3], rightRotations: 3}
	cube[2].up = &CubeSideTransition{nextSide: cube[0], rightRotations: 0}

	cube[3].right = &CubeSideTransition{nextSide: cube[4], rightRotations: 0}
	cube[3].down = &CubeSideTransition{nextSide: cube[5], rightRotations: 0}
	cube[3].left = &CubeSideTransition{nextSide: cube[0], rightRotations: 2}
	cube[3].up = &CubeSideTransition{nextSide: cube[2], rightRotations: 1}

	cube[4].right = &CubeSideTransition{nextSide: cube[1], rightRotations: 2}
	cube[4].down = &CubeSideTransition{nextSide: cube[5], rightRotations: 1}
	cube[4].left = &CubeSideTransition{nextSide: cube[3], rightRotations: 0}
	cube[4].up = &CubeSideTransition{nextSide: cube[2], rightRotations: 0}

	cube[5].right = &CubeSideTransition{nextSide: cube[4], rightRotations: 3}
	cube[5].down = &CubeSideTransition{nextSide: cube[1], rightRotations: 0}
	cube[5].left = &CubeSideTransition{nextSide: cube[0], rightRotations: 3}
	cube[5].up = &CubeSideTransition{nextSide: cube[3], rightRotations: 0}
}

func fillEmptyPositions() {
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			if x >= len(flatBoard[y]) {
				flatBoard[y] = append(flatBoard[y], NoContent)
			}
		}
	}
}

func walkPath() {
	for len(path) > 0 {
		pathElement := path[0]
		path = path[1:]
		updateDirection(pathElement.direction)
		walk(pathElement.stepsToMove, isEndOfBoardFunctions[currentDirIndex])
	}
}

func updateDirection(change Direction) {
	if change == Right {
		currentDirIndex = (currentDirIndex + 1) % len(directions)
	} else if currentDirIndex == 0 {
		currentDirIndex = len(directions) - 1
	} else {
		currentDirIndex--
	}
}

func isTopEndOfBoard() bool {
	return currentPos.y-1 < 0 || getBoard()[currentPos.y-1][currentPos.x] == NoContent
}

func isBottomEndOfBoard() bool {
	return currentPos.y+1 > len(getBoard())-1 || getBoard()[currentPos.y+1][currentPos.x] == NoContent
}

func isRightEndOfBoard() bool {
	return currentPos.x+1 > len(getBoard()[currentPos.y])-1 || getBoard()[currentPos.y][currentPos.x+1] == NoContent
}

func isLeftEndOfBoard() bool {
	return currentPos.x-1 < 0 || getBoard()[currentPos.y][currentPos.x-1] == NoContent
}

func walk(steps int, isAtEndOfBoardFunction IsEndOfBoardFunction) {
	currentDir := directions[currentDirIndex]
	singleStep := 1

	if currentDir == Up || currentDir == Left {
		singleStep = -1
	}

	for j := 1; j <= steps; j++ {
		isAtEndOfBoard := isAtEndOfBoardFunction()
		updateGraphics(currentDir)

		if isAtEndOfBoard {
			if isUsingCube {
				cubeTransitions := [4]*CubeSideTransition{currentCubeSide.right, currentCubeSide.down, currentCubeSide.left, currentCubeSide.up}
				if cubeTransitionOnEnd(steps-j, cubeTransitions[currentDirIndex], cubeTransitionFunctions[currentDirIndex]) == Interrupted {
					return
				}
			} else {
				flatTransitionFunctions[currentDirIndex%2](singleStep)
			}
		} else {
			stepFunctions[currentDirIndex%2](isAtEndOfBoard, singleStep, steps-j)
		}
	}
}

func walkOneStepOnX(isAtEndOfBoard bool, step int, stepsLeft int) {
	if getBoard()[currentPos.y][currentPos.x+step] != Wall {
		currentPos.x += step
	}
}

func walkOneStepOnY(isAtEndOfBoard bool, step int, stepsLeft int) {
	if getBoard()[currentPos.y+step][currentPos.x] != Wall {
		currentPos.y += step
	}
}

func cubeTransitionOnEnd(stepsLeft int, transition *CubeSideTransition, transitionFunction CubeTransitionFunction) WalkResult {
	oldPos := &Position{x: currentPos.x, y: currentPos.y}
	transitionFunction(currentPos.x, currentPos.y, transition.rightRotations)
	if transition.nextSide.board[currentPos.y][currentPos.x] == Wall {
		currentPos.x = oldPos.x
		currentPos.y = oldPos.y
		return Completed
	} else {
		currentCubeSide = transition.nextSide
		currentDirIndex = (currentDirIndex + transition.rightRotations + 1) % len(directions)
		path = append([]*PathElement{{direction: Left, stepsToMove: stepsLeft}}, path...)
		return Interrupted
	}
}

func upCubeTransition(x int, y int, clockwiseRotations int) {
	if clockwiseRotations == 0 {
		rightCubeTransition(currentPos.y, currentPos.x, 3)
	} else {
		rightCubeTransition(currentPos.y, currentPos.x, clockwiseRotations-1)
	}
}

func downCubeTransition(x int, y int, clockwiseRotations int) {
	if clockwiseRotations == 0 {
		leftCubeTransition(currentPos.y, currentPos.x, 3)
	} else {
		leftCubeTransition(currentPos.y, currentPos.x, clockwiseRotations-1)
	}
}

func rightCubeTransition(x int, y int, clockwiseRotations int) {
	if clockwiseRotations == 0 {
		currentPos.x = 0
		currentPos.y = y
	} else if clockwiseRotations == 1 {
		currentPos.x = 0
		currentPos.y = cubeSideLength - 1 - y
	} else if clockwiseRotations == 2 {
		currentPos.x = cubeSideLength - 1
		currentPos.y = cubeSideLength - 1 - y
	} else {
		currentPos.x = y
		currentPos.y = cubeSideLength - 1
	}
}

func leftCubeTransition(x int, y int, clockwiseRotations int) {
	if clockwiseRotations == 0 {
		currentPos.x = cubeSideLength - 1
		currentPos.y = y
	} else if clockwiseRotations == 1 {
		currentPos.x = cubeSideLength - 1 - y
		currentPos.y = cubeSideLength - 1
	} else if clockwiseRotations == 2 {
		currentPos.x = 0
		currentPos.y = cubeSideLength - 1 - y
	} else {
		currentPos.x = y
		currentPos.y = 0
	}
}

func flatTransitionOnXEnd(step int) {
	var oppositeXPos int

	if step > 0 {
		oppositeXPos = findMostLeftXPos(currentPos.y)
	} else {
		oppositeXPos = findMostRightXPos(currentPos.y)
	}

	if oppositeXPos != -1 && getBoard()[currentPos.y][oppositeXPos] != Wall {
		currentPos.x = oppositeXPos
	}
}

func flatTransitionOnYEnd(step int) {
	var oppositeYPos int

	if step < 0 {
		oppositeYPos = findMostDownYPos(currentPos.x)
	} else {
		oppositeYPos = findMostUpYPos(currentPos.x)
	}

	if oppositeYPos != -1 && getBoard()[oppositeYPos][currentPos.x] != Wall {
		currentPos.y = oppositeYPos
	}
}

func findMostLeftXPos(y int) int {
	for x := 0; x < maxX; x++ {
		if flatBoard[y][x] != NoContent {
			return x
		}
	}
	return -1
}

func findMostRightXPos(y int) int {
	for x := maxX - 1; x >= 0; x-- {
		if flatBoard[y][x] != NoContent {
			return x
		}
	}
	return -1
}

func findMostUpYPos(x int) int {
	for y := 0; y < maxY; y++ {
		if flatBoard[y][x] != NoContent {
			return y
		}
	}
	return -1
}

func findMostDownYPos(x int) int {
	for y := maxY - 1; y >= 0; y-- {
		if flatBoard[y][x] != NoContent {
			return y
		}
	}
	return -1
}

func calculateAndPrintResult(x int, y int) {
	password := 1000 * (y + 1)
	password += 4 * (x + 1)
	password += currentDirIndex

	fmt.Println("Final pos: " + strconv.Itoa(x+1) + "," + strconv.Itoa(y+1))
	fmt.Println("Facing: " + string(directions[currentDirIndex]))
	fmt.Println("Password: " + strconv.Itoa(password))
}

func updateGraphics(direction Direction) {
	getBoard()[currentPos.y][currentPos.x] = BoardContent(string(direction))
}

func printFlatBoard() {
	shared.PrintXHeader(1, maxX+1)

	for y := 0; y < maxY; y++ {
		var printedLine string = strconv.Itoa(shared.Abs(y+1)) + shared.GetYPadding(y+1)
		for x := 0; x < maxX; x++ {
			printedLine += string(flatBoard[y][x])
		}
		fmt.Println(printedLine)
	}
}

func printCube() {
	shared.PrintXHeader(1, maxX+1)

	for y := 0; y < maxY; y++ {
		var printedLine string = strconv.Itoa(shared.Abs(y+1)) + shared.GetYPadding(y+1)
		for x := 0; x < maxX; x++ {
			cubeContent := getCubeContent(x, y)
			if cubeContent == "" {
				printedLine += string(NoContent)
			} else {
				printedLine += cubeContent
			}
		}
		fmt.Println(printedLine)
	}
}

func getCubeContent(x int, y int) string {
	for i := 0; i < len(cube); i++ {
		xStart, _ := localToWorldCoords(0, 0, cube[i])
		_, yStart := localToWorldCoords(0, 0, cube[i])
		isInXRange := x >= xStart && x < xStart+cubeSideLength
		isInYRange := y >= yStart && y < yStart+cubeSideLength
		if isInXRange && isInYRange {
			cubeXPos := x - cube[i].boardCoords.x*cubeSideLength
			cubeYPos := y - cube[i].boardCoords.y*cubeSideLength
			return string(cube[i].board[cubeYPos][cubeXPos])
		}
	}
	return ""
}

func getBoard() [][]BoardContent {
	if isUsingCube {
		return currentCubeSide.board
	}
	return flatBoard
}

func localToWorldCoords(x int, y int, localCube *CubeSide) (int, int) {
	absX := x + (localCube.boardCoords.x * cubeSideLength)
	absY := y + (localCube.boardCoords.y * cubeSideLength)
	return absX, absY
}
