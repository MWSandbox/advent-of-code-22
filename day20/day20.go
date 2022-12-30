package main

import (
	"advent-of-code/shared"
	"fmt"
	"strconv"
)

type ListEntry struct {
	originalIndex int
	shiftValue    int
	successor     *ListEntry
}

const decryptionKey int = 811589153
const roundsToMix int = 10

var encryptedList []*ListEntry
var listForPuzzle2 []*ListEntry
var zeroEntry *ListEntry

func main() {
	filePuzzle1 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle1, readEncryptedFile)
	prepareListForPuzzle2()

	solvePuzzle1()
	solvePuzzle2()
}

func solvePuzzle1() {
	mixFile(encryptedList[0])

	fmt.Println("########### PUZZLE 1 ###########")
	fmt.Println("Puzzle 1 result: " + strconv.Itoa(calculateGrovePosition()))
}

func solvePuzzle2() {
	encryptedList = listForPuzzle2
	startEntry := encryptedList[0]
	for i := 0; i < roundsToMix; i++ {
		mixFile(startEntry)
	}

	fmt.Println("")
	fmt.Println("########### PUZZLE 2 ###########")
	fmt.Println("Puzzle 2 result: " + strconv.Itoa(calculateGrovePosition()))
}

func calculateGrovePosition() int {
	zeroIndex := indexOf(zeroEntry)
	fmt.Println("Position of 0: " + strconv.Itoa(zeroIndex))

	puzzleResult := 0
	for i := 1000; i <= 3000; i += 1000 {
		currentValue := encryptedList[(zeroIndex+i)%len(encryptedList)].shiftValue
		puzzleResult += currentValue
		fmt.Println(strconv.Itoa(i) + ". number after 0: " + strconv.Itoa(currentValue))
	}
	return puzzleResult
}

func readEncryptedFile(line string) {
	shiftValue := shared.ConvertStringToInt(line)
	index := len(encryptedList)
	listEntry := &ListEntry{originalIndex: index, shiftValue: shiftValue}

	if index > 0 {
		encryptedList[index-1].successor = listEntry
	}

	if shiftValue == 0 {
		zeroEntry = listEntry
	}

	encryptedList = append(encryptedList, listEntry)
}

func mixFile(startEntry *ListEntry) {
	currentEntry := encryptedList[indexOf(startEntry)]

	for currentEntry != nil {
		sourceIndex := indexOf(currentEntry)
		targetIndex := sourceIndex + currentEntry.shiftValue

		if targetIndex <= 0 {
			targetIndex = len(encryptedList) - 1 + ((sourceIndex + currentEntry.shiftValue) % (len(encryptedList) - 1))
		} else if targetIndex >= len(encryptedList) {
			targetIndex = (sourceIndex + currentEntry.shiftValue) % (len(encryptedList) - 1)
		}

		if targetIndex < sourceIndex {
			shiftLeft(currentEntry, sourceIndex, targetIndex)
		} else {
			shiftRight(currentEntry, sourceIndex, targetIndex)
		}

		currentEntry = currentEntry.successor
	}
}

func shiftLeft(listEntry *ListEntry, sourceIndex int, targetIndex int) {
	listBeforeTarget := encryptedList[:targetIndex]
	listBetweenTargetAndSource := encryptedList[targetIndex:sourceIndex]
	listAfterSource := encryptedList[sourceIndex+1:]

	updatedList := make([]*ListEntry, 0)
	updatedList = append(updatedList, listBeforeTarget...)
	updatedList = append(updatedList, listEntry)
	updatedList = append(updatedList, listBetweenTargetAndSource...)
	updatedList = append(updatedList, listAfterSource...)
	encryptedList = updatedList
}

func shiftRight(listEntry *ListEntry, sourceIndex int, targetIndex int) {
	listBeforeSource := encryptedList[:sourceIndex]
	listBetweenSourceAndTarget := encryptedList[sourceIndex+1 : targetIndex+1]
	listAfterTarget := encryptedList[targetIndex+1:]

	updatedList := make([]*ListEntry, 0)
	updatedList = append(updatedList, listBeforeSource...)
	updatedList = append(updatedList, listBetweenSourceAndTarget...)
	updatedList = append(updatedList, listEntry)
	updatedList = append(updatedList, listAfterTarget...)
	encryptedList = updatedList
}

func prepareListForPuzzle2() {
	for i := 0; i < len(encryptedList); i++ {
		puzzle1Entry := encryptedList[i]
		shiftValue := puzzle1Entry.shiftValue
		puzzle2Entry := &ListEntry{originalIndex: i, shiftValue: shiftValue * decryptionKey}

		if i > 0 {
			listForPuzzle2[len(listForPuzzle2)-1].successor = puzzle2Entry
		}

		listForPuzzle2 = append(listForPuzzle2, puzzle2Entry)
	}
}

func indexOf(entry *ListEntry) int {
	for i := 0; i < len(encryptedList); i++ {
		if encryptedList[i].originalIndex == entry.originalIndex {
			return i
		}
	}
	return -1
}
