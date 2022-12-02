package shared

import (
	"bufio"
	"log"
	"os"
)

type parserFn func(string)

func OpenFile(filePath string) *os.File {
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}
	return file
}

func ReadFileLineByLine(file *os.File, parser parserFn) {
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var currentLine string = scanner.Text()
		parser(currentLine)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}
