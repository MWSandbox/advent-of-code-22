package main

import (
	"advent-of-code/shared"
	"fmt"
	"log"
	"sort"
	"strconv"
)

var totalCalories []int
var currentCalories int = 0

func main() {
	file := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(file, calculateCalories)
	printOutNHighestCalories(totalCalories, 3)
}

func calculateCalories(line string) {

	if line == "" {
		totalCalories = append(totalCalories, currentCalories)
		currentCalories = 0
	} else {
		caloriesAsInt, err := strconv.Atoi(line)
		currentCalories += caloriesAsInt

		if err != nil {
			log.Fatal("Could not parse int: ")
		}
	}
}

func printOutNHighestCalories(totalCalories []int, n int) {
	var sum int = 0
	var currentCalories int = 0

	sort.Ints(totalCalories)

	for i := 1; i <= n; i++ {
		currentCalories = totalCalories[len(totalCalories)-i]
		sum += currentCalories
		fmt.Println(strconv.Itoa(i) + ". highest: " + strconv.Itoa(currentCalories))
	}

	fmt.Println("Total Highest calories: " + strconv.Itoa(sum))
}
