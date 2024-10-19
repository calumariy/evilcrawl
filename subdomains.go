package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func compileSubdomains(mainDomain string, wordlist string, customSubdomainsFile string) {
	// Check that domain exists
	_, err := http.Get(mainDomain)
	handleErr(err)
	passiveDNSRecon(mainDomain)

	if wordlist == "" {
		fmt.Fprintln(os.Stderr, "[!] no wordlist detected, skipping brute force")
	} else {
		activeDNSRecon(mainDomain, wordlist)
	}

	if customSubdomainsFile != "" {
		customSubdomainList, err := os.Open(customSubdomainsFile)
		handleErr(err)
		customSubdomainScanner := bufio.NewScanner(customSubdomainList)
		customSubdomainScanner.Split(bufio.ScanLines)

		for customSubdomainScanner.Scan() {
			customSubdomain := customSubdomainScanner.Text()
			if !strings.Contains(customSubdomain, "http://") {
				customSubdomain = "http://" + customSubdomain
			}
			_, err := http.Get(customSubdomain)
			if err != nil {
				continue
			}

			newSubdomain(customSubdomain)
		}
		handleErr(err)
	} else {
		fmt.Fprintln(os.Stderr, "[!] no custom subdomains supplied")
	}

	newSubdomain(mainDomain)
	return
}

func passiveDNSRecon(mainDomain string) {
	fmt.Fprintln(os.Stderr, "[!] Passive dns recon not supported yet")
	return
}

func activeDNSRecon(mainDomain string, wordlistPath string) {
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
	for wordlistScanner.Scan() {
		subdomainWord := wordlistScanner.Text()
		subdomain := proto + subdomainWord + "." + domain
		_, err := http.Get(subdomain)
		if err != nil {
			continue
		}
		newSubdomain(subdomain)
	}

	return
}
