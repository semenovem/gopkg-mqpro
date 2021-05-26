package mqpro

import (
  "context"
)

// GetByCorrelId Получение сообщения по correlId
func (p *Mqpro) GetByCorrelId(ctx context.Context, correlId []byte, t int) (*Msg, bool, error) {

  if len(p.connPut) == 0 {
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
        msg, ok, err = conn.GetByCorrelId(correlId, t)
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

//// GetById Получение сообщения по его MsgID
//func (p *Mqpro) GetById_(msgId []byte) ([]byte, bool, error) {
//  conn, err := p.getConns(&p.putConn)
//  if err != nil {
//    return nil, false, err
//  }
//
//  return conn.Get(msgId)
//}
