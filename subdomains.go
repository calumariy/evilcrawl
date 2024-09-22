package main

import (
	"fmt"
	"io"
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
