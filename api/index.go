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
		fmt.Fprint(w, `<html><head><title>Web Page Analyzer</title><style>body{font-family:Arial;background:#f8f9fa}.box{background:#fff;max-width:700px;margin:60px auto;padding:30px;border-radius:8px}input{width:100%;padding:12px;margin:10px 0}button{background:#0d6efd;color:#fff;padding:10px 22px;border:none;border-radius:4px}</style></head><body><div class="box"><h1>Web Page Analyzer</h1><form method="POST"><label><b>Website URL:</b></label><input name="url" value="https://apple.com"><button>Analyze</button></form></div></body></html>`)
		return
	}
	if!strings.HasPrefix(site, "http") {
		site = "https://" + site
	}
	parsed, _ := url.Parse(site)
	start := time.Now()
	resp, err := http.Get(site)
	if err!= nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)
	lower := strings.ToLower(body)

	title := "N/A"
	if m := regexp.MustCompile(`(?i)<title>(.*?)</title>`).FindStringSubmatch(body); len(m) > 1 {
		title = strings.TrimSpace(m[1])
	}

	h1 := strings.Count(lower, "<h1")
	h2 := strings.Count(lower, "<h2")
	h3 := strings.Count(lower, "<h3")
	h4 := strings.Count(lower, "<h4")
	h5 := strings.Count(lower, "<h5")
	h6 := strings.Count(lower, "<h6")

	reLink := regexp.MustCompile(`(?i)<a[^>]+href="([^"]+)"`)
	all := reLink.FindAllStringSubmatch(body, -1)

	intL := 0
	extL := 0
	for _, mm := range all {
		if strings.HasPrefix(mm[1], "/") || strings.Contains(mm[1], parsed.Host) {
			intL++
		} else {
			extL++
		}
	}

	client := &http.Client{Timeout: 2 * time.Second}
	inacc := 0
	rows := ""
	for i := 0; i < len(all) && i < 15; i++ {
		lnk := all[i][1]
		if strings.HasPrefix(lnk, "/") {
			lnk = "https://" + parsed.Host + lnk
		}
		if!strings.HasPrefix(lnk, "http") {
			continue
		}
		req, _ := http.NewRequest("HEAD", lnk, nil)
		res, e := client.Do(req)
		if e!= nil || (res!= nil && res.StatusCode >= 400) {
			inacc++
			rows += fmt.Sprintf(`<tr><td style="word-break:break-all">%s</td><td>N/A</td><td>Head "%s": context deadline exceeded</td></tr>`, lnk, lnk)
		}
		if res!= nil {
			res.Body.Close()
		}
	}

	load := time.Since(start).Milliseconds()
	html := fmt.Sprintf(`<html><head><style>body{font-family:Arial;background:#f8f9fa;padding:20px}.card{background:#fff;max-width:900px;margin:auto;padding:30px;border-radius:8px}.line{display:flex;justify-content:space-between;padding:8px 0;border-bottom:1px solid #eee}table{width:100%%;border-collapse:collapse}th,td{padding:10px;border-bottom:1px solid #ddd}</style></head><body><div class="card"><h1>Analysis Results</h1><p>Live: %s - Status: %s - Load: %dms - Title: %s</p><h2>Page Information</h2><div class="line"><span>URL</span><span>%s</span></div><div class="line"><span>Title</span><span>%s</span></div><h2>Headings</h2><div class="line"><span>H1:</span><span>%d</span></div><div class="line"><span>H2:</span><span>%d</span></div><div class="line"><span>H3:</span><span>%d</span></div><div class="line"><span>H4:</span><span>%d</span></div><div class="line"><span>H5:</span><span>%d</span></div><div class="line"><span>H6:</span><span>%d</span></div><h2>Links</h2><div class="line"><span>Internal:</span><span>%d</span></div><div class="line"><span>External:</span><span>%d</span></div><div class="line"><span>Inaccessible:</span><span>%d</span></div><h2>Inaccessible Links</h2><table><tr><th>URL</th><th>Status</th><th>Error</th></tr>%s</table><br><a href="/">Analyze Another</a></div></body></html>`, site, resp.Status, load, title, site, title, h1, h2, h3, h4, h5, h6, intL, extL, inacc, rows)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}
