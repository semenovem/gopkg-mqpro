package mqpro

import (
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

func HasConnBroken(err error) bool {
  mqrc := err.(*ibmmq.MQReturn).MQRC
  return mqrc == ibmmq.MQRC_CONNECTION_BROKEN || mqrc == ibmmq.MQRC_CONNECTION_QUIESCING
}
