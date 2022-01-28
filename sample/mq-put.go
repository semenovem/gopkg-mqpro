package main

import (
  "context"
  "fmt"
  "github.com/semenovem/mqm/v2/queue"
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
  ctx, cancel := context.WithTimeout(rootCtx, time.Second*10)
  defer cancel()

  // Свойства сообщения
  props := map[string]interface{}{
    "foo":   "10101001110110",
    "BAR": "cb31e8610231",
  }

  b := []byte(`{"HoldJetFuelPaymentMsg":{"id":"f021d4ec-27f5-41be-8af3-946e65686902","result":"OK"}}`)

  msg := &queue.Msg{
    Payload:  b,
    Props:   props,
  }

  return msg.MsgId, mqQueFooPut.Put(ctx, msg)
}
