package main

import (
	"advent-of-code/shared"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Coordinates struct {
	x int
	y int
}

const ropeLengthPuzzle1 = 2
const ropeLengthPuzzle2 = 10

var posKnots []Coordinates
var visitedLocations map[Coordinates]bool

func main() {
	initPuzzle(ropeLengthPuzzle1)
	filePuzzle1 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle1, simulateRope)
	fmt.Println("Unique locations visited puzzle 1: " + strconv.Itoa(countUniqueLocations()))

	initPuzzle(ropeLengthPuzzle2)
	filePuzzle2 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle2, simulateRope)
	fmt.Println("Unique locations visited puzzle 2: " + strconv.Itoa(countUniqueLocations()))
}

func initPuzzle(ropeLength int) {
	posKnots = []Coordinates{}
	visitedLocations = make(map[Coordinates]bool)
	visitedLocations[Coordinates{0, 0}] = true

	for i := 0; i < ropeLength; i++ {
		posKnots = append(posKnots, Coordinates{x: 0, y: 0})
	}
}

func simulateRope(line string) {
	headMovement := strings.Split(line, " ")
	direction := headMovement[0]
	amount, err := strconv.Atoi(headMovement[1])

	if err != nil {
		fmt.Println("Couldn't parse head movement")
	}

	for i := 0; i < amount; i++ {
		moveHead(direction)
		moveTails()
	}
}

func moveHead(direction string) {
	posHead := &posKnots[0]

	if direction == "L" {
		posHead.x--
	} else if direction == "R" {
		posHead.x++
	} else if direction == "U" {
		posHead.y++
	} else if direction == "D" {
		posHead.y--
	}
}

func moveTails() {
	for i := 1; i < len(posKnots); i++ {
		prevKnot := posKnots[i-1]
		currentKnot := &posKnots[i]
		moveTail(prevKnot, currentKnot)
	}

	tailEnd := posKnots[len(posKnots)-1]
	visitedLocations[Coordinates{x: tailEnd.x, y: tailEnd.y}] = true
}

func moveTail(prevKnot Coordinates, currentKnot *Coordinates) {
	isMovementRequired := abs(prevKnot.x-currentKnot.x) > 1 || abs(prevKnot.y-currentKnot.y) > 1
	if !isMovementRequired {
		return
	}

	xMovement := (prevKnot.x - currentKnot.x)
	yMovement := (prevKnot.y - currentKnot.y)

	if prevKnot.y == currentKnot.y {
		currentKnot.x += xMovement / 2
	} else if prevKnot.x == currentKnot.x {
		currentKnot.y += yMovement / 2
	} else {
		currentKnot.x += xMovement / abs(xMovement)
		currentKnot.y += yMovement / abs(yMovement)
	}
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func countUniqueLocations() int {
	keys := reflect.ValueOf(visitedLocations).MapKeys()

	return len(keys)
}
