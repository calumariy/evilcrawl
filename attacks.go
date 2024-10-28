package main

import (
	"fmt"
	"golang.org/x/net/html"
)

func attackInput(baseSub string, baseDomain string, n *html.Node) {
	for _, attr := range n.Attr {
		if attr.Key == "type" && attr.Val == "password" && !isAuthorised {
			fmt.Println("found password field on: " + baseDomain)
		}
	}
}
