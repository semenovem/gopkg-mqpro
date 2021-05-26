package mqpro

import (
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "strings"
)

func (c *Mqconn) Get(msgId []byte) ([]byte, bool, error) {
  d, ok, err := c.get(msgId)

  if err != nil {
    go func() {
      c.reqError()
    }()
  }

  return d, ok, err
}

// получение сообщения
func (c *Mqconn) get(msgId []byte) ([]byte, bool, error) {
  var datalen int
  var err error

  getmqmd := ibmmq.NewMQMD()
  gmo := ibmmq.NewMQGMO()
  gmo.Options = ibmmq.MQGMO_NO_SYNCPOINT

  gmo.Options |= ibmmq.MQGMO_WAIT
  gmo.WaitInterval = 3 * 1000 // The WaitInterval is in milliseconds

  gmo.MatchOptions = ibmmq.MQMO_MATCH_MSG_ID
  getmqmd.MsgId = msgId
  getmqmd.Format = ibmmq.MQFMT_STRING

  buffer := make([]byte, 0, 1024)

  buffer, datalen, err = c.que.GetSlice(getmqmd, gmo, buffer)

  if err != nil {
    c.log.Error("004", err)
    mqret := err.(*ibmmq.MQReturn)
    if mqret.MQRC == ibmmq.MQRC_NO_MSG_AVAILABLE {
      // If there's no message available, then I won't treat that as a real error as
      // it's an expected situation
      return nil, false, nil
    }
  }

  c.log.Debugf("Got message of length %d: ", datalen)
  c.log.Debug(strings.TrimSpace(string(buffer)))

  return buffer, true, nil
}
