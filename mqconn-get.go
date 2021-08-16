package mqpro

import (
  "context"
  "fmt"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "github.com/sirupsen/logrus"
  "time"
)

func (c *Mqconn) Get(ctx context.Context) (*Msg, bool, error) {
  l := c.log.WithField("method", "Get")

  msg, ok, err := c.get(ctx, operGet, nil, l)
  if err == ErrConnBroken {
    c.reqError()
  }
  return msg, ok, err
}

// GetByCorrelId Извлекает сообщение из очереди по его CorrelId
func (c *Mqconn) GetByCorrelId(ctx context.Context, correlId []byte) (*Msg, bool, error) {
  l := c.log.WithFields(map[string]interface{}{
    "correlId": fmt.Sprintf("%x", correlId),
    "method":   "GetByCorrelId",
  })

  msg, ok, err := c.get(ctx, operGetByCorrelId, correlId, l)
  if err == ErrConnBroken {
    c.reqError()
  }
  return msg, ok, err
}

// GetByMsgId Извлекает сообщение из очереди по его MsgId
func (c *Mqconn) GetByMsgId(ctx context.Context, msgId []byte) (*Msg, bool, error) {
  l := c.log.WithFields(map[string]interface{}{
    "msgId":  fmt.Sprintf("%x", msgId),
    "method": "GetByMsgId",
  })

  msg, ok, err := c.get(ctx, operGetByMsgId, msgId, l)
  if err == ErrConnBroken {
    c.reqError()
  }

  return msg, ok, err
}

// получение сообщения
func (c *Mqconn) get(ctx context.Context, oper queueOper, id []byte, l *logrus.Entry) (
  *Msg, bool, error) {

  if !c.IsConnected() {
    l.Error(ErrNoConnection)
    return nil, false, ErrNoConnection
  }

  l.Info("Start")

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

  c.mxGet.Lock()
  getMsgHandle, err := c.mgr.CrtMH(cmho)
  c.mxGet.Unlock()

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

  // TODO - если MQFMT_RF_HEADER_2 - при получении не требуется - можно удалить код ниже:
  switch c.h {
  case HeaderRfh2:
    getmqmd.Format = ibmmq.MQFMT_RF_HEADER_2
  default:
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
      c.mxGet.Lock()
      buffer, datalen, err = c.que.GetSlice(getmqmd, gmo, buffer)
      c.mxGet.Unlock()

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

    select {
    case <-time.After(c.msgWaitInterval):
    case <-ctx.Done():
      l.Debug("No message")
      return nil, false, nil
    }
  }

  props, err := properties(getMsgHandle)
  if err != nil {
    l.Errorf("Ошибка получения свойств сообщения: %s", err)
    return nil, false, ErrGetMsg
  }

  l.Info("Success")

  ret := &Msg{
    Payload:  buffer,
    Props:    props,
    CorrelId: getmqmd.CorrelId,
    MsgId:    getmqmd.MsgId,
    Time:     getmqmd.PutDateTime,
    MQRFH2:   make([]*MQRFH2, 0),
  }

  if c.h == HeaderRfh2 {
    headers, err := c.Rfh2Unmarshal(buffer)
    if err != nil {
      c.log.Error(err)
      return nil, false, err
    }
    ret.MQRFH2 = headers

    var ofs int32
    for _, h := range headers {
      unionPropsDeep(ret.Props, h.NameValues)
      ofs += h.StrucLength
    }
    ret.Payload = buffer[ofs:]
  }

  //logMsg(ret, buffer)

  return ret, true, nil
}

// Вывод инфоормации о ibmmq сообщении
func logMsg(msg *Msg, b []byte) {
  fmt.Println("\n--------------------------------")
  fmt.Println("Сообщение:")
  if len(msg.Payload) < 200 {
    fmt.Printf(">>>>> msg.Payload  = %s\n", string(msg.Payload))
  } else {
    fmt.Printf(">>>>> msg.Payload len = %d\n", len(msg.Payload))
  }

  if len(b) > 0 {
    fmt.Printf(">>>>> original buffer  = %+v\n", b)
  }

  fmt.Printf(">>>>> msg.Props    = %+v\n", msg.Props)
  //fmt.Printf(">>>>> msg.CorrelId = %x\n", msg.CorrelId)
  //fmt.Printf(">>>>> msg.MsgId    = %x\n", msg.MsgId)
}
