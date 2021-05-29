package mqpro

import (
  "context"
)

// GetByCorrelId Получение сообщения по correlId
func (p *Mqpro) GetByCorrelId(ctx context.Context, correlId []byte) (*Msg, bool, error) {
  fn := func(c *Mqconn) (*Msg, bool, error) {
    return c.GetByCorrelId(correlId)
  }
  return p.callGet(ctx, fn)
}

// Get Получение очередного сообщения
func (p *Mqpro) Get(ctx context.Context) (*Msg, bool, error) {
  fn := func(c *Mqconn) (*Msg, bool, error) {
    return c.Get()
  }
  return p.callGet(ctx, fn)
}

// GetByMsgId Получение сообщения по его MsgID.
func (p *Mqpro) GetByMsgId(ctx context.Context, msgId []byte) (*Msg, bool, error) {
  fn := func(c *Mqconn) (*Msg, bool, error) {
    return c.GetByMsgId(msgId)
  }
  return p.callGet(ctx, fn)
}

//
func (p *Mqpro) callGet(ctx context.Context, fn func(c *Mqconn) (*Msg, bool, error)) (*Msg, bool, error) {

  if len(p.connGet) == 0 {
    return nil, false, ErrNoConnection
  }

  var (
    msg *Msg
    err = ErrNoEstablishedConnection
    ok  bool
  )

loop:
  for ctx.Err() == nil {
    for _, conn := range p.connGet {
      if conn.IsConnected() {
        msg, ok, err = fn(conn)
        if err == nil || err != ErrConnBroken {
          break loop
        }
      }
    }

    select {
    case <-ctx.Done():
    case <-p.waitConnGet():
    }
  }

  return msg, ok, err
}
