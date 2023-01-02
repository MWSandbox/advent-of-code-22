package shared

import (
	"fmt"
	"strconv"
)

func PrintXHeader(minX int, maxX int) {
	xHeader1 := "    "
	xHeader2 := "    "
	xHeader3 := "    "

	for x := minX; x < maxX; x++ {
		firstDigit := Abs(x) / 100
		secondDigit := (Abs(x) % 100) / 10
		thirdDigit := Abs(x) % 10

		xHeader1 += strconv.Itoa(firstDigit)
		xHeader2 += strconv.Itoa(secondDigit)
		xHeader3 += strconv.Itoa(thirdDigit)
	}

	fmt.Println(xHeader1)
	fmt.Println(xHeader2)
	fmt.Println(xHeader3)
}

func GetYPadding(y int) string {
	padding := " "

	if Abs(y) < 10 {
		padding += " "
	}

	if Abs(y) < 100 {
		padding += " "
	}

	return padding
}
