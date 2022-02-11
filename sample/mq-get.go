package main

import (
  "context"
  "fmt"
  "github.com/semenovem/mqm/v2/queue"
  "net/http"
  "sync"
  "time"
)

// Получает сообщение из очереди
// curl host:port/get
func getMsg(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Получение сообщения из IBM MQ")
  var (
    msg = &queue.Msg{}
    err error
  )

  async := true

  // -----------------------------------------------------------
  // Асинхронное получение множества сообщений
  if async {
    var (
      wg          = sync.WaitGroup{}
      ctx, cancel = context.WithTimeout(rootCtx, time.Second*10)
      countErr    int
      //list [][]byte
    )
    defer cancel()

    //ch, err := mqQueGet.Browse(ctx)
    //if err != nil {
    //  log.Error("Сломан Browse")
    //  _, _ = fmt.Fprintf(w, "[get] Error. async Err: %s\n", err)
    //  return
    //}

    //for msg = range ch {
    //  list = append(list, msg.MsgId)
    //}

    for i := 0; i < 500; i++ {
      wg.Add(1)
      go func() {
        defer wg.Done()
        msg := &queue.Msg{}
        err := mqQueGet.Get(ctx, msg)
        if err != nil {
          countErr++
          fmt.Println(">>>>>>>>>>> err = ", err)
        } else {
          fmt.Println(">> ", formatMsgId(msg.MsgId))
        }
      }()
    }
    wg.Wait()
    fmt.Println("#################################")
    fmt.Println("Ошибок = ", countErr)
    fmt.Println("#################################")
  }
  // -----------------------------------------------------------

  //for i := 0; i < 200; i++ {
  //  go _getMsg()
  //}

  err = _getMsg(msg)
  if err != nil {
    fmt.Println("[ERROR] при получении сообщения: ", err)
    _, _ = fmt.Fprintf(w, "[get] Error: %s\n", err.Error())
    return
  }

  if msg.MsgId == nil {
    fmt.Println("[WARN] нет сообщений")
    _, _ = fmt.Fprintf(w, "[get]. Message queue is empty\n")
    return
  }

  _, _ = fmt.Fprintf(w, "[get] Ok. msgId: %x\n", msg.MsgId)
}

func _getMsg(msg *queue.Msg) error {
  ctx, cancel := context.WithTimeout(rootCtx, time.Second*10)
  defer cancel()

  return mqQueGet.Get(ctx, msg)
}
