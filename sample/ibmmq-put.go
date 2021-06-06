package main

import (
  "context"
  "fmt"
  mqpro "github.com/semenovem/gopkg_mqpro"
  "net/http"
  "time"
)

var flag bool

// Отправляет сообщение в очередь
// curl host:port/put
func putMsg(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Отправка сообщения в IBM MQ")

  ctx, cancel := context.WithTimeout(rootCtx, time.Second*60)
  defer cancel()

  // Свойства сообщения
  props := map[string]interface{}{
    "firstProp":   "this is first prop",
    "anotherProp": "... another prop",
  }

  size := 1024
  b := make([]byte, size)

  for i := 0; i < size; i++ {
    b[i] = byte(i)
  }

  msg := &mqpro.Msg{
    Payload:  b,
    Props:    props,
    CorrelId: correlId,
  }

  flag = !flag
  if flag {
    msg.CorrelId = correlId
  } else {
    msg.CorrelId = correlId2
  }

  msgId, err := ibmmq.Put(ctx, msg)
  if err != nil {
    fmt.Fprintf(w, "put Error: %s\n", err.Error())
    return
  }

  fmt.Fprintf(w, "put Ok. msgId: %x\n", msgId)
}
