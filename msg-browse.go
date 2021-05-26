package mqpro

import (
  "context"
)

func (p *Mqpro) Browse(ctx context.Context) (<-chan *Msg, error) {
  if len(p.connBrowse) == 0 {
    return nil, ErrNoConnection
  }

  var (
    err = ErrNoEstablishedConnection
    ch  <-chan *Msg
  )

loop:
  for ctx.Err() == nil {
    for _, conn := range p.connBrowse {
      if conn.IsConnected() {
        ch, err = conn.Browse(ctx)
        if err == nil || err != ErrConnBroken {
          break loop
        }
      }
    }

    select {
    case <-ctx.Done():
    case <-p.waitConnBrowse():
    }
  }

  return ch, err
}
