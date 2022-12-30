package shared

func Abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func Max(values ...int) int {
	max := values[0]

	for i := 0; i < len(values); i++ {
		if values[i] > max {
			max = values[i]
		}
	}

	return max
}
