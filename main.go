package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	domainPtr := flag.String("d", "", "the ip of target")
	flag.Parse()
	domain := *domainPtr
	if domain == "" {
		fmt.Fprintln(os.Stderr, "error: please specify a domain with -d")
		os.Exit(1)
	}
	subdomains, err := compileSubdomains(domain)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(subdomains)
}