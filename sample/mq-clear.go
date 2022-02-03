package main

import (
  "context"
  "fmt"
  "net/http"
  "time"
)

// Удаляет все сообщения из очереди
// curl host:port/clear
func clearQueue(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Удаление всех сообщений в очереди")

  ctx, cancel := context.WithTimeout(rootCtx, time.Second*10)
  defer cancel()

  ch, err := mqQueGet.Browse(ctx)
  if err != nil {
    fmt.Println("ERROR: ", err)
    return
  }

  i := 0
  for msg := range ch {
    i++

    m, err := mqQueGet.GetByMsgId(ctx, msg.MsgId)
    if err != nil {
      fmt.Println("[ERROR] при получении сообщения из очереди: ", err)
      continue
    }

    if m != nil {
      logMsgDel(m)
    }
  }

  _, _ = fmt.Fprintf(w, "[clear] Ok. %d messages deleted\n", i)
}
