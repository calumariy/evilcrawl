package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func recursivelyAttackDirectory(baseDomain string, domain string) {
	resp, err := http.Get(domain)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// Dont continue if cant access page
	if resp.StatusCode >= 400 {
		fmt.Fprintln(os.Stdout, "[!] Webpage "+domain+" returned a "+strconv.Itoa(resp.StatusCode))
		return
	}

	page, err := html.Parse(resp.Body)
	handleErr(err)

	findNewInfo(removeProtocol(baseDomain), page)
}

func findNewInfo(baseDomain string, n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" && strings.Contains(attr.Val, baseDomain) {
				newDomain := attr.Val
				newDomainArr := strings.Split(newDomain, "/")
				newSub := newDomainArr[0] + "//" + newDomainArr[2]
				if _, exists := subdomains[newSub]; !exists {
					newSubdomain(newSub)
					recursivelyAttackDirectory(baseDomain, newSub)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findNewInfo(baseDomain, c)
	}

	return
}
