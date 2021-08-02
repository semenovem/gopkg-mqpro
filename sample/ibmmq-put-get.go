package main

import (
  "context"
  "fmt"
  mqpro "github.com/semenovem/gopkg_mqpro"
  "net/http"
  "time"
)

// положить / считать сообщение из очереди
// curl host:port/putget
func putGetMsg(w http.ResponseWriter, _ *http.Request) {

  ctx, cancel := context.WithTimeout(rootCtx, time.Second*5)
  defer cancel()

  // Свойства сообщения
  props := map[string]interface{}{
    "firstProp":   "this is first prop",
    "anotherProp": "... another prop",
  }

  // Отправляемое сообщение
  msg := &mqpro.Msg{
    Payload:  []byte("Sending a message to IBM MQ"),
    Props:    props,
    CorrelId: correlId,
  }

  fmt.Println("\n--------------------------------")
  fmt.Printf("Отправляет сообщение в очередь с установленным correlId:\n")
  _, err := ibmmq.Put(ctx, msg)
  if err != nil {
    fmt.Println("ERROR: ", err)
    return
  }

  fmt.Println("\n--------------------------------")
  fmt.Printf("Получает сообщение по correlId = %x\n", correlId)
  msg, ok, err := ibmmq.GetByCorrelId(ctx, correlId)
  if err != nil {
    fmt.Println("ERROR: ", err)
    return
  }

  fmt.Printf(">>>>> ok           = %t\n", ok)
  if ok {
    fmt.Printf(">>>>> msg.Payload  = %s\n", string(msg.Payload))
    fmt.Printf(">>>>> msg.Props    = %+v\n", msg.Props)
    fmt.Printf(">>>>> msg.CorrelId = %x\n", msg.CorrelId)
    fmt.Printf(">>>>> msg.MsgId    = %x\n", msg.MsgId)
  }

  // ----------------------------------------------------------
  // ----------------------------------------------------------
  // ----------------------------------------------------------
  msg = &mqpro.Msg{
    Payload:  []byte("Sending a message to IBM MQ"),
    Props:    props,
    CorrelId: correlId,
  }
  fmt.Println("\n--------------------------------")
  fmt.Printf("Отправляет сообщение в очередь:\n")
  msgId, err := ibmmq.Put(ctx, msg)
  if err != nil {
    fmt.Println("ERROR: ", err)
    return
  }

  fmt.Println("\n--------------------------------")
  fmt.Printf("Получает сообщение по msgId = %x\n", msgId)
  msg, ok, err = ibmmq.GetByMsgId(ctx, msgId)
  if err != nil {
    fmt.Println("ERROR: ", err)
    return
  }

  fmt.Printf(">>>>> ok           = %t\n", ok)
  if ok {
    fmt.Printf(">>>>> msg.Payload  = %s\n", string(msg.Payload))
    fmt.Printf(">>>>> msg.Props    = %+v\n", msg.Props)
    fmt.Printf(">>>>> msg.CorrelId = %x\n", msg.CorrelId)
    fmt.Printf(">>>>> msg.MsgId    = %x\n", msg.MsgId)
  }

  // ----------------------------------------------------------
  // ----------------------------------------------------------
  // ----------------------------------------------------------
  msg = &mqpro.Msg{
    Payload:  []byte("Sending a message to IBM MQ"),
    Props:    props,
    CorrelId: correlId,
  }
  fmt.Println("\n--------------------------------")
  fmt.Printf("Отправляет сообщение в очередь:\n")
  _, err = ibmmq.Put(ctx, msg)
  if err != nil {
    fmt.Println("ERROR: ", err)
    return
  }

  fmt.Println("\n--------------------------------")
  fmt.Println("Получает сообщение первое доступное сообщение")
  msg, ok, err = ibmmq.Get(ctx)
  if err != nil {
    fmt.Println("ERROR: ", err)
    return
  }

  fmt.Printf(">>>>> ok           = %t\n", ok)
  if ok {
    fmt.Printf(">>>>> msg.Payload  = %s\n", string(msg.Payload))
    fmt.Printf(">>>>> msg.Props    = %+v\n", msg.Props)
    fmt.Printf(">>>>> msg.CorrelId = %x\n", msg.CorrelId)
    fmt.Printf(">>>>> msg.MsgId    = %x\n", msg.MsgId)
  }

  fmt.Fprintf(w, "[putget] Ok\n")
}
