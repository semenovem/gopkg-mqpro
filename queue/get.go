package queue

import (
  "context"
  "fmt"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "github.com/pkg/errors"
  "github.com/sirupsen/logrus"
)

func (q *Queue) Get(ctx context.Context, msg *Msg) error {
  var (
    l   *logrus.Entry
    oper queueOper
  )
  if msg.MsgId != nil {
    oper = operGetByMsgId
    l = q.log.WithField("method", "Get")
  } else if msg.CorrelId != nil {
    oper = operGetByCorrelId
    l = q.log.WithFields(map[string]interface{}{
      "correlId": fmt.Sprintf("%x", msg.CorrelId),
      "method":   "GetByCorrelId",
    })
  } else {
    oper = operGet
    l = q.log.WithField("method", "Get")
  }
  return q.get(ctx, oper, msg, l)
}

// GetByCorrelId Извлекает сообщение из очереди по его CorrelId
func (q *Queue) GetByCorrelId(ctx context.Context, correlId []byte) (*Msg, error) {
  l := q.log.WithFields(map[string]interface{}{
    "correlId": fmt.Sprintf("%x", correlId),
    "method":   "GetByCorrelId",
  })
  msg := &Msg{
    CorrelId: correlId,
  }
  err := q.get(ctx, operGetByCorrelId, msg, l)
  return msg, err
}

// GetByMsgId Извлекает сообщение из очереди по его MsgId
func (q *Queue) GetByMsgId(ctx context.Context, msgId []byte) (*Msg, error) {
  l := q.log.WithFields(map[string]interface{}{
    "msgId":  fmt.Sprintf("%x", msgId),
    "method": "GetByMsgId",
  })
  msg := &Msg{
    MsgId: msgId,
  }
  err := q.get(ctx, operGetByMsgId, msg, l)
  return msg, err
}

func (q *Queue) get(ctx context.Context, op queueOper, msg *Msg, l *logrus.Entry) error {
  if q.IsClosed() {
    l.Error(ErrNotOpen)
    return ErrNotOpen
  }

  if q.ctlo != nil {
    l.Error(ErrBusySubsc)
    return ErrBusySubsc
  }

  var (
    conn *mqConn
    err  error
  )

  if q.devMode {
    defer func() {
      if msg == nil {
        return
      }
      logMsg(msg, nil, "get | get by correl id | get msg id")
    }()
  }

  q.mxMsg.Lock()
  defer q.mxMsg.Unlock()

  for {
    select {
    case <-ctx.Done():
      l.Error(ErrInterrupted)
      return ErrInterrupted
    case conn = <-q.RegisterOpen():
    }

    err = q._get(op, msg, conn)
    if err == nil {
      return nil
    }

    if q.errorHandler(errors.Cause(err)) {
      l.Warn(err)
      continue
    }

    l.Error(err)
    return err
  }
}

// Получение сообщения
func (q *Queue) _get(op queueOper, msg *Msg, conn *mqConn) error {
  var (
    datalen      int
    err          error
    mqrc         *ibmmq.MQReturn
    buffer       = make([]byte, 0, 1024)
    getMsgHandle ibmmq.MQMessageHandle
    props        map[string]interface{}
  )

  getmqmd := ibmmq.NewMQMD()
  gmo := ibmmq.NewMQGMO()
  cmho := ibmmq.NewMQCMHO()
  gmo.Options = ibmmq.MQGMO_NO_SYNCPOINT | ibmmq.MQGMO_PROPERTIES_IN_HANDLE

  getMsgHandle, err = conn.m.CrtMH(cmho)
  if err != nil {
    return errors.Wrap(err, msgErrPropCreation)
  }

  defer func() {
    err = dltMh(getMsgHandle)
    if err != nil {
      q.log.WithField("msg", fmt.Sprintf("%+v", msg)).
        Warnf(msgErrPropDeletion, err)
    }
  }()

  gmo.MsgHandle = getMsgHandle

  switch op {
  case operGet:
  case operGetByMsgId:
    gmo.MatchOptions = ibmmq.MQMO_MATCH_MSG_ID
    getmqmd.MsgId = msg.MsgId
  case operGetByCorrelId:
    gmo.MatchOptions = ibmmq.MQMO_MATCH_CORREL_ID
    getmqmd.CorrelId = msg.CorrelId
  case operBrowseFirst:
    gmo.Options |= ibmmq.MQGMO_BROWSE_FIRST
  case operBrowseNext:
    gmo.Options |= ibmmq.MQGMO_BROWSE_NEXT
  default:
    q.log.WithField("msg", fmt.Sprintf("%+v", msg)).
      Panicf("Unknown operation. queueOper = %v", op)
  }

  for i := 0; i < 2; i++ {
    buffer, datalen, err = conn.q.GetSlice(getmqmd, gmo, buffer)
    if err == nil {
      break
    }

    mqrc = err.(*ibmmq.MQReturn)
    switch mqrc.MQRC {
    case ibmmq.MQRC_TRUNCATED_MSG_FAILED:
      buffer = make([]byte, 0, datalen)
      continue
    case ibmmq.MQRC_NO_MSG_AVAILABLE:
      msg.Erase()
      err = nil
      return nil
    }

    return err
  }

  props, err = properties(getMsgHandle)
  if err != nil {
    return errors.Wrap(err, msgErrPropGetting)
  }

  msg.MsgId = getmqmd.MsgId
  msg.CorrelId = getmqmd.CorrelId
  msg.Props = props
  msg.Payload = buffer
  msg.Time = getmqmd.PutDateTime
  msg.MQRFH2 = nil

  return nil
}
