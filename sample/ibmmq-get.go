package main

import (
  "context"
  "fmt"
  "net/http"
  "time"
)

// Получает сообщение из очереди
// Пример:
// curl localhost:8080/get-by-correl/414d5120514d3120202020202020202005b3b060019f0240
func getMsg(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Получение сообщения из IBM MQ")

  ctx, cancel := context.WithTimeout(rootCtx, time.Second*10)
  defer cancel()

  msgId, ok, err := ibmmq.Get(ctx, nil)
  if err != nil {
    fmt.Fprintf(w, "get by correlId Error: %s", err.Error())
    return
  }

  if !ok {
    fmt.Fprintf(w, "get by correlId. Message not found")
    return
  }

  fmt.Fprintf(w, "get by correlId Ok. msgId: %x", msgId)
}
