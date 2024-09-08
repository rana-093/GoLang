package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Measurement struct {
	Station string
	Temp    float64
}

type Stats struct {
	Min, Max, Mean float64
}

func min2(x float64, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

func max2(x float64, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

func calc() error {
	start := time.Now()
	file, err := os.Open("measurements.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	stationTemps := make(map[string][]float64)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ";")
		temp, _ := strconv.ParseFloat(parts[1], 64)
		stationName := parts[0]
		stationTemps[stationName] = append(stationTemps[stationName], temp)
	}

	statOfStations := make(map[string]Stats)

	var stations []string

	for key, value := range stationTemps {
		min, max, sum := value[0], value[0], 0.0
		for _, t := range value {
			min = min2(min, t)
			max = max2(max, t)
			sum += t
		}
		mean := sum / float64(len(value))
		statOfStations[key] = Stats{Max: max, Mean: mean, Min: min}
		stations = append(stations, key)
	}

	sort.Strings(stations)
	fmt.Print("{")

	for idx, name := range stations {
		result := statOfStations[name]
		fmt.Printf("%s=%.1f/%.1f/%.1f", name, result.Min, result.Mean, result.Max)
		if idx < len(stations)-1 {
			fmt.Print(", ")
		}
		fmt.Println("}")
	}

	end := time.Now()

	fmt.Printf("Total time taken is %s", end.Sub(start))
	return nil
}
