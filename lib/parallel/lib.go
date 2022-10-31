package parallel

import "sync"

func reduce(arr *[]int, fn func(int32, int32) int32) int32 {
	acc := int32(0)
	for _, ele := range *arr {
		acc = fn(int32(ele), acc)
	}
	return acc
}

func stepSum(ch chan int32, xs *[]int, wg *sync.WaitGroup, stepFn func(int32, int32) int32) {
	defer wg.Done()
	stepSum := reduce(xs, stepFn)
	ch <- stepSum
}

func summation(ch chan int32, outCh chan int32) {
	sum := int32(0)
	for step := range ch {
		sum += step
	}
	outCh <- sum
}

func Process(nSteps int, series *[]int, stepFn func(int32, int32) int32) int32 {
	stepCh := make(chan int32)
	outCh := make(chan int32)
	wg := &sync.WaitGroup{}
	// instead of this for loop, we can just multiply step-sum with nSteps
	// I have done this to retain the logic of original program
	for step := 1; step < nSteps; step++ {
		wg.Add(1)
		go stepSum(stepCh, series, wg, stepFn)
	}
	go summation(stepCh, outCh)
	wg.Wait()
	close(stepCh)
	sum := <-outCh
	close(outCh)
	return sum
}
