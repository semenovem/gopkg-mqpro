package main

import (
  "context"
  "fmt"
  mqpro "github.com/semenovem/gopkg_mqpro"
  "net/http"
  "time"
)

// Отправляет сообщение в очередь
//
// curl localhost:8080/put
func putMsg(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Отправка сообщения в IBM MQ")

  ctx, cancel := context.WithTimeout(rootCtx, time.Second*60)
  defer cancel()

  // Свойства сообщения
  props := map[string]interface{}{
    "firstProp":   "this is first prop",
    "anotherProp": "... another prop",
  }

  msg := &mqpro.Msg{
    Payload: []byte("Sending a message to IBM MQ"),
    Props:   props,
  }

  msgId, err := ibmmq.Put(ctx, msg)
  if err != nil {
    fmt.Fprintf(w, "put Error: %s\n", err.Error())
    return
  }

  fmt.Fprintf(w, "put Ok. msgId: %x\n", msgId)
}
