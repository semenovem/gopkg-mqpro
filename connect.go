package mqpro

import (
  "context"
  "sync"
  "time"
)

func (p *Mqpro) Connect() error {
  p.mx.Lock()
  defer p.mx.Unlock()

  ctx, cancel := context.WithCancel(p.rootCtx)
  p.ctx = ctx
  p.ctxCancel = cancel

  p.log.Trace("Request to establish connection to IBM MQ...")

  if len(p.conns) == 0 {
    p.log.Error(ErrNoData)
    return ErrNoData
  }

  for _, conn := range p.conns {
    conn.Connect(p.reconnDelay)
  }

  go func() {
    <-p.ctx.Done()
    p.Disconnect()
  }()

  var wg sync.WaitGroup
  p.waitConnAll(&wg, p.connBrowse)
  p.waitConnAll(&wg, p.connGet)
  p.waitConnAll(&wg, p.connPut)
  wg.Wait()

  return nil
}

func (p *Mqpro) waitConnAll(wg *sync.WaitGroup, conns []*Mqconn) {
  if len(conns) != 0 {
    wg.Add(1)
    go func() {
      select {
      case <-p.waitConn(conns):
      case <-p.ctx.Done():
      }
      wg.Done()
    }()
  }
}

func (p *Mqpro) Disconnect() {
  p.log.Trace("Request disconnect...")

  if p.ctx == nil {
    p.log.Trace("Already disconnected")
    return
  }
  if p.ctx.Err() != nil {
    p.ctxCancel()
  }

  p.mx.Lock()
  defer p.mx.Unlock()

  p.ctx = nil
  p.ctxCancel = nil

  for _, conn := range p.conns {
    conn.Disconnect()
  }

  select {
  case <-p.waitDisconn():
  case <-time.After(p.delayBeforeDisconnect):
  }
  p.log.Info("Disconnected")
}

func (p *Mqpro) waitConnPut() <-chan struct{} {
  return p.waitConn(p.connPut)
}

func (p *Mqpro) waitConnGet() <-chan struct{} {
  return p.waitConn(p.connGet)
}

func (p *Mqpro) waitConnBrowse() <-chan struct{} {
  return p.waitConn(p.connBrowse)
}

func (p *Mqpro) waitConn(conns []*Mqconn) <-chan struct{} {
  if len(conns) == 0 {
    p.log.Panic(ErrNoConnection)
  }

  cc := make(chan struct{})
  is := true
  var mx sync.Mutex

  for _, c := range conns {
    go func(ch <-chan struct{}) {
      <-ch
      mx.Lock()
      if is {
        is = false
        close(cc)
      }
      mx.Unlock()
    }(c.RegisterEventConn())
  }

  return cc
}

func (p *Mqpro) waitDisconn() <-chan struct{} {
  cc := make(chan struct{})
  var wg sync.WaitGroup

  for _, c := range p.conns {
    go func(ch <-chan struct{}) {
      wg.Add(1)
      <-ch
      wg.Done()
    }(c.RegisterEventDisconn())
  }

  go func() {
    wg.Wait()
    close(cc)
  }()

  return cc
}
