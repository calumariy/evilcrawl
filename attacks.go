package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
)

func attackInput(baseSub string, baseDomain string, n *html.Node) {
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
	ctx, cancel = context.WithTimeout(ctx, time.Second)
	defer cancel()

	// navigate to a page, wait for an element, click
	err := chromedp.Run(ctx,
		setCookie(baseDomain),
		// Navigate to domain
		chromedp.Navigate(baseDomain),
		// Load webpage
		chromedp.WaitReady("body", chromedp.ByQuery),
		// find and click "Example" link
		chromedp.Click(`button[type="submit"]`),
	)
	if err != nil {
		if err != context.DeadlineExceeded {
			handleErr(err)
		}
		return
	}

	// Button found, try to do more
	fmt.Println("[!] submit buuton found at ")
}

func setCookie(domain string) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {

		if auth == "" {
			return nil
		}

		domainPath := strings.Split(domain, "/")

		mainDomain := domainPath[2]
		mainDomainArr := strings.Split(mainDomain, ".")
		cookieDomain := "." + mainDomainArr[len(mainDomainArr)-2] + "." + mainDomainArr[len(mainDomainArr)-1]

		cookieStr := strings.Split(auth, ":")

		err := network.SetCookie(cookieStr[0], cookieStr[1]).WithDomain(cookieDomain).Do(ctx)
		if err != nil {
			fmt.Errorf("error: %w")
		}
		return err
	})
}
