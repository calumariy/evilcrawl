package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	domainPtr := flag.String("d", "", "the ip of target")
	wordlistPtr := flag.String("w", "", "the wordlist for subdomain enumeration")
	customSubdomainsPtr := flag.String("c", "", "the wordlist of subdomains you have found")
	outFilePtr := flag.String("o", "", "output file of program")
	flag.Parse()

	domain := *domainPtr
	wordlist := *wordlistPtr
	customSubdomains := *customSubdomainsPtr
	outFileName := *outFilePtr

	if domain == "" {
		fmt.Fprintln(os.Stderr, "[x] ERROR: please specify a domain or ip address with -d")
		os.Exit(1)
	}

	outFile := os.Stdout
	if outFileName == "" {
		fmt.Fprintln(os.Stderr, "[!] No out file specified - output will be set to stdout")
	} else {
		var err error
		outFile, err = os.OpenFile(outFileName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		handleErr(err)
	}

	subdomains, err := compileSubdomains(domain, wordlist, customSubdomains)
	handleErr(err)
	amntSubdomains := len(subdomains)
	for _, subdomain := range subdomains {
		fmt.Fprintln(outFile, subdomain)
	}
	fmt.Fprintln(os.Stderr, "[.] found "+strconv.Itoa(amntSubdomains)+" subdomains! Launching workers...")
}
