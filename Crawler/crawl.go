package Crawler

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly/v2"
	"log"
	"strings"
)

type ProductDetails struct {
	Category                               string
	ProductName                            string
	Price                                  float64
	ImageUrl                               []string
	TitleOfDescription                     string
	GeneralDescriptionOfTheProduct         string
	GeneralDescriptionOfTheProductItemized string
}

func Crawler(baseUrl string) error {
	c := colly.NewCollector()
	//details := make(map[string]ProductDetails)

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	var nodes []*cdp.Node

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://shop.adidas.jp/products/IU0964/"),
		chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`, nil),
		//chromedp.Sleep(3*time.Second), // Wait for the page to load initially
		//chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`, nil),
		//chromedp.Sleep(5*time.Second), // Wait for additional content to load after scrolling
		//chromedp.Nodes("div.sizeChart.test-sizeChart", &nodes, chromedp.ByQueryAll),
	)

	if len(nodes) == 0 {
		fmt.Println("No nodes found")
	}

	fmt.Printf("Found %d nodes\n", len(nodes))

	// Output node content (for debugging)
	for _, node := range nodes {
		fmt.Printf("Node: %v\n", node)
	}

	if err != nil {
		fmt.Printf("err is : %v\n", err)
		return err
	}
	fmt.Printf("Ok I am evaluating!")
	//
	//log.Println("Starting navigation...")
	//err = chromedp.Run(ctx,
	//  chromedp.Navigate("https://shop.adidas.jp/products/IU0964/"),
	//  chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`, nil),
	//  chromedp.Sleep(3*time.Second),
	//  chromedp.WaitReady("body", chromedp.ByQuery),
	//)
	//if err != nil {
	//  log.Fatal("Navigation error: ", err)
	//}
	//

	detailsCollector := colly.NewCollector()

	log.Printf("Base Url is %s\n", baseUrl)

	//err := detailsCollector.Limit(&colly.LimitRule{
	//	DomainGlob:  "*",
	//	Parallelism: 1,
	//	Delay:       15 * time.Second, // Adjust the delay based on the page load time
	//})
	//if err != nil {
	//	fmt.Printf("err is: %v", err)
	//	return err
	//}

	detailsCollector.OnRequest(func(request *colly.Request) {
		request.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
		log.Printf("Visiting URL is: %s\n", request.URL.String())

	})

	//detailsCollector.OnRequest(func(request *colly.Request) {
	//	//log.Printf("Visiting URL is: %s\n", request.URL.String())
	//})

	detailsCollector.OnHTML("h4.heading.itemFeature.test-commentItem-subheading", func(e *colly.HTMLElement) {
		//titleOfDescription := e.Text
		//fmt.Println("titleOfDescription is: ", titleOfDescription)
	})

	detailsCollector.OnHTML("div.sizeChart.test-sizeChart", func(e *colly.HTMLElement) {
		fmt.Printf("chart is: %s\n", e.Text)
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

	detailsCollector.OnHTML("ul.articleFeatures.description_part", func(e *colly.HTMLElement) {
		generalDescriptionOfTheProductItemized := e.Text
		fmt.Printf("generalDescriptionOfTheProductItemized: %s\n", generalDescriptionOfTheProductItemized)
	})

	detailsCollector.OnHTML("div.commentItem-mainText.test-commentItem-mainText", func(e *colly.HTMLElement) {
		generalDescriptionOfTheProduct := e.Text
		fmt.Printf("generalDescriptionOfTheProduct: %s\n", generalDescriptionOfTheProduct)
	})

	detailsCollector.OnHTML("div.articleImageWrapper img", func(e *colly.HTMLElement) {
		// Extract the src attribute from the img tag
		imgSrc := e.Attr("src")
		if strings.HasPrefix(imgSrc, "/static/images/tools/itemCard_dummy") {
			return
		}
		// Print the image source
		//imgSrc = baseUrl + imgSrc
		//fmt.Printf("imgSrc is: %s, url is: %s\n", imgSrc, e.Request.URL.String())
	})

	detailsCollector.OnHTML("span.price-value.test-price-value", func(element *colly.HTMLElement) {
		//url := element.Request.URL.String()
		//if !details[url] {
		//	details[url] = ProductDetails{}
		//}
		//price := element.Text
		//fmt.Println("======> ", price, url)
	})

	detailsCollector.OnHTML("h1.itemTitle.test-itemTitle", func(e *colly.HTMLElement) {
		// Extract the text inside the <h1> tag
		//text := e.Text
		//fmt.Println("Category:", text)
	})

	// Go to the next page
	c.OnHTML("a.image_link.test-image_link", func(e *colly.HTMLElement) {
		//url := e.Attr("href")
		//urlToVisit := baseUrl + url
		//log.Printf("urlToVisit is: %s\n", urlToVisit)
		//err := detailsCollector.Visit(urlToVisit)
		//if err != nil {
		//	//log.Printf("Failed to visit: %s\n, got error: %s", urlToVisit, err)
		//	return
		//}
	})

	c.OnHTML("ul.lpc-ukLocalNavigation_itemList li:last-child a", func(e *colly.HTMLElement) {
		relativeURL := e.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(relativeURL)
		newUrlToVisit := baseUrl + absoluteURL
		err := c.Visit(newUrlToVisit)
		//log.Printf("New Url found to visit is: %s\n", newUrlToVisit)
		if err != nil {
			log.Println("Failed to visit:", newUrlToVisit, err)
		}
	})

	err2 := detailsCollector.Visit("https://shop.adidas.jp/products/IU0964/")
	if err2 != nil {
		return err2
	}

	//err := c.Visit(baseUrl + "/men")
	//if err != nil {
	//	return err
	//}

	return nil

}
