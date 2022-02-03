package queue

import (
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "time"
)

func (q *Queue) workerState() {
  var (
    l   = q.log.WithField("mod", "workerState")
    err error
    st  state
    mgr *ibmmq.MQQueueManager
    que *ibmmq.MQObject
    i   int8
  )

worker:
  for st = range q.chState {
    //l.Debug(stateMapByKey[q.state], " >>> ", stateMapByKey[st])

    if q.state == st {
      continue
    }

    if st == stateErr {
      st = stateOpen

      q.close()
      q.manager.Reconnect()
    }

    switch st {
    case stateOpen:
      if q.state == stateOpen || q.state == stateConnecting {
        continue
      }
      q.state = stateConnecting

      for {
        select {
        case mgr = <-q.manager.RegisterConn():
        case <-q.ctx.Done():
          continue worker
        }

        if que, err = q.open(mgr); err == nil {
          q.conn = &mqConn{q: que, m: mgr}

          for i = 0; i < 3; i++ {
            if err = q.subscInMsg(q.conn); err == nil {
              q.state = stateOpen
              q.fireConn()
              continue worker
            }
          }
        }

        l.WithField("oper", "open").Warn(err)

        q.close()
        q.manager.Reconnect()

        select {
        case <-q.ctx.Done():
          continue worker
        case <-time.After(q.reconnectDelay):
        }
      }

    case stateClosed:
      q.state = stateClosed
      q.close()
    }
  }
}

func (q *Queue) stateOpen() {
  q.chState <- stateOpen
}

func (q *Queue) stateClose() {
  q.state = stateTransitional
  q.chState <- stateClosed
}

func (q *Queue) IsOpen() bool {
  return q.state == stateOpen
}

func (q *Queue) IsClosed() bool {
  return q.state == stateClosed
}

// bool - планируется ли рестарт из-за представленной проблемы
func (q *Queue) errorHandler(err error) bool {
  if err == nil {
    return false
  }

  isNeedRestart := true

  switch p := err.(type) {
  case *ibmmq.MQReturn:
    switch p.MQRC {
    case ibmmq.MQRC_CALL_IN_PROGRESS,
      ibmmq.MQRC_NOT_OPEN_FOR_INPUT,
      ibmmq.MQRC_NOT_OPEN_FOR_OUTPUT,
      ibmmq.MQRC_NOT_OPEN_FOR_BROWSE,
      ibmmq.MQRC_SECURITY_ERROR:
      isNeedRestart = false
    }
  case error:
    switch p {
    case ErrBusySubsc:
      isNeedRestart = false
    }
  }

  if isNeedRestart && q.state == stateOpen {
    q.state = stateTransitional
    q.chState <- stateErr
    return true
  }

  return false
}
