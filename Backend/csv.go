package Backend

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"sort"
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
}

func millisecondsToHours(ms int64) float32 {
	return float32(ms) / 3_600_000
}

func millisecondsToSeconds(ms int64) float32 {
	return float32(ms) / 1_000
}

func ParseXSLX(fileName string) (error, *map[string]AccelerationResult) {
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
	var lock sync.Mutex

	const Layout, Workers, Interval = "02-01-2006 15:04:05", 50, 20
	const AccelerateThresholdInKMPerSecond = 3
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

	processChunk := func(data [][]string, deviceId string) {
		for _, row := range data {
			if len(row) < 18 || strings.EqualFold(row[0], "Device") {
				continue
			}
			timeStamp, err := time.Parse(Layout, row[2])
			if err != nil {
				panic(err)
			}
			speed, _ := strconv.ParseFloat(row[7], 32)
			distance, _ := strconv.ParseFloat(row[16], 32)
			totalDistance, _ := strconv.ParseFloat(row[17], 32)
			lat, _ := strconv.ParseFloat(row[3], 32)
			lon, _ := strconv.ParseFloat(row[4], 32)
			motion, _ := strconv.ParseBool(row[15])

			speedInfo := SpeedInfo{
				TimeStamp:     timeStamp,
				Speed:         float32(speed),
				Distance:      float32(distance),
				TotalDistance: float32(totalDistance),
				Lat:           float32(lat),
				Lon:           float32(lon),
				Motion:        motion,
			}

			key := timeStamp.Truncate(time.Duration(Interval) * time.Second)
			lock.TryLock()
			speedMap[deviceId][key] = append(speedMap[deviceId][key], speedInfo)
			lock.Unlock()
		}
	}

	handleTimeStampWiseAggregation := func(data []SpeedInfo) (bool, bool, float32) {
		sort.Slice(data, func(i, j int) bool {
			return data[i].TimeStamp.Before(data[j].TimeStamp)
		})
		N := len(data)
		isAcceleration, isNegativeAcceleration := true, true
		totalDistanceCoveredInKM := data[0].Distance
		for i := 1; i < N; i++ {
			if data[i].Speed > data[i-1].Speed {
				isNegativeAcceleration = false
			} else if data[i].Speed < data[i-1].Speed {
				isAcceleration = false
			}
			totalDistanceCoveredInKM += data[i].Distance
		}
		totalMilliSeconds := data[N-1].TimeStamp.Sub(data[0].TimeStamp).Milliseconds()
		velocityInKMPerSecond := totalDistanceCoveredInKM / millisecondsToSeconds(totalMilliSeconds)
		if totalMilliSeconds == 0 || velocityInKMPerSecond == 0 {
			return false, false, totalDistanceCoveredInKM
		}
		if isAcceleration && velocityInKMPerSecond >= AccelerateThresholdInKMPerSecond {
			return true, false, totalDistanceCoveredInKM
		} else if isNegativeAcceleration && velocityInKMPerSecond >= AccelerateThresholdInKMPerSecond {
			return false, true, totalDistanceCoveredInKM
		}
		return false, false, totalDistanceCoveredInKM
	}

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
			processChunk(chunk, device)
		}(chunkWiseDevices[i], chunk)
	}

	wg.Wait()

	const NoOfThreads = 15
	totalLength := 0
	for _, val := range speedMap {
		totalLength += len(val)
	}
	eachThread := totalLength / NoOfThreads

	chunkWiseData := make(chan map[string][]AccelerationResult, NoOfThreads)

	fmt.Println("eachThread is:", eachThread, "speedMap: ", len(speedMap))

	var keysGroup [][]SpeedInfo
	idx := 0

	for i, timeMap := range speedMap {
		for _, val := range timeMap {
			idx++
			keysGroup = append(keysGroup, val)
			if (idx+1)%eachThread == 0 || idx == len(speedMap)-1 {
				wg.Add(1)
				go func(group [][]SpeedInfo, deviceId string) {
					defer wg.Done()
					positiveAcc, negativeAcc, totalDistance := 0, 0, float32(0)
					for _, row := range group {
						if len(row) < 2 {
							continue
						}
						posAcc, negAcc, totalDistanceInKM := handleTimeStampWiseAggregation(row)
						if posAcc {
							positiveAcc++
						} else if negAcc {
							negativeAcc++
						}
						totalDistance += totalDistanceInKM
					}
					data := map[string][]AccelerationResult{
						i: {
							{
								PositiveAcceleration: positiveAcc,
								NegativeAcceleration: negativeAcc,
								TotalDistanceInKM:    totalDistance,
							},
						},
					}

					chunkWiseData <- data
				}(keysGroup, i)
				keysGroup = [][]SpeedInfo{}
			}
		}
	}

	go func() {
		wg.Wait()
		close(chunkWiseData)
	}()

	result := make(map[string]AccelerationResult)

	for data := range chunkWiseData {
		for key, accResults := range data {
			totalPositiveAccelerationCase, totalNegativeAccelerationCase, totalDistanceCovered := 0, 0, float32(0)
			for _, accResult := range accResults {
				totalPositiveAccelerationCase += accResult.PositiveAcceleration
				totalNegativeAccelerationCase += accResult.NegativeAcceleration
				totalDistanceCovered += accResult.TotalDistanceInKM
			}
			if existing, exists := result[key]; exists {
				existing.PositiveAcceleration += totalPositiveAccelerationCase
				existing.NegativeAcceleration += totalNegativeAccelerationCase
				existing.TotalDistanceInKM += totalDistanceCovered
				result[key] = existing
			} else {
				result[key] = AccelerationResult{
					PositiveAcceleration: totalPositiveAccelerationCase,
					NegativeAcceleration: totalNegativeAccelerationCase,
					TotalDistanceInKM:    totalDistanceCovered,
				}
			}
		}
	}

	for key, val := range result {
		fmt.Printf("key: %s , value: %v", key, val)
	}

	endTime := time.Now()
	fmt.Println("total time: ", endTime.Sub(startTime).Milliseconds())
	return nil, &result
}
