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

  q.log.Trace("Запрос на открытие очереди")

  if q.IsOpen() {
    q.log.Warn(ErrAlreadyOpen)
    return ErrAlreadyOpen
  }

  if !q.IsConfigured() {
    q.log.Error(ErrNotConfigured)
    return ErrNotConfigured
  }

  if !q.manager.IsConfigured() {
    q.log.Error(ErrManagerNotConfigured)
    return ErrManagerNotConfigured
  }

  q.UpdateBaseCfg()

  q.ctx, q.ctxCanc = context.WithCancel(q.rootCtx)

  q.stateOpen()

  select {
  case <-q.RegisterOpen():
  case <-q.ctx.Done():
  }

  q.log.Info("Очередь готова к работе")

  return nil
}

func (q *Queue) Close() error {
  q.log.Debug("Запрос на закрытие очереди")

  if q.IsClosed() {
    q.log.Warn(ErrNotOpen)
    return ErrNotOpen
  }

  q.ctxCanc()
  q.stateClose()

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
// Вызов из горутины изменения состояния объекта очереди
func (q *Queue) open(mgr *ibmmq.MQQueueManager) (*ibmmq.MQObject, error) {
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

  que, err := mgr.Open(mqod, flag)
  if err != nil {
    return nil, err
  }

  q.log.WithFields(map[string]interface{}{
    "mod":  "open",
    "perm": strings.Join(q.convPermToVal(), ","),
  }).Debug("Opened")

  return &que, nil
}

// Вызов из горутины изменения состояния объекта очереди
func (q *Queue) close() {
  q.unsubscInMsg()

  if q.conn == nil {
    return
  }

  err := q.conn.q.Close(0)
  if err != nil {
    q.log.WithField("mod", "close").Warn(err)
  }

  q.conn = nil
}
