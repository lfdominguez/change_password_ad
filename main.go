package main

import (
	"fmt"
  "net/http"
  "strings"
  "log"
  "bytes"
//   "os"
//   "os/exec"

	"github.com/gobuffalo/packr"
)

var box = packr.NewBox("./template")


func showIndex (w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  fmt.Println(r.Form)
  fmt.Println("path", r.URL.Path)
  fmt.Println("scheme", r.URL.Scheme)
  fmt.Println(r.Form["url_long"])
  for k, v := range r.Form {
      fmt.Println("key:", k)
      fmt.Println("val:", strings.Join(v, ""))
  }
  
  if _, ok := r.Form["user"]; ok {
    user := r.Form["user"]
    old_pass := r.Form["old_pass"]
    new_pass := r.Form["new_pass"]
    retype_pass := r.Form["retype_pass"]
    
    fmt.Fprintf(w, "%s", user)
    fmt.Fprintf(w, "%s", old_pass)
    fmt.Fprintf(w, "%s", new_pass)
    fmt.Fprintf(w, "%s", retype_pass)
  }

  
  w.Header().Set("Content-type", "text/html")
  
  str, _ := box.MustString("index.html")
  
  fmt.Fprintf(w, "%s", str)
}

func sendLogo (w http.ResponseWriter, r *http.Request) {
  bytesRead, _ := box.MustBytes("logo.png")
  b := bytes.NewBuffer(bytesRead)
  
  w.Header().Set("Content-type", "image/png")

  if _, err := b.WriteTo(w); err != nil {
    fmt.Fprintf(w, "%s", err)
  }
}

func sendStyle (w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-type", "text/css")
  
  str, _ := box.MustString("style.css")
  
  fmt.Fprintf(w, "%s", str)
}


func main() {
  http.HandleFunc("/", showIndex)
  http.HandleFunc("/logo.png", sendLogo)
  http.HandleFunc("/style.css", sendStyle)
  
  err := http.ListenAndServe(":9090", nil) // set listen port
  if err != nil {
      log.Fatal("ListenAndServe: ", err)
  }
} 
