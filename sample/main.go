package main

import (
  "context"
  "encoding/hex"
  "fmt"
  mqpro "github.com/semenovem/gopkg_mqpro/v2"
  "github.com/sirupsen/logrus"
  "net/http"
  "os"
  "os/signal"
  "strconv"
  "syscall"
  "time"
)

var log = logrus.NewEntry(logrus.New())
var rootCtx, rootCtxCancel = context.WithCancel(context.Background())
var ibmmq = mqpro.New(rootCtx, log)
var correlId []byte
var correlId2 []byte

func init() {
  go func() {
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-sig
    rootCtxCancel()
  }()

  correlId, _ = hex.DecodeString("414d5120514d3120202020202020202005b3b06029480440")
  correlId2, _ = hex.DecodeString("414d5120514d3120202020202020202005b3b06029480444")
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

  // api
  go func() {
    p, err := strconv.Atoi(os.Getenv("ENV_API_PORT"))
    if err != nil {
      fmt.Println("not set correct ENV_API_PORT: ", err)
      panic("not set correct ENV_API_PORT")
    }

    err = http.ListenAndServe(fmt.Sprintf(":%d", p), nil)
    fmt.Println("ListenAndServe: ", err)
  }()

  if cfg.SimpleSubscriber || cfg.Mirror {
    subscr()
  }

  <-rootCtx.Done()
  time.Sleep(time.Second * 1)
}
