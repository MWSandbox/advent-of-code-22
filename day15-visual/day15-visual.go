package main

import (
	"advent-of-code/shared"
	"fmt"
	"strconv"
	"strings"
)

type Location string

const (
	Sensor          Location = "S"
	Beacon                   = "B"
	ScannedLocation          = "#"
	UnknownLocation          = "."
)

const yPosToCheck int = 10

type Position struct {
	x int
	y int
}

var minX int = 9999999
var minY int = 9999999
var maxX int = -9999999
var maxY int = -9999999

var sensorToBeacon map[*Position]*Position = make(map[*Position]*Position)
var tunnelMap [][]Location

// Does not work with full puzzle input since solution does not scale very well
func main() {
	filePuzzle1 := shared.OpenFile("./input-for-visual.txt")
	shared.ReadFileLineByLine(filePuzzle1, readSensorData)
	initializeTunnelMap()
	scanSurroundings()
	printMap()
	resultPuzzle1 := countKnownLocationsFor(yPosToCheck)
	fmt.Println("Known locations on yPos=" + strconv.Itoa(yPosToCheck) + ": " + strconv.Itoa(resultPuzzle1))
}

func readSensorData(line string) {
	split := strings.Split(line, ": closest beacon is at ")
	sensorString := strings.Split(split[0], "Sensor at ")[1]
	beaconString := split[1]
	sensorPos := extractCoordinatesFromPosition(sensorString)
	beaconPos := extractCoordinatesFromPosition(beaconString)
	sensorToBeacon[sensorPos] = beaconPos

	updateMinAndMax(sensorPos, beaconPos)
}

func extractCoordinatesFromPosition(position string) *Position {
	coordinates := strings.Split(position, ", ")
	posX, errX := strconv.Atoi(coordinates[0][2:])
	posY, errY := strconv.Atoi(coordinates[1][2:])

	if errX != nil || errY != nil {
		fmt.Println("Couldn't read coordinates: " + position)
	}
	return &Position{x: posX, y: posY}
}

func updateMinAndMax(sensorPos *Position, beaconPos *Position) {
	updateSingleMinAndMax(sensorPos)
	updateSingleMinAndMax(beaconPos)
}

func updateSingleMinAndMax(pos *Position) {
	if pos.x < minX {
		minX = pos.x
	}

	if pos.y < minY {
		minY = pos.y
	}

	if pos.x > maxX {
		maxX = pos.x
	}

	if pos.y > maxY {
		maxY = pos.y
	}
}

func initializeTunnelMap() {
	tunnelMap = make([][]Location, maxX-minX+1)

	for x := 0; x <= maxX-minX; x++ {
		tunnelMap[x] = make([]Location, maxY-minY+1)

		for y := 0; y <= maxY-minY; y++ {
			tunnelMap[x][y] = UnknownLocation
		}
	}
}

func scanSurroundings() {
	for sensor, beacon := range sensorToBeacon {
		tunnelMap[sensor.x+abs(minX)][sensor.y+abs(minY)] = Sensor
		tunnelMap[beacon.x+abs(minX)][beacon.y+abs(minY)] = Beacon

		distX := abs(sensor.x - beacon.x)
		distY := abs(sensor.y - beacon.y)
		totalDist := distX + distY

		for x := sensor.x - totalDist; x < sensor.x+totalDist; x++ {

			distLeftForY := totalDist - abs(x-sensor.x)

			for y := sensor.y - distLeftForY; y < sensor.y+distLeftForY; y++ {

				if x < minX {
					extendTunnelMapXValues(x - minX)
					minX = x
				} else if x > maxX {
					extendTunnelMapXValues(x - maxX)
					maxX = x
				}

				if y < minY {
					extendTunnelMapYValues(y - minY)
					minY = y
				} else if y > maxY {
					extendTunnelMapYValues(y - minY)
					maxY = y
				}

				if tunnelMap[x+abs(minX)][y+abs(minY)] == UnknownLocation {
					tunnelMap[x+abs(minX)][y+abs(minY)] = ScannedLocation
				}
			}
		}
	}
}

func extendTunnelMapXValues(extension int) {
	for i := 0; i < abs(extension); i++ {
		sliceToAdd := make([]Location, maxY-minY+1)

		for y := 0; y <= maxY-minY; y++ {
			sliceToAdd[y] = UnknownLocation
		}

		tunnelMap = append(tunnelMap, sliceToAdd)

		if extension < 0 {
			tunnelMap = append([][]Location{sliceToAdd}, tunnelMap...)
		}
	}
}

func extendTunnelMapYValues(extension int) {
	for i := 0; i < abs(extension); i++ {
		for j := 0; j < len(tunnelMap); j++ {
			tunnelMap[j] = append(tunnelMap[j], UnknownLocation)

			if extension < 0 {
				tunnelMap[j] = append([]Location{UnknownLocation}, tunnelMap[j]...)
			}
		}
	}
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func printMap() {
	shared.PrintXHeader(minX, maxX)

	for y := 0; y < maxY-minY; y++ {
		var printedLine string = strconv.Itoa(abs(y+minY)) + shared.GetYPadding(y+minY)
		for x := 0; x < maxX-minX; x++ {
			printedLine += string(tunnelMap[x][y])
		}
		fmt.Println(printedLine)
	}
}

func countKnownLocationsFor(yPosToCheck int) int {
	count := 0

	for x := 0; x < len(tunnelMap); x++ {
		if tunnelMap[x][yPosToCheck+abs(minY)] != UnknownLocation {
			count++
		}
	}

	return count
}
