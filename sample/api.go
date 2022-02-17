package main

import (
  "fmt"
  "net/http"
)

func init() {
  http.HandleFunc("/", api404)
  http.HandleFunc("/ready", apiReady)

  http.HandleFunc("/get", getMsg)
  http.HandleFunc("/correl/", getMsgByCorrelId)
  http.HandleFunc("/msgid/", getMsgByMsgId)

  http.HandleFunc("/put", putMsg)
  http.HandleFunc("/browse", onBrowse)

  http.HandleFunc("/sub", onRegisterInMsg)
  http.HandleFunc("/unsub", offRegisterInMsg)

  http.HandleFunc("/ping", mqPing)

  http.HandleFunc("/on-mirror", onMirror)
  http.HandleFunc("/off-mirror", offMirror)

  http.HandleFunc("/on-dev-mode", onDevMode)
  http.HandleFunc("/off-dev-mode", offDevMode)

  http.HandleFunc("/config", apiPrintCfg)
  http.HandleFunc("/clear", clearQueue)

  http.HandleFunc("/conn", apiConn)
  http.HandleFunc("/disconn", apiDisconn)

  http.HandleFunc("/open", apiOpen)
  http.HandleFunc("/close", apiClose)
}

func api404(w http.ResponseWriter, _ *http.Request) {
  _, _ = fmt.Fprintf(w, "404\nuse: [/ping, /on-mirror, /off-mirror,"+
    " /on-dev-mode, /off-dev-mode,"+
    "/get, /correl, /msgid, /config, /clear, /sub, /unsub, /browse, /put]\n")
}

func apiReady(w http.ResponseWriter, _ *http.Request) {
  _, _ = fmt.Fprintf(w, "[ready] Ok. pipe:%s=%t. "+
    "queueGet:%s=%t. "+
    "queuePut:%s=%t.\n",
    mqBar.Alias(), mqBar.Ready(),
    mqQueFooGet.Alias(), mqQueFooGet.Ready(),
    mqQueFooPut.Alias(), mqQueFooPut.Ready())
}

func apiConn(w http.ResponseWriter, _ *http.Request) {
  fmt.Fprint(w, "start ibmmq connect:\n")
  err := mq.Connect()
  if err != nil {
    fmt.Fprintf(w, "ERROR: %s\n", err.Error())
  }
  fmt.Fprintf(w, "end\n")
}

func apiDisconn(w http.ResponseWriter, _ *http.Request) {
  fmt.Fprint(w, "start ibmmq disconnect:\n")
  err := mq.Disconnect()
  if err != nil {
    fmt.Fprintf(w, "ERROR: %s\n", err.Error())
  }
  fmt.Fprintf(w, "end\n")
}

func apiOpen(w http.ResponseWriter, _ *http.Request) {
  fmt.Fprint(w, "opening ibmmq queue:\n")
  err := mqQueGet.Open()
  if err != nil {
    fmt.Fprintf(w, "ERROR: %s\n", err.Error())
  }
  fmt.Fprintf(w, "end\n")
}

func apiClose(w http.ResponseWriter, _ *http.Request) {
  fmt.Fprint(w, "closing ibmmq queue:\n")
  err := mqQueGet.Close()
  if err != nil {
    fmt.Fprintf(w, "ERROR: %s\n", err.Error())
  }
  fmt.Fprintf(w, "end\n")
}

// Включает режим DevMode для библиотеки mqm
// curl host:port/on-dev-mode
func onDevMode(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Включает режим DevMode для библиотеки mqm")
  mq.SetDevMode(true)
  _, _ = fmt.Fprint(w, "[on-dev-mode] Ok\n")
}

// Выключает режим DevMode для библиотеки mqm
// curl host:port/off-dev-mode
func offDevMode(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Выключает режим DevMode для библиотеки mqm")
  mq.SetDevMode(false)
  _, _ = fmt.Fprint(w, "[off-dev-mode] Ok\n")
}

// Вывести текущие настройки
// curl host:port/config
func apiPrintCfg(w http.ResponseWriter, _ *http.Request) {
  printCfg()
}
