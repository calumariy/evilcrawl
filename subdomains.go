package main

import (
	"fmt"
	"net/http"
	"os"
)

func compileSubdomains(mainDomain string) ([]string, error) {
	// Check that domain exists
	_, err := http.Get(mainDomain)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	subdomains := []string{mainDomain}

	return subdomains, nil
}

func passiveDNSRecon(mainDomain string) ([]string, error) {
	return []string{mainDomain}, nil
}

func activeDNSRecon(mainDomain string) ([]string, error) {
	return []string{mainDomain}, nil
}

func joinSubdomains(subdomains1 []string, subdomains2 []string) ([]string, error) {
	return []string{}, nil
}
