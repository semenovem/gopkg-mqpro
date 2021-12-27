package queue

import (
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "sync/atomic"
)


func (c *Conn) regConnMgr() <-chan *ibmmq.MQQueueManager {
  ch := make(chan *ibmmq.MQQueueManager)
  if c.state == stateConn {
    ch <- c.mgr
    close(ch)
  } else {
    c.subConnMgr <- ch
  }
  return ch
}





// deprecated
var ind uint32

func (c *Conn) fireEventConnected() {
  for i, ch := range c.fnsConn {
    delete(c.fnsConn, i)
    close(ch)
  }
}

func (c *Conn) fireEventDisconnected() {
  for i, ch := range c.fnsDisconn {
    delete(c.fnsDisconn, i)
    close(ch)
  }
}

func (c *Conn) RegisterEventConn() <-chan struct{} {
  ch := make(chan struct{})
  if c.IsConnected() {
    close(ch)
  } else {
    i := atomic.AddUint32(&ind, 1)
    c.fnsConn[i] = ch
  }
  return ch
}

func (c *Conn) RegisterEventDisconn() <-chan struct{} {
  ch := make(chan struct{})
  if c.IsDisconnected() {
    close(ch)
  } else {
    i := atomic.AddUint32(&ind, 1)
    c.fnsDisconn[i] = ch
  }
  return ch
}
