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
	subMu.Lock()
	defer subMu.Unlock()

	if _, exists := subdomains[domain]; !exists {
		subdomains[domain] = struct{}{}

		fmt.Fprintln(outFile, domain)
	}
}

func newDirectory(domain string) {
	dirMu.Lock()
	defer dirMu.Unlock()

	if _, exists := directories[domain]; !exists {
		directories[domain] = struct{}{}

		fmt.Fprintln(outFile, domain)
	}
}
