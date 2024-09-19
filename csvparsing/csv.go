package csvparsing

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
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
	PositiveAcceleration int
	NegativeAcceleration int
}

func millisecondsToHours(ms int64) float32 {
	return float32(ms) / 3_600_000
}

func millisecondsToSeconds(ms int64) float32 {
	return float32(ms) / 1_000
}

func ParseCSV(fileName string) error {
	startTime := time.Now()
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var wg sync.WaitGroup
	var lock sync.Mutex

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	const Layout, Workers, Interval = "2006-01-02 15:04:05", 50, 5
	const AccelerateThresholdInKMPerSecond = 6
	chunkSize := len(data) / Workers
	fmt.Printf("chunksize is: %d\n", chunkSize)
	speedMap := make(map[time.Time][]SpeedInfo)

	processChunk := func(data [][]string, wg *sync.WaitGroup) {
		for _, row := range data {
			timeStamp, err := time.Parse(Layout, row[0])
			if err != nil {
				panic(err)
			}
			speed, _ := strconv.ParseFloat(row[5], 32)
			distance, _ := strconv.ParseFloat(row[13], 32)
			totalDistance, _ := strconv.ParseFloat(row[14], 32)
			lat, _ := strconv.ParseFloat(row[1], 32)
			lon, _ := strconv.ParseFloat(row[2], 32)
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
			speedMap[key] = append(speedMap[key], speedInfo)
			lock.Unlock()
		}
	}

	handleTimeStampWiseAggregation := func(data []SpeedInfo) (bool, bool) {
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
		if isAcceleration && velocityInKMPerSecond >= AccelerateThresholdInKMPerSecond {
			return true, false
		} else if isNegativeAcceleration && velocityInKMPerSecond >= AccelerateThresholdInKMPerSecond {
			return false, true
		}
		return false, false
	}

	processChunk(data[1:len(data)], &wg)

	totalPositiveAccelerationCase, totalNegativeAccelerationCase := 0, 0

	const NoOfThreads = 15
	eachThread := len(speedMap) / NoOfThreads

	chunkWiseData := make(chan AccelerationResult, NoOfThreads)

	var keysGroup [][]SpeedInfo
	idx := 0
	for _, val := range speedMap {
		idx++
		keysGroup = append(keysGroup, val)
		if (idx+1)%eachThread == 0 || idx == len(speedMap)-1 {
			wg.Add(1)
			go func(group [][]SpeedInfo) {
				defer wg.Done()
				positiveAcc, negativeAcc := 0, 0
				for _, row := range group {
					if len(row) < 2 {
						continue
					}
					posAcc, negAcc := handleTimeStampWiseAggregation(row)
					if posAcc {
						positiveAcc++
					} else if negAcc {
						negativeAcc++
					}
				}
				chunkWiseData <- AccelerationResult{
					PositiveAcceleration: positiveAcc,
					NegativeAcceleration: negativeAcc,
				}
			}(keysGroup)
			keysGroup = [][]SpeedInfo{}
		}
	}

	go func() {
		wg.Wait()
		close(chunkWiseData)
	}()

	for result := range chunkWiseData {
		totalPositiveAccelerationCase += result.PositiveAcceleration
		totalNegativeAccelerationCase += result.NegativeAcceleration
	}

	endTime := time.Now()
	fmt.Println("total time: ", endTime.Sub(startTime).Milliseconds())
	fmt.Printf("Total Harsh Brakes: %d. Total Harsh Acceleration: %d\n", totalNegativeAccelerationCase, totalPositiveAccelerationCase)

	return nil
}
