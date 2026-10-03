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
	site := r.FormValue("url")
	if site == "" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<html><head><title>Web Page Analyzer</title><style>body{font-family:Arial;background:#f8f9fa}.box{background:#fff;max-width:700px;margin:60px auto;padding:30px;border-radius:8px;box-shadow:0 2px 8px #ddd}input{width:100%;padding:12px;border:1px solid #ccc;border-radius:4px;margin:10px 0}button{background:#0d6efd;color:#fff;padding:10px 22px;border:none;border-radius:4px;cursor:pointer}</style></head><body><div class="box"><h1>Web Page Analyzer</h1><p>Enter a URL to analyze its HTML structure and links.</p><form method="POST"><label><b>Website URL:</b></label><input type="text" name="url" value="https://apple.com" placeholder="https://example.com"><button type="submit">Analyze</button></form></div></body></html>`)
		return
	}
	if!strings.HasPrefix(site, "http") { site = "https://" + site }
	parsed, _ := url.Parse(site)
	start := time.Now()
	resp, err := http.Get(site)
	if err!= nil { fmt.Fprintf(w, "Error: %v", err); return }
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	body := string(b)
	lower := strings.ToLower(body)
	title := "N/A"
	if m := regexp.MustCompile(`(?i)<title>(.*?)</title>`).FindStringSubmatch(body); len(m) > 1 { title = strings.TrimSpace(m[1]) }
	h1 := strings.Count(lower, "<h1")
	h2 := strings.Count(lower, "<h2")
	h3 := strings.Count(lower, "<h3")
	h4 := strings.Count(lower, "<h4")
	h5 := strings.Count(lower, "<h5")
	h6 := strings.Count(lower, "<h6")
	re := regexp.MustCompile(`(?i)<a[^>]+href="([^"]+)"`)
	all := re.FindAllStringSubmatch(body, -1)
	intL, extL := 0, 0
	for _, mm := range all {
		if strings.HasPrefix(mm[1], "/") || strings.Contains(mm[1], parsed.Host) { intL++ } else { extL++ }
	}
	client := &http.Client{Timeout: 2 * time.Second}
	inacc := 0
	rows := ""
	for i := 0; i < len(all) && i < 15; i++ {
		lnk := all[i][1]
		if strings.HasPrefix(lnk, "/") { lnk = "https://" + parsed.Host + lnk }
		if!strings.HasPrefix(lnk, "http") { continue }
		req, _ := http.NewRequest("HEAD",
