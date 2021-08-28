package mqpro

import (
  "context"
)


// Get2 Получение очередного сообщения
//func (p *Mqpro) Get2(ctx context.Context, s ...string) (*Msg, bool, error) {
//
//}

// GetByCorrelId Получение сообщения по correlId
func (p *Mqpro) GetByCorrelId(ctx context.Context, correlId []byte) (*Msg, bool, error) {
  fn := func(c *Mqconn) (*Msg, bool, error) {
    return c.GetByCorrelId(ctx, correlId)
  }
  return p.callGet(ctx, fn)
}

// Get Получение очередного сообщения
func (p *Mqpro) Get(ctx context.Context) (*Msg, bool, error) {
  fn := func(c *Mqconn) (*Msg, bool, error) {
    return c.Get(ctx)
  }
  return p.callGet(ctx, fn)
}

// GetByMsgId Получение сообщения по его MsgID.
func (p *Mqpro) GetByMsgId(ctx context.Context, msgId []byte) (*Msg, bool, error) {
  fn := func(c *Mqconn) (*Msg, bool, error) {
    return c.GetByMsgId(ctx, msgId)
  }
  return p.callGet(ctx, fn)
}

//
func (p *Mqpro) callGet(ctx context.Context, fn func(c *Mqconn) (*Msg, bool, error)) (
  *Msg, bool, error) {

  if len(p.connGet) == 0 {
    p.log.Error(ErrNoConnection)
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
