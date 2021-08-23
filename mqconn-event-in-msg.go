package mqpro

import (
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

// RegisterEventInMsg подписка на входящие сообщения
func (c *Mqconn) RegisterEventInMsg(fn func(*Msg)) {
  c.fnInMsg = fn

  if c.typeConn != TypeGet {
    c.log.Panic("The connection must be of type: 'TypeGet'")
  }

  if c.IsConnected() {
    if c.registerEventInMsg() != nil {
      c.reqError()
    }
  }
}

func (c *Mqconn) registerEventInMsg() error {
  c.log.Trace("Subscribing to incoming messages...")
  err := c._registerEventInMsg()
  if err != nil {
    c.log.Error("Failed to subscribe to incoming messages")
    return err
  }
  c.log.Info("Subscribed to incoming messages")
  return nil
}

// Подписка на входящие сообщения
func (c *Mqconn) _registerEventInMsg() error {
  cmho := ibmmq.NewMQCMHO()
  mh, err := c.mgr.CrtMH(cmho)
  if err != nil {
    c.log.Error("Ошибка создания объекта свойств сообщения", err)
    return err
  }

  cbd := ibmmq.NewMQCBD()
  gmo := ibmmq.NewMQGMO()
  getmqmd := ibmmq.NewMQMD()
  ctlo := ibmmq.NewMQCTLO()

  gmo.Options = ibmmq.MQGMO_NO_SYNCPOINT | ibmmq.MQGMO_PROPERTIES_IN_HANDLE
  gmo.MsgHandle = mh

  cbd.CallbackFunction = c.handlerInMsg

  err = c.que.CB(ibmmq.MQOP_REGISTER, cbd, getmqmd, gmo)
  if err != nil {
    return err
  }

  err = c.mgr.Ctl(ibmmq.MQOP_START, ctlo)
  if err != nil {
    return err
  }
  c.ctlo = ctlo

  return nil
}

func (c *Mqconn) UnregisterInMsg() {
  c.log.Trace("Unsubscribing from incoming messages...")
  c._unregisterInMsg()
  c.log.Info("unsubscribed from incoming messages")
}

func (c *Mqconn) _unregisterInMsg() {
  if c.ctlo != nil {
    c.isWarn(c.mgr.Ctl(ibmmq.MQOP_STOP, c.ctlo))
    c.ctlo = nil
  }
}

func (c *Mqconn) handlerInMsg(
  _ *ibmmq.MQQueueManager,
  _ *ibmmq.MQObject,
  md *ibmmq.MQMD,
  gmo *ibmmq.MQGMO,
  buffer []byte,
  _ *ibmmq.MQCBC,
  err *ibmmq.MQReturn) {

  if err.MQRC == ibmmq.MQRC_NO_MSG_AVAILABLE {
    return
  }

  if err.MQRC == ibmmq.MQRC_CONNECTION_BROKEN {
    c.log.Warnf("Ошибка подключения: %v", err)
    c.reqError()
    return
  }

  if err.MQCC != ibmmq.MQCC_OK {
    c.log.Warnf("Subscription error: %v", err)
    return
  }

  props, err1 := properties(gmo.MsgHandle)
  if err1 != nil {
    c.log.Error("Ошибка при получении свойств сообщения: ", err)
    return
  }

  msg := &Msg{
    MsgId:    md.MsgId,
    CorrelId: md.CorrelId,
    Payload:  buffer,
    Props:    props,
    Time:     md.PutDateTime,
  }

  if c.DevMode {
    c.devMsg = *msg
    f := devMode(&c.devMsg, buffer, "subscribe")
    defer func() {
      f()
    }()
  }

  if c.h == HeaderRfh2 {
    headers, err := c.Rfh2Unmarshal(buffer)
    if err != nil {
      c.log.Warn(err)
      return
    }
    msg.MQRFH2 = headers

    var ofs int32
    for _, h := range headers {
      unionPropsDeep(msg.Props, h.NameValues)
      ofs += h.StructLength
    }
    msg.Payload = buffer[ofs:]

    if c.DevMode {
      c.devMsg.Payload = buffer[ofs:]
      c.devMsg.MQRFH2 = headers
      c.devMsg.Props = msg.Props
    }
  }

  c.log.Info("Получено сообщение")

  go c.fnInMsg(msg)
}
