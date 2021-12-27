package queue

import "github.com/ibm-messaging/mq-golang/v5/ibmmq"

func (c *Conn) IsConnected() bool {
  return c.state == stateConn
}

func (c *Conn) IsDisconnected() bool {
  return c.state == stateDisconn
}

func (c *Conn) isWarnConn(err error) {
  if err != nil {
    mqret := err.(*ibmmq.MQReturn)
    if mqret == nil || mqret.MQRC != ibmmq.MQRC_CONNECTION_BROKEN {
      c.log.Warn(err)
    }
  }
}

func (c *Conn) isWarn(err error) {
  if err != nil {
    c.log.Warn(err)
  }
}
