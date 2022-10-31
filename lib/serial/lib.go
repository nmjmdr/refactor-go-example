package serial

func reduce(arr *[]int, fn func(int32, int32) int32) int32 {
	acc := int32(0)
	for _, ele := range *arr {
		acc = fn(int32(ele), acc)
	}
	return acc
}

func series(start int, end int) *[]int {
	xs := []int{}
	for i := start; i <= end; i++ {
		if i%2 == 0 {
			i += 99
		}
		xs = append(xs, i)
	}
	return &xs
}

func stepSum(xs *[]int, stepFn func(int32, int32) int32) int32 {
	return reduce(xs, stepFn)
}

func Process(nSteps int, series *[]int, stepFn func(int32, int32) int32) int32 {
	sum := int32(0)
	// instead of this for loop, we can just multiply step-sum with nSteps
	// I have done this to retain the logic of original program
	for step := 1; step < nSteps; step++ {
		sum += stepSum(series, stepFn)
	}
	return sum
}

func Processv2(nSteps int, series *[]int, stepFn func(int32, int32) int32) int32 {
	// just multiply directly
	return stepSum(series, stepFn) * int32(nSteps)
}
