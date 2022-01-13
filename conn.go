package mqpro

import (
  "context"
  "github.com/semenovem/gopkg_mqpro/v2/queue"
  "sync"
  "time"
)

func (m *Mqpro) Connect() error {
  m.log.Trace("Request to establish connection to IBM MQ...")

  m.mx.Lock()
  defer m.mx.Unlock()

  if !m.isConfigured() {
    m.log.Warn(ErrNoConfig)
    return ErrNoConfig
  }

  if m.isConnected {
    m.log.Warn(ErrAlreadyConnected)
    return ErrAlreadyConnected
  }

  m.ctx, m.ctxCanc = context.WithCancel(m.rootCtx)

  // Открытие очередей
  select {
  case <-m.ctx.Done():
  case err := <-m.openQues():
    if err != nil {
      return err
    }
  }

  return nil
}

func (m *Mqpro) openQues() <-chan error {
  var (
    ch   = make(chan error)
    wg   = sync.WaitGroup{}
    err1 error
  )

  // Запуск открытия очередей
  for _, q := range m.queues {
    if !q.IsConfigured() {
      m.log.Warnf("Очередь {%s} без конфигурации не будет открыта", q.Alias())
      continue
    }

    wg.Add(1)
    go func(q *queue.Queue) {
      defer wg.Done()

      err := q.Open()
      if err != nil {
        err1 = err
      }
    }(q)
  }

  // Ожидание открытия очередей
  go func() {
    defer close(ch)
    wg.Wait()
    if err1 != nil {
      ch <- err1
      return
    }
  }()

  return ch
}

func (m *Mqpro) Disconnect() error {
  m.log.Trace("Request to disconnect from IBM MQ...")

  if !m.isConnected {
    m.log.Warn(ErrNoEstablishedConnection)
    return ErrNoEstablishedConnection
  }

  m.ctxCanc()

  m.mx.Lock()
  defer m.mx.Unlock()

  select {
  case <-m.rootCtx.Done():
  case <-time.After(m.disconnDelay):
  }

  m.log.Info("Connection dropped")

  return nil
}
