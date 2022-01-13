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
  log                  = logger()
  rootCtx, rootCtxCanc = context.WithCancel(context.Background())
  logIbmmq             = log.WithField("sys", "mq")
  mq                   = mqpro.New(rootCtx, logIbmmq)
  mqOper1Put           = mq.Queue("payPut")
  mqOper1Get           = mq.Queue("payGet")
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
    rootCtxCanc()
  }()
  log.Logger.SetFormatter(formatter())
}

func main() {
  log.Info("Старт тестового приложения работы с IBM MQ")
  defer log.Info("Остановка приложения")

  go func() {
    err := mq.Connect()
    if err == nil {
      log.Info(">>>>> Подключение к IBMMQ успешно")
    } else {
      log.Error("Err: ошибка запуска приложения:", err)
      rootCtxCanc()
    }
  }()

  mqOper1Get.RegisterInMsg(hndIncomingMsg)

  // api
  if cfg.ApiPort != 0 {
    go func() {
      err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ApiPort), nil)
      if err != nil {
        log.Error("ListenAndServe: ", err)
      }
    }()
  }

  if cfg.SimpleSubscriber || cfg.Mirror {
    //subscr()
  }

  <-rootCtx.Done()
  time.Sleep(time.Millisecond * 100)
}
