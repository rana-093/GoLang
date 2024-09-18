package csvparsing

import (
	"encoding/csv"
	"fmt"
	"os"
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
	TimeStamp     string
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
	const layout, workers, interval = "2006-01-02 15:04:05", 50, 5
	chunkSize := len(data) / workers
	fmt.Printf("chunksize is: %d\n", chunkSize)
	speedMap := make(map[time.Time][]SpeedInfo)

	processChunk := func(data [][]string, wg *sync.WaitGroup) {
		defer wg.Done()
		for _, row := range data {
			timeStamp := row[0]
			speed, _ := strconv.ParseFloat(row[5], 32)
			distance, _ := strconv.ParseFloat(row[13], 32)
			totalDistance, _ := strconv.ParseFloat(row[14], 32)
			lat, _ := strconv.ParseFloat(row[1], 32)
			lon, _ := strconv.ParseFloat(row[2], 32)

			//fmt.Printf("time: %s, speed: %0.1f, distance: %0.1f\n", timeStamp, speed, distance)

			speedInfo := SpeedInfo{
				TimeStamp:     timeStamp,
				Speed:         float32(speed),
				Distance:      float32(distance),
				TotalDistance: float32(totalDistance),
				Lat:           float32(lat),
				Lon:           float32(lon),
			}
			datetime, _ := time.Parse(layout, timeStamp)
			key := datetime.Truncate(time.Duration(interval) * time.Second)
			//fmt.Println("key: ", key)
			lock.Lock()
			speedMap[key] = append(speedMap[key], speedInfo)
			lock.Unlock()
		}
	}

	for i := 0; i < workers; i++ {
		startIndex := i * chunkSize
		if startIndex == 0 {
			startIndex = 1
		}
		endIndex := (i + 1) * chunkSize
		if i == workers-1 {
			endIndex = len(data)
		}
		wg.Add(1)
		go processChunk(data[startIndex:endIndex], &wg)
	}

	wg.Wait()
	for key, value := range speedMap {
		fmt.Printf("key: %s, value: %v\n", key, value)
	}
	endTime := time.Now()
	fmt.Println("total time: ", endTime.Sub(startTime).Milliseconds())

	return nil
}
