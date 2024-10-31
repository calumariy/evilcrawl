package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func recursivelyAttackDirectory(baseSub string, baseDomain string, domain string, client *http.Client) {
	resp, err := client.Get(domain)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	// Try LFI

	attemptLFI(domain, client)

	// Dont continue if cant access page
	if resp.StatusCode >= 400 {
		fmt.Fprintln(os.Stdout, "[!] Webpage "+domain+" returned a "+strconv.Itoa(resp.StatusCode))
		return
	}

	page, err := html.Parse(resp.Body)
	handleErr(err)

	findNewInfo(baseSub, domain, page, client)
}

func findNewInfo(baseSub string, baseDomain string, n *html.Node, client *http.Client) {
	if n.Type == html.ElementNode && n.Data == "input" {
		attackInput(baseSub, baseDomain, n)
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
					recursivelyAttackDirectory(baseSub, baseDomain, newSub, client)
				} else if _, exists := directories[newDomain]; !exists {
					newDirectory(newDomain)
					recursivelyAttackDirectory(baseSub, baseDomain, newDomain, client)
				}
			} else if attr.Key == "href" && strings.HasPrefix(attr.Val, "/") && !strings.HasSuffix(baseDomain, attr.Val) {

				newDomain := attr.Val
				newSubArr := strings.Split(baseSub, "/")

				domain := newSubArr[0] + "//" + newSubArr[2] + newDomain

				if _, exists := directories[domain]; !exists {
					newDirectory(domain)
					recursivelyAttackDirectory(baseSub, baseDomain, domain, client)
				}
			}

		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findNewInfo(baseSub, baseDomain, c, client)
	}

	return
}
