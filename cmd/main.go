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
