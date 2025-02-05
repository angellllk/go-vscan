package main

import (
	"flag"
	"fmt"
	"go-vscan/core"
	"strings"
	"sync"
)

func main() {
	targetPtr := flag.String("target", "", "target url")
	flag.Parse()

	if len(*targetPtr) == 0 {
		fmt.Printf("error: target must include a valid url")
		return
	}

	domain := strings.Replace(strings.Replace(*targetPtr, "https://", "", 1), "http://", "", 1)
	domain = strings.Split(domain, "/")[0]

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		core.ScanSubdomains(domain)
	}()

	go func() {
		defer wg.Done()
		core.ResolveDNSSubdomains(domain)
	}()

	wg.Wait()

	if err := core.ScanTarget(*targetPtr); err != nil {
		return
	}

	core.CrawlWebsite(*targetPtr, *targetPtr, 2)

	for link, _ := range core.VisitedLinks {
		core.TestForXSS(link)
		core.TestForSQLInjection(link)
	}
}
