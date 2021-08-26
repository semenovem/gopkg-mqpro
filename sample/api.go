package main

import (
  "fmt"
  "net/http"
)

func init() {
  http.HandleFunc("/", api404)
  http.HandleFunc("/put", putMsg)
  http.HandleFunc("/get", getMsg)
  http.HandleFunc("/putget", putGetMsg)
  http.HandleFunc("/sub", onRegisterInMsg)
  http.HandleFunc("/unsub", offRegisterInMsg)
  http.HandleFunc("/browse", onBrowse)
  http.HandleFunc("/correl", getMsgByCorrelId)
}

func api404(w http.ResponseWriter, r *http.Request) {
  _, _ = fmt.Fprintf(w, "404\n")
}
