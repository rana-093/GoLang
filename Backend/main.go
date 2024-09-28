package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	region = "ap-south-1"
	bucket = "report-sl"
)

type VehicleReport struct {
	TotalAcceleration    int
	TotalDeceleration    int
	AccelerationPer100KM float32
	DecelerationPer100KM float32
	MaxSpeed             float32
	TotalDistanceCovered float32
	DeviceId             string
}

func downloadAndGetFileName(svc *s3.S3, key string) (string, error) {
	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", fmt.Errorf("Error: %w\n", err)
	}

	defer result.Body.Close()
	newUUID := uuid.New()
	tempFile, err := os.CreateTemp("", fmt.Sprintf("downloaded-%s-*.xlsx", newUUID.String()))

	if err != nil {
		return "", fmt.Errorf("Error: %w\n", err)
	}

	if _, err := io.Copy(tempFile, result.Body); err != nil {
		return "", fmt.Errorf("Error: %w\n", err)
	}

	if err := tempFile.Close(); err != nil {
		return "", fmt.Errorf("Error: %w\n", err)
	}
	return tempFile.Name(), nil
}

func calculateDailyUsageReport(svc *s3.S3, key string) map[string]float64 {
	fileName, err := downloadAndGetFileName(svc, key)
	defer os.Remove(fileName)
	if err != nil {
		return nil
	}
	_, totalDistanceByDevice := ParseXSLXFromDailyUsageReport(fileName)
	return totalDistanceByDevice
}

func calculateObjectHistoryForGettingAccAndDcc(svc *s3.S3, key string) map[string]AccelerationResult {
	fileName, err := downloadAndGetFileName(svc, key)
	defer os.Remove(fileName)
	if err != nil {
		return nil
	}
	_, accAndDccReport := ParseXSLXFromObjectHistoryReport(fileName)
	return accAndDccReport
}

func handleDownLoadFile(w http.ResponseWriter, r *http.Request) {
	objectHistoryUrl := r.URL.Query().Get("objectHistoryUrl")
	dailyUsageUrl := r.URL.Query().Get("dailyUsageUrl")

	slog.Info("Object History URL is: ", objectHistoryUrl, " and daily usage url is: ", dailyUsageUrl)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}

	svc := s3.New(sess)

	var vechicleReport []VehicleReport

	objectHistoryReport := calculateObjectHistoryForGettingAccAndDcc(svc, objectHistoryUrl)
	dailyUsageReport := calculateDailyUsageReport(svc, dailyUsageUrl)

	for key, value := range objectHistoryReport {
		totalDistance := dailyUsageReport[key]
		totalPositiveAcc := value.PositiveAcceleration
		totalNegativeAcc := value.NegativeAcceleration
		positiveAccPer100KM, negativeAccPer100KM := 0.0, 0.0
		if totalDistance != 0.0 {
			positiveAccPer100KM = (float64(totalPositiveAcc) / totalDistance) * float64(100)
			negativeAccPer100KM = (float64(totalNegativeAcc) / totalDistance) * float64(100)
		}
		report := VehicleReport{
			TotalAcceleration:    totalPositiveAcc,
			TotalDeceleration:    totalNegativeAcc,
			TotalDistanceCovered: float32(totalDistance),
			AccelerationPer100KM: float32(positiveAccPer100KM),
			DecelerationPer100KM: float32(negativeAccPer100KM),
			MaxSpeed:             value.MaxSpeed,
			DeviceId:             key,
		}
		vechicleReport = append(vechicleReport, report)
	}

	for key, value := range dailyUsageReport {
		fmt.Printf("key: %s, value: %f\n", key, value)
	}

	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(objectHistoryReport)

	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	slog.Info("Listening on port 8088")
	http.HandleFunc("/", handleDownLoadFile)
	http.ListenAndServe(":8088", nil)
}
