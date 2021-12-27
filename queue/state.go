package queue

import (
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "time"
)

func (c *Conn) workerState() {
  var (
    err error
    mgr *ibmmq.MQQueueManager
  )

  for cmd := range c.chState {
    if c.state == cmd {
      continue
    }

    switch c.state {
    case stateConn:
      if cmd != stateDisconn {
        continue
      }
      c.state = stateConn
      mgr, err = c.connToMgr()
      if err != nil {
        c.log.Errorf("Failed connection attempt to IBM MQ manager: %s", err)
        go c.reqError()
        continue
      }
      c.log.Info("Connected to IBM MQ manager")
      c.mgr = mgr

    case stateDisconn:
      c.state = stateDisconn
      c.ctxCns()
      if c.mgr != nil {
        _ = c.mgr.Disc()
      }

    case stateErr:
      if c.state == stateDisconn {
        continue
      }
      c.state = stateErr

      if c.mgr != nil {
        _ = c.mgr.Disc()
      }

      select {
      case <-c.ctx.Done():
        continue
      case <-time.After(c.reconnectDelay):
        go c.reqConn()
      }
    }
  }
}

func (c *Conn) reqConn() {
  c.chState <- stateConn
}

func (c *Conn) reqError() {
  c.chState <- stateErr
}

func (c *Conn) reqDisconn() {
  c.chState <- stateDisconn
}
