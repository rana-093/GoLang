package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
)

// LoginResponseDTO holds the response and the CSRF token
type LoginResponseDTO struct {
	CSRFToken string
}

// RetrieveCSRFTokenForUser retrieves the CSRF token from the login page using GoColly
func RetrieveCSRFTokenForUser(baseURL string, loginURLSuffix string) (*LoginResponseDTO, error) {
	// Create a new collector
	c := colly.NewCollector()

	var csrfToken string

	// On HTML callback
	c.OnHTML("input[name=_token]", func(e *colly.HTMLElement) {
		csrfToken = e.Attr("value")
	})

	// Error handling for the request
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	// Make a GET request to the login page
	err := c.Visit(baseURL + loginURLSuffix)
	if err != nil {
		return nil, err
	}

	if csrfToken == "" {
		return nil, fmt.Errorf("CSRF token not found")
	}

	return &LoginResponseDTO{
		CSRFToken: csrfToken,
	}, nil
}

func main() {
	baseURL := "https://gps.carcopolo.com"
	loginURLSuffix := "/login"

	loginResponse, err := RetrieveCSRFTokenForUser(baseURL, loginURLSuffix)
	if err != nil {
		log.Fatalf("Error retrieving CSRF token: %v", err)
	}

	fmt.Printf("CSRF Token: %s\n", loginResponse.CSRFToken)
}
