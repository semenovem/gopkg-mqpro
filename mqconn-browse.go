package mqpro

import (
  "context"
)

func (c *Mqconn) Browse(ctx context.Context) (<-chan *Msg, error) {
  ch, err := c.browse(ctx)
  if err == ErrConnBroken {
    c.reqError()
  }

  return ch, err
}

func (c *Mqconn) browse(ctx context.Context) (<-chan *Msg, error) {
  l := c.log.WithField("method", "BrowseOpen")

  c.mxBrowse.Lock()
  defer c.mxBrowse.Unlock()

  if !c.IsConnected() {
    c.log.Error(ErrNoConnection)
    return nil, ErrNoConnection
  }

  l.Info("Start open BROWSE")

  var (
    ch   = make(chan *Msg)
    wait = make(chan struct{})
    err  error
    ok   bool
  )

  go func(w chan struct{}) {
    var msg *Msg
    cx, cancel := context.WithCancel(ctx)
    cancel()
    ll := l.WithField("method", "BrowseGet")
    oper := operBrowseFirst

    for ctx.Err() == nil {
      msg, ok, err = c.get(cx, oper, nil, ll)
      if err != nil || !ok {
        break
      }

      if w != nil {
        close(w)
        w = nil
      }
      ch <- msg
      oper = operBrowseNext
    }
    if w != nil {
      close(w)
    }
    close(ch)
    l.Info("Закрытие канала обзора сообщений BROWSE")
  }(wait)

  select {
  case <-ctx.Done():
  case <-wait:
  }

  if err != nil {
    return nil, err
  }

  l.Info("Success open for BROWSE")

  return ch, nil
}
