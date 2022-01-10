package queue

import (
  "context"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "strings"
  "time"
)

func (q *Queue) Open() error {
  q.mx.Lock()
  defer q.mx.Unlock()

  q.log.Debug("Запрос на открытие очереди")

  if q.IsConn() {
    q.log.Warn(ErrAlreadyOpen)
    return ErrAlreadyOpen
  }

  if q.queueName == "" || len(q.perm) == 0 {
    q.log.Error(ErrNoConfig)
    return ErrNoConfig
  }

  mainCfg := q.manager.GetQueueConfig()
  q.devMode = mainCfg.DevMode
  q.h = mainCfg.Header

  q.ctx, q.ctxCanc = context.WithCancel(q.rootCtx)

  q.stateConn()

  select {
  case <-q.RegisterOpen():
  case <-q.ctx.Done():
  }

  q.log.Info("Очередь готова к работе")

  return nil
}

func (q *Queue) Close() error {
  q.log.Debug("Запрос на закрытие очереди")

  if q.IsDisconn() {
    q.log.Warn(ErrClosed)
    return ErrClosed
  }

  q.ctxCanc()
  q.stateDisconn()

  q.mx.Lock()
  defer q.mx.Unlock()

  select {
  case <-q.rootCtx.Done():
  case <-time.After(q.delayClose):
  }

  q.log.Info("Очередь закрыта")

  return nil
}

// Открывает очередь
func (q *Queue) open(m *ibmmq.MQQueueManager) error {
  mqod := ibmmq.NewMQOD()
  mqod.ObjectType = ibmmq.MQOT_Q
  mqod.ObjectName = q.queueName

  var flag int32
  for _, v := range q.perm {
    switch v {
    case permGet:
      flag |= ibmmq.MQOO_INPUT_SHARED
    case permBrowse:
      flag |= ibmmq.MQOO_BROWSE
    case permPut:
      flag |= ibmmq.MQOO_OUTPUT
    }
  }

  que, err := m.Open(mqod, flag)
  if err != nil {
    return err
  }

  q.log.WithFields(map[string]interface{}{
    "mod":  "open",
    "perm": strings.Join(q.permString(), ","),
  }).Info("Opened")

  q.que = &que
  return nil
}

func (q *Queue) close() {
  o := q.que
  if o != nil {
    q.que = nil
    err := o.Close(0)
    if err != nil {
      q.log.WithField("mod", "close").Warn(err)
    }
  }
}
