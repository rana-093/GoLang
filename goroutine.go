package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

type PaymentMethod struct {
	UserId        string  `json:"userId"`
	BankName      string  `json:"bankName"`
	BankAccNo     *string `json:"bankAccNo"`
	BranchName    *string `json:"branchName"`
	AccHolderName *string `json:"accHolderName"`
	PaymentMethod *string `json:"paymentMethod"`
}

type UserData struct {
	Id                           int            `json:"id"`
	UserID                       string         `json:"userId"`
	Name                         string         `json:"name"`
	Phone                        string         `json:"phone"`
	Email                        *string        `json:"email"`
	UserShortCode                *string        `json:"userShortCode"`
	Enabled                      bool           `json:"enabled"`
	Roles                        []string       `json:"roles"`
	Token                        *string        `json:"token"`
	DriverFilesUrls              *string        `json:"driverFilesUrls"`
	Designation                  *string        `json:"designation"`
	CreatedOn                    string         `json:"createdOn"`
	ContactPerson                *string        `json:"contactPerson"`
	TruckStand                   *string        `json:"truckStand"`
	PaymentMethod                *PaymentMethod `json:"paymentMethod"`
	Customer                     *string        `json:"customer"`
	DrivingLicenseExpiryDate     *string        `json:"drivingLicenseExpiryDate"`
	CompletedTripCountsForDriver *int           `json:"completedTripCountsForDriver"`
	CreatedBy                    *string        `json:"createdBy"`
	UpdatedBy                    *string        `json:"updatedBy"`
	Status                       *string        `json:"status"`
}

type Response struct {
	Status     string     `json:"status"`
	StatusCode int        `json:"statusCode"`
	Message    string     `json:"message"`
	Data       []UserData `json:"data"`
}

func prepareAndCallEndpoint(url string, headers map[string]string, wg *sync.WaitGroup) (*Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	fmt.Println("Got new Request")

	defer wg.Done()

	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing response body: %v", err)
	}
	return &response, nil
}

func checkGoRoutinesAndChannels() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	url, token := os.Getenv("URL"), os.Getenv("TOKEN")

	ch := make(chan *Response)

	fmt.Println("url is : %s, token is : %s\n", url, token)

	headers := map[string]string{
		"Authorization": token,
	}

	start := time.Now()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			response, err := prepareAndCallEndpoint(url, headers, &wg)
			if err != nil {
				fmt.Errorf("Error preparing endpoint: %v", err)
				ch <- nil
				return
			}
			ch <- response
		}()
	}

	wg.Wait()
	response := <-ch
	end := time.Now()
	fmt.Println("Time taken is : ", end.Sub(start))
	fmt.Println("total respose is , ", response)
}
