package main

import (
	"advent-of-code/shared"
	"fmt"
	"strconv"
)

type Graph struct {
	startNode *shared.GraphNode
	endNode   *shared.GraphNode
}

var graph *Graph = &Graph{}
var yPos int = 0
var heightMap [][]*shared.GraphNode = make([][]*shared.GraphNode, 0)
var nodeIdToHeight map[int]int = make(map[int]int)
var shortestPath int = 9999999
var nodeIdToX map[int]int = make(map[int]int)
var nodeIdToY map[int]int = make(map[int]int)

func main() {
	filePuzzle := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle, parseGraph)

	solvePuzzle1()
	solvePuzzle2()
}

func solvePuzzle1() {
	shortestPath = searchShortestPath([]*shared.GraphNode{graph.startNode})
	fmt.Println("Shortest path from S to E: " + strconv.Itoa(shortestPath))
}

func solvePuzzle2() {
	var startNodes []*shared.GraphNode

	for y := 0; y < len(heightMap); y++ {
		for x := 0; x < len(heightMap[y]); x++ {
			node := heightMap[y][x]
			if nodeIdToHeight[node.Index] == int('a') {
				startNodes = append(startNodes, node)
			}
		}
	}

	shortestPath = searchShortestPath(startNodes)
	fmt.Println("Shortest path from any a spot: " + strconv.Itoa(shortestPath))
}

func parseGraph(line string) {
	heightMap = append(heightMap, make([]*shared.GraphNode, 0))

	for i := 0; i < len(line); i++ {
		currentChar := rune(line[i])
		var currentHeight int
		isStart := currentChar == 'S'
		isEnd := currentChar == 'E'

		if isStart {
			currentHeight = int('a')
		} else if isEnd {
			currentHeight = int('z')
		} else {
			currentHeight = int(currentChar)
		}

		nodeId := shared.CantorHashTwoFields(i, yPos)
		nodeIdToX[nodeId] = i
		nodeIdToY[nodeId] = yPos

		nodeIdToHeight[nodeId] = currentHeight
		currentNode := &shared.GraphNode{Index: nodeId, Neighbors: make([]*shared.GraphNode, 0)}
		heightMap[yPos] = append(heightMap[yPos], currentNode)

		if isStart {
			graph.startNode = currentNode
		} else if isEnd {
			graph.endNode = currentNode
		}

		updateNeighbors(i, yPos, currentNode)
	}
	yPos++
}

func updateNeighbors(currentX int, currentY int, currentNode *shared.GraphNode) {
	if isNeighborReachable(currentX, currentY, currentX-1, currentY) {
		leftNode := heightMap[currentY][currentX-1]
		currentNode.Neighbors = append(currentNode.Neighbors, leftNode)
	}

	if isNeighborReachable(currentX-1, currentY, currentX, currentY) {
		leftNode := heightMap[currentY][currentX-1]
		leftNode.Neighbors = append(leftNode.Neighbors, currentNode)
	}

	if isNeighborReachable(currentX, currentY, currentX, currentY-1) {
		upperNode := heightMap[currentY-1][currentX]
		currentNode.Neighbors = append(currentNode.Neighbors, upperNode)
	}

	if isNeighborReachable(currentX, currentY-1, currentX, currentY) {
		upperNode := heightMap[currentY-1][currentX]
		upperNode.Neighbors = append(upperNode.Neighbors, currentNode)
	}
}

func isNeighborReachable(currentX int, currentY, targetX int, targetY int) bool {
	if isInRange(currentX, currentY) && isInRange(targetX, targetY) {
		targetId := heightMap[targetY][targetX].Index
		currentId := heightMap[currentY][currentX].Index
		neighborHeight := nodeIdToHeight[targetId]
		currentHeight := nodeIdToHeight[currentId]

		if neighborHeight-currentHeight <= 1 {
			return true
		}
	}
	return false
}

func isInRange(x int, y int) bool {
	return x >= 0 && y >= 0 && y < len(heightMap) && x < len(heightMap[0])
}

func searchShortestPath(startNodes []*shared.GraphNode) int {
	var allNodes []*shared.GraphNode

	for i := 0; i < len(heightMap); i++ {
		allNodes = append(allNodes, heightMap[i]...)
	}

	return shared.DijkstraWithMultiStart(allNodes, startNodes, graph.endNode)
}
