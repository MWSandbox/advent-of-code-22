package main

import (
	"advent-of-code/shared"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Path struct {
	nodes            []*shared.GraphNode
	releasedPressure int
	closedValves     []*shared.GraphNode
	minutesLeft      int
}

const pathsMinValue int = 914 //half of first puzzles result to improve puzzle 2 performance

var totalMinutes int = 30
var rootNode *shared.GraphNode
var valveToNode map[string]*shared.GraphNode = make(map[string]*shared.GraphNode)
var nodeToFlowRate map[*shared.GraphNode]int = make(map[*shared.GraphNode]int)
var nodeToValve map[*shared.GraphNode]string = make(map[*shared.GraphNode]string)
var valves []string
var closedValves []*shared.GraphNode
var maxWaterOutput int
var distancesBetweenNodes [][]int
var highestPressureReleased int
var allPaths []*Path

func main() {
	filePuzzle1 := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle1, parseGraph)
	setIndices()
	calculateMaxWaterOutput()
	findShortestDistancesBetweenRelevantNodes()
	dfs()
	fmt.Println("Pressure released in 30 min.: " + strconv.Itoa(highestPressureReleased))

	resetForPuzzle2()
	dfs()
	findBestPathToWalkWithElephant()
}

func resetForPuzzle2() {
	highestPressureReleased = 0
	totalMinutes = 26
	allPaths = make([]*Path, 0)
}

func parseGraph(line string) {
	valve := line[6:8]
	valves = append(valves, valve)
	flowRateString := strings.Split(line[23:], ";")[0]

	flowRate, err := strconv.Atoi(flowRateString)

	if err != nil {
		fmt.Println("Couldn't convert flow rate: " + flowRateString)
	}

	neighborSplit := strings.Split(line, "to valves ")
	if len(neighborSplit) == 1 {
		neighborSplit = strings.Split(line, "to valve ")
	}
	neighbors := strings.Split(neighborSplit[1], ", ")
	graphNode := valveToNode[valve]

	if graphNode == nil {
		graphNode = &shared.GraphNode{}
		graphNode.Neighbors = append(graphNode.Neighbors, graphNode)
		valveToNode[valve] = graphNode
	}
	nodeToFlowRate[graphNode] = flowRate
	nodeToValve[graphNode] = valve

	if flowRate > 0 {
		closedValves = append(closedValves, graphNode)
	}

	for i := 0; i < len(neighbors); i++ {
		neighborValve := neighbors[i]
		neighborNode := valveToNode[neighborValve]

		if neighborNode == nil {
			neighborNode = &shared.GraphNode{}
			neighborNode.Neighbors = append(neighborNode.Neighbors, neighborNode)
			valveToNode[neighborValve] = neighborNode
		}

		graphNode.Neighbors = append(graphNode.Neighbors, neighborNode)
	}

	if rootNode == nil {
		rootNode = graphNode
	}
}

func setIndices() {
	sort.Strings(valves)

	for i := 0; i < len(valves); i++ {
		valveToNode[valves[i]].Index = i
	}
}

func calculateMaxWaterOutput() {
	for i := 0; i < len(closedValves); i++ {
		maxWaterOutput += 30 * nodeToFlowRate[closedValves[i]]
	}
}

func dfs() {
	var startNodes []*shared.GraphNode
	startNodes = append(startNodes, valveToNode["AA"])
	var closedStartValves []*shared.GraphNode
	closedStartValves = append(closedStartValves, closedValves...)
	startPath := &Path{nodes: startNodes, releasedPressure: 0, closedValves: closedStartValves, minutesLeft: totalMinutes}

	var pathStack []*Path
	pathStack = append(pathStack, startPath)

	for len(pathStack) > 0 {
		currentPath := pathStack[len(pathStack)-1]
		pathStack = pathStack[:len(pathStack)-1]

		if currentPath.releasedPressure > highestPressureReleased {
			highestPressureReleased = currentPath.releasedPressure
		}

		for i := 0; i < len(currentPath.closedValves); i++ {
			valveToConsider := currentPath.closedValves[i]
			fromStartNode := currentPath.nodes[len(currentPath.nodes)-1]
			minutesToOpenValve := distancesBetweenNodes[fromStartNode.Index][valveToConsider.Index] + 1

			if minutesToOpenValve <= currentPath.minutesLeft {
				var newNodes []*shared.GraphNode
				newNodes = append(newNodes, currentPath.nodes...)
				newNodes = append(newNodes, valveToConsider)

				waterOutputUntilEnd := currentPath.releasedPressure + ((currentPath.minutesLeft - minutesToOpenValve) * nodeToFlowRate[valveToConsider])

				var newClosedValves []*shared.GraphNode
				newClosedValves = append(newClosedValves, currentPath.closedValves...)

				indexOfValve := -1
				for i := 0; i < len(newClosedValves); i++ {
					if newClosedValves[i] == valveToConsider {
						indexOfValve = i
					}
				}

				newClosedValves = append(newClosedValves[:indexOfValve], newClosedValves[indexOfValve+1:]...)

				newPath := &Path{nodes: newNodes, releasedPressure: waterOutputUntilEnd, closedValves: newClosedValves, minutesLeft: currentPath.minutesLeft - minutesToOpenValve}
				pathStack = append(pathStack, newPath)
				allPaths = append(allPaths, newPath)
			}
		}
	}
}

func printPath(path *Path) {
	printablePath := ""

	for i := 0; i < len(path.nodes); i++ {
		printablePath += nodeToValve[path.nodes[i]]

		if i != len(path.nodes)-1 {
			printablePath += " -> "
		}
	}

	fmt.Println("Path: " + printablePath)
}

func findShortestDistancesBetweenRelevantNodes() {
	distancesBetweenNodes = make([][]int, len(valves))

	for i := 0; i < len(valves); i++ {
		distancesBetweenNodes[i] = make([]int, len(valves))
		for j := 0; j < len(valves); j++ {
			nodeI := valveToNode[valves[i]]
			nodeJ := valveToNode[valves[j]]

			if i != j && (isNodeClosedValve(nodeI) && isNodeClosedValve(nodeJ)) || (i == 0) {
				distance := findShortestPathFromTo(nodeI, nodeJ)
				distancesBetweenNodes[i][j] = distance
			}
		}
	}
}

func isNodeClosedValve(node *shared.GraphNode) bool {
	for i := 0; i < len(closedValves); i++ {
		if closedValves[i] == node {
			return true
		}
	}
	return false
}

func findShortestPathFromTo(start *shared.GraphNode, end *shared.GraphNode) int {
	var allNodes []*shared.GraphNode

	for i := 0; i < len(valves); i++ {
		allNodes = append(allNodes, valveToNode[valves[i]])
	}

	return shared.Dijkstra(allNodes, start, end)
}

func findBestPathToWalkWithElephant() {
	maxPressure := 0
	var ownPath *Path
	var elephantPath *Path

	filterAllPaths()

	for i := 0; i < len(allPaths); i++ {
		for j := i + 1; j < len(allPaths); j++ {
			pressureSum := allPaths[i].releasedPressure + allPaths[j].releasedPressure
			if pressureSum > maxPressure && !doPathsContainSameValves(allPaths[i], allPaths[j]) {
				maxPressure = pressureSum
				ownPath = allPaths[i]
				elephantPath = allPaths[j]
			}
		}
	}

	fmt.Println("Total pressure released with elephant: " + strconv.Itoa(maxPressure))

	fmt.Println("Own Path:")
	printPath(ownPath)

	fmt.Println("Elephant Path:")
	printPath(elephantPath)
}

func doPathsContainSameValves(path1 *Path, path2 *Path) bool {
	var nodesPath1 map[int]bool = make(map[int]bool)

	for i := 1; i < len(path1.nodes); i++ {
		nodesPath1[path1.nodes[i].Index] = true
	}

	for i := 1; i < len(path2.nodes); i++ {
		if nodesPath1[path2.nodes[i].Index] {
			return true
		}
	}

	return false
}

func filterAllPaths() {
	var trimmedPaths []*Path

	for i := 0; i < len(allPaths); i++ {
		if allPaths[i].releasedPressure > pathsMinValue {
			trimmedPaths = append(trimmedPaths, allPaths[i])
		}
	}

	allPaths = trimmedPaths
}
