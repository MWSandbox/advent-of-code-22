package main

import (
	"advent-of-code/shared"
	"fmt"
	"math"
	"strconv"
)

var snafuNumbers [][]int = make([][]int, 0)
var decimalNumbers []int

func main() {
	filePuzzle := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle, parseFuelRequirements)

	decimalSum := convertAllSnafuNumbers()
	printTable(decimalSum)
}

func parseFuelRequirements(line string) {
	var digits []int

	for i := len(line) - 1; i >= 0; i-- {
		digitAsString := line[i : i+1]
		var digit int

		if digitAsString == "-" {
			digit = -1
		} else if digitAsString == "=" {
			digit = -2
		} else {
			digit = shared.ConvertStringToInt(digitAsString)
		}

		digits = append(digits, digit)
	}

	snafuNumbers = append(snafuNumbers, digits)
}

func convertAllSnafuNumbers() int {
	sum := 0

	for i := 0; i < len(snafuNumbers); i++ {
		decimalNumber := convertSnafuToDecimal(snafuNumbers[i])
		decimalNumbers = append(decimalNumbers, decimalNumber)
		sum += decimalNumber
	}

	return sum
}

func convertSnafuToDecimal(snafuNumber []int) int {
	decimalNumber := 0

	for i := 0; i < len(snafuNumber); i++ {
		digit := snafuNumber[i]
		decimalValue := digit * int(math.Pow(5, float64(i)))
		decimalNumber += decimalValue
	}

	return decimalNumber
}

func convertDecimalToSnafu(decimalNumber int) []int {
	var snafuNumber []int
	baseFiveString := strconv.FormatInt(int64(decimalNumber), 5)

	for i := len(baseFiveString) - 1; i >= 0; i-- {
		snafuNumber = append(snafuNumber, shared.ConvertStringToInt(baseFiveString[i:i+1]))
	}

	for i := 0; i < len(snafuNumber); i++ {
		digit := snafuNumber[i]

		if digit > 2 {
			diff := 5 - digit
			snafuNumber[i] = -diff

			if i == len(snafuNumber)-1 {
				snafuNumber = append(snafuNumber, 0)
			}

			snafuNumber[i+1]++
		}
	}
	return snafuNumber
}

func snafuNumberAsString(snafuNumber []int) string {
	text := ""

	for i := len(snafuNumber) - 1; i >= 0; i-- {
		digit := snafuNumber[i]

		if digit == -1 {
			text += "-"
		} else if digit == -2 {
			text += "="
		} else {
			text += strconv.Itoa(digit)
		}
	}
	return text
}

func printTable(sum int) {
	snafuSum := convertDecimalToSnafu(sum)
	padding := len(snafuSum) + 1

	sumRow := getTableRow(convertDecimalToSnafu(sum), sum, padding)
	separator := ""
	for i := 0; i < len(sumRow); i++ {
		separator += "-"
	}

	header := "SNAFU"
	header = addPaddingToText(header, padding)
	header += "DECIMAL"
	fmt.Println(header)
	fmt.Println(separator)

	for i := 0; i < len(snafuNumbers); i++ {
		fmt.Println(getTableRow(snafuNumbers[i], decimalNumbers[i], padding))
	}

	fmt.Println(separator)
	fmt.Println(sumRow)
}

func addPaddingToText(text string, maxPadding int) string {
	padding := ""

	for i := 0; i < maxPadding-len(text); i++ {
		padding += " "
	}
	return text + padding
}

func getTableRow(snafuNumber []int, decimalNumber int, padding int) string {
	row := snafuNumberAsString(snafuNumber)
	row = addPaddingToText(row, padding)
	row += strconv.Itoa(decimalNumber)
	return row
}
