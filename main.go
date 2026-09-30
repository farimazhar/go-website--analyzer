package main

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)
//go:embed index.html
var content embed.FS

func handler(w http.ResponseWriter, r *http.Request){
	if r.Method=="POST" || r.URL.Path=="/analyze"{
		r.ParseForm()
		url:=r.FormValue("url")
		if!strings.HasPrefix(url,"http"){url="https://"+url}
		resp,err:=http.Get(url)
		if err!=nil{fmt.Fprintf(w,"Error: %v",err); return}
		defer resp.Body.Close()
		b,_:=io.ReadAll(resp.Body); body:=string(b)
		title:="Not Found"
		re:=regexp.MustCompile(`(?i)<title>(.*?)</title>`)
		if m:=re.FindStringSubmatch(body); len(m)>1{title=m[1]}
		hasLogin:="No"
		if strings.Contains(strings.ToLower(body),`type="password"`){hasLogin="Yes"}
		h1:=strings.Count(strings.ToLower(body),"<h1")
		fmt.Fprintf(w,`<html><body style="font-family:Arial;padding:20px"><h1>Analysis Results</h1><hr><p><b>URL:</b> %s</p><p><b>Title:</b> %s</p><p><b>Login Form:</b> %s</p><p><b>H1 Count:</b> %d</p><a href="/">Back</a></body></html>`,url,title,hasLogin,h1)
		return
	}
	data,_:=content.ReadFile("index.html")
	w.Write(data)
}
func main(){
	http.HandleFunc("/",handler)
	http.HandleFunc("/analyze",handler)
	http.ListenAndServe(":8080",nil)
}
