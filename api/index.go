package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	siteURL := r.FormValue("url")

	// Input Page - koi bhi link daal sakte ho
	if siteURL == "" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<html><head><meta name="viewport" content="width=device-width,initial-scale=1"><style>body{background:#f5f5f5;font-family:Arial}.card{background:white;max-width:800px;margin:60px auto;padding:35px;border-radius:10px;box-shadow:0 2px 10px #0001}input{width:100%;padding:12px;border:1px solid #ccc;border-radius:5px;margin-top:8px}button{background:#0d6efd;color:#fff;padding:12px 22px;border:0;border-radius:5px;margin-top:12px;cursor:pointer}</style></head><body><div class="card"><h1>Web Page Analyzer</h1><p>Enter a URL to analyze its HTML structure and links.</p><form><label><b>Website URL:</b></label><input name="url" placeholder="https://google.com or https://youtube.com" required><button>Analyze</button></form></div></body></html>`)
		return
	}

	if!strings.HasPrefix(siteURL, "http") {
		siteURL = "https://" + siteURL
	}
	parsed, _ := url.Parse(siteURL)

	resp, err := http.Get(siteURL)
	if err!= nil {
		fmt.Fprintf(w, "Error fetching %s : %v <br><a href='/'>Back</a>", siteURL, err)
		return
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	body := string(b)
	lower := strings.ToLower(body)

	title := "N/A"
	if m := regexp.MustCompile(`(?i)<title>(.*?)</title>`).FindStringSubmatch(body); len(m) > 1 {
		title = m[1]
	}
	htmlVer := "HTML5"
	login := "No"
	if strings.Contains(lower, `type="password"`) {
		login = "Yes"
	}

	h1 := strings.Count(lower, "<h1")
	h2 := strings.Count(lower, "<h2")
	h3 := strings.Count(lower, "<h3")
	h4 := strings.Count(lower, "<h4")
	h5 := strings.Count(lower, "<h5")
	h6 := strings.Count(lower, "<h6")

	re := regexp.MustCompile(`(?i)<a[^>]+href="([^"]+)"`)
	matches := re.FindAllStringSubmatch(body, -1)
	intC, extC := 0, 0
	var links []string
	for _, mm := range matches {
		l := mm[1]
		if l == "" || strings.HasPrefix(l, "#") || strings.HasPrefix(l, "mailto:") {
			continue
		}
		links = append(links, l)
		if strings.HasPrefix(l, "/") || strings.Contains(l, parsed.Host) {
			intC++
		} else {
			extC++
		}
	}

	client := &http.Client{Timeout: 2 * time.Second}
	var rows string
	inacc := 0
	for i := 0; i < len(links) && i < 10; i++ {
		t := links[i]
		if strings.HasPrefix(t, "/") {
			t = "https://"
