package mqpro

import (
  "fmt"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

// Put отправка сообщения в очередь
func (c *Mqconn) Put(msg *Msg) ([]byte, error) {
  d, err := c.put(msg)
  if err != nil {
    if HasConnBroken(err) {
      c.reqError()
      err = ErrConnBroken
    } else {
      err = ErrPutMsg
    }
  }

  return d, err
}

func (c *Mqconn) put(msg *Msg) ([]byte, error) {
  l := c.log
  if msg.CorrelId != nil {
    l = c.log.WithField("correlId", fmt.Sprintf("%x", msg.CorrelId))
  }

  c.mxPut.Lock()
  defer c.mxPut.Unlock()

  if !c.IsConnected() {
    return nil, ErrNoConnection
  }

  l.Tracef("Start PUT message")

  msgId, err := c._put(msg)
  if err != nil {
    l.Errorf("Failed to PUT message: %v", err)
    return nil, err
  }

  l.Debugf("Success PUT message. MsgId: %x", msgId)

  return msgId, nil
}

func (c *Mqconn) _put(msg *Msg) ([]byte, error) {
  cmho := ibmmq.NewMQCMHO()
  putMsgHandle, err := c.mgr.CrtMH(cmho)
  if err != nil {
    return nil, err
  }

  err = c.setProps(&putMsgHandle, msg.Props)
  if err != nil {
    return nil, err
  }

  putmqmd := ibmmq.NewMQMD()
  pmo := ibmmq.NewMQPMO()

  if msg.CorrelId != nil {
    putmqmd.CorrelId = msg.CorrelId
  }

  pmo.Options = ibmmq.MQPMO_NO_SYNCPOINT
  pmo.OriginalMsgHandle = putMsgHandle
  putmqmd.Format = ibmmq.MQFMT_STRING

  var d []byte
  if msg.Payload == nil {
    d = make([]byte, 0)
  } else {
    d = msg.Payload
  }

  err = c.que.Put(putmqmd, pmo, d)
  if err != nil {
    return nil, err
  }

  return putmqmd.MsgId, nil
}

func (c *Mqconn) setProps(h *ibmmq.MQMessageHandle, props map[string]interface{}) error {
  var err error
  smpo := ibmmq.NewMQSMPO()
  pd := ibmmq.NewMQPD()

  for k, v := range props {
    err = h.SetMP(smpo, k, pd, v)
    if err != nil {
      c.log.Errorf("Failed to set message property '%s' value '%v': %v", k, v, err)
      return err
    }
  }

  return nil
}
