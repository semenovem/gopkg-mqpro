package main

import (
  "context"
  "fmt"
  "github.com/semenovem/gopkg_mqpro/v2/queue"
  "net/http"
  "time"
)

// Отправляет сообщение в очередь и ждет ответа
// curl host:port/ping
// Requirements:
// - client
// - server ENV_MIRROR=true
// не должно быть других client с активной подпиской на ту же очередь
func mqPing(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Отправка сообщения ping в IBM MQ")

  ctx, cancel := context.WithTimeout(rootCtx, time.Second*5)
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

  msgId, err := ibmmqOper1In.Put(ctx, msg)
  if err != nil {
    _, _ = fmt.Fprintf(w, "[ping] Error: %s\n", err.Error())
    return
  }

  msg.MsgId = msgId
  logMsgOut(msg)

  fmt.Println()
  fmt.Println("Ждем ответа: ")

  reply, ok, err := ibmmqOper1In.GetByCorrelId(ctx, msgId)

  if err != nil {
    fmt.Println("[ERROR] ошибка при получении сообщения по CorrelID")
    _, _ = fmt.Fprintf(w, "[ping] Error receiving response: %s\n", err.Error())
    return
  }

  if !ok {
    fmt.Println("[WARN] нет ответного сообщения")
    _, _ = fmt.Fprintf(w, "[ping] Warn no response was received\n")
    return
  }

  logMsgIn(reply)

  _, _ = fmt.Fprintf(w, "[ping] Ok. msgId: %x\n", msgId)
}
