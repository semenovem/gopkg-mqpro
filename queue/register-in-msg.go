package queue

import (
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "time"
)

// RegisterInMsg регистрирует подписку на сообщения
func (q *Queue) RegisterInMsg(hnd func(*Msg)) {
  if q.hndInMsg != nil {
    q.log.Panic(ErrRegisterInMsg)
  }

  q.hndInMsg = hnd

  go q._subscInMsg()
}

// UnregisterInMsg удаляет подписку на сообщения
func (q *Queue) UnregisterInMsg() {
  if q.hndInMsg == nil {
    q.log.Panic("Нет подписки на приходящие сообщения для отмены")
  }
  q.unsubscInMsg()
  q.hndInMsg = nil
}

func (q *Queue) IsSubscribed() bool {
  return q.hndInMsg != nil
}

func (q *Queue) _subscInMsg() {
  var (
    err error
  )
  for q.IsOpen() {
    if err = q.subscInMsg(q.conn); err == nil {
      q.log.Info("Добавлена подписка на сообщения")
      return
    }

    q.errorHandler(err)

    select {
    case <-q.ctx.Done():
      return
    case <-time.After(q.reconnectDelay):
    }
  }
}

func (q *Queue) subscInMsg(conn *mqConn) error {
  if q.hndInMsg == nil {
    return nil
  }

  if !q.HasPermQueue(permGet) {
    q.log.Panic(ErrNotGetOpen)
  }

  cbd := ibmmq.NewMQCBD()
  gmo := ibmmq.NewMQGMO()
  getmqmd := ibmmq.NewMQMD()
  ctlo := ibmmq.NewMQCTLO()
  cmho := ibmmq.NewMQCMHO()

  cbd.CallbackFunction = q.handlerInMsg

  q.mxMsg.Lock()
  defer q.mxMsg.Unlock()

  mh, err := conn.m.CrtMH(cmho)
  if err != nil {
    q.log.Errorf(msgErrPropCreation, err)
    return err
  }

  gmo.Options = ibmmq.MQGMO_NO_SYNCPOINT | ibmmq.MQGMO_PROPERTIES_IN_HANDLE
  gmo.MsgHandle = mh

  err = conn.q.CB(ibmmq.MQOP_REGISTER, cbd, getmqmd, gmo)
  if err != nil {
    return err
  }

  err = conn.m.Ctl(ibmmq.MQOP_START, ctlo)
  if err != nil {
    return err
  }
  q.ctlo = ctlo

  return nil
}

func (q *Queue) handlerInMsg(
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
    q.log.Warnf("Ошибка подключения: %v", err)
    q.errorHandler(err)
    return
  }

  if err.MQCC != ibmmq.MQCC_OK {
    q.log.Warnf("Subscription error: %v", err)
    return
  }

  props, err1 := properties(gmo.MsgHandle)
  if err1 != nil {
    q.log.Errorf(msgErrPropGetting, err)
    return
  }

  msg := &Msg{
    MsgId:    md.MsgId,
    CorrelId: md.CorrelId,
    Payload:  buffer,
    Props:    props,
    Time:     md.PutDateTime,
  }

  var devMsg Msg
  if q.devMode {
    devMsg = *msg
    f1 := devMode(&devMsg, buffer, "subscribe")
    defer func() {
      f1()
    }()
  }

  if q.h == HeaderRfh2 {
    headers, err := q.Rfh2Unmarshal(buffer)
    if err != nil {
      q.log.Warn(err)
      return
    }
    msg.MQRFH2 = headers

    var ofs int32
    for _, h := range headers {
      unionPropsDeep(msg.Props, h.NameValues)
      ofs += h.StructLength
    }
    msg.Payload = buffer[ofs:]

    if q.devMode {
      devMsg.Payload = buffer[ofs:]
      devMsg.MQRFH2 = headers
      devMsg.Props = msg.Props
    }
  }

  q.log.Info("Получено сообщение")

  go q.hndInMsg(msg)
}

func (q *Queue) unsubscInMsg() {
  q.mxSubsc.Lock()
  defer q.mxSubsc.Unlock()

  if q.ctlo == nil {
    return
  }

  conn := q.conn
  if conn != nil {
    q.isWarn(conn.m.Ctl(ibmmq.MQOP_STOP, q.ctlo))
  }
  q.ctlo = nil
}
