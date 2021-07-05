package mqpro

import (
  "fmt"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "github.com/sirupsen/logrus"
)

// Put отправка сообщения в очередь
func (c *Mqconn) Put(msg *Msg) ([]byte, error) {
  l := c.log.WithField("method", "Put")

  if msg.CorrelId != nil {
    l = l.WithField("correlId", fmt.Sprintf("%x", msg.CorrelId))
  }

  d, err := c.put(msg, l)
  if err == ErrConnBroken {
    c.reqError()
  }

  return d, err
}

func (c *Mqconn) put(msg *Msg, l *logrus.Entry) ([]byte, error) {
  if !c.IsConnected() {
    return nil, ErrNoConnection
  }

  c.mxPut.Lock()
  defer c.mxPut.Unlock()

  l.Info("Start")

  putmqmd := ibmmq.NewMQMD()
  pmo := ibmmq.NewMQPMO()
  cmho := ibmmq.NewMQCMHO()

  putMsgHandle, err := c.mgr.CrtMH(cmho)
  if err != nil {
    l.Errorf("Ошибка создания объекта свойств сообщения: %s", err)

    if IsConnBroken(err) {
      err = ErrConnBroken
    } else {
      err = ErrPutMsg
    }

    return nil, err
  }

  err = setProps(&putMsgHandle, msg.Props, l)
  if err != nil {
    return nil, ErrPutMsg
  }

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
    l.Error("Ошибка отправки сообщения: ", err)

    if IsConnBroken(err) {
      err = ErrConnBroken
    } else {
      err = ErrPutMsg
    }

    return nil, err
  }

  l.Infof("Success. MsgId: %x", putmqmd.MsgId)

  return putmqmd.MsgId, nil
}
