package queue

import "github.com/ibm-messaging/mq-golang/v5/ibmmq"

func IsConnBroken(err error) bool {
  mqrc := err.(*ibmmq.MQReturn).MQRC
  return mqrc == ibmmq.MQRC_CONNECTION_BROKEN || mqrc == ibmmq.MQRC_CONNECTION_QUIESCING
}

func unionProps(dst map[string]interface{}, src map[string]interface{}) {
  for n, v := range src {
    dst[n] = v
  }
}

func unionPropsDeep(dst map[string]interface{}, src []map[string]interface{}) {
  for _, a := range src {
    for n, v := range a {
      dst[n] = v
    }
  }
}

func tailFour(n int) int {
  r := n % 4
  if r == 0 {
    return 0
  }
  return 4 - r
}

func tailFour32(n int32) int32 {
  r := n % 4
  if r == 0 {
    return 0
  }
  return 4 - r
}
