package main

import (
	"advent-of-code/shared"
	"strconv"
	"fmt"
)

type IsTreeVisibleChecker func(int, int) (bool, int)

var treeSizes [][]int
var currentRow int = 0
var visibleTreeCount int = 0
var highestScenicScore int = 0
var currentTreeSize int = 0
var visibilityCheckers [4]IsTreeVisibleChecker = [4]IsTreeVisibleChecker{isTreeVisibleFromLeft, isTreeVisibleFromRight, isTreeVisibleFromTop, isTreeVisibleFromBottom}


func main() {
	filePuzzle := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle, parseTreeMap)

	calculateVisibleTreeCount()
	fmt.Println("Visible tree count: " + strconv.Itoa(visibleTreeCount))
	fmt.Println("Highest scenic score: " + strconv.Itoa(highestScenicScore))
}

func parseTreeMap(line string) {
	treeSizes = append(treeSizes, []int{})

	for i := 0; i < len(line); i++ {
		currentTreeSizeAsString := line[i:i+1]
		currentTreeSize, err := strconv.Atoi(currentTreeSizeAsString)

		if err != nil {
			fmt.Println("Couldn't convert tree size to int: " + currentTreeSizeAsString)
		}

		treeSizes[currentRow] = append(treeSizes[currentRow], currentTreeSize)
	}
	currentRow++
}

func calculateVisibleTreeCount() {
	for i := 0; i < len(treeSizes); i ++ {
		for j := 0; j < len(treeSizes[i]); j++ {
			isVisible, scenicScore := calculateVisibility(j, i)
			
			if isVisible {
				visibleTreeCount++
			}

			if scenicScore > highestScenicScore {
				highestScenicScore = scenicScore
			}
		}
	}
}

func calculateVisibility(x int, y int) (bool, int) {
	var isTreeVisible bool = false
	var scenicScore int = 1
	currentTreeSize = treeSizes[y][x]

	for i := 0; i < len(visibilityCheckers); i++ {
		isVisible, visibleTreeCount := visibilityCheckers[i](x, y)
		scenicScore *= visibleTreeCount
		isTreeVisible = isTreeVisible || isVisible
	}

	return isTreeVisible, scenicScore
}

func isTreeVisibleFromLeft(x int, y int) (bool, int) {
	var visibleTrees int = 0

	for i := x-1; i >= 0; i-- {
		visibleTrees++
		if treeSizes[y][i] >= currentTreeSize {
			return false, visibleTrees
		}
	}
	return true, visibleTrees
}

func isTreeVisibleFromRight(x int, y int) (bool, int) {
	var visibleTrees int = 0

	for i := x+1; i < len(treeSizes[y]); i++ {
		visibleTrees++
		if treeSizes[y][i] >= currentTreeSize {
			return false, visibleTrees
		}
	}
	return true, visibleTrees
}

func isTreeVisibleFromTop(x int, y int) (bool, int) {
	var visibleTrees int = 0

	for i := y-1; i >= 0; i-- {
		visibleTrees++
		if treeSizes[i][x] >= currentTreeSize {
			return false, visibleTrees
		}
	}
	return true, visibleTrees
}

func isTreeVisibleFromBottom(x int, y int) (bool, int) {
	var visibleTrees int = 0

	for i := y+1; i < len(treeSizes); i++ {
		visibleTrees++
		if treeSizes[i][x] >= currentTreeSize {
			return false, visibleTrees
		}
	}
	return true, visibleTrees
}
