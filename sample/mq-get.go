package main

import (
  "context"
  "fmt"
  "github.com/semenovem/mqm/v2/queue"
  "net/http"
  "time"
)

// Получает сообщение из очереди
// curl host:port/get
func getMsg(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Получение сообщения из IBM MQ")

  //for i := 0; i < 200; i++ {
  //  go _getMsg()
  //}

  msg, err := _getMsg()
  if err != nil {
    fmt.Println("[ERROR] при получении сообщения: ", err)
    _, _ = fmt.Fprintf(w, "[get] Error: %s\n", err.Error())
    return
  }

  if msg == nil {
    fmt.Println("[WARN] нет сообщений")
    _, _ = fmt.Fprintf(w, "[get]. Message queue is empty\n")
    return
  }

  _, _ = fmt.Fprintf(w, "[get] Ok. msgId: %x\n", msg.MsgId)
}

func _getMsg() (*queue.Msg, error) {
  ctx, cancel := context.WithTimeout(rootCtx, time.Second*10)
  defer cancel()

  return mqQueFooGet.Get(ctx)
}
