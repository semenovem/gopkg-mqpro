package queue

// RegisterOpen TODO добавить отмену отправки в канал данных, когда они уже не нужны
func (q *Queue) RegisterOpen() <-chan *mqConn {
  ch := make(chan *mqConn, 1)
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
    store    = make([]chan *mqConn, capacity)
    ch       chan *mqConn
    conn     *mqConn
  )

  fire := func() {
    conn := q.conn
    if conn == nil || !q.IsOpen() {
      return
    }
    for i := int32(0); i < ind; i++ {
      go func(i int32) {
        store[i] <- conn
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

    conn = q.conn
    if conn != nil && q.IsOpen() {
      ch <- conn
      close(ch)
      continue
    }

    if ind >= capacity {
      l.Warnf("Exceeding the waiting queue. Increasing the queue size +%d", inc)
      capacity += inc
      store = append(store, make([]chan *mqConn, inc)...)
    }
    store[ind] = ch
    ind++
  }
}
