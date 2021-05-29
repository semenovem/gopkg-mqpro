package main

import (
  "context"
  "encoding/hex"
  "fmt"
  "net/http"
  "time"
)

// Получает сообщение из очереди
// Пример:
// curl localhost:8080/get-by-correl/414d5120514d3120202020202020202005b3b060019f0240
func getMsgByCorrelId(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Получение сообщения из IBM MQ")

  ctx, cancel := context.WithTimeout(rootCtx, time.Second*10)
  defer cancel()

  _, err := hex.DecodeString(r.URL.Path[1:])
  if err != nil {
    fmt.Fprintf(w, "get by correlId Error. need to pass correlId: %s", err.Error())
    return
  }

  msgId, ok, err := ibmmq.GetByCorrelId(ctx, nil, 100)
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
