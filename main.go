package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		siteURL := r.FormValue("url")
		var fullReport string

		if siteURL!= "" {
			if!strings.HasPrefix(siteURL, "http") {
				siteURL = "https://" + siteURL
			}
			start := time.Now()
			resp, err := http.Get(siteURL)
			if err == nil {
				defer resp.Body.Close()
				bodyBytes, _ := io.ReadAll(resp.Body)
				body := string(bodyBytes)
				elapsed := time.Since(start).Milliseconds()
				parsedURL, _ := url.Parse(siteURL)

				// Title
				titleRe := regexp.MustCompile(`(?i)<title>(.*?)</title>`)
				titleMatch := titleRe.FindStringSubmatch(body)
				title := "Not Found"
				if len(titleMatch) > 1 { title = titleMatch[1] }

				// Meta Description
				descRe := regexp.MustCompile(`(?i)<meta[^>]*name="description"[^>]*content="([^"]*)"`)
				descMatch := descRe.FindStringSubmatch(body)
				description := "Not Found"
				if len(descMatch) > 1 { description = descMatch[1] }

				h1Count := strings.Count(strings.ToLower(body), "<h1")
				h2Count := strings.Count(strings.ToLower(body), "<h2")
				h3Count := strings.Count(strings.ToLower(body), "<h3")
				imgCount := strings.Count(strings.ToLower(body), "<img")
				linkCount := strings.Count(strings.ToLower(body), "<a ")
				wordCount := len(strings.Fields(body))

				// Internal / External Links
				internalLinks := strings.Count(body, parsedURL.Host)

				fullReport = fmt.Sprintf(`
				<div style="text-align:left; background:white; border:1px solid #ddd; border-radius:12px; padding:20px; margin-top:20px;">
					<h3 style="color:green;">✅ Complete Report for %s</h3>
					<hr>
					<p><b>🌐 URL:</b> %s</p>
					<p><b>📡 Status:</b> %s (%d)</p>
					<p><b>⏱️ Load Time:</b> %d ms</p>
					<p><b>📝 Title:</b> %s</p>
					<p><b>📄 Meta Description:</b> %s</p>
					<hr>
					<p><b>🔍 SEO Structure:</b></p>
					<p>H1 Tags: %d | H2 Tags: %d | H3 Tags: %d</p>
					<p><b>🔗 Links:</b> Total %d (Internal ~ %d)</p>
					<p><b>🖼️ Images:</b> %d</p>
					<p><b>📖 Word Count:</b> %d words</p>
					<p><b>📦 Page Size:</b> %d KB</p>
					<hr>
					<p style="color:#2b6cb0;"><b>By: Fari Mazhar - Alhamdulillah</b></p>
				</div>
				`, parsedURL.Host, siteURL, resp.Status, resp.StatusCode, elapsed, title, description, h1Count, h2Count, h3Count, linkCount, internalLinks, imgCount, wordCount, len(bodyBytes)/1024)
			} else {
				fullReport = fmt.Sprintf("<p style='color:red;'>Error: %v</p>", err)
			}
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
		<html><head><meta name="viewport" content="width=device-width, initial-scale=1"><title>Complete Website Analyzer</title></head>
		<body style="font-family:Arial; max-width:700px; margin:30px auto; background:#f8fafc; padding:20px;">
		<h2 style="text-align:center;">Go Website Analyzer - Full Audit</h2>
		<p style="text-align:center;">Bismillah Hir Rahman Nir Raheem</p>
		<form method="POST" style="background:white; padding:20px; border-radius:12px; box-shadow:0 2px 10px #ddd;">
			<label><b>Enter Website URL:</b></label><br><br>
			<input type="text" name="url" placeholder="https://adidas.com or https://apple.com" style="width:100%%; padding:12px; border:2px solid #2b6cb0; border-radius:8px;" required>
			<br><br>
			<button style="width:100%%; background:#2b6cb0; color:white; padding:12px; border:none; border-radius:8px; font-weight:bold; font-size:16px;">Get Full Details</button>
		</form>
		%s
		</body></html>
		`, fullReport)
	})
	http.ListenAndServe(":8080", nil)
}
