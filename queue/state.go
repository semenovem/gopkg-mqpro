package queue

import (
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "time"
)

func (q *Queue) workerState() {
  var (
    l     = q.log.WithField("mod", "workerState")
    err   error
    st    state
    mgrMq *ibmmq.MQQueueManager
  )

worker:
  for st = range q.chState {
    l.Debug(stateKey[q.state], " >>> ", stateKey[st])

    if q.state == st {
      continue
    }

    switch st {
    case stateConn:
      if q.state == stateConn || q.state == stateConnecting {
        continue
      }
      q.state = stateConnecting

      for {
        select {
        case mgrMq = <-q.manager.RegisterConn():
        case <-q.ctx.Done():
          continue worker
        }

        err = q.open(mgrMq)
        if err == nil {
          q.state = stateConn
          q.fireConn()
          continue worker
        }

        l.WithField("oper", "open").Warn(err)

        select {
        case <-q.ctx.Done():
          continue worker
        case <-time.After(q.reconnectDelay):
        }
      }

    case stateDisconn:
      q.state = stateDisconn
      q.close()

    case stateErr:
      q.state = stateErr
      q.close()
      go q.stateConn()
    }
  }
}

func (q *Queue) stateConn() {
  q.chState <- stateConn
}

func (q *Queue) stateError() {
  if q.state == stateConn {
    q.chState <- stateErr
  }
}

func (q *Queue) stateDisconn() {
  q.chState <- stateDisconn
}

func (q *Queue) IsConn() bool {
  return q.state == stateConn
}

func (q *Queue) IsDisconn() bool {
  return q.state == stateDisconn
}
