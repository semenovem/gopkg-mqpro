package main

import (
  "context"
  "fmt"
  "github.com/google/uuid"
  "github.com/semenovem/mqm/v2/queue"
  "net/http"
  "strconv"
  "time"
)

var mqPutFooCounter uint64 = 10000000000

// Отправляет сообщение в очередь
// curl host:port/put
func putMsg(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Отправка сообщения в IBM MQ")

  for i := 0; i < 500; i++ {
  go func() {
    _, err := _putMsg()
    if err != nil {
      fmt.Println(">>>>>> ERR: ", err)
    }
  }()
  }

  msg, err := _putMsg()
  if err != nil {
    _, _ = fmt.Fprintf(w, "put Error: %s\n", err.Error())
    return
  }

  _, _ = fmt.Fprintf(w, "put Ok. msgId: %x\n", msg.MsgId)
}

func _putMsg() (*queue.Msg, error) {
  ctx, cancel := context.WithTimeout(rootCtx, time.Second*5)
  defer cancel()

  mqPutFooCounter++

  // Свойства сообщения
  props := map[string]interface{}{
    "foo":  strconv.FormatUint(mqPutFooCounter, 10),
    "BAR": "cb31e8610231",
  }

  id := uuid.New().String()
  b := []byte(`{"HoldJetFuelPaymentMsg":{"id":"` + id + `","result":"OK"}}`)

  msg := &queue.Msg{
    Payload: b,
    Props:   props,
  }

  return msg, mqQuePut.Put(ctx, msg)
}
