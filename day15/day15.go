package main

import (
	"advent-of-code/shared"
	"fmt"
	"strconv"
	"strings"
)

const targetPosY int = 2000000
const distressSignalMin int = 0
const distressSignalMax int = 4000000

type Position struct {
	x int
	y int
}

var sensorToSignalStrength map[*Position]int = make(map[*Position]int)
var vectorsCovered []*Position
var beacons []*Position

func main() {
	filePuzzle1 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle1, readSensorData)
	count := countFieldsWithoutBeaconOnTarget()
	fmt.Println("Known location on yPos=" + strconv.Itoa(targetPosY) + ": " + strconv.Itoa(count))

	filePuzzle2 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle2, readSensorData)
	distressX, distressY := checkFieldsForDistressSignal()
	fmt.Println("Tuning frequency: " + strconv.Itoa(distressSignalMax*distressX+distressY))
}

func readSensorData(line string) {
	split := strings.Split(line, ": closest beacon is at ")
	sensorString := strings.Split(split[0], "Sensor at ")[1]
	beaconString := split[1]
	sensorPos := extractCoordinatesFromPosition(sensorString)
	beaconPos := extractCoordinatesFromPosition(beaconString)

	if !isBeaconAlreadyPresent(beaconPos) {
		beacons = append(beacons, beaconPos)
	}

	distX := abs(sensorPos.x - beaconPos.x)
	distY := abs(sensorPos.y - beaconPos.y)
	signalStrength := distX + distY
	sensorToSignalStrength[sensorPos] = signalStrength
}

func isBeaconAlreadyPresent(beaconPos *Position) bool {
	for i := 0; i < len(beacons); i++ {
		if beaconPos.x == beacons[i].x && beaconPos.y == beacons[i].y {
			return true
		}
	}
	return false
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

func countFieldsWithoutBeaconOnTarget() int {
	scanVectors(-9999999, 9999999, targetPosY)
	beaconsOnTarget := countScannedBeaconsOnTarget()
	scannedFieldsOnTarget := countScannedFieldsOnTarget()
	return scannedFieldsOnTarget - beaconsOnTarget
}

func checkFieldsForDistressSignal() (int, int) {
	for i := distressSignalMin; i <= distressSignalMax; i++ {
		scanVectors(distressSignalMin, distressSignalMax, i)
		scannedFieldsOnTarget := countScannedFieldsOnTarget()

		if scannedFieldsOnTarget != distressSignalMax+1 {
			return findDistressSignalXPos(), i
		}
	}
	return 0, 0
}

func scanVectors(minX int, maxX int, yPos int) {
	vectorsCovered = make([]*Position, 0)

	for sensor, signalStrength := range sensorToSignalStrength {
		distToYPos := abs(yPos - sensor.y)

		if distToYPos <= signalStrength {
			remainingRange := signalStrength - distToYPos
			newVectorMin := sensor.x - remainingRange
			newVectorMax := sensor.x + remainingRange

			if newVectorMin < minX {
				newVectorMin = minX
			}

			if newVectorMax > maxX {
				newVectorMax = maxX
			}

			newVector := &Position{newVectorMin, newVectorMax}
			vectorsCovered = append(vectorsCovered, newVector)
		}
	}
	removeVectorOverlaps()
}

func removeVectorOverlaps() {
	isOverlapDetected := false

	for i := 0; i < len(vectorsCovered) && !isOverlapDetected; i++ {
		isOverlapDetected = checkOverlapOfVector(vectorsCovered[i], i)
	}

	if isOverlapDetected {
		removeVectorOverlaps()
	}
}

func checkOverlapOfVector(vector *Position, index int) bool {
	isOverlap := false

	for i := 0; i < len(vectorsCovered) && !isOverlap; i++ {
		if vector != vectorsCovered[i] {
			otherVector := vectorsCovered[i]
			isYInOtherVector := vector.y >= otherVector.x && vector.y <= otherVector.y
			isXInOtherVector := vector.x >= otherVector.x && vector.x <= otherVector.y
			isLeftOverlap := isYInOtherVector && !isXInOtherVector
			isRightOverlap := isXInOtherVector && !isYInOtherVector
			isFullyIncluded := isXInOtherVector && isYInOtherVector
			doesFullyInclude := !isXInOtherVector && !isYInOtherVector && vector.x <= otherVector.x && vector.y >= otherVector.y
			isOverlap = isLeftOverlap || isRightOverlap || isFullyIncluded || doesFullyInclude

			if isLeftOverlap {
				otherVector.x = vector.x
				removeVectorAtIndex(index)
			} else if isRightOverlap {
				otherVector.y = vector.y
				removeVectorAtIndex(index)
			} else if isFullyIncluded {
				removeVectorAtIndex(index)
			} else if doesFullyInclude {
				removeVectorAtIndex(i)
			}
		}
	}
	return isOverlap
}

func removeVectorAtIndex(index int) {
	vectorsCovered = append(vectorsCovered[:index], vectorsCovered[index+1:]...)
}

func countScannedBeaconsOnTarget() int {
	count := 0

	for i := 0; i < len(beacons); i++ {
		if beacons[i].y == targetPosY {
			for j := 0; j < len(vectorsCovered); j++ {
				if beacons[i].x >= vectorsCovered[j].x && beacons[i].x <= vectorsCovered[j].y {
					count++
					break
				}
			}
		}
	}
	return count
}

func countScannedFieldsOnTarget() int {
	count := 0

	for i := 0; i < len(vectorsCovered); i++ {
		count += vectorsCovered[i].y - vectorsCovered[i].x + 1
	}
	return count
}

func findDistressSignalXPos() int {
	if len(vectorsCovered) == 1 {
		if vectorsCovered[0].x > distressSignalMin {
			return distressSignalMin
		} else {
			return distressSignalMax
		}
	}

	if vectorsCovered[0].x < vectorsCovered[1].x {
		return vectorsCovered[0].y + 1
	} else {
		return vectorsCovered[1].y + 1
	}
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
