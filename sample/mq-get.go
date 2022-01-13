package main

import (
  "context"
  "fmt"
  "github.com/semenovem/gopkg_mqpro/v2/queue"
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

  msg, ok, err := _getMsg()
  if err != nil {
    fmt.Println("[ERROR] при получении сообщения: ", err)
    _, _ = fmt.Fprintf(w, "[get] Error: %s\n", err.Error())
    return
  }

  if !ok {
    fmt.Println("[WARN] нет сообщений")
    _, _ = fmt.Fprintf(w, "[get]. Message queue is empty\n")
    return
  }

  _, _ = fmt.Fprintf(w, "[get] Ok. msgId: %x\n", msg.MsgId)
}

func _getMsg() (*queue.Msg, bool, error) {
  ctx, cancel := context.WithTimeout(rootCtx, time.Second*10)
  defer cancel()

  return mqOper1Get.Get(ctx)
}
