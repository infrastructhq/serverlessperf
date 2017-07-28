package main

import (
	"log"

	"github.com/mthenw/faasperf/awslambda"
)

func main() {
	results, err := awslambda.Benchmark()
	if err != nil {
		log.Fatal(err.Error())
	}
}
