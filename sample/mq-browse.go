package main

import (
  "context"
  "fmt"
  "net/http"
  "time"
)

// Просмотр сообщений в очереди
// curl host:port/browse
func onBrowse(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Просмотр сообщений в очереди")

  ctx, canc := context.WithTimeout(rootCtx, time.Second*10)
  defer canc()

  ch, err := mqOper1Get.Browse(ctx)
  if err != nil {
    fmt.Println("ERROR: ", err)
    return
  }

  i := 0
  for msg := range ch {
    i++
    logMsgIn(msg)
  }

  if i == 0 {
    fmt.Println("Нет сообщений в очереди", i)
  } else {
    fmt.Printf("Кол-во ообщений в очереди: %d\n", i)
  }

  _, _ = fmt.Fprintf(w, "[browse] Ok. %d messages viewed\n", i)
}
