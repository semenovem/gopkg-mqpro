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

  ctx, cancel := context.WithTimeout(rootCtx, time.Second*10)
  defer cancel()

  msg, ok, err := ibmmqOper1In.Get(ctx)
  if err != nil {
    fmt.Println("[ERROR] при получении сообщения: ", err)
    _, _ = fmt.Fprintf(w, "[get] Error: %s\n", err.Error())
    return
  }

  if !ok {
    fmt.Println("[WARN] нет сообщений")
    _, _ = fmt.Fprintf(w, "[get]. Message queue is empty\n")
    return
  }

  _, _ = fmt.Fprintf(w, "[get] Ok. msgId: %x\n", msg.MsgId)
}
