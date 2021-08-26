package main

import (
  "context"
  "encoding/hex"
  "fmt"
  mqpro "github.com/semenovem/gopkg_mqpro"
  "net/http"
  "os"
  "os/signal"
  "strconv"
  "syscall"
  "time"
)

var rootCtx, rootCtxCancel = context.WithCancel(context.Background())
var ibmmq = mqpro.New(rootCtx)
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
      panic("")
    }

    fmt.Println()
    err = http.ListenAndServe(fmt.Sprintf(":%d", p), nil)
    fmt.Println("ListenAndServe: ", err)
  }()

  if cfg.SimpleSubscriber {
    subscr()
  }


  <-rootCtx.Done()
  time.Sleep(time.Second * 1)
}


func logMsg(msg *mqpro.Msg) {
  fmt.Println("\n--------------------------------")
  fmt.Println("Получили сообщение:")
  if len(msg.Payload) < 300 {
    fmt.Printf(">>>>> msg.Payload  = %s\n", string(msg.Payload))
  } else {
    fmt.Printf(">>>>> len msg.Payload  = %d\n", len(msg.Payload))
  }
  fmt.Printf(">>>>> msg.Props    = %+v\n", msg.Props)
  fmt.Printf(">>>>> msg.CorrelId = %x\n", msg.CorrelId)
  fmt.Printf(">>>>> msg.MsgId    = %x\n", msg.MsgId)
  fmt.Printf(">>>>> msg.Time     = %s\n", msg.Time.Format(time.RFC822))
}
