package mqpro

import "sync/atomic"

func (c *Mqconn) fireEventConnected() {
  for i, ch := range c.fnsConn {
    delete(c.fnsConn, i)
    close(ch)
  }
}

func (c *Mqconn) fireEventDisconnected() {
  for i, ch := range c.fnsDisconn {
    delete(c.fnsDisconn, i)
    close(ch)
  }
}

func (c *Mqconn) RegisterEventConn() <-chan struct{} {
  ch := make(chan struct{})
  if c.IsConnected() {
    close(ch)
  } else {
    i := atomic.AddUint32(&c.ind, 1)
    c.fnsConn[i] = ch
  }
  return ch
}

func (c *Mqconn) RegisterEventDisconn() <-chan struct{} {
  ch := make(chan struct{})
  if c.IsDisconnected() {
    close(ch)
  } else {
    i := atomic.AddUint32(&c.ind, 1)
    c.fnsDisconn[i] = ch
  }
  return ch
}
