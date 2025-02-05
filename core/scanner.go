package core

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

var securityHeaders = []string{
	"Content-Security-Policy",
	"Strict-Transport-Security",
	"X-Frame-Options",
	"X-XSS-Protection",
	"X-Content-Type-Options",
	"Referrer-Policy",
	"Content-Type",
	"Set-Cookie",
	"Access-Control-Allow-Origin",
}

var commonSubdomains = []string{
	"www", "api", "mail", "blog", "shop", "cdn", "dev", "test", "staging",
	"secure", "portal", "vpn", "admin", "support", "status", "files", "assets",
	"ucp", "acp", "admin", "user", "client", "forum",
}

func ScanTarget(target string) error {
	fmt.Println("\nAccessing target...")
	client := http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(target)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	defer resp.Body.Close()

	for _, header := range securityHeaders {
		if resp.Header.Get(header) == "" {
			fmt.Printf("⚠️ Missing: %s\n", header)
		}
	}

	return nil
}

func ScanSubdomains(domain string) {
	client := http.Client{Timeout: 5 * time.Second}

	for _, sub := range commonSubdomains {
		subdomain := fmt.Sprintf("https://%s.%s", sub, domain)
		resp, err := client.Get(subdomain)
		if err == nil && resp.StatusCode < 400 {
			fmt.Printf("\n✅  Found subdomain: %s (Status: %d)\n", subdomain, resp.StatusCode)
		}
	}
}

func ResolveDNSSubdomains(domain string) {
	for _, sub := range commonSubdomains {
		fullSubdomain := fmt.Sprintf("%s.%s", sub, domain)
		ips, err := net.LookupHost(fullSubdomain)
		if err == nil {
			fmt.Printf("\n🌍 DNS Resolved: %s -> %v\n", fullSubdomain, ips)
		}
	}
}
