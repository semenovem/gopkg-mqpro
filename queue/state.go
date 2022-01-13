package queue

import (
  "fmt"
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
        q.errorHandler(err)

        select {
        case <-q.ctx.Done():
          continue worker
        case <-time.After(q.reconnectDelay):
        }
      }

    case stateClosed:
      q.state = stateClosed
      q.close()

    case stateErr:
      q.state = stateErr
      q.close()
      q.manager.Reconnect()
      go q.stateOpen()
    }
  }
}

func (q *Queue) stateOpen() {
  q.chState <- stateOpen
}

func (q *Queue) stateError() {
  q.chState <- stateErr
}

func (q *Queue) stateClose() {
  q.chState <- stateClosed
}

func (q *Queue) IsOpen() bool {
  return q.state == stateOpen
}

func (q *Queue) IsClosed() bool {
  return q.state == stateClosed
}

func (q *Queue) errorHandler(err error) {
  if err == nil {
    return
  }

  isNeedRestart := true

  switch p := err.(type) {
  case *ibmmq.MQReturn:
    switch p.MQRC {
    case ibmmq.MQRC_CONNECTION_BROKEN:
    case ibmmq.MQRC_CALL_IN_PROGRESS:
      isNeedRestart = false
    }
  case error:
    switch p {
    case ErrConnBroken:
    case ErrBusySubsc:
      isNeedRestart = false
    }
  }

  fmt.Println(">>>>>>>>>>>>>> debug errorHandler:: ", err)

  if isNeedRestart && q.state == stateOpen {
    q.stateError()
  }
}
