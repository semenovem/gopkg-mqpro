package main

import (
  "context"
  "encoding/hex"
  "fmt"
  mqpro "github.com/semenovem/gopkg_mqpro"
  "github.com/sirupsen/logrus"
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

func init() {
  go func() {
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-sig
    rootCtxCancel()
  }()

  correlId, _ = hex.DecodeString("414d5120514d3120202020202020202005b3b06029480440")

  http.HandleFunc("/", api404)
  http.HandleFunc("/put", putMsg)
  http.HandleFunc("/get", getMsg)
  http.HandleFunc("/putget", putGetMsg)
  http.HandleFunc("/sub", onRegisterInMsg)
  http.HandleFunc("/unsub", offRegisterInMsg)
  http.HandleFunc("/browse", onBrowse)

  lev, err := logrus.ParseLevel(os.Getenv("ENV_LOG_LEVEL"))
  if err == nil {
    l := logrus.New()
    l.SetLevel(lev)

    mqpro.Log = logrus.NewEntry(l)

    //mqpro.SetLogger(logrus.NewEntry(l))
  }
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

  <-rootCtx.Done()
  time.Sleep(time.Second * 1)
}

func api404(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "404\n")
}

func logMsg(msg *mqpro.Msg) {
  fmt.Println("\n--------------------------------")
  fmt.Println("Получили сообщение:")
  fmt.Printf(">>>>> msg.Payload  = %s\n", string(msg.Payload))
  fmt.Printf(">>>>> msg.Props    = %+v\n", msg.Props)
  fmt.Printf(">>>>> msg.CorrelId = %x\n", msg.CorrelId)
  fmt.Printf(">>>>> msg.MsgId    = %x\n", msg.MsgId)
}
