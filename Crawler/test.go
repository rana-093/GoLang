package Crawler

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly/v2"
	"log"
	"time"
)

func Main() {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create a new ChromeDP context
	ctx, cancel := chromedp.NewContext(
		timeoutCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	var htmlContent string
	// Run the ChromeDP tasks
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://shop.adidas.jp/products/IU0964/"),
		chromedp.Sleep(3*time.Second), // Initial wait
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Here I am with some scrolling - 1!!")

	for i := 0; i < 3; i++ {
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`window.scrollBy(0, window.innerHeight)`, nil),
			chromedp.Sleep(3*time.Second), // Wait for new content to load
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Here I am with some scrolling - 2!!")

	err = chromedp.Run(ctx,
		chromedp.WaitVisible(`div.sizeChart.test-sizeChart`, chromedp.ByQuery),
		chromedp.OuterHTML(`div.sizeChart.test-sizeChart`, &htmlContent),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Output the captured HTML content
	fmt.Println("Captured HTML content:")
	fmt.Println(htmlContent)

	c := colly.NewCollector()

	c.OnHTML("div.sizeChart.test-sizeChart", func(e *colly.HTMLElement) {
		fmt.Println("Found size chart!")
		e.ForEach("thead tr", func(_ int, row *colly.HTMLElement) {
			row.ForEach("th", func(_ int, header *colly.HTMLElement) {
				fmt.Printf("Header: %s\n", header.Text)
			})
		})
		e.ForEach("tbody tr", func(_ int, row *colly.HTMLElement) {
			row.ForEach("td", func(_ int, cell *colly.HTMLElement) {
				fmt.Printf("Cell: %s\n", cell.Text)
			})
		})
	})

	fmt.Println("html: ", htmlContent)

	// Pass HTML content to colly
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "text/html")
	})

	// Process the HTML content with colly
	//err = c.Parse(htmlContent)
	//if err != nil {
	//	log.Fatal(err)
	//}
}
