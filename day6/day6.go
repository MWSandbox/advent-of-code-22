package main

import (
	"advent-of-code/shared"
	"fmt"
	"strconv"
	"reflect"
)

var markerLength int = 4
var startOfSignalPos int = 0

func main() {
	filePuzzle1 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle1, calculateStartOfSignal)
	fmt.Println("Start of signal at position: " + strconv.Itoa(startOfSignalPos))

	markerLength = 14
	filePuzzle2 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle2, calculateStartOfSignal)
	fmt.Println("Start of message at position: " + strconv.Itoa(startOfSignalPos))
}

func calculateStartOfSignal(line string) {
	for i := 0; i < len(line) - 3; i++ {
		uniqueKeyMap := make(map[byte]bool)
		uniqueKeyMap[line[i]] = true

		for j := 1; j < markerLength; j++ {
			currentByte := line[i+j]
			if !uniqueKeyMap[currentByte] {
				fmt.Println(string(currentByte))
				uniqueKeyMap[currentByte] = true
			}
		}

		keys := reflect.ValueOf(uniqueKeyMap).MapKeys()

		if len(keys) == markerLength {
			startOfSignalPos = i + markerLength
			return
		}
	}
}