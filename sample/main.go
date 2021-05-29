package main

import (
  "context"
  "fmt"
  mqpro "github.com/semenovem/gopkg_mqpro"
  "net/http"
  "os"
  "os/signal"
  "syscall"
  "time"
)

var rootCtx, rootCtxCancel = context.WithCancel(context.Background())
var ibmmq = mqpro.New(rootCtx)
var correlId = []byte{65, 77, 81, 32, 81, 77, 49, 32, 32, 32, 32, 32, 32, 32, 32, 32, 5,
  179, 176, 96, 1, 183, 2, 6, 4}

func init() {
  go func() {
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-sig
    rootCtxCancel()
  }()

  http.HandleFunc("/", api404)
  http.HandleFunc("/put", putMsg)
}

func main() {
  fmt.Println("Старт тестового приложения работы с IBM MQ")
  defer fmt.Println("Остановка приложения")

  go func() {
    ibmmq.UseDefEnv()

    err := ibmmq.Connect()
    if err != nil {
      fmt.Println("Err: ошибка запуска приложения:", err)
      rootCtxCancel()
    }
  }()

  go func() {
    err := http.ListenAndServe(":8080", nil)
    fmt.Println("ListenAndServe: ", err)
  }()

  <-rootCtx.Done()
  time.Sleep(time.Second * 2)
}

func api404(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "404\n")
}
