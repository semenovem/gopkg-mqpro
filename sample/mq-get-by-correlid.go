package main

import (
  "context"
  "encoding/hex"
  "fmt"
  "net/http"
  "strings"
  "time"
)

// Получает сообщение из очереди по CorrelID
// curl host:port/correl
func getMsgByCorrelId(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Получение сообщения по CorrelID")

  ctx, cancel := context.WithTimeout(rootCtx, time.Second*5)
  defer cancel()

  var s string

  if strings.HasPrefix(r.URL.Path, "/") && len(r.URL.Path) > 2 {
    s = r.URL.Path[1:]
  } else {
    s = r.URL.Path
  }

  vv := strings.Split(s, "/")

  if len(vv) < 2 {
    fmt.Println("[ERROR] ошибка в url")
    _, _ = fmt.Fprintf(w, "[correl] Error: on url\n")
    return
  }

  id := vv[1]
  if !regMsgId.MatchString(id) {
    fmt.Println("[ERROR] передан не валидный CorrelID. " +
      "Ожидается: [414d5120514d31202020202020202020607527610b9a0340]")
    _, _ = fmt.Fprintf(w, "[correl] Error: invalid CorrelID passed\n")
    return
  }

  b, err := hex.DecodeString(id)
  if err != nil {
    fmt.Println("[ERROR] передан не валидный CorrelID. " +
      "Ожидается: [414d5120514d31202020202020202020607527610b9a0340]")
    _, _ = fmt.Fprintf(w, "[correl] Error: invalid CorrelID passed\n")
    return
  }

  msg, ok, err := ibmmqOper1Get.GetByCorrelId(ctx, b)
  if err != nil {
    fmt.Println("[ERROR] ошибка при получении сообщения по CorrelID: ", id)
    _, _ = fmt.Fprintf(w, "[correl] Error: %s\n", err.Error())
    return
  }

  if !ok {
    fmt.Println("[WARN] нет сообщения с CorrelID: ", id)
    _, _ = fmt.Fprintf(w, "[correl] Warn: message not found\n")
    return
  }

  fmt.Println("[INFO] получено сообщение с CorrelID: ", id)

  logMsgIn(msg)

  _, _ = fmt.Fprintf(w, "[correl] Ok. msgId: %x\n", msg.MsgId)
}
