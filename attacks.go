package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
	"log"
	"time"
)

func attackInput(baseSub string, baseDomain string, n *html.Node) {
	fmt.Println("[!] Input found at " + baseDomain)
	for _, attr := range n.Attr {
		if attr.Key == "type" && attr.Val == "password" && !isAuthorised {
			fmt.Println("[!] found password field on: " + baseDomain + "\nWant to use authorisation? Use the -a flag!")
		}
	}

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// navigate to a page, wait for an element, click
	var example string
	err := chromedp.Run(ctx,
		chromedp.Navigate(baseDomain),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.WaitVisible(`body > footer`),
		// find and click "Example" link
		chromedp.Click(`#example-After`, chromedp.NodeVisible),
		// retrieve the text of the textarea
		chromedp.Value(`#example-After textarea`, &example),
	)
	if err != nil {
		log.Fatal(err)
	}
}
