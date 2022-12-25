package main

import (
	"advent-of-code/shared"
	"fmt"
	"sort"
	"strconv"
)

type Packet struct {
	value    int
	children []*Packet
	parent   *Packet
}

type OrderedPacket struct {
	packet       *Packet
	originalLine string
}

type ComparisonResult int

const (
	Larger ComparisonResult = iota
	Smaller
	Equal
)

const dividerPacket1 string = "[[2]]"
const dividerPacket2 string = "[[6]]"

var lineNumber int = 0
var firstTree *Packet
var secondTree *Packet
var indexSum int = 0
var orderedPackets []*OrderedPacket

func main() {
	filePuzzle1 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle1, detectOrdering)
	compareAndAddToIndex()
	fmt.Println("Final index sum: " + strconv.Itoa(indexSum))

	filePuzzle2 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle2, addLinesForOrdering)
	addDividerPackets()
	sortPackets()
	index1, index2 := provideDividerIndices()
	printOrderedPackets()

	fmt.Println("Divider indices: " + strconv.Itoa(index1) + ", " + strconv.Itoa(index2))
	fmt.Println("Distress signal: " + strconv.Itoa(index1*index2))
}

func detectOrdering(line string) {
	isFirstTree := lineNumber%3 == 0
	isSecondTree := lineNumber%3 == 1

	if isFirstTree {
		firstTree = readTreeInput(line)
	} else if isSecondTree {
		secondTree = readTreeInput(line)
	} else {
		compareAndAddToIndex()
	}

	lineNumber++
}

func compareAndAddToIndex() {
	if compare(firstTree, secondTree) != 0 {
		indexSum += lineNumber/3 + 1
	}
}

func addLinesForOrdering(line string) {
	if line != "" {
		currentPacket := readTreeInput(line)
		orderedPackets = append(orderedPackets, &OrderedPacket{packet: currentPacket, originalLine: line})
	}

	lineNumber++
}

func readTreeInput(line string) *Packet {

	var currentPacket *Packet
	var rootNode *Packet
	var currentValue string = ""

	for i := 0; i < len(line); i++ {
		currentChar := line[i : i+1]

		if currentChar == "[" {
			isRootNode := currentPacket == nil
			currentPacket = createNewChild(currentPacket, -1)
			if isRootNode {
				rootNode = currentPacket
			}
		} else if currentChar == "," {
			if currentValue != "" {
				appendCurrentValue(currentPacket, currentValue)
				currentValue = ""
			}
		} else if currentChar == "]" {
			if currentValue != "" {
				appendCurrentValue(currentPacket, currentValue)
				currentValue = ""
			}
			currentPacket = currentPacket.parent
		} else {
			currentValue += currentChar
		}
	}

	return rootNode
}

func createNewChild(currentPacket *Packet, value int) *Packet {
	child := &Packet{value: value, children: make([]*Packet, 0), parent: currentPacket}

	if currentPacket != nil {
		currentPacket.children = append(currentPacket.children, child)
	}
	return child
}

func appendCurrentValue(currentElement *Packet, strValue string) {
	intValue, err := strconv.Atoi(strValue)

	if err != nil {
		fmt.Println("Couldn't convert to int: " + strValue)
	}

	createNewChild(currentElement, intValue)
}

func compare(leftStart *Packet, rightStart *Packet) ComparisonResult {
	var currentLeftPacket *Packet = leftStart
	var currentRightPacket *Packet = rightStart

	for i := 0; i < len(leftStart.children); i++ {
		currentLeftPacket = leftStart.children[i]

		if len(rightStart.children) > i {
			currentRightPacket = rightStart.children[i]
		} else {
			return Larger
		}

		leftPacketContainsValue := currentLeftPacket.value != -1
		rightPacketContainsValue := currentRightPacket.value != -1
		bothPacketsContainLists := !leftPacketContainsValue && !rightPacketContainsValue
		bothPacketsHaveValue := leftPacketContainsValue && rightPacketContainsValue

		if bothPacketsHaveValue {
			if currentLeftPacket.value > currentRightPacket.value {
				return Larger
			} else if currentLeftPacket.value < currentRightPacket.value {
				return Smaller
			}
		} else if bothPacketsContainLists {
			result := compare(currentLeftPacket, currentRightPacket)
			if result != Equal {
				return result
			}
		} else if leftPacketContainsValue {
			leftListElement := wrapValueInList(currentLeftPacket.value)
			result := compare(leftListElement, currentRightPacket)
			if result != Equal {
				return result
			}
		} else if rightPacketContainsValue {
			rightListElement := wrapValueInList(currentRightPacket.value)
			result := compare(currentLeftPacket, rightListElement)
			if result != Equal {
				return result
			}
		}
	}

	if len(leftStart.children) < len(rightStart.children) {
		return Smaller
	}
	return Equal
}

func wrapValueInList(valueToWrap int) *Packet {
	listElement := &Packet{value: -1, children: make([]*Packet, 0), parent: nil}
	listChild := &Packet{value: valueToWrap, children: make([]*Packet, 0), parent: listElement}
	listElement.children = append(listElement.children, listChild)
	return listElement
}

func addDividerPackets() {
	addLinesForOrdering(dividerPacket1)
	addLinesForOrdering(dividerPacket2)
}

func sortPackets() {
	sort.Slice(orderedPackets, func(i, j int) bool {
		if orderedPackets[j].originalLine == "[[2]]" && orderedPackets[i].originalLine == "[[2],[9,8,[]]]" {
			fmt.Println("breakpoint")
		}

		comparisonValue := compare(orderedPackets[i].packet, orderedPackets[j].packet)
		return comparisonValue != 0
	})
}

func provideDividerIndices() (int, int) {
	var indexDivider1 int
	var indexDivider2 int

	for i := 0; i < len(orderedPackets); i++ {
		currentLine := orderedPackets[i].originalLine
		if currentLine == dividerPacket1 {
			indexDivider1 = i + 1
		} else if currentLine == dividerPacket2 {
			indexDivider2 = i + 1
		}
		fmt.Println(orderedPackets[i].originalLine)
	}

	return indexDivider1, indexDivider2
}

func printOrderedPackets() {
	for i := 0; i < len(orderedPackets); i++ {
		fmt.Println(orderedPackets[i].originalLine)
	}
}
