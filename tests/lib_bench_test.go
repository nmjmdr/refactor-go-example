package tests

import (
	"kablamo/lib/functions"
	"kablamo/lib/parallel"
	"kablamo/lib/serial"
	"testing"
)

func BenchmarkParallel(b *testing.B) {
	series := functions.Series(1, 1000000000)
	nSteps := 1000
	for n := 0; n < b.N; n++ {
		parallel.Process(nSteps, series, functions.StepFunction)
	}
}

func BenchmarkSerial(b *testing.B) {
	series := functions.Series(1, 1000000000)
	nSteps := 1000
	for n := 0; n < b.N; n++ {
		serial.Process(nSteps, series, functions.StepFunction)
	}
}

func BenchmarkSerialv2(b *testing.B) {
	series := functions.Series(1, 1000000000)
	nSteps := 1000
	for n := 0; n < b.N; n++ {
		serial.Processv2(nSteps, series, functions.StepFunction)
	}
}
