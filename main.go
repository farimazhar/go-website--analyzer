package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func analyzerHandler(w http.ResponseWriter, r *http.Request) {
	siteURL := r.FormValue("url")
	if siteURL == "" {
		siteURL = r.URL.Query().Get("url")
	}

	w.Header().Set("Content-Type", "text/html")
	
	// Home Page
	if siteURL == "" {
		fmt.Fprint(w, `
		<h2 style="font-family:sans-serif;">Web Page Analyzer</h2>
		<p>Enter a URL to analyze its HTML structure and links.</p>
		<form>
			<label>Website URL:</label><br>
			<input type="text" name="url" placeholder="https://apple.com" style="width:400px;padding:8px;" required>
			<br><br>
			<button style="background:#2b6cb0;color:white;padding:8px 20px;border:none;border-radius:4px;">Analyze</button>
		</form>
		`)
		return
	}

	// Fix URL if missing https://
	if !strings.HasPrefix(siteURL, "http") {
		siteURL = "https://" + siteURL
	}

	start := time.Now()
	resp, err := http.Get(siteURL)
	duration := time.Since(start)

	if err != nil {
		fmt.Fprintf(w, "<h3>Error: %v</h3><a href='/'>Go Back</a>", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	u, _ := url.Parse(siteURL)
	linksCount := strings.Count(string(body), "<a ")

	fmt.Fprintf(w, `
		<h2>✅ Result</h2>
		<p><b>URL:</b> %s</p>
		<p><b>Status:</b> Live (%d %s)</p>
		<p><b>Time:</b> %d ms</p>
		<p><b>Links Found:</b> %d</p>
		<p><b>Host:</b> %s</p>
		<br>
		<p>By: Fari Mazhar - Alhamdulillah</p>
		<p>Bismillah Hir Rahman Nir Raheem</p>
		<a href="/">Analyze Another</a>
	`, siteURL, resp.StatusCode, resp.Status, duration.Milliseconds(), linksCount, u.Host)
}

func main() {
	http.HandleFunc("/", analyzerHandler)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
