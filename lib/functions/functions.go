package functions

func Series(start int, end int) *[]int {
	xs := []int{}
	for i := start; i <= end; i++ {
		if i%2 == 0 {
			i += 99
		}
		xs = append(xs, i)
	}
	return &xs
}

func StepFunction(i int32, acc int32) int32 {
	r := i % 3
	if r == 1 {
		return acc + int32(i*4)
	}
	if r == 2 {
		return acc + int32(i*3)
	}
	return acc
}
