
package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"regexp"
)

func Handler(w http.ResponseWriter, r *http.Request) {
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
			body, _ := io.ReadAll(resp.Body)
			loadTime := time.Since(start)
			fullReport += fmt.Sprintf("✅ Live: %s\nStatus: %s\nLoad Time: %s\n", siteURL, resp.Status, loadTime)
			
			re := regexp.MustCompile(`<title>(.*?)</title>`)
			m := re.FindStringSubmatch(string(body))
			if len(m) > 1 {
				fullReport += fmt.Sprintf("Title: %s\n", m[1])
			}
		} else {
			fullReport = fmt.Sprintf("❌ Error: %v", err)
		}
	}

	html := fmt.Sprintf(`<html><body style="font-family:sans-serif;padding:20px;"><h1>Go Website Analyzer</h1><form><input name="url" placeholder="google.com" style="padding:8px;width:300px;"><button>Analyze</button></form><pre style="background:#f4f4f4;padding:15px;margin-top:20px;">%s</pre></body></html>`, fullReport)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}
