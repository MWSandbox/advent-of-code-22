package main

import (
	"advent-of-code/shared"
	"fmt"
	"strconv"
	"container/heap"
	"sort"
)

type Item struct {
	value *GraphNode
	priority int
	index int
	// predecessor *Item
}

type PriorityQueue []*Item

func (priorityQueue PriorityQueue) Len() int { 
	return len(priorityQueue) 
}

func (priorityQueue PriorityQueue) Less(i, j int) bool {
	return priorityQueue[i].priority < priorityQueue[j].priority
}

func (priorityQueue PriorityQueue) Swap(i, j int) {
	priorityQueue[i], priorityQueue[j] = priorityQueue[j], priorityQueue[i]
	priorityQueue[i].index = i
	priorityQueue[j].index = j
}

func (priorityQueue *PriorityQueue) Push(x interface{}) {
	n := len(*priorityQueue)
	item := x.(*Item)
	item.index = n
	*priorityQueue = append(*priorityQueue, item)
}

func (priorityQueue *PriorityQueue) Pop() interface{} {
	old := *priorityQueue
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*priorityQueue = old[0 : n-1]
	return item
}

func (priorityQueue *PriorityQueue) update(item *Item, value *GraphNode, priority int) {
	item.value = value
	item.priority = priority
	// item.predecessor = predecessor
	heap.Fix(priorityQueue, item.index)
}

type Graph struct {
	startNode *GraphNode
	endNode *GraphNode
}

type GraphNode struct {
	id *GraphNodeId
	height int
	edges map[*GraphNodeId]int
}

type GraphNodeId struct {
	x int
	y int
}

type Interface interface {
	sort.Interface
	Push(x interface{}) // add x as element Len()
	Pop() interface{}   // remove and return element Len() - 1.
}

var graph *Graph = &Graph{}
var yPos int = 0
var heightMap [][]*GraphNode

func main() {
	filePuzzle := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle, parseGraph)
	printHeightMap()
	dijkstraSearch()
}

func parseGraph(line string) {
	emptyArray := []*GraphNode{}
	heightMap = append(heightMap, emptyArray)

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
			currentHeight	= int(currentChar)
		}

		currentNode := &GraphNode{height: currentHeight, id: &GraphNodeId{x: i, y: yPos}, edges: make(map[*GraphNodeId]int)}
		heightMap[yPos] = append(heightMap[yPos], currentNode)

		if isStart {
			graph.startNode = currentNode
		} else if isEnd {
			graph.endNode = currentNode
		}

		if isLeftNeighborReachable(i) {
			leftNode := heightMap[yPos][i-1]
			currentNode.edges[leftNode.id] = 1
			leftNode.edges[currentNode.id] = 1
		}

		if isUpperNeighborReachable(i) {
			upperNode := heightMap[yPos-1][i]
			currentNode.edges[upperNode.id] = 1
			upperNode.edges[currentNode.id] = 1
		}
	}
	yPos++
}

func isLeftNeighborReachable(currentXPos int) bool {
	if currentXPos > 0 {
		leftNeighborHeight := heightMap[yPos][currentXPos-1].height
		currentHeight := heightMap[yPos][currentXPos].height

		if leftNeighborHeight - currentHeight == 1 {
			return true
		}
	}
	return false
}

func isUpperNeighborReachable(currentXPos int) bool {
	if yPos > 0 {
		upperNeighborHeight := heightMap[yPos-1][currentXPos].height
		currentHeight := heightMap[yPos][currentXPos].height

		if upperNeighborHeight - currentHeight == 1 {
			return true
		}
	}
	return false
}

func printHeightMap() {
	var line string = ""

	for y := 0; y < len(heightMap); y++ {
		for x := 0; x < len(heightMap[y]); x++ {
			line += strconv.Itoa(heightMap[y][x].height) + " "
		}
		fmt.Println(line)
		line = ""
	}
}

func dijkstraSearch() {
	totalNodes := len(heightMap) * len(heightMap[0])
	priorityQueue := make(PriorityQueue, 1)
	priorityQueue[0] = &Item{value: nil, priority: 9999, index: 0}
	fmt.Println("Init Start")
	heap.Init(&priorityQueue)
	fmt.Println("Init End")

	fmt.Println("totalnodes: " + strconv.Itoa(totalNodes))

	for y := 0; y < len(heightMap); y++ {
		for x := 0; x < len(heightMap[y]); x++ {
			currentNode := heightMap[y][x]
			priority := 9999
			index := (y*len(heightMap[y]))+x

			if currentNode == graph.startNode {
				priority = 0
			}

			queueItem := &Item{
				value: currentNode,
				priority: priority,
				index: index,
				// predecessor: nil,
			}

			fmt.Println("Initializing: " + strconv.Itoa(currentNode.id.x) + "," + strconv.Itoa(currentNode.id.y) + " " + strconv.Itoa(priority) + " " + strconv.Itoa(index))

			priorityQueue[y*x] = queueItem
		}
	}

	heap.Init(&priorityQueue)

	for len(priorityQueue) > 0 {
		currentItem := heap.Pop(&priorityQueue).(*Item)
		currentNode := currentItem.value

		for key, value := range currentNode.edges {
			fmt.Println(strconv.Itoa(value))
			neighborNode := heightMap[key.y][key.x]
			neighborItem := getNodeFromQueue(priorityQueue, neighborNode)
			if neighborItem != nil {
				updateDistance(priorityQueue, currentItem, neighborItem)
			}
		}

	}
}

func updateDistance(priorityQueue PriorityQueue, currentItem *Item, neighborItem *Item) {
	newPriority := currentItem.priority + currentItem.value.edges[neighborItem.value.id]

	if newPriority < neighborItem.priority {
		priorityQueue.update(neighborItem, neighborItem.value, newPriority)
	}
}

func shortestPath(priorityQueue PriorityQueue) {
	var endItem *Item

	for i := 0; i < len(priorityQueue); i++ {
		if priorityQueue[i].value == graph.endNode {
			endItem = priorityQueue[i]
		}
	}

	var currentItem *Item = endItem

	for endItem != nil {
		fmt.Println(currentItem.value.id)
		// currentItem = currentItem.predecessor
	}
}

func getNodeFromQueue(queue PriorityQueue, node *GraphNode) *Item {
	for i := 0; i < len(queue); i++ {
		if queue[i].value == node {
			return queue[i]
		}
	}

	return nil
}
