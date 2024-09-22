package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func compileSubdomains(mainDomain string, wordlist string) ([]string, error) {
	// Check that domain exists
	_, err := http.Get(mainDomain)
	handleErr(err)
	passiveSubdomains, err := passiveDNSRecon(mainDomain)
	var activeSubdomains []string
	if wordlist == "" {
		fmt.Fprintln(os.Stderr, "no wordlist detected, skipping brute force")
	} else {
		activeSubdomains, err = activeDNSRecon(mainDomain, wordlist)
		handleErr(err)
	}

	joinSubdomains(passiveSubdomains, activeSubdomains)

	subdomains := []string{mainDomain}

	return subdomains, nil
}

func passiveDNSRecon(mainDomain string) ([]string, error) {
	return []string{mainDomain}, nil
}

func activeDNSRecon(mainDomain string, wordlistPath string) ([]string, error) {
	wordlist, err := os.Open(wordlistPath)
	handleErr(err)
	wordlistScanner := bufio.NewScanner(wordlist)
	wordlistScanner.Split(bufio.ScanLines)
	isHttp := false
	isHttps := false
	domain := mainDomain

	// Remove protocol prefix
	if strings.HasPrefix(mainDomain, "http://") {
		isHttp = true
		domain = strings.TrimPrefix(domain, "http://")
	} else if strings.HasPrefix(mainDomain, "https://") {
		isHttps = true
		domain = strings.TrimPrefix(domain, "https://")
	}

	for wordlistScanner.Scan() {
		subdomainWord := wordlistScanner.Text()

	}
	return []string{mainDomain}, nil
}

func joinSubdomains(subdomains1 []string, subdomains2 []string) ([]string, error) {
	return []string{}, nil
}
