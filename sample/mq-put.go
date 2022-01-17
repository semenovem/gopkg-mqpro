package main

import (
  "context"
  "fmt"
  "github.com/semenovem/gopkg_mqpro/v2/queue"
  "net/http"
  "time"
)

// Отправляет сообщение в очередь
// curl host:port/put
func putMsg(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Отправка сообщения в IBM MQ")

  //for i := 0; i < 200; i++ {
  // go _putMsg()
  //}

  msgId, err := _putMsg()
  if err != nil {
    _, _ = fmt.Fprintf(w, "put Error: %s\n", err.Error())
    return
  }

  _, _ = fmt.Fprintf(w, "put Ok. msgId: %x\n", msgId)
}

func _putMsg() ([]byte, error) {
  ctx, cancel := context.WithTimeout(rootCtx, time.Second*7)
  defer cancel()

  // Свойства сообщения
  props := map[string]interface{}{
    "firstProp":   "this is first prop",
    "anotherProp": "... another prop",
  }

  size := 8 * 1
  b := make([]byte, size)

  for i := 0; i < size; i++ {
    b[i] = byte(i)
  }

  msg := &queue.Msg{
    Payload: b,
    Props:   props,
  }

  return mqQueFooPut.Put(ctx, msg)
}
