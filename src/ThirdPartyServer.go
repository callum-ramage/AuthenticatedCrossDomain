package main

import (
  "net/http"
  "text/template"
  "crypto/sha512"
  "crypto/hmac"
  "net/url"
  "RequestStore"
  "strconv"
)

type Substitutions struct {
  URL string
}

func index(w http.ResponseWriter, r *http.Request) {
  queryParams := r.URL.Query()

  if _, domainOk := queryParams["domain"]; !domainOk || len(queryParams["domain"]) != 1 {
    return
  }
  domain := queryParams["domain"][0]
  request := "domain=" + domain

  maxRepeats := 1;
  if repeats, repeatsOk := queryParams["repeats"]; repeatsOk && len(repeats) == 1 {
    request = request + "&repeats=" + repeats[0]
    maxRepeats, _ = strconv.Atoi(repeats[0])
  }

  if expiry, expiryOk := queryParams["expiry"]; expiryOk && len(expiry) == 1 {
    request = request + "&expiry=" + url.QueryEscape(expiry[0])
  }

  if _, authOk := queryParams["auth"]; !authOk || len(queryParams["auth"]) != 1 {
    return
  }
  auth := queryParams["auth"][0]

  key, err := RequestStore.RequestKey(domain)
  if err != nil {
    return
  }

  mac := hmac.New(sha512.New, key)
  mac.Write([]byte(request))

  if hmac.Equal(mac.Sum(nil), []byte(auth)) {
    if RequestStore.RequestUsed(auth) >= maxRepeats {
      indexTemplate, _ := template.ParseFiles("consumed.html")
      indexTemplate.Execute(w, nil)
    } else {
      RequestStore.RequestAdd(auth)
      substitute := Substitutions{domain}
      indexTemplate, _ := template.ParseFiles("response.html")
      indexTemplate.Execute(w, substitute)
    }
  }
}

func main() {
  http.HandleFunc("/", index)
  http.ListenAndServe(":8001", nil)
}
