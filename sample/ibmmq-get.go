package main

import (
  "context"
  "fmt"
  "net/http"
  "time"
)

// Получает сообщение из очереди
// curl host:port/get
func getMsg(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Получение сообщения из IBM MQ")
  ctx, cancel := context.WithTimeout(rootCtx, time.Second*60)
  defer cancel()

  msg, ok, err := ibmmq.Get(ctx)
  if err != nil {
    fmt.Fprintf(w, "[get] Error: %s\n", err.Error())
    return
  }

  if !ok {
    fmt.Fprintf(w, "[get]. Message queue is empty\n")
    return
  }

  logMsg(msg)

  fmt.Fprintf(w, "[get] Ok. msgId: %x\n", msg.MsgId)
}
