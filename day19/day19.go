package main

import (
	"advent-of-code/shared"
	"fmt"
	"strconv"
	"strings"
)

type Blueprint struct {
	oreRobot                *Robot
	clayRobot               *Robot
	obsidianRobot           *Robot
	geodeRobot              *Robot
	maxOreRequirements      int
	maxClayRequirements     int
	maxObsidianRequirements int
}

type Robot struct {
	ore      int
	clay     int
	obsidian int
}

type Inventory struct {
	oreRobots      int
	clayRobots     int
	obsidianRobots int
	geodeRobots    int
	ore            int
	clay           int
	obsidian       int
	geode          int
}

const ore string = "ore"
const clay string = "clay"
const obsidian string = "obsidian"

var timeLimit int = 24
var blueprints []*Blueprint
var maxGeodesPerBlueprint []int
var earliestGeode int = timeLimit

func main() {
	filePuzzle := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle, parseBlueprints)

	maxGeodesPerBlueprint = make([]int, len(blueprints))

	solvePuzzle1()
	solvePuzzle2()
}

func parseBlueprints(line string) {
	blueprintContent := strings.Split(line, ": ")[1]
	recipes := strings.Split(blueprintContent, ". ")

	oreRobot := parseOreRobot(recipes[0])
	clayRobot := parseClayRobot(recipes[1])
	obsidianRobot := parseObsidianRobot(recipes[2])
	geodeRobot := parseGeodeRobot(recipes[3])

	blueprint := &Blueprint{oreRobot: oreRobot, clayRobot: clayRobot, obsidianRobot: obsidianRobot, geodeRobot: geodeRobot}
	blueprint.maxOreRequirements = oreRobot.ore + clayRobot.ore + obsidianRobot.ore + geodeRobot.ore
	blueprint.maxClayRequirements = obsidianRobot.clay
	blueprint.maxObsidianRequirements = geodeRobot.obsidian
	blueprints = append(blueprints, blueprint)
}

func parseOreRobot(recipe string) *Robot {
	oreCosts := parseCostsOfMaterial(recipe, ore)
	return &Robot{ore: oreCosts, clay: 0, obsidian: 0}
}

func parseClayRobot(recipe string) *Robot {
	oreCosts := parseCostsOfMaterial(recipe, ore)
	return &Robot{ore: oreCosts, clay: 0, obsidian: 0}
}

func parseObsidianRobot(recipe string) *Robot {
	oreCosts := parseCostsOfMaterial(recipe, ore)
	clayCosts := parseCostsOfMaterial(recipe, clay)
	return &Robot{ore: oreCosts, clay: clayCosts, obsidian: 0}
}

func parseGeodeRobot(recipe string) *Robot {
	oreCosts := parseCostsOfMaterial(recipe, ore)
	obsidianCosts := parseCostsOfMaterial(recipe, obsidian)
	return &Robot{ore: oreCosts, clay: 0, obsidian: obsidianCosts}
}

func parseCostsOfMaterial(recipe string, material string) int {
	materialSplit := strings.Split(recipe, " "+material)
	start := materialSplit[len(materialSplit)-2]
	parts := strings.Split(start, " ")
	costs := parts[len(parts)-1]
	return shared.ConvertStringToInt(costs)
}

func solvePuzzle1() {
	qualityLevel := 0
	for i := 0; i < len(blueprints); i++ {
		calculateMaxOutputOfBlueprint(i)
		qualityLevel += (i + 1) * maxGeodesPerBlueprint[i]
		fmt.Println("Max output of blueprint " + strconv.Itoa(i+1) + ": " + strconv.Itoa(maxGeodesPerBlueprint[i]))
	}
	fmt.Println("Quality level: " + strconv.Itoa(qualityLevel))
}

func solvePuzzle2() {
	blueprints = blueprints[:3]
	timeLimit = 32
	result := 1
	earliestGeode = timeLimit
	for i := 0; i < len(blueprints); i++ {
		calculateMaxOutputOfBlueprint(i)
		result *= maxGeodesPerBlueprint[i]
		fmt.Println("Max output of blueprint " + strconv.Itoa(i+1) + ": " + strconv.Itoa(maxGeodesPerBlueprint[i]))
	}

	fmt.Println("Result: " + strconv.Itoa(result))
}

func calculateMaxOutputOfBlueprint(blueprintIndex int) {
	inventory := &Inventory{oreRobots: 1}
	earliestGeode = timeLimit
	runSimulation(blueprintIndex, inventory, 0)
}

func runSimulation(blueprintIndex int, inventory *Inventory, minutesPassed int) {
	var robotsToIgnore []*Robot
	blueprint := blueprints[blueprintIndex]

	for time := minutesPassed; time < timeLimit; time++ {
		if time > earliestGeode && inventory.geodeRobots == 0 {
			return
		}

		buildableRobots := getBuildableRobots(inventory, blueprint, robotsToIgnore)
		robotsToIgnore = append(robotsToIgnore, buildableRobots...)
		isGeodeRobotBuildable := isRobotIgnored(blueprint.geodeRobot, robotsToIgnore)
		lastTimeForClayRobot := calculateLastTimeToBuildRobot(inventory.clayRobots, time, blueprint.obsidianRobot.clay)
		lastTimeForObsidianRobot := calculateLastTimeToBuildRobot(inventory.obsidianRobots, time, blueprint.geodeRobot.obsidian)

		for j := 0; j < len(buildableRobots); j++ {
			if shouldBuildRobot(time, lastTimeForClayRobot, lastTimeForObsidianRobot, buildableRobots[j], blueprint, isGeodeRobotBuildable, inventory) {
				newInventory := copyInventory(inventory)
				buildRobot(newInventory, buildableRobots[j])

				if buildableRobots[j] == blueprint.geodeRobot && earliestGeode > time {
					earliestGeode = time
				}

				collectResources(newInventory)
				addRobotToInventory(newInventory, buildableRobots[j], blueprint)
				runSimulation(blueprintIndex, newInventory, time+1)
			}
		}
		collectResources(inventory)
	}

	if inventory.geode > maxGeodesPerBlueprint[blueprintIndex] {
		maxGeodesPerBlueprint[blueprintIndex] = inventory.geode
	}
}

func calculateLastTimeToBuildRobot(currentRobotCount int, minutesPassed int, costsForNextBetterRobot int) int {
	currentClayOutput := currentRobotCount * (timeLimit - minutesPassed)
	maxClayOutput := currentClayOutput

	for i := minutesPassed + 1; i < timeLimit; i++ {
		obsidianRobotsToBuild := maxClayOutput / costsForNextBetterRobot
		timeToBuildObsidianRobots := obsidianRobotsToBuild + 1
		if timeLimit-i < timeToBuildObsidianRobots {
			return i
		}
		for j := i; j < timeLimit; j++ {
			maxClayOutput++
		}
	}
	return 0
}

func getBuildableRobots(inventory *Inventory, blueprint *Blueprint, robotsToIgnore []*Robot) []*Robot {
	var buildableRobots []*Robot
	var allRobots [4]*Robot = [4]*Robot{blueprint.oreRobot, blueprint.clayRobot, blueprint.obsidianRobot, blueprint.geodeRobot}

	for i := 0; i < len(allRobots); i++ {
		robot := allRobots[i]
		if !isRobotIgnored(robot, robotsToIgnore) && hasResourcesToBuild(inventory, robot) {
			buildableRobots = append(buildableRobots, robot)
		}
	}

	return buildableRobots
}

func shouldBuildRobot(timePassed int, lastTimeForClayRobot int, lastTimeForObsidianRobot int, robotToBuild *Robot, blueprint *Blueprint, isGeodeRobotBuildable bool, inventory *Inventory) bool {
	isEnoughTimeLeft := timePassed < timeLimit-1
	isLastGeodeRobotor := timePassed >= timeLimit-1 && robotToBuild == blueprint.geodeRobot
	isOreOverproduction := blueprint.maxOreRequirements <= inventory.ore && robotToBuild == blueprint.oreRobot
	isClayOverproduction := blueprint.maxClayRequirements <= inventory.clay && robotToBuild == blueprint.clayRobot
	isObsidianOverproduction := blueprint.maxObsidianRequirements <= inventory.obsidian && robotToBuild == blueprint.obsidianRobot
	isLastTimeForClayReached := timePassed >= lastTimeForClayRobot && robotToBuild == blueprint.clayRobot
	isLastTimeForObsidianReached := timePassed >= lastTimeForObsidianRobot && robotToBuild != blueprint.geodeRobot

	isRobotBuilt := (isEnoughTimeLeft || isLastGeodeRobotor)
	isRobotBuilt = isRobotBuilt && !isOreOverproduction && !isClayOverproduction && !isObsidianOverproduction
	isRobotBuilt = isRobotBuilt && !isLastTimeForClayReached && !isLastTimeForObsidianReached
	return isRobotBuilt
}

func isRobotIgnored(robot *Robot, robotsToIgnore []*Robot) bool {
	for i := 0; i < len(robotsToIgnore); i++ {
		if robotsToIgnore[i] == robot {
			return true
		}
	}
	return false
}

func hasResourcesToBuild(inventory *Inventory, robot *Robot) bool {
	return inventory.ore >= robot.ore && inventory.clay >= robot.clay && inventory.obsidian >= robot.obsidian
}

func buildRobot(inventory *Inventory, robot *Robot) {
	inventory.ore -= robot.ore
	inventory.clay -= robot.clay
	inventory.obsidian -= robot.obsidian
}

func addRobotToInventory(inventory *Inventory, robot *Robot, blueprint *Blueprint) {
	if robot == blueprint.oreRobot {
		inventory.oreRobots++
	} else if robot == blueprint.clayRobot {
		inventory.clayRobots++
	} else if robot == blueprint.obsidianRobot {
		inventory.obsidianRobots++
	} else {
		inventory.geodeRobots++
	}
}

func collectResources(inventory *Inventory) {
	inventory.ore += inventory.oreRobots
	inventory.clay += inventory.clayRobots
	inventory.obsidian += inventory.obsidianRobots
	inventory.geode += inventory.geodeRobots
}

func copyInventory(inventory *Inventory) *Inventory {
	newInventory := &Inventory{ore: inventory.ore, clay: inventory.clay, obsidian: inventory.obsidian, geode: inventory.geode}
	newInventory.oreRobots = inventory.oreRobots
	newInventory.clayRobots = inventory.clayRobots
	newInventory.obsidianRobots = inventory.obsidianRobots
	newInventory.geodeRobots = inventory.geodeRobots
	return newInventory
}
