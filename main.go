package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	domainPtr := flag.String("d", "", "the ip of target")
	wordlistPtr := flag.String("w", "", "the wordlist for subdomain enumeration")
	flag.Parse()
	domain := *domainPtr
	wordlist := *wordlistPtr
	if domain == "" {
		fmt.Fprintln(os.Stderr, "error: please enter a domain with -d")
		os.Exit(1)
	}

	subdomains, err := compileSubdomains(domain, wordlist)
	handleErr(err)
	amntSubdomains := len(subdomains)
	for i := 0; i < amntSubdomains; i++ {
		fmt.Println(subdomains[i])
	}
}
