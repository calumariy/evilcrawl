package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

var subdomains = make(map[string]struct{})
var subMu sync.Mutex

var directories = make(map[string]struct{})
var dirMu sync.Mutex

var outFile *os.File

var isAuthorised bool
var auth string

func main() {

	domainPtr := flag.String("d", "", "the ip of target")
	wordlistPtr := flag.String("w", "", "the wordlist for subdomain enumeration")
	customSubdomainsPtr := flag.String("c", "", "the wordlist of subdomains you have found")
	outFilePtr := flag.String("o", "", "output file of program")
	authPtr := flag.String("a", "", "cooke for auth in the form name:cookie")
	flag.Parse()

	domain := *domainPtr
	wordlist := *wordlistPtr
	customSubdomains := *customSubdomainsPtr
	outFileName := *outFilePtr
	authStr := *authPtr

	var wg sync.WaitGroup
	if domain == "" {
		fmt.Fprintln(os.Stderr, "[x] ERROR: please specify a domain or ip address with -d")
		os.Exit(1)
	}

	jar, err := cookiejar.New(nil)
	handleErr(err)

	var cookies []*http.Cookie
	auth = authStr
	if authStr != "" {
		isAuthorised = true
		domainPath := strings.Split(domain, "/")
		path := ""
		for i := 3; i < len(domainPath)-1; i++ {
			path += "/" + domainPath[i]
		}

		mainDomain := domainPath[2]
		mainDomainArr := strings.Split(mainDomain, ".")
		cookieDomain := "." + mainDomainArr[len(mainDomainArr)-2] + "." + mainDomainArr[len(mainDomainArr)-1]
		authCookie := &http.Cookie{
			Name:   strings.Split(authStr, ":")[0],
			Value:  strings.Split(authStr, ":")[1],
			Path:   path,
			Domain: cookieDomain,
		}
		cookies = append(cookies, authCookie)
	} else {
		isAuthorised = false
	}

	url, _ := url.Parse(domain)
	jar.SetCookies(url, cookies)

	client := &http.Client{
		Jar: jar,
	}

	outFile = os.Stdout
	if outFileName == "" {
		fmt.Fprintln(os.Stderr, "[!] No out file specified - output will be set to stdout")
	} else {
		var err error
		outFile, err = os.OpenFile(outFileName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		handleErr(err)
	}

	compileSubdomains(domain, wordlist, customSubdomains)
	fmt.Fprintln(os.Stderr, "[.] found "+strconv.Itoa(len(subdomains))+" subdomains! Launching workers...")

	for subdomain := range subdomains {
		recursivelyAttackDirectory(subdomain, domain, subdomain, client, &wg)
	}

	wg.Wait()
}
