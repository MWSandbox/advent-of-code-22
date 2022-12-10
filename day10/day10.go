package main

import (
	"advent-of-code/shared"
	"strconv"
	"fmt"
	"strings"
)

const lineLength = 40
const tick int = 40
const minTick int = 20
const maxCycle int = 220

var cycle int = 0
var register int = 1
var signalStrength int = 0
var currentPrintRow string = ""
var currentRowIndex = 0

func main() {
	filePuzzle := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle, processSignal)
	fmt.Println("Signal strength: " + strconv.Itoa(signalStrength))
}

func processSignal(line string) {
	command := strings.Split(line, " ")
	instruction := command[0]

	if instruction == "addx" {
		for i := 0; i < 2; i++ {
			runCycle()
		}

		changeRegister(command[1])
	} else {
		runCycle()
	}
}

func runCycle() {
	drawPixel()
	cycle++
	printIfEndOfLine()
	checkTick()
}

func drawPixel() {
	drawPosition := cycle - (currentRowIndex * lineLength)
	isSpriteDrawn := drawPosition >= register - 1 && drawPosition <= register + 1

	if isSpriteDrawn {
		currentPrintRow += "#"
	} else {
		currentPrintRow += "."
	}
}

func printIfEndOfLine() {
	if cycle % lineLength == 0 {
		fmt.Println(currentPrintRow)
		currentPrintRow = ""
		currentRowIndex++
	}
}

func checkTick() {
	if (cycle % tick) - minTick == 0 && cycle <= maxCycle {
		signalStrength += (cycle * register)
	}
}

func changeRegister(parameter string) {
	value, err := strconv.Atoi(parameter)

	if err != nil {
		fmt.Println("Couldn't convert addx value " + parameter)
	}
	register += value
}