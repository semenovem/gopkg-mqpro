package main

import (
  "context"
  "fmt"
  "github.com/semenovem/mqm/v2"
  "github.com/sirupsen/logrus"
  "net/http"
  "os"
  "os/signal"
  "syscall"
  "time"
)

var (
  log                 = logger()
  rootCtx, rootCtxEsc = context.WithCancel(context.Background())
  logIbmmq            = log.WithField("sys", "mq")
  mq                  = mqm.New(rootCtx, logIbmmq)
  mqQueFooPut         = mq.NewQueue("aliasQueueFooPut")
  mqQueFooGet         = mq.NewQueue("aliasQueueFooGet")
  mqBar               = mq.NewPipe("aliasQueueBar")
  mqQuePut            mqm.Queue
  mqQueGet            mqm.Queue
)

func logger() *logrus.Entry {
  l := logrus.NewEntry(logrus.New())
  l.Logger.SetFormatter(formatter())
  return l
}

func init() {
  mqQuePut = mqQueFooPut
  mqQueGet = mqQueFooGet

  mqQuePut = mqBar
  mqQueGet = mqBar

  go func() {
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-sig
    rootCtxEsc()
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
      rootCtxEsc()
    }
  }()

  //mqQueFooGet.RegisterInMsg(hndIncomingMsg)

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
