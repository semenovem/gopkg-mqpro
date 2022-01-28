package manager

import (
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

func (m *Manager) RegisterConn() <-chan *ibmmq.MQQueueManager {
  if m.IsDisconn() {
    _ = m.Connect()
  }
  return m.registerConn()
}

func (m *Manager) registerConn() <-chan *ibmmq.MQQueueManager {
  ch := make(chan *ibmmq.MQQueueManager)
  m.chRegisterConn <- ch
  return ch
}

func (m *Manager) fireConn() {
  m.chRegisterConn <- nil
}

func (m *Manager) workerRegisterConn() {
  var (
    l        = m.log.WithField("fn", "workerRegisterConn")
    origCap  = int32(100)
    capacity = origCap
    inc      = origCap
    ind      = int32(0)
    store    = make([]chan *ibmmq.MQQueueManager, capacity)
    ch       chan *ibmmq.MQQueueManager
    mgr      *ibmmq.MQQueueManager
  )

  fire := func() {
    mgr := m.mgr
    if mgr == nil || !m.IsConn() {
      return
    }
    for i := int32(0); i < ind; i++ {
      go func(i int32) {
        store[i] <- mgr
        close(store[i])
      }(i)
    }
    ind = 0

    if capacity != origCap {
      capacity = origCap
      store = store[:capacity]
    }
  }

  for ch = range m.chRegisterConn {
    if ch == nil {
      fire()
      continue
    }

    mgr = m.mgr
    if mgr != nil && m.IsConn() {
      ch <- mgr
      close(ch)
      continue
    }

    if ind >= capacity {
      l.Warnf("Exceeding the waiting queue. Increasing the queue size +%d", inc)
      capacity += inc
      store = append(store, make([]chan *ibmmq.MQQueueManager, inc)...)
    }
    store[ind] = ch
    ind++
  }
}
