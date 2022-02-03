package queue

import "time"

// Msg отправляемое / получаемое сообщение
type Msg struct {
  MsgId    []byte
  CorrelId []byte
  Payload  []byte
  Props    map[string]interface{}
  Time     time.Time
  MQRFH2   []*MQRFH2
}

var msgInitTime time.Time

func (m *Msg) Erase() {
  m.MsgId = nil
  m.CorrelId = nil
  m.Props = nil
  m.Payload = nil
  m.MQRFH2 = nil
  m.Time = msgInitTime
}
