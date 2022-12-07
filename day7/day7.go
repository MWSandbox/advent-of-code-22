package main

import (
	"advent-of-code/shared"
	"fmt"
	"strconv"
	"strings"
)

type FileSystemTree struct {
	root *Directory
}

type FileSystemNode struct {
	name   string
	parent *Directory
}

type Directory struct {
	FileSystemNode
	childDirectories []*Directory
	childFiles       []*File
}

type File struct {
	FileSystemNode
	size int
}

type DirSizeCallbackFn func(*Directory, int)

const totalDiskSpace int = 70000000
const requiredDiskSpace int = 30000000
const dirMaxSize int = 100000

var currentDirectory *Directory
var fileSystemTree FileSystemTree
var MaxDirSizeSum int = 0
var sizeOfDirToBeDeleted int = 999999999
var missingDiskSpace int = 0
var totalSpaceUsed int = 0

func main() {
	filePuzzle := shared.OpenFile("./input.txt")
	shared.ReadFileLineByLine(filePuzzle, createFileSystemTree)
	calculateDirSizes(fileSystemTree.root, callBackForHandlingMaxDirSize)
	fmt.Println("Total directory size of all directory smaller than 100000: " + strconv.Itoa(MaxDirSizeSum))

	missingDiskSpace = requiredDiskSpace - (totalDiskSpace - totalSpaceUsed)
	fmt.Println("Missing disk space: " + strconv.Itoa(missingDiskSpace))
	calculateDirSizes(fileSystemTree.root, callBackForHandlingDirToBeDeleted)
	fmt.Println("Directory size to be deleted: " + strconv.Itoa(sizeOfDirToBeDeleted))
}

func createFileSystemTree(line string) {
	var isCdCommand bool = line[0:4] == "$ cd"
	var isLsCommand bool = line[0:4] == "$ ls"
	var isDirectory bool = line[0:3] == "dir"

	if isCdCommand {
		parameter := line[5:len(line)]
		executeCdCommand(parameter)
	} else if isDirectory {
		parameter := line[4:len(line)]
		addDirectoryIfNotExists(parameter)
	} else if !isLsCommand {
		fileDescription := strings.Split(line, " ")
		addFileIfNotExists(fileDescription[0], fileDescription[1])
	}
}

func executeCdCommand(parameter string) {
	if parameter == ".." {
		currentDirectory = currentDirectory.parent

	} else {
		if parameter == "/" {
			if currentDirectory == nil {
				fileSystemTree.root = createNewDir("/", nil)
			}
			currentDirectory = fileSystemTree.root
		} else {
			currentDirectory = addDirectoryIfNotExists(parameter)
		}
	}
}

func addDirectoryIfNotExists(name string) *Directory {
	var child *Directory

	for i := 0; i < len(currentDirectory.childDirectories); i++ {
		if name == currentDirectory.childDirectories[i].name {
			child = currentDirectory.childDirectories[i]
		}
	}
	if child == nil {
		child = createNewDir(name, currentDirectory)
		currentDirectory.childDirectories = append(currentDirectory.childDirectories, child)
	}

	return child
}

func createNewDir(name string, parent *Directory) *Directory {
	return &Directory{
		FileSystemNode:   FileSystemNode{name: name, parent: parent},
		childDirectories: nil,
		childFiles:       nil,
	}
}

func addFileIfNotExists(sizeAsString string, name string) *File {
	var child *File

	size, err := strconv.Atoi(sizeAsString)

	if err != nil {
		fmt.Println("Couldn't convert file size to int: " + sizeAsString)
	}

	for i := 0; i < len(currentDirectory.childFiles); i++ {
		if name == currentDirectory.childFiles[i].name {
			child = currentDirectory.childFiles[i]
		}
	}
	if child == nil {
		child = &File{
			size:           size,
			FileSystemNode: FileSystemNode{name: name, parent: currentDirectory},
		}
		currentDirectory.childFiles = append(currentDirectory.childFiles, child)
	}

	return child
}

func calculateDirSizes(dir *Directory, dirSizeCallback DirSizeCallbackFn) int {
	var size int = 0

	for i := 0; i < len(dir.childFiles); i++ {
		fileSize := dir.childFiles[i].size
		size += fileSize
		totalSpaceUsed += fileSize
	}

	for i := 0; i < len(dir.childDirectories); i++ {
		size += calculateDirSizes(dir.childDirectories[i], dirSizeCallback)
	}

	dirSizeCallback(dir, size)

	return size
}

func callBackForHandlingMaxDirSize(dir *Directory, size int) {
	if size < dirMaxSize {
		MaxDirSizeSum += size
	}
}

func callBackForHandlingDirToBeDeleted(dir *Directory, size int) {
	if size >= missingDiskSpace && size < sizeOfDirToBeDeleted {
		sizeOfDirToBeDeleted = size
	}
}
