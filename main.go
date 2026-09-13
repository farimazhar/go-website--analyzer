package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Go Website Analyzer v1.0.0")
	
	websites := []string{
		"https://google.com",
		"https://github.com",
	}

	for _, url := range websites {
		start := time.Now()
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("❌ %s is DOWN - %v\n", url, err)
			continue
		}
		duration := time.Since(start)
		fmt.Printf("✅ %s is UP - Status: %d - Time: %v\n", url, resp.StatusCode, duration)
		resp.Body.Close()
	}
}
