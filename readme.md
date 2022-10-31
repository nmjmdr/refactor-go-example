**Refactoring go code example**

_Original code:_
```
package main

import "fmt"

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		go func() {
			switch j % 3 {
			case 0:
				j = j * 1
			case 1:
				j = j * 2
				results <- j * 2
			case 2:
				results <- j * 3
				j = j * 3
			}
		}()
	}
}

func main() {
	jobs := make(chan int)
	results := make(chan int)
	for i := 1; i <= 1000000000; i++ {
		go func() {
			if i%2 == 0 {
				i += 99
			}
			jobs <- i
		}()
	}
	close(jobs)
	jobs2 := []int{}
	for w := 1; w < 1000; w++ {
		jobs2 = append(jobs2, w)
	}
	for i, w := range jobs2 {
		go worker(w, jobs, results)
		i = i + 1
	}
	close(results)
	var sum int32 = 0
	for r := range results {
		sum += int32(r)
	}
	fmt.Println(sum)
}
```
_Things that are wrong in above code:_

. The intention of the code is not clear. From reading it is not immeditaely apparent what the code is trying to accomplish. 
Good code makes the intention apparent. The next set of points will eloborate this further

. Function main seems to be doing a lot of things. It could be broken down into smaller functions

. From the name "worker" it is not immediately apparent, what the function is trying to do. 
That name is generic and does not tell us anything about its functionality

. The code seems to be adding jobs to a job channel, but then immediately closes the channel. 

. It is not clear i%2 check should be performed in a go routine. 
```
for i := 1; i <= 1000000000; i++ {
		go func() {
```
This can be avoided. The second issue is that the value "i" wont be consistent. 
Two go routines in the above loop could recieve the same value of "i". This can be avoided by passing "i" to the go routine. 
Irrespective of this I would get of the go rountine at this stage

. The variable "jobs2" is quite wrong. It is apparent from the code that it is more of "x steps" to perform. Neither is it a channel

. In the following steps, i = i+1 is performed, but its value is never used
```
for i, w := range jobs2 {
		go worker(w, jobs, results)
		i = i + 1
	}
```

. A lot of unnessary and superfolous operations seems to be carried out in worker function
```
for _, j := range arr {
		switch j % 3 {
		case 0:
			j = j * 1 //computation performed, but not used
		case 1:
			// j = j * 2 //unnessary computation performed, can just use append(results, j*4)
			results = append(results, j*2)
		case 2:
			results = append(results, j*3)
			j = j * 3 //computation performed, but not used
		}
	}
```

. The channel "results" is closed but then an attempt is made to iterate over it.

. Overall it took some time to understand the "goal" of the code.

_How I refactored it_

. The first step of any refactoring is to understand the "intention of the code" or what the code is trying to accomplish

. In order to do this, I first simplied the code to a simple iterative version below:
```
package main

import "fmt"

func stepFunction(arr *[]int) int32 {
	sum := int32(0)
	for _, v := range *arr {
		if v%3 == 1 {
			sum = sum + int32(v*2)
		} else if v%3 == 0 {
			sum = sum + int32(v*3)
		}
	}
	return sum
}

func series(start int, stop int) *[]int {
	xs := make([]int, 0)
	for i := start; i <= stop; i++ {
		if i%2 == 0 {
			i += 99
		}
		xs = append(xs, i)
	}
	return &xs
}

func main() {
	start := 1
	stop := 200
	xs := series(start, stop)
	nSteps := 5
	sum := int32(0)
	for step := 0; step < nSteps; step++ {
		r := stepFunction(xs)
		sum += int32(r)
	}
	fmt.Println(sum)
}

```

Now the intetion of the code is bit more clear. It first generates a series (based on check for i%2 == 0) 
and then transforms that series, based on next set of criterion, and then adds the values to generate a sum. Such sums are added over nSteps to get a final sum

We can now think about making the cod parallel. Also I have refactored the code to remove the logic of stepFunction and series function into dependencies that are injected.
This way we can change the series and stepFunction and rest of the code wont be impacted.
Here is the parallel version:
```
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
```
Here is the serial version:
```
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
```
Here are the functions (stepFunction and series) as seperate library functions:
```
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
```

Here is how we set it up in a main function:
```
package main

import (
	"fmt"
	"kablamo/lib/functions"
	"kablamo/lib/parallel"
	"kablamo/lib/serial"
	"time"
)

func main() {
	start := 1
	end := 1000000000
	series := functions.Series(start, end)
	nSteps := 1000

	startTime := time.Now()
	sumParallel := parallel.Process(nSteps, series, functions.StepFunction)
	elapsedParallel := time.Since(startTime)

	startTime = time.Now()
	sumSerial := serial.Process(nSteps, series, functions.StepFunction)
	elapsedSerial := time.Since(startTime)

	startTime = time.Now()
	sumSerialv2 := serial.Processv2(nSteps, series, functions.StepFunction)
	elapsedSerialv2 := time.Since(startTime)

	fmt.Printf("Parallel, value: %d, Time taken: %d ns\n", sumParallel, elapsedParallel.Nanoseconds())
	fmt.Printf("Serial v1, value: %d, Time taken: %d ns\n", sumSerial, elapsedSerial.Nanoseconds())
	fmt.Printf("Serial v2, value: %d, Time taken: %d ns\n", sumSerialv2, elapsedSerialv2.Nanoseconds())
}
```

Also note that it is unnecessary to repetedly invoke stepFunction to perform the same computation. We can just evaluate it once and then multiply nSteps.
The shortest time is taken by the version which performs this. 

```
go run cmd/main.go                                                            
Parallel, value: 751852706, Time taken: 4479434500 ns
Serial v1, value: 751852706, Time taken: 24585894792 ns
Serial v2, value: 834291376, Time taken: 25093208 ns
```

I have also added a few test cases and benchmark tests.


