package mqpro

import "github.com/ibm-messaging/mq-golang/v5/ibmmq"

func (c *Mqconn) Type() TypeConn {
  return c.typeConn
}

func (c *Mqconn) IsConnected() bool {
  return c.stateConn == stateConnect
}

func (c *Mqconn) IsDisconnected() bool {
  return c.stateConn == stateDisconnect
}

func (c *Mqconn) isWarnConn(err error) {
  if err != nil {
    mqret := err.(*ibmmq.MQReturn)
    if mqret == nil || mqret.MQRC != ibmmq.MQRC_CONNECTION_BROKEN {
      c.log.Warn(err)
    }
  }
}

func (c *Mqconn) isWarn(err error) {
  if err != nil {
    c.log.Warn(err)
  }
}

func (c *Mqconn) GetRootTag() string {
  return c.rfh2RootTag
}

func (c *Mqconn) SetRootTag(tag string) {
  c.rfh2RootTag = tag
}
