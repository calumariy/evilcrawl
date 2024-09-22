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
		fmt.Println("error: please enter a domain with -d")
		os.Exit(1)
	}
	fmt.Println("domain:", domain)
}
