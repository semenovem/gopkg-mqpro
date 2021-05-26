package mqpro

import "github.com/ibm-messaging/mq-golang/v5/ibmmq"

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

func HasConnBroken(err error) bool {
  mqrc := err.(*ibmmq.MQReturn).MQRC
  return mqrc == ibmmq.MQRC_CONNECTION_BROKEN || mqrc == ibmmq.MQRC_CONNECTION_QUIESCING
}
