package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/mthenw/faasperf/awslambda"
	"github.com/sanity-io/litter"
	chart "github.com/wcharczuk/go-chart"
)

func main() {
	fmt.Println("Running benchmark for AWS Lambda provider...")

	results, err := awslambda.Benchmark()
	if err != nil {
		log.Fatal(err.Error())
	}

	x := []float64{}
	y := []float64{}
	for _, memory := range awslambda.MemorySizes {
		x = append(x, float64(memory))
		y = append(y, float64(results[memory].CPUMs))
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{Show: true},
			Ticks: ticks(),
		},
		YAxis: chart.YAxis{
			Style: chart.Style{Show: true},
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 20,
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "AWS Lambda CPU",
				XValues: x,
				YValues: y,
			},
		},
	}

	buf := bytes.NewBuffer([]byte{})
	err = graph.Render(chart.PNG, buf)
	if err != nil {
		log.Fatalf("error creating image: %s", err.Error())
	}

	err = ioutil.WriteFile("results/awslambda.png", buf.Bytes(), 0644)
	if err != nil {
		log.Fatalf("error saving image: %s", err.Error())
	}

	litter.Dump(results)
}

func ticks() (ticks []chart.Tick) {
	for _, size := range awslambda.MemorySizes {
		ticks = append(ticks, chart.Tick{Value: float64(size), Label: strconv.Itoa(int(size))})
	}
	return
}
