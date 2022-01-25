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
