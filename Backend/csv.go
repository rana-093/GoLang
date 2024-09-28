package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SpeedInfo struct {
	Lat           float32
	Lon           float32
	Speed         float32
	Distance      float32
	TotalDistance float32
	TimeStamp     time.Time
	Motion        bool
}

type GroupWiseAggregatedResult struct {
	DurationInSeconds    float32
	TotalCoveredDistance float32
	Accelerate           bool
	Decelerate           bool
	TotalVelocity        float32
}

type AccelerationResult struct {
	PositiveAcceleration         int
	NegativeAcceleration         int
	AccelerationPer100KM         float32
	NegativeAccelerationPer100KM float32
	TotalDistanceInKM            float32
	DeviceId                     string
	MaxSpeed                     float32
}

func millisecondsToHours(ms int64) float32 {
	return float32(ms) / 3_600_000
}

func millisecondsToSeconds(ms int64) float32 {
	return float32(ms) / 1_000
}

func ParseXSLXFromDailyUsageReport(fileName string) (error, map[string]float64) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, nil
	}

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	totalDistanceByDevice := make(map[string]float64)
	for _, row := range rows {
		if strings.EqualFold(row[0], "Driver") {
			continue
		}
		distanceParts := strings.Split(row[3], " ")
		deviceName := row[1]
		distance, _ := strconv.ParseFloat(distanceParts[0], 16)
		if _, exists := totalDistanceByDevice[deviceName]; exists {
			totalDistanceByDevice[deviceName] += distance
		} else {
			totalDistanceByDevice[deviceName] = distance
		}
	}
	return nil, totalDistanceByDevice
}

func ParseXSLXFromObjectHistoryReport(fileName string) (error, map[string]AccelerationResult) {
	startTime := time.Now()
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, nil
	}

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		fmt.Println("Error reading rows:", err)
		return nil, nil
	}

	var wg sync.WaitGroup
	var _ sync.Mutex

	const Layout, Workers = "02-01-2006 15:04:05", 50
	const AccelerateThresholdInKMPerSecond = 6
	chunkSize := len(rows) / Workers
	fmt.Printf("chunksize is: %d. Total rows: %d\n", chunkSize, len(rows))

	chunkWiseRows := make([]int, 0)
	chunkWiseDevices := make([]string, 0)

	for i, row := range rows {
		for _, cell := range row {
			if strings.EqualFold(cell, "Device:") {
				if len(row) > 1 {
					chunkWiseRows = append(chunkWiseRows, i)
					chunkWiseDevices = append(chunkWiseDevices, row[1])
				}
			}
		}
	}

	speedMap := make(map[string]map[time.Time][]SpeedInfo)

	for _, val := range chunkWiseDevices {
		speedMap[val] = make(map[time.Time][]SpeedInfo)
	}

	processChunk := func(data [][]string, deviceId string) (int, int, float32) {
		totalPositiveAccCase, totalNegativeAccCase := 0, 0
		maxSpeed := 0.0
		for i := 1; i < len(data); i++ {
			if len(data[i]) < 18 || strings.EqualFold(data[i][0], "Device") || len(data[i-1]) < 18 {
				continue
			}
			curTimeStamp, err := time.Parse(Layout, data[i][2])
			prevTimestamp, err := time.Parse(Layout, data[i-1][2])
			timeBetween := curTimeStamp.Sub(prevTimestamp).Milliseconds()
			if timeBetween == 0 || err != nil {
				continue
			}

			curSpeedParts := strings.Split(data[i][7], " ")
			prevSpeedParts := strings.Split(data[i-1][7], " ")

			curSpeed, _ := strconv.ParseFloat(curSpeedParts[0], 32)
			prevSpeed, _ := strconv.ParseFloat(prevSpeedParts[0], 32)

			maxSpeed = math.Max(maxSpeed, curSpeed)
			maxSpeed = math.Max(maxSpeed, prevSpeed)

			acceleration := ((curSpeed - prevSpeed) / float64(timeBetween)) * 1000.0

			if math.Abs(acceleration) >= AccelerateThresholdInKMPerSecond {
				if acceleration > 0 {
					totalPositiveAccCase++
				} else {
					totalNegativeAccCase++
				}
			}
		}
		return totalPositiveAccCase, totalNegativeAccCase, float32(maxSpeed)
	}

	chunkWiseData := make(chan map[string]AccelerationResult, 10)

	for i := 0; i < len(chunkWiseRows); i++ {
		start := chunkWiseRows[i]
		var end int
		if i+1 < len(chunkWiseRows) {
			end = chunkWiseRows[i+1]
		} else {
			end = len(rows)
		}
		chunk := rows[start:end]
		wg.Add(1)
		go func(device string, chunk [][]string) {
			defer wg.Done()
			positiveCase, negativeCase, maxSpeed := processChunk(chunk, device)
			chunkWiseData <- map[string]AccelerationResult{
				device: {
					PositiveAcceleration: positiveCase,
					NegativeAcceleration: negativeCase,
					MaxSpeed:             maxSpeed,
				},
			}
		}(chunkWiseDevices[i], chunk)
	}

	go func() {
		wg.Wait()
		close(chunkWiseData)
	}()

	result := make(map[string]AccelerationResult)

	for data := range chunkWiseData {
		for key, value := range data {

			totalPositiveAccelerationCase, totalNegativeAccelerationCase, totalDistanceCovered :=
				value.PositiveAcceleration, value.NegativeAcceleration, float32(0)
			maxSpeed := value.MaxSpeed

			if existing, exists := result[key]; exists {
				existing.PositiveAcceleration += totalPositiveAccelerationCase
				existing.NegativeAcceleration += totalNegativeAccelerationCase
				existing.TotalDistanceInKM += totalDistanceCovered
				existing.MaxSpeed = float32(math.Max(float64(existing.MaxSpeed), float64(maxSpeed)))
				result[key] = existing
			} else {
				result[key] = AccelerationResult{
					PositiveAcceleration: totalPositiveAccelerationCase,
					NegativeAcceleration: totalNegativeAccelerationCase,
					TotalDistanceInKM:    totalDistanceCovered,
					MaxSpeed:             maxSpeed,
				}
			}
		}
	}

	for key, val := range result {
		fmt.Printf("key: %s , value: %v", key, val)
	}

	endTime := time.Now()
	fmt.Println("total time: ", endTime.Sub(startTime).Milliseconds())
	return nil, result
}
