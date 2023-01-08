package shared

func CantorHashThreeFields(a int, b int, c int) int {
	return (CantorHashTwoFields(a, b)+c)*(CantorHashTwoFields(a, b)+c+1)/2 + c
}

func CantorHashTwoFields(a int, b int) int {
	return (a+b)*(a+b+1)/2 + b
}
