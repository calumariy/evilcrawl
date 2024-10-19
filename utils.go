package main

import (
	"fmt"
	"os"
	"strings"
)

func handleErr(err error) {
	if err == nil {
		return
	}

	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func removeProtocol(domain string) string {
	if strings.HasPrefix(domain, "http://") {
		domain = strings.TrimPrefix(domain, "http://")
	} else if strings.HasPrefix(domain, "https://") {
		domain = strings.TrimPrefix(domain, "https://")
	}

	return domain
}

func newSubdomain(domain string) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := subdomains[domain]; !exists {
		subdomains[domain] = struct{}{}

		fmt.Fprintln(outFile, domain)
	}
}
