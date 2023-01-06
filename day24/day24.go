package main

import (
	"advent-of-code/shared"
	"fmt"
	"sort"
	"strconv"
)

type Direction string

const (
	Right  Direction = ">"
	Bottom           = "v"
	Left             = "<"
	Top              = "^"
)

type Position struct {
	x int
	y int
}

type Path struct {
	positions []*Position
	blizzards []*Blizzard
}

type Blizzard struct {
	pos       *Position
	direction Direction
}

var blizzards []*Blizzard
var yParserPos int = 0
var currentPos *Position = &Position{x: 1, y: 0}
var currentTarget *Position
var maxX int
var maxY int
var minutesAndPosHashesToIsVisited map[int]bool = make(map[int]bool)
var bestPaths [3]*Path
var targets [2]*Position = [2]*Position{{x: 1, y: 0}}

func main() {
	filePuzzle := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle, parseValley)
	maxY = yParserPos - 1
	targets[1] = &Position{x: maxX - 1, y: maxY}
	currentTarget = targets[1]

	solvePuzzle1()
	solvePuzzle2()
}

func solvePuzzle1() {
	bestPaths[0] = simulateWalkThroughValley()
	fmt.Println("Minutes for escape through valley: " + strconv.Itoa(len(bestPaths[0].positions)-1))
}

func solvePuzzle2() {
	resetForTarget(targets[0], bestPaths[0])
	bestPaths[1] = simulateWalkThroughValley()
	fmt.Println("Minutes to return for a snack: " + strconv.Itoa(len(bestPaths[1].positions)-1))

	resetForTarget(targets[1], bestPaths[1])
	bestPaths[2] = simulateWalkThroughValley()
	fmt.Println("Minutes to bring the snack back to the elves: " + strconv.Itoa(len(bestPaths[2].positions)-1))

	sum := 0
	for i := 0; i < len(bestPaths); i++ {
		sum += len(bestPaths[i].positions) - 1
	}
	fmt.Println("Total minutes in the blizzard valley: " + strconv.Itoa(sum))
}

func resetForTarget(target *Position, lastPath *Path) {
	currentPos.x = currentTarget.x
	currentPos.y = currentTarget.y
	currentTarget = target
	minutesAndPosHashesToIsVisited = make(map[int]bool)
	blizzards = copyPath(lastPath).blizzards
}

func parseValley(line string) {
	maxX = len(line) - 1

	for x := 0; x < len(line); x++ {
		currentChar := line[x : x+1]

		if currentChar != "#" && currentChar != "." {
			blizzards = append(blizzards, &Blizzard{pos: &Position{x: x, y: yParserPos}, direction: Direction(currentChar)})
		}
	}

	yParserPos++
}

func simulateWalkThroughValley() *Path {
	var paths []*Path = make([]*Path, 0)
	startPositions := []*Position{{x: currentPos.x, y: currentPos.y}}
	paths = append(paths, &Path{positions: startPositions, blizzards: blizzards})

	return walkPath(paths)
}

func walkPath(paths []*Path) *Path {
	for len(paths) > 0 {
		currentPath := paths[len(paths)-1]
		paths = paths[:len(paths)-1]
		lastStep := currentPath.positions[len(currentPath.positions)-1]
		currentPos.x = lastStep.x
		currentPos.y = lastStep.y
		// printMap()
		// printPath(currentPath)

		if isNextToTarget() {
			printMap()
			currentPath.positions = append(currentPath.positions, &Position{x: currentTarget.x, y: currentTarget.y})
			for i := 0; i < len(currentPath.blizzards); i++ {
				updateBlizzardPosition(currentPath.blizzards[i])
			}
			return currentPath
		} else {
			surroundingBlockers := getAvailableTilesAfterBlizzardMovement(currentPath.blizzards)
			if !surroundingBlockers[1][1] {
				paths = addPath(paths, currentPath, currentPos.x, currentPos.y)
			}

			for y := -1; y <= 1; y += 2 {
				if !surroundingBlockers[1+y][1] {
					paths = addPath(paths, currentPath, currentPos.x, currentPos.y+y)
				}
			}

			for x := -1; x <= 1; x += 2 {
				if !surroundingBlockers[1][1+x] {
					paths = addPath(paths, currentPath, currentPos.x+x, currentPos.y)
				}
			}
		}
	}
	return nil
}

func isNextToTarget() bool {
	diffX := currentPos.x - currentTarget.x
	diffY := currentPos.y - currentTarget.y
	isNextToTarget := shared.Abs(diffX)+shared.Abs(diffY) == 1
	return isNextToTarget
}

func addPath(paths []*Path, oldPath *Path, x int, y int) []*Path {
	hashValue := cantorHashMinutesAndPos(len(oldPath.positions)+1, &Position{x: x, y: y})

	if minutesAndPosHashesToIsVisited[hashValue] {
		return paths
	}

	minutesAndPosHashesToIsVisited[hashValue] = true
	newPath := copyPath(oldPath)
	nextPos := &Position{x: x, y: y}
	newPath.positions = append(newPath.positions, nextPos)
	paths = append(paths, newPath)

	sort.Slice(paths, func(i, j int) bool {

		isPathEquals := len(paths[i].positions) == len(paths[j].positions)
		if isPathEquals {
			lastStepI := paths[i].positions[len(paths[i].positions)-1]
			lastStepJ := paths[j].positions[len(paths[j].positions)-1]
			diffI := shared.Abs(currentTarget.x-lastStepI.x) + shared.Abs(currentTarget.y-lastStepI.y)
			diffJ := shared.Abs(currentTarget.x-lastStepJ.x) + shared.Abs(currentTarget.y-lastStepJ.y)
			return diffI > diffJ
		}
		return len(paths[i].positions) > len(paths[j].positions)
	})

	return paths
}

func copyPath(oldPath *Path) *Path {
	var newPositions []*Position
	var newBlizzards []*Blizzard

	for i := 0; i < len(oldPath.positions); i++ {
		newPositions = append(newPositions, &Position{x: oldPath.positions[i].x, y: oldPath.positions[i].y})
	}

	for i := 0; i < len(oldPath.blizzards); i++ {
		oldBlizzard := oldPath.blizzards[i]
		newBlizzards = append(newBlizzards, &Blizzard{pos: &Position{x: oldBlizzard.pos.x, y: oldBlizzard.pos.y}, direction: oldBlizzard.direction})
	}

	return &Path{positions: newPositions, blizzards: newBlizzards}
}

func getAvailableTilesAfterBlizzardMovement(currentBlizzards []*Blizzard) [][]bool {
	surroundingBlockers := createFreeSurrounding()
	updateSurroundingOnBoundaries(surroundingBlockers)

	for i := 0; i < len(currentBlizzards); i++ {
		updateBlizzardPosition(currentBlizzards[i])
		updateSurroundingBasedOnBlizzardPos(surroundingBlockers, currentBlizzards[i])
	}

	return surroundingBlockers
}

func createFreeSurrounding() [][]bool {
	var surroundingBlockers [][]bool = make([][]bool, 3)

	for i := 0; i < 3; i++ {
		surroundingBlockers[i] = make([]bool, 3)
	}

	return surroundingBlockers
}

func updateSurroundingOnBoundaries(surroundingBlockers [][]bool) {
	if currentPos.x == 1 && currentPos.y == 0 {
		surroundingBlockers[1][2] = true
	}

	if currentPos.x == 1 {
		setXBlockers(surroundingBlockers, 0)
	}

	if currentPos.x == maxX-1 {
		setXBlockers(surroundingBlockers, 2)
	}

	if currentPos.y == 0 || currentPos.y == 1 {
		setYBlockers(surroundingBlockers, 0)
	}

	if currentPos.y == maxY-1 {
		setYBlockers(surroundingBlockers, 2)
	}

	if currentPos.x == maxX-1 && currentPos.y == maxY {
		surroundingBlockers[1][0] = true
		surroundingBlockers[1][2] = true
		setYBlockers(surroundingBlockers, 2)
	}
}

func setXBlockers(surroundingBlockers [][]bool, x int) {
	for y := 0; y < 3; y++ {
		surroundingBlockers[y][x] = true
	}
}

func setYBlockers(surroundingBlockers [][]bool, y int) {
	for x := 0; x < 3; x++ {
		surroundingBlockers[y][x] = true
	}
}

func updateSurroundingBasedOnBlizzardPos(surroundingBlockers [][]bool, blizzard *Blizzard) {
	isInXRange := shared.IsInRange(blizzard.pos.x, currentPos.x-1, currentPos.x+1)
	isInYRange := shared.IsInRange(blizzard.pos.y, currentPos.y-1, currentPos.y+1)

	if isInXRange && isInYRange {
		xDiff := blizzard.pos.x - currentPos.x
		yDiff := blizzard.pos.y - currentPos.y
		surroundingBlockers[1+yDiff][1+xDiff] = true
	}
}

func updateBlizzardPosition(blizzard *Blizzard) {
	if blizzard.direction == Right {
		if blizzard.pos.x == maxX-1 {
			blizzard.pos.x = 1
		} else {
			blizzard.pos.x++
		}
	} else if blizzard.direction == Left {
		if blizzard.pos.x == 1 {
			blizzard.pos.x = maxX - 1
		} else {
			blizzard.pos.x--
		}
	} else if blizzard.direction == Top {
		if blizzard.pos.y == 1 {
			blizzard.pos.y = maxY - 1
		} else {
			blizzard.pos.y--
		}
	} else {
		if blizzard.pos.y == maxY-1 {
			blizzard.pos.y = 1
		} else {
			blizzard.pos.y++
		}
	}
}

func printMap() {
	for y := 0; y <= maxY; y++ {
		lineToPrint := ""
		for x := 0; x <= maxX; x++ {
			blizzardsOnPos := getBlizzardForPos(x, y)
			isTargetPos := x == currentTarget.x && y == currentTarget.y
			isStartPos := x == 1 && y == 0

			if currentPos.x == x && currentPos.y == y {
				lineToPrint += "E"
			} else if isTargetPos || isStartPos {
				lineToPrint += "."
			} else if len(blizzardsOnPos) == 0 {
				if x == 0 || x == maxX || y == 0 || y == maxY {
					lineToPrint += "#"
				} else {
					lineToPrint += "."
				}
			} else if len(blizzardsOnPos) == 1 {
				lineToPrint += string(blizzardsOnPos[0].direction)
			} else {
				lineToPrint += strconv.Itoa(len(blizzardsOnPos))
			}
		}
		fmt.Println(lineToPrint)
	}
}

func printPath(path *Path) {
	printedPath := ""

	for i := 0; i < len(path.positions); i++ {
		printedPath += strconv.Itoa(path.positions[i].y) + "," + strconv.Itoa(path.positions[i].x) + " "
	}

	fmt.Println(printedPath)
}

func getBlizzardForPos(x int, y int) []*Blizzard {
	var blizzardsOnPos []*Blizzard

	for i := 0; i < len(blizzards); i++ {
		blizzard := blizzards[i]
		if blizzard.pos.x == x && blizzard.pos.y == y {
			blizzardsOnPos = append(blizzardsOnPos, blizzard)
		}
	}

	return blizzardsOnPos
}

func cantorHashMinutesAndPos(minutes int, pos *Position) int {
	return (cantorHashPos(pos)+minutes)*(cantorHashPos(pos)+minutes+1)/2 + minutes
}

func cantorHashPos(pos *Position) int {
	return (pos.x+pos.y)*(pos.x+pos.y+1)/2 + pos.y
}
