package queue

import (
  "context"
  "fmt"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "github.com/sirupsen/logrus"
)

func (q *Queue) Get(ctx context.Context) (*Msg, bool, error) {
  l := q.log.WithField("method", "Get")

  msg, ok, err := q.get(ctx, operGet, nil, l)
  q.errorHandler(err)

  return msg, ok, err
}

// GetByCorrelId Извлекает сообщение из очереди по его CorrelId
func (q *Queue) GetByCorrelId(ctx context.Context, correlId []byte) (*Msg, bool, error) {
  l := q.log.WithFields(map[string]interface{}{
    "correlId": fmt.Sprintf("%x", correlId),
    "method":   "GetByCorrelId",
  })

  msg, ok, err := q.get(ctx, operGetByCorrelId, correlId, l)
  q.errorHandler(err)

  return msg, ok, err
}

// GetByMsgId Извлекает сообщение из очереди по его MsgId
func (q *Queue) GetByMsgId(ctx context.Context, msgId []byte) (*Msg, bool, error) {
  l := q.log.WithFields(map[string]interface{}{
    "msgId":  fmt.Sprintf("%x", msgId),
    "method": "GetByMsgId",
  })

  msg, ok, err := q.get(ctx, operGetByMsgId, msgId, l)
  q.errorHandler(err)

  return msg, ok, err
}

// Получение сообщения
func (q *Queue) get(ctx context.Context, oper queueOper, id []byte, l *logrus.Entry) (
  *Msg, bool, error) {

  if q.IsClosed() {
    l.Error(ErrNotOpen)
    return nil, false, ErrNotOpen
  }

  if q.ctlo != nil {
    return nil, false, ErrBusySubsc
  }

  var conn *mqConn

  select {
  case <-ctx.Done():
    return nil, false, ErrInterrupted
  case conn = <-q.RegisterOpen():
  }

  l.Trace("Start")

  var (
    datalen int
    err     error
    mqrc    *ibmmq.MQReturn
    buffer  = make([]byte, 0, 1024)
  )

  getmqmd := ibmmq.NewMQMD()
  gmo := ibmmq.NewMQGMO()
  cmho := ibmmq.NewMQCMHO()
  gmo.Options = ibmmq.MQGMO_NO_SYNCPOINT | ibmmq.MQGMO_PROPERTIES_IN_HANDLE

  q.mxMsg.Lock()
  defer q.mxMsg.Unlock()

  getMsgHandle, err := conn.m.CrtMH(cmho)
  if err != nil {
    l.Errorf("Ошибка создания объекта свойств сообщения: %s", err)

    if IsConnBroken(err) {
      err = ErrConnBroken
    } else {
      err = ErrGetMsg
    }
    return nil, false, err
  }
  defer func() {
    err := dltMh(getMsgHandle)
    if err != nil {
      l.Warnf("Ошибка удаления объекта свойств сообщения: %s", err)
    }
  }()

  gmo.MsgHandle = getMsgHandle

  switch q.h {
  case HeaderRfh2:
    getmqmd.Format = ibmmq.MQFMT_RF_HEADER_2
  default:
    // TODO код, получения стандартных заголовков перенести сюда
    getmqmd.Format = ibmmq.MQFMT_STRING
  }

  switch oper {
  case operGet:
  case operGetByMsgId:
    gmo.MatchOptions = ibmmq.MQMO_MATCH_MSG_ID
    getmqmd.MsgId = id
  case operGetByCorrelId:
    gmo.MatchOptions = ibmmq.MQMO_MATCH_CORREL_ID
    getmqmd.CorrelId = id
  case operBrowseFirst:
    gmo.Options |= ibmmq.MQGMO_BROWSE_FIRST
  case operBrowseNext:
    gmo.Options |= ibmmq.MQGMO_BROWSE_NEXT

  default:
    l.Panicf("Unknown operation. queueOper = %v", oper)
  }

loopCtx:
  for {
  loopGet:
    for i := 0; i < 2; i++ {
      buffer, datalen, err = conn.q.GetSlice(getmqmd, gmo, buffer)

      if err == nil {
        break loopCtx
      }

      mqrc = err.(*ibmmq.MQReturn)

      switch mqrc.MQRC {
      case ibmmq.MQRC_TRUNCATED_MSG_FAILED:
        buffer = make([]byte, 0, datalen)
        continue
      case ibmmq.MQRC_NO_MSG_AVAILABLE:
        err = nil
        break loopGet
      }

      l.Error(err)

      if IsConnBroken(err) {
        err = ErrConnBroken
      } else {
        err = ErrGetMsg
      }

      return nil, false, err
    }

    l.Debug("No message")

    return nil, false, nil
  }

  props, err := properties(getMsgHandle)
  if err != nil {
    l.Errorf("Ошибка получения свойств сообщения: %s", err)
    return nil, false, ErrGetMsg
  }

  l.Debug("Success")

  msg := &Msg{
    Payload:  buffer,
    Props:    props,
    CorrelId: getmqmd.CorrelId,
    MsgId:    getmqmd.MsgId,
    Time:     getmqmd.PutDateTime,
    MQRFH2:   make([]*MQRFH2, 0),
  }

  var devMsg Msg
  if q.devMode {
    devMsg = *msg
    f := devMode(&devMsg, buffer, "get")
    defer func() {
      f()
    }()
  }

  if q.h == HeaderRfh2 {
    headers, err := q.Rfh2Unmarshal(buffer)
    if err != nil {
      q.log.Error(err)
      return nil, false, err
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

  return msg, true, nil
}
