package queue

import (
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

func (q *Queue) RegisterOpen() <-chan *ibmmq.MQObject {
  ch := make(chan *ibmmq.MQObject)
  q.chRegisterOpen <- ch
  return ch
}

func (q *Queue) fireConn() {
  q.chRegisterOpen <- nil
}

func (q *Queue) workerRegisterOpen() {
  var (
    l        = q.log.WithField("fn", "workerRegisterOpen")
    origCap  = int32(100)
    capacity = origCap
    inc      = origCap
    ind      = int32(0)
    store    = make([]chan *ibmmq.MQObject, capacity)
    ch       chan *ibmmq.MQObject
    que      *ibmmq.MQObject
  )

  fire := func() {
    que := q.que
    if que == nil || !q.IsConn() {
      return
    }
    for i := int32(0); i < ind; i++ {
      go func(i int32) {
        store[i] <- que
        close(store[i])
      }(i)
    }
    ind = 0

    if capacity != origCap {
      capacity = origCap
      store = store[:capacity]
    }
  }

  for ch = range q.chRegisterOpen {
    if ch == nil {
      fire()
      continue
    }

    que = q.que
    if que != nil && q.IsConn() {
      ch <- que
      close(ch)
      continue
    }

    if ind >= capacity {
      l.Warnf("Exceeding the waiting queue. Increasing the queue size +%d", inc)
      capacity += inc
      store = append(store, make([]chan *ibmmq.MQObject, inc)...)
    }
    store[ind] = ch
    ind++
  }
}
