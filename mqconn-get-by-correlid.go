package mqpro

import (
  "fmt"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

// GetByCorrelId Извлекает сообщение из очереди по его CorrelId
func (c *Mqconn) GetByCorrelId(correlId []byte) (*Msg, bool, error) {
  msg, ok, err := c.getByCorrelId(correlId)
  if err != nil {
    if HasConnBroken(err) {
      c.reqError()
      err = ErrConnBroken
    } else {
      err = ErrGetMsg
    }
  }

  return msg, ok, err
}

func (c *Mqconn) getByCorrelId(correlId []byte) (*Msg, bool, error) {
  l := c.log.WithField("correlId", fmt.Sprintf("%x", correlId))

  c.mxPut.Lock()
  defer c.mxPut.Unlock()

  l.Tracef("Start GET by correlId")

  if !c.IsConnected() {
    return nil, false, ErrNoConnection
  }

  msg, ok, err := c._getByCorrelId(correlId)
  if err != nil {
    c.log.Errorf("Failed to GET by correlId: %v", err)
    return nil, false, err
  }

  c.log.Debugf("Success GET by correlId of length %d", len(msg.Payload))

  return msg, ok, err
}

func (c *Mqconn) _getByCorrelId(correlId []byte) (*Msg, bool, error) {
  var datalen int
  var err error

  getmqmd := ibmmq.NewMQMD()
  gmo := ibmmq.NewMQGMO()
  getmqmd.Format = ibmmq.MQFMT_STRING
  gmo.Options = ibmmq.MQGMO_NO_SYNCPOINT

  gmo.Options |= ibmmq.MQGMO_WAIT
  gmo.WaitInterval = int32(3 * 1000) // The WaitInterval is in milliseconds
  gmo.MatchOptions = ibmmq.MQMO_MATCH_CORREL_ID
  getmqmd.CorrelId = correlId

  cmho := ibmmq.NewMQCMHO()

  c.mxGet.Lock()
  defer c.mxGet.Unlock()

  getMsgHandle, err := c.mgr.CrtMH(cmho)
  if err != nil {
    c.log.Error("Ошибка создания объекта свойств сообщения: ", err)
    return nil, false, err
  }
  defer DltMh(getMsgHandle, c.log)
  gmo.MsgHandle = getMsgHandle
  gmo.Options |= ibmmq.MQGMO_PROPERTIES_IN_HANDLE

  buffer := make([]byte, 0, 1024)

  for i := 0; i < 2; i++ {
    buffer, datalen, err = c.que.GetSlice(getmqmd, gmo, buffer)
    if err != nil {
      mqret := err.(*ibmmq.MQReturn)

      if mqret.MQRC == ibmmq.MQRC_TRUNCATED_MSG_FAILED {
        buffer = make([]byte, 0, datalen)
        continue
      }

      if mqret.MQRC == ibmmq.MQRC_NO_MSG_AVAILABLE {
        err = nil
      } else {
        c.log.Error("Ошибка получения сообщения: ", err, "  len: ", datalen)
      }

      return nil, false, err
    }

    break
  }

  props, err := Properties(getMsgHandle, c.log)
  if err != nil {
    return nil, false, err
  }

  ret := &Msg{
    Payload:  buffer,
    Props:    props,
    CorrelId: correlId,
    MsgId:    getmqmd.MsgId,
  }

  return ret, true, nil
}
