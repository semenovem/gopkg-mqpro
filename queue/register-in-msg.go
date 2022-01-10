package queue

func (q *Queue) RegisterInMsg(hnd func(*Msg)) error {
  if q.hndInMsg != nil {
    return ErrRegisterEventInMsg
  }
  q.hndInMsg = hnd
  return nil
}

func (q *Queue) UnregisterInMsg() {
  q.hndInMsg = nil
}
