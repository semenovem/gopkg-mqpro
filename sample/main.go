package main

import (
  "context"
  "fmt"
  mqpro "github.com/semenovem/gopkg_mqpro/v2"
  "github.com/semenovem/gopkg_mqpro/v2/queue"
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

var mqPayOut = queue.New(rootCtx, log.WithField("a", "PayOut"))

func init() {
  go func() {
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-sig
    rootCtxCancel()
  }()
}

func main() {
  fmt.Println("Старт тестового приложения работы с IBM MQ")
  defer fmt.Println("Остановка приложения")

  go func() {
    ibmmq.UseDefEnv()

    //err := ibmmq.Connect()
    err := ibmmq.Connect2()
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
  time.Sleep(time.Millisecond * 300)
}
