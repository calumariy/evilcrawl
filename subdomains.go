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
	domain := mainDomain
	proto := ""

	// Remove protocol prefix
	if strings.HasPrefix(mainDomain, "http://") {
		domain = strings.TrimPrefix(domain, "http://")
		proto = "http://"
	} else if strings.HasPrefix(mainDomain, "https://") {
		domain = strings.TrimPrefix(domain, "https://")
		proto = "https://"
	}

	// List of subdomains that give a response
	subdomains := []string{}
	for wordlistScanner.Scan() {
		subdomainWord := wordlistScanner.Text()
		subdomain := proto + subdomainWord + "." + domain
		_, err := http.Get(subdomain)
		if err != nil {
			continue
		}
		subdomains = append(subdomains, subdomain)
	}

	return []string{mainDomain}, nil
}

func joinSubdomains(subdomains1 []string, subdomains2 []string) ([]string, error) {
	return []string{}, nil
}
