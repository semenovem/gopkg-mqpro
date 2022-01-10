package main

import (
  "context"
  "fmt"
  mqpro "github.com/semenovem/gopkg_mqpro/v2"
  "github.com/sirupsen/logrus"
  "net/http"
  "os"
  "os/signal"
  "syscall"
  "time"
)

var (
  log                    = logger()
  rootCtx, rootCtxCancel = context.WithCancel(context.Background())
  logIbmmq               = log.WithField("sys", "mq")
  ibmmq                  = mqpro.New(rootCtx, logIbmmq)
  ibmmqOper1In           = ibmmq.Queue("payIn")
  //ibmmqOper1Out          = ibmmq.Queue("payOut")
)

func logger() *logrus.Entry {
  l := logrus.NewEntry(logrus.New())
  l.Logger.SetFormatter(formatter())
  return l
}

func init() {
  go func() {
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-sig
    rootCtxCancel()
  }()
  log.Logger.SetFormatter(formatter())
}

func main() {
  log.Info("Старт тестового приложения работы с IBM MQ")
  defer log.Info("Остановка приложения")

  go func() {
    err := ibmmq.Connect2()
    if err != nil {
      log.Error("Err: ошибка запуска приложения:", err)
      rootCtxCancel()
    }
  }()

  // api
  if cfg.ApiPort != 0 {
    go func() {
      err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ApiPort), nil)
      if err != nil {
        log.Error("ERR: ListenAndServe: ", err)
      }
    }()
  }

  if cfg.SimpleSubscriber || cfg.Mirror {
    subscr()
  }

  <-rootCtx.Done()
  time.Sleep(time.Millisecond * 200)
}
