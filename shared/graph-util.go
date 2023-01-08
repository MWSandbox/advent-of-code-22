package shared

import (
	"container/heap"
)

type Item struct {
	value    *GraphNode
	priority int
	index    int
}

type PriorityQueue []*Item

type GraphNode struct {
	Neighbors []*GraphNode
	Index     int
}

const infinity int = 9999999

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
	heap.Fix(priorityQueue, item.index)
}

func DijkstraWithMultiStart(allNodes []*GraphNode, startNodes []*GraphNode, endNode *GraphNode) int {
	shortestPath := infinity

	for i := 0; i < len(startNodes); i++ {
		result := Dijkstra(allNodes, startNodes[i], endNode)

		if result < shortestPath {
			shortestPath = result
		}
	}
	return shortestPath
}

func DijkstraWithLimit(allNodes []*GraphNode, startNode *GraphNode, endNode *GraphNode, limit int) int {
	nodeToDistance := initializeDistances(allNodes, startNode)

	priorityQueue := make(PriorityQueue, len(allNodes))

	for i := 0; i < len(allNodes); i++ {
		priorityQueue[i] = &Item{value: allNodes[i], priority: nodeToDistance[allNodes[i]], index: i}
	}
	heap.Init(&priorityQueue)

	for priorityQueue.Len() > 0 {
		item := heap.Pop(&priorityQueue).(*Item)

		if item.priority > limit {
			return infinity
		}

		for i := 0; i < len(item.value.Neighbors); i++ {
			neighborItem := getItemFromQueue(priorityQueue, item.value.Neighbors[i])

			if neighborItem != nil {
				updateDistance(nodeToDistance, priorityQueue, item, neighborItem)
			}
		}
	}

	return nodeToDistance[endNode]
}

func Dijkstra(allNodes []*GraphNode, startNode *GraphNode, endNode *GraphNode) int {
	return DijkstraWithLimit(allNodes, startNode, endNode, infinity)
}

func initializeDistances(allNodes []*GraphNode, startNode *GraphNode) map[*GraphNode]int {
	var nodeToDistance map[*GraphNode]int = make(map[*GraphNode]int)
	for i := 0; i < len(allNodes); i++ {
		nodeToDistance[allNodes[i]] = infinity
	}
	nodeToDistance[startNode] = 0
	return nodeToDistance
}

func getItemFromQueue(queue PriorityQueue, node *GraphNode) *Item {
	for i := 0; i < len(queue); i++ {
		if queue[i].value == node {
			return queue[i]
		}
	}

	return nil
}

func updateDistance(nodeToDistance map[*GraphNode]int, priorityQueue PriorityQueue, currentItem *Item, neighborItem *Item) {
	newPriority := currentItem.priority + 1

	if newPriority < neighborItem.priority {
		nodeToDistance[neighborItem.value] = newPriority
		priorityQueue.update(neighborItem, neighborItem.value, newPriority)
	}
}
