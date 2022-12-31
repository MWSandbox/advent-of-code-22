package main

import (
	"advent-of-code/shared"
	"fmt"
	"strconv"
	"strings"
)

type TreeNode struct {
	id           string
	operator     string
	leftOperand  *TreeNode
	rightOperand *TreeNode
	result       int
}

const rootIdentifier string = "root"
const noResult int = -1
const human string = "humn"

var idToTreeNode map[string]*TreeNode = make(map[string]*TreeNode)
var depthToTreeNode map[int][]*TreeNode = make(map[int][]*TreeNode)
var idToOperands map[string][2]string = make(map[string][2]string)
var rootNode *TreeNode
var humanNode *TreeNode

func main() {
	filePuzzle := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle, parseTreeNodes)
	buildBinaryTree(rootNode, 0)

	solvePuzzle1()
	solvePuzzle2()
}

func parseTreeNodes(line string) {
	var operator string
	var staticNumber int

	parts := strings.Split(line, ": ")
	name := parts[0]
	equationParts := strings.Split(parts[1], " ")

	if len(equationParts) > 2 {
		leftOperand := equationParts[0]
		operator = equationParts[1]
		rightOperand := equationParts[2]
		staticNumber = noResult
		idToOperands[name] = [2]string{leftOperand, rightOperand}
	} else {
		staticNumber = shared.ConvertStringToInt(equationParts[0])
	}

	player := &TreeNode{id: name, operator: operator, result: staticNumber}
	idToTreeNode[name] = player

	if name == rootIdentifier {
		rootNode = player
	}

	if name == human {
		humanNode = player
	}
}

func buildBinaryTree(currentNode *TreeNode, depth int) {
	if currentNode.operator != "" {
		leftOperandId := idToOperands[currentNode.id][0]
		rightOperandId := idToOperands[currentNode.id][1]
		currentNode.leftOperand = idToTreeNode[leftOperandId]
		currentNode.rightOperand = idToTreeNode[rightOperandId]

		buildBinaryTree(currentNode.leftOperand, depth+1)
		buildBinaryTree(currentNode.rightOperand, depth+1)
	}
	depthToTreeNode[depth] = append(depthToTreeNode[depth], currentNode)
}

func solvePuzzle1() {
	for depth := len(depthToTreeNode) - 1; depth >= 0; depth-- {
		for i := 0; i < len(depthToTreeNode[depth]); i++ {
			currentNode := depthToTreeNode[depth][i]

			if currentNode.result == noResult {
				currentNode.result = solveEquation(currentNode)
			}
		}
	}

	fmt.Println("Result for root monkey: " + strconv.Itoa(rootNode.result))
}

func solveEquation(nodeWithEquation *TreeNode) int {
	leftValue := nodeWithEquation.leftOperand.result
	rightValue := nodeWithEquation.rightOperand.result

	if nodeWithEquation.operator == "+" {
		return leftValue + rightValue
	} else if nodeWithEquation.operator == "-" {
		return leftValue - rightValue
	} else if nodeWithEquation.operator == "*" {
		return leftValue * rightValue
	} else if nodeWithEquation.operator == "/" {
		return leftValue / rightValue
	}
	return leftValue
}

func solvePuzzle2() {
	rootNode.operator = "="
	missingValue := findMissingValue(rootNode)

	fmt.Println("I need to yell: " + strconv.Itoa(missingValue))
}

func findMissingValue(currentNode *TreeNode) int {
	var nextNode *TreeNode
	isHumanInLeftSubtree := isHumanInSubTree(currentNode.leftOperand)

	if isHumanInLeftSubtree {
		currentNode.leftOperand.result = noResult
		nextNode = currentNode.leftOperand
	} else {
		currentNode.rightOperand.result = noResult
		nextNode = currentNode.rightOperand
	}

	setUnknownInEquation(currentNode)

	if nextNode == humanNode {
		return nextNode.result
	}

	return findMissingValue(nextNode)
}

func isHumanInSubTree(currentNode *TreeNode) bool {
	isHuman := currentNode.id == human
	hasOperands := currentNode.leftOperand != nil
	return isHuman || (hasOperands && (isHumanInSubTree(currentNode.leftOperand) || isHumanInSubTree(currentNode.rightOperand)))
}

func setUnknownInEquation(nodeWithEquation *TreeNode) {
	var unknownOperand *TreeNode
	var determinedOperand *TreeNode

	if nodeWithEquation.leftOperand.result == noResult {
		unknownOperand = nodeWithEquation.leftOperand
		determinedOperand = nodeWithEquation.rightOperand
	} else {
		unknownOperand = nodeWithEquation.rightOperand
		determinedOperand = nodeWithEquation.leftOperand
	}

	if nodeWithEquation.operator == "+" {
		unknownOperand.result = nodeWithEquation.result - determinedOperand.result
	} else if nodeWithEquation.operator == "-" && unknownOperand == nodeWithEquation.leftOperand {
		unknownOperand.result = determinedOperand.result + nodeWithEquation.result
	} else if nodeWithEquation.operator == "-" {
		unknownOperand.result = determinedOperand.result - nodeWithEquation.result
	} else if nodeWithEquation.operator == "*" {
		unknownOperand.result = nodeWithEquation.result / determinedOperand.result
	} else if nodeWithEquation.operator == "/" && unknownOperand == nodeWithEquation.leftOperand {
		unknownOperand.result = determinedOperand.result * nodeWithEquation.result
	} else if nodeWithEquation.operator == "/" {
		unknownOperand.result = determinedOperand.result / nodeWithEquation.result
	} else {
		unknownOperand.result = determinedOperand.result
	}
}
