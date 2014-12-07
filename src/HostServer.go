package main

import (
	"net/http"
  "text/template"
  "crypto/sha512"
  "crypto/hmac"
  "net/url"
  "RequestStore"
  "time"
)

type Substitutions struct {
  URL string
}

func index(w http.ResponseWriter, r *http.Request) {
  domain := "localhost:8000"
  repeats := "2"
  expiry := time.Now().Add(time.Hour * 24).UTC().String()
  request := "domain=" + domain + "&repeats=" + repeats + "&expiry=" + url.QueryEscape(expiry)

  key, err := RequestStore.RequestKey(domain)
  if err != nil {
    print(err)
    return
  }

  mac := hmac.New(sha512.New, key)
  urlString := "http://localhost:8001/?" + request
  mac.Write([]byte(request))
  urlString = urlString + "&auth=" + url.QueryEscape(string(mac.Sum(nil)))

  substitute := Substitutions{urlString}
  indexTemplate, _ := template.ParseFiles("index.html")
  indexTemplate.Execute(w, substitute)
}

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8000", nil)
}
