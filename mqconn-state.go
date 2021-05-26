package mqpro

import (
  "time"
)

func (c *Mqconn) state() {
  for cmd := range c.chMgr {

    switch c.stateConn {
    case stateConnect:
      if cmd != reqDisconnect {
        continue
      }
    case stateDisconnect:
      if cmd != reqConnect {
        continue
      }
      c.stateConn = stateErr
    case stateErr:
    default:
      c.log.Panic("Unknown state")
    }

    switch cmd {
    case reqConnect, reqReconnect:
      if c.connect() == nil {
        c.stateConn = stateConnect
        go c.fireEventConnected()
      } else {
        c.stateConn = stateErr
        c._disconnect()
        c.reconnect()

        <-time.After(c.reconnectDelay)
      }
    case reqDisconnect:
      c.disconnect()
      c.stateConn = stateDisconnect
      go c.fireEventDisconnected()
      close(c.chMgr)
    default:
      c.log.Panic("Unknown state")
    }
  }
}

func (c *Mqconn) reqConnect() {
  go func() {
    c.chMgr <- reqConnect
  }()
}

func (c *Mqconn) reqError() {
  if c.stateConn == stateDisconnect {
    return
  }
  if c.stateConn == stateConnect {
    c._disconnect()
  }
  c.stateConn = stateErr
  c.reconnect()
}

func (c *Mqconn) reqDisconnect() {
  go func() {
    c.chMgr <- reqDisconnect
  }()
}

func (c *Mqconn) reconnect() {
  go func() {
    c.chMgr <- reqReconnect
  }()
}
