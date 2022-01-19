package queue

import (
  "context"
  "fmt"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "github.com/sirupsen/logrus"
)

// Put отправка сообщения в очередь
func (q *Queue) Put(ctx context.Context, msg *Msg) error {
  l := q.log.WithField("method", "Put")
  if msg.CorrelId != nil {
    l = l.WithField("correlId", fmt.Sprintf("%x", msg.CorrelId))
  }
  return q.put(ctx, msg, l)
}

func (q *Queue) put(ctx context.Context, msg *Msg, l *logrus.Entry) error {
  if q.IsClosed() {
    return ErrNotOpen
  }

  var (
    conn    *mqConn
    err     error
    hd      []byte
    payload []byte
  )

  if q.devMode {
    defer func() {
      if hd != nil {
        msg.MQRFH2, err = q.Rfh2Unmarshal(hd)
        if err != nil {
          l.Warn("DevMode: ", err)
        }
      }

      logMsg(msg, payload, "put")
    }()
  }

  if q.h == HeaderRfh2 {
    hd, err = q.Rfh2Marshal(msg.Props)
    if err != nil {
      l.Error("Не удалось подготовить сообщение с заголовками rfh2: ", err)
      return err
    }
    if msg.Payload == nil {
      payload = hd
    } else {
      payload = append(hd, msg.Payload...)
    }
  } else {
    if msg.Payload == nil {
      payload = make([]byte, 0)
    } else {
      payload = make([]byte, len(msg.Payload))
      copy(payload, msg.Payload)
    }
  }

  q.mxMsg.Lock()
  defer q.mxMsg.Unlock()

  for {
    select {
    case <-ctx.Done():
      return ErrInterrupted
    case conn = <-q.RegisterOpen():
    }

    err = q._put(conn, msg, payload, l)
    if err == nil {
      return nil
    }

    if q.errorHandler(err) {
      continue
    }
    return err
  }
}

func (q *Queue) _put(conn *mqConn, msg *Msg, payload []byte, l *logrus.Entry) error {
  var (
    err     error
    putmqmd = ibmmq.NewMQMD()
    pmo     = ibmmq.NewMQPMO()
    cmho    = ibmmq.NewMQCMHO()
  )

  pmo.Options = ibmmq.MQPMO_NO_SYNCPOINT

  if msg.CorrelId != nil {
    putmqmd.CorrelId = msg.CorrelId
  }

  switch q.h {
  case HeaderRfh2:
    putmqmd.Format = ibmmq.MQFMT_RF_HEADER_2

  default:
    var putMsgHandle ibmmq.MQMessageHandle

    putmqmd.Format = ibmmq.MQFMT_STRING
    putMsgHandle, err = conn.m.CrtMH(cmho)
    if err != nil {
      l.Errorf(msgErrPropCreation, err)
      return err
    }

    err = setProps(&putMsgHandle, msg.Props, l)
    if err != nil {
      l.Errorf(msgErrPropSetting, err)
      return err
    }
    pmo.OriginalMsgHandle = putMsgHandle

    defer func() {
      err = dltMh(putMsgHandle)
      if err != nil {
        l.Warnf(msgErrPropDeletion, err)
      }
    }()
  }

  err = conn.q.Put(putmqmd, pmo, payload)
  if err == nil {
    msg.MsgId = putmqmd.MsgId
    msg.Time = putmqmd.PutDateTime
  }

  return err
}
