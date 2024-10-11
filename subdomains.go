package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func compileSubdomains(mainDomain string, wordlist string, customSubdomainsFile string) ([]string, error) {
	// Check that domain exists
	_, err := http.Get(mainDomain)
	handleErr(err)
	passiveSubdomains, err := passiveDNSRecon(mainDomain)

	var activeSubdomains []string
	if wordlist == "" {
		fmt.Fprintln(os.Stderr, "[!] no wordlist detected, skipping brute force")
	} else {
		activeSubdomains, err = activeDNSRecon(mainDomain, wordlist)
		handleErr(err)
	}

	passiveAndActiveDomains, err := joinSubdomains(passiveSubdomains, activeSubdomains)
	handleErr(err)

	customSubdomains := []string{}
	customSubdomainList, err := os.Open(customSubdomainsFile)
	handleErr(err)

	subdomains := []string{}

	if customSubdomainsFile != "" {
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
			customSubdomains = append(customSubdomains, customSubdomain)
		}
		subdomains, err = joinSubdomains(passiveAndActiveDomains, customSubdomains)
		handleErr(err)
	} else {
		fmt.Fprintln(os.Stderr, "[!] no custom subdomains supplied")
	}

	foundDomain := false
	for _, domainName := range subdomains {
		if domainName == mainDomain {
			foundDomain = true
			break
		}
	}

	if !foundDomain {
		subdomains = append(subdomains, mainDomain)
	}

	return subdomains, nil
}

func passiveDNSRecon(mainDomain string) ([]string, error) {
	fmt.Fprintln(os.Stderr, "[!] Passive dns recon not supported yet")
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

	return subdomains, nil
}

func joinSubdomains(subdomains1 []string, subdomains2 []string) ([]string, error) {
	m := make(map[string]bool)
	compiledSubdomains := []string{}

	for _, domain1 := range subdomains1 {
		m[domain1] = true
		compiledSubdomains = append(compiledSubdomains, domain1)
	}

	for _, domain2 := range subdomains2 {
		_, hasDomain := m[domain2]
		if !hasDomain {
			compiledSubdomains = append(compiledSubdomains, domain2)
			m[domain2] = true
		}
	}

	return compiledSubdomains, nil
}
