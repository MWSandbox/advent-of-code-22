package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	file := openFile("./input.txt")
	totalCalories := readCaloriesFromFile(file)
	printOutNHighestCalories(totalCalories, 3)
}

func openFile(filePath string) *os.File {
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}
	return file
}

func readCaloriesFromFile(file *os.File) []int {
	scanner := bufio.NewScanner(file)
	var currentCalories int = 0
	var totalCalories []int

	for scanner.Scan() {

		var currentLine string = scanner.Text()

		if currentLine == "" {
			totalCalories = append(totalCalories, currentCalories)
			currentCalories = 0
		} else {
			caloriesAsInt, err := strconv.Atoi(currentLine)
			currentCalories += caloriesAsInt

			if err != nil {
				log.Fatal("Could not parse int: ")
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	return totalCalories
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
