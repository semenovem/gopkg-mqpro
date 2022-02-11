package queue

import (
  "context"
  "fmt"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "github.com/pkg/errors"
  "github.com/sirupsen/logrus"
  "time"
)

var pMqdmho = ibmmq.NewMQDMHO()

func (q *Queue) Get(ctx context.Context, msg *Msg) error {
  var (
    l    *logrus.Entry
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

// Получение сообщения
func (q *Queue) get(ctx context.Context, op queueOper, msg *Msg, l *logrus.Entry) (err error) {
  var (
    datalen int
    buffer  = make([]byte, 0, 1024)
    props   map[string]interface{}
    conn    *mqConn
    getmqmd = ibmmq.NewMQMD()
    gmo     = ibmmq.NewMQGMO()
    cmho    = ibmmq.NewMQCMHO()
    off     = false
  )
  gmo.Options = ibmmq.MQGMO_NO_SYNCPOINT | ibmmq.MQGMO_PROPERTIES_IN_HANDLE

  if q.IsClosed() {
    l.Error(ErrNotOpen)
    return ErrNotOpen
  }

  if q.ctlo != nil {
    l.Error(ErrBusySubsc)
    return ErrBusySubsc
  }

  if q.devMode {
    defer func() {
      if msg == nil {
        return
      }
      logMsg(msg, nil, "get | get by correl id | get msg id")
    }()
  }

  defer func() {
    if err != nil {
      l.Error(err)
    }

    if off {
      err1 := gmo.MsgHandle.DltMH(pMqdmho)
      if err1 != nil {
        l.WithField("msg", fmt.Sprintf("%+v", msg)).Warnf(msgErrPropDeletion, err1)
      }
    }
  }()

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
    l.WithField("msg", fmt.Sprintf("%+v", msg)).
      Panicf("Unknown operation. queueOper = %v", op)
  }

  q.mxMsg.Lock()
  defer q.mxMsg.Unlock()

loopMain:
  for {
    select {
    case <-ctx.Done():
      return ErrInterrupted

    case conn = <-q.RegisterOpen():
      gmo.MsgHandle, err = conn.m.CrtMH(cmho)
      off = true
      if err != nil {
        return errors.Wrap(err, msgErrPropCreation)
      }

    loopGet:
      for {
        buffer, datalen, err = conn.q.GetSlice(getmqmd, gmo, buffer)

        if err == nil {
          break loopMain
        }

        switch err.(*ibmmq.MQReturn).MQRC {

        // Не достаточен размер выделенной памяти
        case ibmmq.MQRC_TRUNCATED_MSG_FAILED:
          buffer = make([]byte, 0, datalen)
          continue

        // Нет сообщений
        case ibmmq.MQRC_NO_MSG_AVAILABLE:
          if op == operGet || op == operBrowseNext {
            msg.Erase()
            return nil
          }

          select {
          case <-time.After(q.waitInterval):
          case <-ctx.Done():
            msg.Erase()
            return nil
          }

        // Ошибка получения сообщения
        default:
          if q.errorHandler(err) {
            l.Warn(err)
            break loopGet
          } else {
            return err
          }
        }
      }
    }

    off = false
    err1 := gmo.MsgHandle.DltMH(pMqdmho)
    if err1 != nil {
      l.Warnf(msgErrPropDeletion, err1)
    }
  }

  props, err = properties(gmo.MsgHandle)
  if err != nil {
    err = errors.Wrap(err, msgErrPropGetting)
    return err
  }

  // -------------
  off = false
  err1 := gmo.MsgHandle.DltMH(pMqdmho)
  if err1 != nil {
    l.Warnf(msgErrPropDeletion, err1)
  }

  msg.MsgId = getmqmd.MsgId
  msg.CorrelId = getmqmd.CorrelId
  msg.Props = props
  msg.Payload = buffer
  msg.Time = getmqmd.PutDateTime
  msg.MQRFH2 = nil

  return nil
}
