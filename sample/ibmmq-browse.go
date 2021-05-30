package main

import (
  "fmt"
  "net/http"
)

// Просмотр сообщений в очереди
func onBrowse(w http.ResponseWriter, _ *http.Request) {
  ch, err := ibmmq.Browse(rootCtx)

  if err != nil {
    fmt.Println("ERROR: ", err)
    return
  }

  i := 0
  for msg := range ch {
    i++
    logMsg(msg)
  }

  fmt.Fprintf(w, "[browse] Ok. %d messages viewed\n", i)
}
