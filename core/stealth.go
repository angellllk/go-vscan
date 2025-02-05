package core

import (
	"math/rand"
	"time"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 14_2 like Mac OS X) AppleWebKit/537.36 (KHTML, like Gecko) Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 10; SM-G975F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
}

func GetRandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}
func RandomDelay() {
	delay := time.Duration(rand.Intn(2500)+500) * time.Millisecond
	time.Sleep(delay)
}
