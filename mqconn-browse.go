package mqpro

import (
  "context"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

func (c *Mqconn) Browse(ctx context.Context) (<-chan *Msg, error) {
  ch, err := c.browse(ctx)
  if err != nil {
    if HasConnBroken(err) {
      c.reqError()
      err = ErrConnBroken
    } else {
      err = ErrBrowseMsg
    }
  }

  return ch, err
}

func (c *Mqconn) browse(ctx context.Context) (<-chan *Msg, error) {
  c.mxBrowse.Lock()
  defer c.mxBrowse.Unlock()

  if !c.IsConnected() {
    return nil, ErrNoConnection
  }

  c.log.Trace("Start open BROWSE")

  var (
    ch   = make(chan *Msg)
    wait = make(chan struct{})
    err  error
  )

  go func(w chan struct{}) {
    var msg *Msg
    browseOption := ibmmq.MQGMO_BROWSE_FIRST
    for ctx.Err() == nil {
      msg, err = c._browse(browseOption)
      if err != nil {
        if err.(*ibmmq.MQReturn).MQRC == ibmmq.MQRC_NO_MSG_AVAILABLE {
          err = nil
        } else {
          c.log.Error("Ошибка при обзоре сообщений: ", err)
        }
        break
      }
      if w != nil {
        close(w)
        w = nil
      }
      ch <- msg
      browseOption = ibmmq.MQGMO_BROWSE_NEXT
    }
    if w != nil {
      close(w)
    }
    close(ch)
    c.log.Trace("Закрытие канала обзора сообщений BROWSE")
  }(wait)

  select {
  case <-ctx.Done():
  case <-wait:
  }

  if err != nil {
    return nil, err
  }

  c.log.Debugf("Success open for BROWSE")

  return ch, nil
}

func (c *Mqconn) _browse(browseOption int32) (*Msg, error) {
  l := c.log.WithField("method", "Browse messages")

  var (
    err     error
    datalen int
  )

  buffer := make([]byte, 0, 1024)
  getmqmd := ibmmq.NewMQMD()
  gmo := ibmmq.NewMQGMO()

  gmo.Options = ibmmq.MQGMO_NO_SYNCPOINT | browseOption | ibmmq.MQGMO_WAIT
  gmo.Options |= ibmmq.MQGMO_PROPERTIES_IN_HANDLE
  gmo.WaitInterval = c.cfg.WaitInterval

  cmho := ibmmq.NewMQCMHO()
  getMsgHandle, err := c.mgr.CrtMH(cmho)
  if err != nil {
    c.log.Error("Ошибка создания объекта свойств сообщения: ", err)
    return nil, err
  }
  defer func() {
    err := dltMh(getMsgHandle)
    if err != nil {
      l.Warnf("Ошибка удаления объекта свойств сообщения: %s", err.Error())
    }
  }()
  gmo.MsgHandle = getMsgHandle

  for i := 0; i < 2; i++ {
    buffer, datalen, err = c.que.GetSlice(getmqmd, gmo, buffer)

    if err != nil {
      if err.(*ibmmq.MQReturn).MQRC == ibmmq.MQRC_TRUNCATED_MSG_FAILED {
        buffer = make([]byte, 0, datalen)
        continue
      }
    }
    break
  }

  if err != nil {
    return nil, err
  }

  props, err := properties(getMsgHandle)
  if err != nil {
    c.log.Error("Ошибка получения свойств сообщения: ", err)
    return nil, err
  }

  m := &Msg{
    Payload:  buffer,
    MsgId:    getmqmd.MsgId,
    CorrelId: getmqmd.CorrelId,
    Props:    props,
  }

  return m, nil
}
