package queue

import (
  "context"
  "fmt"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "github.com/pkg/errors"
  "github.com/sirupsen/logrus"
)

func (q *Queue) Get(ctx context.Context) (*Msg, error) {
  l := q.log.WithField("method", "Get")
  return q.get(ctx, operGet, nil, l)
}

// GetByCorrelId Извлекает сообщение из очереди по его CorrelId
func (q *Queue) GetByCorrelId(ctx context.Context, correlId []byte) (*Msg, error) {
  l := q.log.WithFields(map[string]interface{}{
    "correlId": fmt.Sprintf("%x", correlId),
    "method":   "GetByCorrelId",
  })
  return q.get(ctx, operGetByCorrelId, correlId, l)
}

// GetByMsgId Извлекает сообщение из очереди по его MsgId
func (q *Queue) GetByMsgId(ctx context.Context, msgId []byte) (*Msg, error) {
  l := q.log.WithFields(map[string]interface{}{
    "msgId":  fmt.Sprintf("%x", msgId),
    "method": "GetByMsgId",
  })
  return q.get(ctx, operGetByMsgId, msgId, l)
}

func (q *Queue) get(ctx context.Context, oper queueOper, id []byte, l *logrus.Entry) (
  *Msg, error) {

  if q.IsClosed() {
    l.Error(ErrNotOpen)
    return nil, ErrNotOpen
  }

  if q.ctlo != nil {
    l.Error(ErrBusySubsc)
    return nil, ErrBusySubsc
  }

  var (
    conn *mqConn
    err  error
    msg  *Msg
  )

  if q.devMode {
    defer func() {
      //if hd != nil {
      //  msg.MQRFH2, err = q.Rfh2Unmarshal(hd)
      //  if err != nil {
      //    l.Warn("DevMode: ", err)
      //  }
      //}

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
      return nil, ErrInterrupted
    case conn = <-q.RegisterOpen():
    }

    msg, err = q._get(oper, id, conn)
    if err == nil {
      return msg, nil
    }

    if q.errorHandler(errors.Cause(err)) {
      l.Warn(err)
      continue
    }

    l.Error(ErrGetMsg)
    return nil, ErrGetMsg
  }
}

// Получение сообщения
func (q *Queue) _get(oper queueOper, id []byte, conn *mqConn) (
  *Msg, error) {

  var (
    datalen      int
    err          error
    mqrc         *ibmmq.MQReturn
    buffer       = make([]byte, 0, 1024)
    getMsgHandle ibmmq.MQMessageHandle
    props        map[string]interface{}
    msg          *Msg
  )

  getmqmd := ibmmq.NewMQMD()
  gmo := ibmmq.NewMQGMO()
  cmho := ibmmq.NewMQCMHO()
  gmo.Options = ibmmq.MQGMO_NO_SYNCPOINT | ibmmq.MQGMO_PROPERTIES_IN_HANDLE

  getMsgHandle, err = conn.m.CrtMH(cmho)
  if err != nil {
    return nil, errors.Wrap(err, msgErrPropCreation)
  }

  defer func() {
    err = dltMh(getMsgHandle)
    if err != nil {
      q.log.WithField("id", fmt.Sprintf("%x", id)).Warnf(msgErrPropDeletion, err)
    }
  }()

  gmo.MsgHandle = getMsgHandle

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
    q.log.WithField("id", fmt.Sprintf("%x", id)).
      Panicf("Unknown operation. queueOper = %v", oper)
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
      err = nil
      return nil, nil
    }

    return nil, err
  }

  props, err = properties(getMsgHandle)
  if err != nil {
    return nil, errors.Wrap(err, msgErrPropGetting)
  }

  msg = &Msg{
    Payload:  buffer,
    Props:    props,
    CorrelId: getmqmd.CorrelId,
    MsgId:    getmqmd.MsgId,
    Time:     getmqmd.PutDateTime,
    MQRFH2:   make([]*MQRFH2, 0),
  }

  return msg, nil
}
