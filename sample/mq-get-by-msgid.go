package main

import (
  "context"
  "encoding/hex"
  "fmt"
  "net/http"
  "strings"
  "time"
)

// Получает сообщение из очереди по MsgID
// curl host:port/msgid
func getMsgByMsgId(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Получение сообщения по MsgId")

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
    _, _ = fmt.Fprintf(w, "[msgid] Error: on url\n")
    return
  }

  id := vv[1]
  if !regMsgId.MatchString(id) {
    fmt.Println("[ERROR] передан не валидный MsgID. " +
      "Ожидается: [414d5120514d31202020202020202020607527610b9a0340]")
    _, _ = fmt.Fprintf(w, "[msgid] Error: invalid MsgID passed\n")
    return
  }

  b, err := hex.DecodeString(id)
  if err != nil {
    fmt.Println("[ERROR] передан не валидный MsgID. " +
      "Ожидается: [414d5120514d31202020202020202020607527610b9a0340]")
    _, _ = fmt.Fprintf(w, "[msgid] Error: invalid MsgID passed\n")
    return
  }

  msg, ok, err := mqQueFooGet.GetByMsgId(ctx, b)
  if err != nil {
    fmt.Println("[ERROR] ошибка при получении сообщения по MsgID: ", id)
    _, _ = fmt.Fprintf(w, "[msgid] Error: %s\n", err.Error())
    return
  }

  if !ok {
    fmt.Println("[WARN] нет сообщения с MsgID: ", id)
    _, _ = fmt.Fprintf(w, "[msgid] Warn: message not found\n")
    return
  }

  fmt.Println("[INFO] получено сообщение с MsgID: ", id)

  logMsgIn(msg)

  _, _ = fmt.Fprintf(w, "[msgid] Ok. msgId: %x\n", msg.MsgId)
}
