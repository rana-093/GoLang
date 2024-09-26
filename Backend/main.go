package Backend

import (
	"GoLang/csvparsing"
	"encoding/json"
	"fmt"
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

func searchFileWithKeyAndDownloadFile(key string) {

}

func handleDownLoadFile(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("objectHistoryUrl")
	slog.Info("Object History URL is: ", key)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}

	svc := s3.New(sess)

	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Println("Error getting object:", err)
		return
	}
	defer result.Body.Close()

	tempFile, err := os.CreateTemp("", "downloaded-*.xlsx")

	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return
	}

	defer os.Remove(tempFile.Name())

	if _, err := io.Copy(tempFile, result.Body); err != nil {
		fmt.Println("Error saving the file:", err)
		return
	}

	if err := tempFile.Close(); err != nil {
		fmt.Println("Error closing temp file:", err)
		return
	}

	err = csvparsing.ParseXSLX(tempFile.Name())
	if err != nil {
		slog.Error("Error parsing CSV file:", err)
		return
	}
	err, output := ParseXSLX(tempFile.Name())
	if err != nil {
		slog.Error("Error parsing CSV file:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(output)

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
