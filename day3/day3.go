package main

import (
	"advent-of-code/shared"
	"fmt"
	"strconv"
	"unicode"
)

var totalPriorityOfDuplicates int = 0
var totalPriorityOfBadges int = 0
var currentMemberIndex int = 0
var sortedItemsOfGroupMembers [3][53]bool
var groupSize = 3

func main() {
	filePuzzle1 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle1, calculatePriorityOfDuplicates)

	filePuzzle2 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle2, calculatePriorityOfBadges)

	fmt.Println("Total priority of duplicates: " + strconv.Itoa(totalPriorityOfDuplicates))
	fmt.Println("Total priority of badges: " + strconv.Itoa(totalPriorityOfBadges))
}

func calculatePriorityOfDuplicates(line string) {
	half := len(line) / 2

	firstCompartment := line[0:half]
	secondCompartment := line[half:len(line)]

	itemsInFirstCompartment := sortItemsOfCompartment(firstCompartment)
	itemsInSecondCompartment := sortItemsOfCompartment(secondCompartment)

	totalPriorityOfDuplicates += getDuplicateItemPriority(itemsInFirstCompartment, itemsInSecondCompartment)
}

func calculatePriorityOfBadges(line string) {
	sortedItems := sortItemsOfCompartment(line)
	isLastGroupMember := currentMemberIndex % 3 == 2

	sortedItemsOfGroupMembers[currentMemberIndex % 3] = sortedItems

	if isLastGroupMember {
		totalPriorityOfBadges += getBadgePriority()
	}

	currentMemberIndex++
}

func sortItemsOfCompartment(compartment string) [53]bool {
	var sortedItems [53]bool
	var bonusPriorityForUpperCase int = 26
	var lowerCaseStartIndex int = int('a') - 1
	var upperCaseStartIndex int = int('A') - 1

	for i := 0; i < len(compartment); i++ {
		var currentItem rune = rune(compartment[i])
		var itemPriority int
		
		if unicode.IsLower(currentItem) {
			itemPriority = int(currentItem) - lowerCaseStartIndex
		} else {
			itemPriority = int(currentItem) - upperCaseStartIndex + bonusPriorityForUpperCase
		}
		sortedItems[itemPriority] = true
	}

	return sortedItems
}

func getDuplicateItemPriority(itemsInFirstCompartment [53]bool, itemsInSecondCompartment [53]bool) int {
	for i := 0; i < len(itemsInFirstCompartment); i++ {
		if itemsInFirstCompartment[i] && itemsInSecondCompartment[i] {
			return i
		}
	}
	return 0
}

func getBadgePriority() int {
	for i := 0; i < len(sortedItemsOfGroupMembers[0]); i++ {
		var isCurrentItemBadge bool = true

		for j := 0; j < groupSize; j++ {
			isCurrentItemBadge = isCurrentItemBadge && sortedItemsOfGroupMembers[j][i]
		}

		if (isCurrentItemBadge) {
			return i
		}
	}
	
	return 0
}