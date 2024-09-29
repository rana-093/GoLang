package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jackc/pgx/v4"
)

const (
	region = "ap-south-1"
	bucket = "report-sl"
)

type TravelMetaData struct {
	DeviceId            string
	From                time.Time
	To                  time.Time
	ObjectHistoryUrl    string
	DailyUsageReportUrl string
	CompanyId           *string
	TrackerId           *string
	MonthOfReport       string
}

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

func saveReportToDB(vehicleReport []VehicleReport, travelMetaData TravelMetaData) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	fmt.Printf("dsn is %s\n", dsn)
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v\n", err)
	}
	defer conn.Close(context.Background())

	fmt.Println("Successfully connected to the database!")

	for _, report := range vehicleReport {
		query := `INSERT INTO cargo_reports 
    			(harsh_acc_count, harsh_dcc_count, acc_per_100_km, 
    			 dcc_per_100_km, distance, from_date, to_date, 
    			 device_id, month_of_report) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
		_, err := conn.Exec(
			context.Background(),
			query,
			report.TotalAcceleration,
			report.TotalDeceleration,
			report.AccelerationPer100KM,
			report.DecelerationPer100KM,
			report.TotalDistanceCovered,
			travelMetaData.From,
			travelMetaData.To,
			travelMetaData.DeviceId,
			travelMetaData.MonthOfReport,
		)
		if err != nil {
			log.Fatalf("Error executing insert query: %v\n", err)
		}
	}

	var greeting string
	err = conn.QueryRow(context.Background(), "SELECT 'Hello, PostgreSQL!'").Scan(&greeting)
	if err != nil {
		log.Fatalf("QueryRow failed: %v\n", err)
	}

	fmt.Println(greeting)
}

func handleDownLoadFile(w http.ResponseWriter, r *http.Request) {
	objectHistoryUrl := r.URL.Query().Get("objectHistoryUrl")
	dailyUsageUrl := r.URL.Query().Get("dailyUsageUrl")
	fromDate := r.URL.Query().Get("fromDate")
	toDate := r.URL.Query().Get("toDate")
	deviceId := r.URL.Query().Get("deviceId")
	trackerId := r.URL.Query().Get("trackerId")

	layout := "2006-01-02"
	parsedFromDate, err := time.Parse(layout, fromDate)
	if err != nil {
		_ = fmt.Errorf("Error parsing fromDate: %v\n", err)
	}
	parsedToDate, err := time.Parse(layout, toDate)
	if err != nil {
		_ = fmt.Errorf("Error parsing toDate: %v\n", err)
	}

	slog.Info("Object History URL is: ", objectHistoryUrl, " and daily usage url is: ", dailyUsageUrl, " from date: ", fromDate, " to date: ", toDate)

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

	response, err := json.Marshal(vechicleReport)

	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	travelMetaData := TravelMetaData{
		From:                parsedFromDate,
		To:                  parsedToDate,
		DeviceId:            deviceId,
		ObjectHistoryUrl:    objectHistoryUrl,
		DailyUsageReportUrl: dailyUsageUrl,
		TrackerId:           &trackerId,
		MonthOfReport:       parsedFromDate.Month().String(),
	}

	saveReportToDB(vechicleReport, travelMetaData)

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	slog.Info("Listening on port 8088")
	http.HandleFunc("/", handleDownLoadFile)
	http.ListenAndServe(":8088", nil)
}
