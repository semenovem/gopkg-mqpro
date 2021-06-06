package mqpro

import (
  "fmt"
  "github.com/sirupsen/logrus"
)

func MqconnNew(tc TypeConn, l *logrus.Entry, c *Cfg) *Mqconn {
  o := &Mqconn{
    cfg:             c,
    fnsConn:         map[uint32]chan struct{}{},
    fnsDisconn:      map[uint32]chan struct{}{},
    reconnectDelay:  defReconnectDelay,
    stateConn:       stateDisconnect,
    msgWaitInterval: defMsgWaitInterval,
  }

  m := map[string]interface{}{
    "conn": fmt.Sprintf("%s|%s|%s|%s|%s",
      o.endpoint(), c.MgrName, c.ChannelName, c.QueueName, typeConnTxt[tc]),
  }

  o.log = l.WithFields(m)

  if tc != TypePut && tc != TypeGet && tc != TypeBrowse {
    o.log.Panic("Unknown connection type")
  }

  o.typeConn = tc

  if c.MaxMsgLength == 0 {
    c.MaxMsgLength = defMaxMsgLength
  }

  return o
}
