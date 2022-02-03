package mqm

import (
  "context"
  "sync"
  "time"
)

func (m *Mqm) Connect() error {
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

  m.ctx, m.ctxEsc = context.WithCancel(m.rootCtx)

  // Открытие очередей
  select {
  case <-m.ctx.Done():
  case err := <-m.openQues():
    if err != nil {
      return err
    }
  }

  m.isConnected = true

  return nil
}

func (m *Mqm) openQues() <-chan error {
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
    go func(q Queue) {
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

func (m *Mqm) Disconnect() error {
  m.log.Trace("Request to disconnect from IBM MQ...")

  if !m.isConnected {
    m.log.Warn(ErrNoConnection)
    return ErrNoConnection
  }

  m.ctxEsc()

  m.mx.Lock()
  defer m.mx.Unlock()

  m.isConnected = false

  select {
  case <-m.rootCtx.Done():
  case <-time.After(m.disconnDelay):
  }

  m.log.Info("Connection dropped")

  return nil
}
