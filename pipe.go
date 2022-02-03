package mqm

import (
  "context"
  "github.com/semenovem/mqm/v2/queue"
  "time"
)

type Pipe struct {
  get   Queue
  put   Queue
  alias string
}

func (m *Mqm) cfgPipe(cfg *PipeCfg) error {
  var (
    p = m.findPipeByAlias(cfg.Alias)
  )
  if p == nil {
    return ErrNotFoundPipe
  }
  p.put = m.NewQueue(cfg.Alias + "_put")
  p.get = m.NewQueue(cfg.Alias + "_get")

  err := p.put.CfgByStr(cfg.Put+":put")
  if err != nil {
    return err
  }

  err = p.get.CfgByStr(cfg.Get+":get,browse")
  if err != nil {
    return err
  }

  return nil
}

func (m *Mqm) findPipeByAlias(a string) *Pipe {
  for _, p := range m.pipes {
    if p.alias == a {
      return p
    }
  }
  return nil
}

// NewPipe Объект канала (имеет в своем составе две очереди: отправка/получение)
func (m *Mqm) NewPipe(a string) Queue {
  p := &Pipe{alias: a}
  m.pipes = append(m.pipes, p)
  return p
}

func (c *Pipe) Put(ctx context.Context, msg *queue.Msg) error {
  return c.put.Put(ctx, msg)
}

func (c *Pipe) Get(ctx context.Context, msg *queue.Msg) error {
  return c.get.Get(ctx, msg)
}

func (c *Pipe) GetByMsgId(ctx context.Context, msgId []byte) (*queue.Msg, error) {
  return c.get.GetByMsgId(ctx, msgId)
}

func (c *Pipe) GetByCorrelId(ctx context.Context, correlId []byte) (*queue.Msg, error) {
  return c.get.GetByCorrelId(ctx, correlId)
}

func (c *Pipe) Browse(ctx context.Context) (<-chan *queue.Msg, error) {
  return c.get.Browse(ctx)
}

func (c *Pipe) Alias() string {
  return c.alias
}

// CfgByStr конфигурирование через строку не поддерживается
func (c *Pipe) CfgByStr(_ string) error {
  return ErrChannelCfgNotSup
}

func (c *Pipe) IsConfigured() bool {
  return c.put.IsConfigured() && c.get.IsConfigured()
}

func ttt() error {
  time.Sleep(time.Millisecond * 500)
  return nil
}

func ttt2() error {
  time.Sleep(time.Millisecond * 200)
  return nil
}

func (c *Pipe) Open() error {
  var (
    ch  = make(chan error)
    err error
  )
  go func() { ch <- c.put.Open() }()
  go func() { ch <- c.get.Open() }()
  err = <-ch
  if err != nil {
    go func() { <-ch; close(ch) }()
    return err
  }
  err = <-ch
  close(ch)
  return err
}

func (c *Pipe) Close() error {
  var (
    ch  = make(chan error)
    err error
  )
  go func() { ch <- c.put.Close() }()
  go func() { ch <- c.get.Close() }()
  err = <-ch
  if err != nil {
    go func() { <-ch; close(ch) }()
    return err
  }
  err = <-ch
  close(ch)
  return err
}

func (c *Pipe) UpdateBaseCfg() {
  c.put.UpdateBaseCfg()
  c.get.UpdateBaseCfg()
}
func (c *Pipe) IsSubscribed() bool {
  return c.get.IsSubscribed()
}

func (c *Pipe) RegisterInMsg(hnd func(*queue.Msg)) {
  c.get.RegisterInMsg(hnd)
}
func (c *Pipe) UnregisterInMsg() {
  c.get.UnregisterInMsg()
}
