package mqpro

import (
  "context"
)

// Put Отправка сообщения
// return: ErrNoConnection | ErrNoEstablishedConnection | ErrConnBroken | ErrPutMsg
func (p *Mqpro) Put(
  ctx context.Context, msg *Msg) ([]byte, error) {

  if len(p.connPut) == 0 {
    return nil, ErrNoConnection
  }

  var (
    b   []byte = nil
    err        = ErrNoEstablishedConnection
  )

loop:
  for ctx.Err() == nil {
    for _, conn := range p.connPut {
      if conn.IsConnected() {
        b, err = conn.Put(msg)
        if err == nil || err != ErrConnBroken {
          break loop
        }
      }
    }

    select {
    case <-ctx.Done():
    case <-p.waitConnPut():
    }
  }

  return b, err
}
