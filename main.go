package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Website Analyzer by Fari Mazhar")
	fmt.Println("-------------------------------")

	// Check karnay wali websites
	urls := []string{
		"https://google.com",
		"https://github.com",
	}

	for _, url := range urls {
		start := time.Now()
		
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(url, "-> DOWN hai")
			continue
		}
		
		speed := time.Since(start)
		fmt.Printf("%s -> UP hai | Time: %v\n", url, speed)
		resp.Body.Close()
	}
}
