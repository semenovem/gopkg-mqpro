package main

import (
  "context"
  "fmt"
  "github.com/semenovem/mqm/v2/queue"
  "sync"
  "time"
)

// Параллельное получение сообщений по MsgId/CorrelId
func testing1 () error {
  var (
    wg          = sync.WaitGroup{}
    ctx, cancel = context.WithTimeout(rootCtx, time.Second*30)
    countErr    int
    msg = &queue.Msg{}
    list [][]byte

  )
  defer cancel()

  ch, err := mqQueGet.Browse(ctx)
  if err != nil {
    log.Error("Сломан Browse")
    return err
  }

  for msg = range ch {
    list = append(list, msg.MsgId)
  }

  wg.Add(len(list))

  for i := 0; i < len(list); i++ {
    go func(b []byte) {
      defer wg.Done()
      //msg := &queue.Msg{MsgId: b}
      msg := &queue.Msg{CorrelId: b}

      err := mqQueGet.Get(ctx, msg)
      if err != nil {
        countErr++
        log.Error(">>>>>>>>>>> err = ", err)
        return
      }
    }(list[i])
  }
  wg.Wait()

  fmt.Println("#################################")
  fmt.Println("Ошибок = ", countErr)
  fmt.Println("#################################")

  return nil
}
