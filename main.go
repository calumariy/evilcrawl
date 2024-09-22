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
		fmt.Fprintln(os.Stderr, "error: please enter a domain with -d")
		os.Exit(1)
	}

	subdomains, err := compileSubdomains(domain)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	amntSubdomains := len(subdomains)
	for i := 0; i < amntSubdomains; i++ {
		fmt.Println(subdomains[i])
	}
}
