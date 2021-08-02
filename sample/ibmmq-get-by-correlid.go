package main

import (
  "context"
  "fmt"
  "net/http"
  "time"
)

// Получает сообщение из очереди
// curl host:port/get
func getMsgByCorrelId(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Получение сообщения из IBM MQ по CorrelId (ожидание в течение 5 сек)")
  ctx, cancel := context.WithTimeout(rootCtx, time.Second*5)
  defer cancel()

  go func() {
    ctx, cancel := context.WithTimeout(rootCtx, time.Second*5)
    defer cancel()
    msg, ok, err := ibmmq.GetByCorrelId(ctx, correlId2)
    if err != nil {
      fmt.Printf("[get] Error: %s\n", err.Error())
      return
    }
    if !ok {
      fmt.Printf("[get] Message queue is empty\n")
      return
    }

    logMsg(msg)
  }()

  msg, ok, err := ibmmq.GetByCorrelId(ctx, correlId)

  if err != nil {
    fmt.Fprintf(w, "[get] Error: %s\n", err.Error())
    return
  }

  if !ok {
    fmt.Fprintf(w, "[get] Message queue is empty\n")
    return
  }

  logMsg(msg)

  fmt.Fprintf(w, "[get] Ok. msgId: %x\n", msg.MsgId)
}
