package main

import (
	"advent-of-code/shared"
	"strconv"
	"strings"
	"fmt"
	"reflect"
	"sort"
)

type WorryLevelIncreaser func(int, int) (int)

type Monkey struct {
	items []uint64
	operationParts []string
	divisor int
	nextMonkeyOnTrue int
	nextMonkeyOnFalse int
	inspectionCount int
}

const startingItemsIdLength int = 14
const operationIdLength int = 9
const oldWorryLevelId string = "old"
const worryLevelReducer int = 1
const rounds int = 10000

var indexToMonkey map[int]*Monkey = make(map[int]*Monkey)
var operatorToFunction map[string]WorryLevelIncreaser = make(map[string]WorryLevelIncreaser)
var currentMonkey int = -1

func main() {
	initialize()
	filePuzzle := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle, parseMonkeyState)
	playMonkeyGame()
}

func initialize() {
	operatorToFunction["*"] = multiplyValues
	operatorToFunction["+"] = addValues
}

func parseMonkeyState(line string) {
	trimmedLine := strings.TrimSpace(line)
	isMonkeyStartPos := len(line) > 6 && line[0:6] == "Monkey"
	isStartingItems := len(trimmedLine) > startingItemsIdLength && trimmedLine[0:startingItemsIdLength] == "Starting items"
	isOperation := len(trimmedLine) > operationIdLength && trimmedLine[0:operationIdLength] == "Operation"
	isTest := len(trimmedLine) > 4 && trimmedLine[0:4] == "Test"
	isTrueCondition := len(trimmedLine) > 7 && trimmedLine[0:7] == "If true"
	isFalseCondition := len(trimmedLine) > 8 && trimmedLine[0:8] == "If false"

	if isMonkeyStartPos {
		currentMonkey++
		indexToMonkey[currentMonkey] = &Monkey{items: nil, operationParts: nil, divisor: 0, inspectionCount: 0}
	} else if isStartingItems {
		itemsAsString := strings.Split(trimmedLine[startingItemsIdLength+2:len(trimmedLine)], ", ")
		if itemsAsString != nil {
			var items []uint64
			for i := 0; i < len(itemsAsString); i++ {
				item, err := strconv.Atoi(itemsAsString[i])

				if err != nil {
					fmt.Println("Couldn't parse item: " + itemsAsString[i])
				}
				items = append(items, item)
			}
			indexToMonkey[currentMonkey].items = items
		}
	} else if isOperation {
		fullOperation := trimmedLine[operationIdLength+8:len(trimmedLine)]
		operationParts := strings.Split(fullOperation, " ")
		indexToMonkey[currentMonkey].operationParts = operationParts
	} else if isTest {
		divisorAsString := trimmedLine[19: len(trimmedLine)]
		divisor, err := strconv.Atoi(divisorAsString)

		if err != nil {
			fmt.Println("Couldn't parse divisor: " + divisorAsString)
		}
		indexToMonkey[currentMonkey].divisor = divisor
	} else if isTrueCondition {
		nextMonkeyOnTrue := getNextMonkey(trimmedLine)
		indexToMonkey[currentMonkey].nextMonkeyOnTrue = nextMonkeyOnTrue
	} else if isFalseCondition {
		nextMonkeyOnFalse := getNextMonkey(trimmedLine)
		indexToMonkey[currentMonkey].nextMonkeyOnFalse = nextMonkeyOnFalse
	}
}

func getNextMonkey(line string) int {
	nextMonkeyAsString := line[len(line)-1:len(line)]
	nextMonkey, err := strconv.Atoi(nextMonkeyAsString)

	if err != nil {
		fmt.Println("Couldn't parse next monkey: " + nextMonkeyAsString)
	}
	return nextMonkey
}

func executeOperation(oldWorryLevel int, operationParts []string) int {
	firstValue := getDynamicFunctionValue(operationParts[0], oldWorryLevel)
	secondValue := getDynamicFunctionValue(operationParts[2], oldWorryLevel)
	worryLevelIncreaserFunction := operatorToFunction[operationParts[1]]

	return worryLevelIncreaserFunction(firstValue, secondValue)
}

func getDynamicFunctionValue(valueAsString string, oldWorryLevel int) int {
	isOldWorryLevel := valueAsString == oldWorryLevelId

	if isOldWorryLevel {
		return oldWorryLevel
	}

	value, err := strconv.Atoi(valueAsString)

	if err != nil {
		fmt.Println("Couldn't convert function value: " + valueAsString)
	}
	return value
}

func multiplyValues(value1 int, value2 int) int {
	return value1 * value2
}

func addValues(value1 int, value2 int) int {
	return value1 + value2
}

func playMonkeyGame() {
	var inspectedItemCounts []uint64
	monkeyIndices := reflect.ValueOf(indexToMonkey).MapKeys()
	totalMonkeys := len(monkeyIndices)

	for i := 0; i < rounds; i++ {
		fmt.Println("Round: ")
		playOneRound(totalMonkeys)
	}	

	for i := 0; i < totalMonkeys; i++ {
		inspectedItemCounts = append(inspectedItemCounts, indexToMonkey[i].inspectionCount)
		fmt.Println("Activity: " + strconv.Itoa(indexToMonkey[i].inspectionCount))
	}

	sort.Ints(inspectedItemCounts)
	mostActive := inspectedItemCounts[len(inspectedItemCounts)-1]
	secondMostActive := inspectedItemCounts[len(inspectedItemCounts)-2]

	fmt.Println("Most active: " + strconv.Itoa(mostActive))
	fmt.Println("Second most active: " + strconv.Itoa(secondMostActive))
	fmt.Println("Monkey business: " + strconv.Itoa(mostActive * secondMostActive))
}

func playOneRound(totalMonkeys int) {
	for i := 0; i < totalMonkeys; i++ {
		fmt.Println("Monkey " + strconv.Itoa(i) + " with currently " + strconv.Itoa(len(indexToMonkey[i].items)) + " items:")
		for j := 0; j < len(indexToMonkey[i].items); j++ {
			indexToMonkey[i].inspectionCount++
			currentItem := indexToMonkey[i].items[j]
			fmt.Println("  Monkey inspects an item with a worry level of " + strconv.Itoa(currentItem) + ".")
			currentItem = executeOperation(currentItem, indexToMonkey[i].operationParts)
			fmt.Println("  Worry level of item increased to " + strconv.Itoa(currentItem) + ".")
			currentItem = currentItem / worryLevelReducer
			fmt.Println("  Worry level is reduced to " + strconv.Itoa(currentItem) + ".")

			if currentItem % indexToMonkey[i].divisor == 0 {
				throwItemToNextMonkey(indexToMonkey[i].nextMonkeyOnTrue, currentItem)
				fmt.Println("  Test is true, throwing to monkey " + strconv.Itoa(indexToMonkey[i].nextMonkeyOnTrue) + ".")
			} else {
				throwItemToNextMonkey(indexToMonkey[i].nextMonkeyOnFalse, currentItem)
				fmt.Println("  Test is false, throwing to monkey " + strconv.Itoa(indexToMonkey[i].nextMonkeyOnFalse) + ".")
			}
		}
		indexToMonkey[i].items = nil
	}
}

func throwItemToNextMonkey(nextMonkey int, item int) {
	indexToMonkey[nextMonkey].items = append(indexToMonkey[nextMonkey].items, item)
}