package main

import (
  "context"
  "fmt"
  mqpro "github.com/semenovem/gopkg_mqpro"
  "net/http"
  "time"
)

// Отправляет сообщение в очередь c CorlId
// Все тоже самое что и в putMsg️, но с установкой ему correlId
func putMsgCorrel(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Отправка сообщения в IBM MQ")

  ctx, cancel := context.WithTimeout(rootCtx, time.Second * 10)
  defer cancel()

  msg := &mqpro.Msg{
    CorrelId: correlId,
    Payload:  []byte("Sending a message to IBM MQ with correlId set"),
  }

  msgId, err := ibmmq.Put(ctx, msg)
  if err != nil {
    fmt.Fprintf(w, "put with correlId Error: %s\n", err.Error())
    return
  }

  fmt.Fprintf(w, "put with correlId Ok. msgId: %x\n", msgId)
}

