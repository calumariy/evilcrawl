package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func recursivelyAttackDirectory(baseSub string, domain string, client *http.Client, wg *sync.WaitGroup) {
	resp, err := client.Get(domain)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	// Try LFI
	wg.Add(1)
	go func() {
		defer wg.Done()
		attemptURLXSS(domain)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		attemptLFI(domain, client)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		attemptSSTI(domain, client)
	}()

	// Dont continue if cant access page
	if resp.StatusCode >= 400 {
		fmt.Fprintln(os.Stdout, "[!] Webpage "+domain+" returned a "+strconv.Itoa(resp.StatusCode))
	}

	page, err := html.Parse(resp.Body)
	handleErr(err)

	findNewInfo(baseSub, domain, page, client, wg)
}

func findNewInfo(baseSub string, baseDomain string, n *html.Node, client *http.Client, wg *sync.WaitGroup) {
	if n.Type == html.ElementNode && n.Data == "input" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			attackInput(baseDomain, n)
		}()
	}

	if n.Type == html.ElementNode && n.Data == "form" {
		action := ""
		method := ""

		for _, attr := range n.Attr {
			if attr.Key == "action" {
				action = attr.Val
			}
			if attr.Key == "method" {
				method = attr.Val
			}
		}

		if (method == "GET" || method == "get") && strings.HasPrefix(action, "/") {
			newDomain := action
			newSubArr := strings.Split(baseSub, "/")

			domain := newSubArr[0] + "//" + newSubArr[2] + newDomain

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "input" {
					inputType := "invalid"
					inputName := "invalid"
					for _, attr := range c.Attr {
						if attr.Key == "type" {
							inputType = attr.Val
						}
						if attr.Key == "name" {
							inputName = attr.Val
						}
					}
					if inputType == "text" {
						if _, exists := directories[domain+"?"+inputName+"="]; !exists {
							newDirectory(domain + "?" + inputName + "=")
							recursivelyAttackDirectory(baseSub, domain+"?"+inputName+"=", client, wg)
						}
					}
				}
			}
		}
		if method == "post" || method == "POST" {
			var newDomain string
			newSubArr := strings.Split(baseSub, "/")
			if action != "" {
				if strings.HasPrefix(action, "/") {
					newDomain = newSubArr[0] + "//" + newSubArr[2] + action
				} else if !strings.HasPrefix(action, "http") {
					newDomain = baseDomain + action
				}
			} else {
				newDomain = baseDomain
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "input" {
					inputType := "invalid"
					inputName := "invalid"
					for _, attr := range c.Attr {
						if attr.Key == "type" {
							inputType = attr.Val
						}
						if attr.Key == "name" {
							inputName = attr.Val
						}
					}
					if inputType == "text" {
						attemptSSTIPost(newDomain, client, inputName)
					}
				}
			}
		}
	}
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {

			if attr.Key == "href" && strings.Contains(attr.Val, baseSub) {
				newDomain := attr.Val
				newDomainArr := strings.Split(newDomain, "/")
				newSub := newDomainArr[0] + "//" + newDomainArr[2]

				if _, exists := subdomains[newSub]; !exists {
					newSubdomain(newSub)
					newDirectory(newSub)
					recursivelyAttackDirectory(baseSub, newSub, client, wg)
				} else if _, exists := directories[newDomain]; !exists {
					newDirectory(newDomain)
					recursivelyAttackDirectory(baseSub, newDomain, client, wg)
				}
			} else if attr.Key == "href" && strings.HasPrefix(attr.Val, "/") && !strings.HasSuffix(baseDomain, attr.Val) {

				newDomain := attr.Val
				newSubArr := strings.Split(baseSub, "/")

				domain := newSubArr[0] + "//" + newSubArr[2] + newDomain

				if _, exists := directories[domain]; !exists {
					newDirectory(domain)
					recursivelyAttackDirectory(baseSub, domain, client, wg)
				}
			} else if attr.Key == "href" && len(strings.Split(attr.Val, ".")) == 1 && !strings.HasPrefix(attr.Val, "/") {

				newDomain := attr.Val

				domain := baseDomain + "/" + newDomain

				if _, exists := directories[domain]; !exists {
					newDirectory(domain)
					recursivelyAttackDirectory(baseSub, domain, client, wg)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findNewInfo(baseSub, baseDomain, c, client, wg)
	}

	return
}
