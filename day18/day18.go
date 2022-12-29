package main

import (
	"advent-of-code/shared"
	"fmt"
	"strconv"
	"strings"
)

const cubeDirections int = 6
const infinity int = 9999999

var xToYToZCubes map[int]map[int]map[int]bool = make(map[int]map[int]map[int]bool)
var totalSurfaceArea int = 0
var cubeHashes map[int]bool = make(map[int]bool)
var trappedAirHashes map[int]bool = make(map[int]bool)
var maxCoordinates [3]int = [3]int{0, 0, 0}

func main() {
	filePuzzle := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle, parseCubes)

	countSurfaceArea()
	fmt.Println("Total cubes: " + strconv.Itoa(len(cubeHashes)))
	fmt.Println("Surface area: " + strconv.Itoa(totalSurfaceArea))

	totalSurfaceArea = 0
	findTrappedAir()
	countSurfaceArea()
	fmt.Println("Total trapped air drops: " + strconv.Itoa(len(trappedAirHashes)))
	fmt.Println("External surface area: " + strconv.Itoa(totalSurfaceArea))
}

func parseCubes(line string) {
	stringCoordinates := strings.Split(line, ",")
	intCoordinates := coordinatesAsInt(stringCoordinates)

	if xToYToZCubes[intCoordinates[0]] == nil {
		xToYToZCubes[intCoordinates[0]] = make(map[int]map[int]bool)
	}

	if xToYToZCubes[intCoordinates[0]][intCoordinates[1]] == nil {
		xToYToZCubes[intCoordinates[0]][intCoordinates[1]] = make(map[int]bool)
	}

	xToYToZCubes[intCoordinates[0]][intCoordinates[1]][intCoordinates[2]] = true
	cubeHashes[hashCoordinates(intCoordinates)] = true

	for i := 0; i < len(maxCoordinates); i++ {
		if intCoordinates[i] > maxCoordinates[i] {
			maxCoordinates[i] = intCoordinates[i]
		}
	}
}

func coordinatesAsInt(stringCoordinates []string) [3]int {
	var coordinatesAsInt [3]int

	for i := 0; i < 3; i++ {
		singleCoordinate, err := strconv.Atoi(stringCoordinates[i])

		if err != nil {
			fmt.Println("Couldn't parse coordinates: " + stringCoordinates[i])
		}
		coordinatesAsInt[i] = singleCoordinate
	}

	return coordinatesAsInt
}

func countSurfaceArea() {
	for x, _ := range xToYToZCubes {
		for y, _ := range xToYToZCubes[x] {
			for z, _ := range xToYToZCubes[x][y] {
				cube := [3]int{x, y, z}
				connectedAirDrops := countAllConnectedTrappedAirDropsOf(cube)
				neighborCount := len(findAllNeighbors(cube))
				totalSurfaceArea += cubeDirections - neighborCount - connectedAirDrops
			}
		}
	}
}

func countAllConnectedTrappedAirDropsOf(cube [3]int) int {
	count := 0
	neighborPositions := generateAllNeighborPositionsOfCube(cube)

	for i := 0; i < len(neighborPositions); i++ {
		if trappedAirHashes[hashCoordinates(neighborPositions[i])] {
			count++
		}
	}

	return count
}

func generateAllNeighborPositionsOfCube(cube [3]int) [][3]int {
	var neighborPositions [][3]int

	for i := 0; i < len(cube); i++ {
		for j := -1; j < 2; j += 2 {
			neighborCube := [3]int{cube[0], cube[1], cube[2]}
			neighborCube[i] += j
			neighborPositions = append(neighborPositions, neighborCube)
		}
	}

	return neighborPositions
}

func findAllNeighbors(cube [3]int) [][3]int {
	var actualNeighbors [][3]int
	neighborPositions := generateAllNeighborPositionsOfCube(cube)

	for i := 0; i < len(neighborPositions); i++ {
		neighbor := neighborPositions[i]
		if xToYToZCubes[neighbor[0]][neighbor[1]][neighbor[2]] {
			actualNeighbors = append(actualNeighbors, neighbor)
		}
	}

	return actualNeighbors
}

func findTrappedAir() {
	for x, _ := range xToYToZCubes {
		for y, _ := range xToYToZCubes[x] {
			for z, _ := range xToYToZCubes[x][y] {
				checkIfAirDropsAroundCubeAreTrapped([3]int{x, y, z})
			}
		}
	}
}

func checkIfAirDropsAroundCubeAreTrapped(cube [3]int) {
	neighborPositions := generateAllNeighborPositionsOfCube(cube)

	for i := 0; i < len(neighborPositions); i++ {
		neighbor := neighborPositions[i]
		if !xToYToZCubes[neighbor[0]][neighbor[1]][neighbor[2]] {
			checkIfAirDropIsTrapped(neighborPositions[i])
		}
	}
}

func checkIfAirDropIsTrapped(airDrop [3]int) {
	hashValue := hashCoordinates(airDrop)
	neighborsOfAirDrop := findAllNeighbors(airDrop)
	isSurroundedByCubes := len(neighborsOfAirDrop) == cubeDirections

	if isSurroundedByCubes {
		if !trappedAirHashes[hashValue] {
			trappedAirHashes[hashValue] = true
		}
	} else if len(neighborsOfAirDrop) >= 3 && !trappedAirHashes[hashValue] {
		isAirChamber, airChamber := checkForAirChamber(airDrop, neighborsOfAirDrop, make(map[int]bool))

		if isAirChamber {
			fmt.Println("Found Air chamber of size " + strconv.Itoa(len(airChamber)))

			for hash, _ := range airChamber {
				trappedAirHashes[hash] = true
			}
		}
	}
}

func checkForAirChamber(airDrop [3]int, neighbors [][3]int, airChamber map[int]bool) (bool, map[int]bool) {
	boundaries := findAllBoundaries(airDrop, neighbors)
	isAirDropTrapped := isAirDropTrapped(boundaries)

	if isAirDropTrapped {
		airChamber[hashCoordinates(airDrop)] = true
		var diffsPerAxis []int = []int{airDrop[0] - boundaries[0], boundaries[1] - airDrop[0]}
		diffsPerAxis = append(diffsPerAxis, []int{airDrop[1] - boundaries[2], boundaries[3] - airDrop[1]}...)
		diffsPerAxis = append(diffsPerAxis, []int{airDrop[2] - boundaries[4], boundaries[5] - airDrop[2]}...)

		for i := 0; i < len(diffsPerAxis); i++ {
			for j := 1; diffsPerAxis[i]-j > 0; j++ {

				nextAirDrop := [3]int{airDrop[0], airDrop[1], airDrop[2]}

				if i%2 == 0 {
					nextAirDrop[i/2] -= j
				} else {
					nextAirDrop[i/2] += j
				}

				if !airChamber[hashCoordinates(nextAirDrop)] {
					isNextAirDropTrapped, nextAirChamber := checkForAirChamber(nextAirDrop, findAllNeighbors(nextAirDrop), airChamber)

					for hash, _ := range nextAirChamber {
						airChamber[hash] = true
					}
					isAirDropTrapped = isAirDropTrapped && isNextAirDropTrapped

					if !isAirDropTrapped {
						return false, nil
					}
				}
			}
		}
	}

	return isAirDropTrapped, airChamber
}

func findAllBoundaries(airDrop [3]int, neighbors [][3]int) [6]int {
	var boundaries [6]int = [6]int{infinity, infinity, infinity, infinity, infinity, infinity}

	for i := 0; i < len(neighbors); i++ {
		var diffs [3]int = [3]int{airDrop[0] - neighbors[i][0], airDrop[1] - neighbors[i][1], airDrop[2] - neighbors[i][2]}

		for j := 0; j < len(boundaries); j++ {
			if j%2 == 0 && diffs[j/2] > 0 {
				boundaries[j] = neighbors[i][j/2]
			} else if j%2 == 1 && diffs[j/2] < 0 {
				boundaries[j] = neighbors[i][j/2]
			}
		}
	}

	for i := 0; i < len(boundaries); i++ {
		if boundaries[i] == infinity {
			var startValue int
			newCoordinates := [3]int{airDrop[0], airDrop[1], airDrop[2]}
			if i%2 == 0 {
				startValue = airDrop[i/2] - 1
				for adjustedCoordinate := startValue; adjustedCoordinate >= 0 && boundaries[i] == infinity; adjustedCoordinate-- {
					boundaries[i] = findNextBoundary(newCoordinates, i/2, adjustedCoordinate)
				}
			} else {
				startValue = airDrop[i/2] + 1
				for adjustedCoordinate := startValue; adjustedCoordinate <= maxCoordinates[i/2] && boundaries[i] == infinity; adjustedCoordinate++ {
					boundaries[i] = findNextBoundary(newCoordinates, i/2, adjustedCoordinate)
				}
			}
		}
	}

	return boundaries
}

func findNextBoundary(newCoordinates [3]int, coordinateIndex int, adjustedCoordinate int) int {
	newCoordinates[coordinateIndex] = adjustedCoordinate
	if xToYToZCubes[newCoordinates[0]][newCoordinates[1]][newCoordinates[2]] {
		return newCoordinates[coordinateIndex]
	}
	return infinity
}

func isAirDropTrapped(boundaries [6]int) bool {
	isAirDropTrapped := true
	for i := 0; i < len(boundaries); i++ {
		if boundaries[i] == infinity {
			isAirDropTrapped = false
		}
	}
	return isAirDropTrapped
}

func hashCoordinates(coordinates [3]int) int {
	hashValue := 3 * 987 * coordinates[0]
	hashValue += 7 * 654 * coordinates[1]
	hashValue += 13 * 854 * coordinates[2]
	return hashValue
}
