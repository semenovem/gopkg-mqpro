package queue

import (
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "github.com/sirupsen/logrus"
)

func properties(getMsgHandle ibmmq.MQMessageHandle) (map[string]interface{}, error) {
  impo := ibmmq.NewMQIMPO()
  pd := ibmmq.NewMQPD()
  props := make(map[string]interface{})

  impo.Options = ibmmq.MQIMPO_CONVERT_VALUE | ibmmq.MQIMPO_INQ_FIRST
  for {
    name, value, err := getMsgHandle.InqMP(impo, pd, "%")
    impo.Options = ibmmq.MQIMPO_CONVERT_VALUE | ibmmq.MQIMPO_INQ_NEXT
    if err != nil {
      mqret := err.(*ibmmq.MQReturn)
      if mqret.MQRC != ibmmq.MQRC_PROPERTY_NOT_AVAILABLE {
        return nil, err
      }
      break
    }
    props[name] = value
  }
  return props, nil
}

func setProps(h *ibmmq.MQMessageHandle, props map[string]interface{}, l *logrus.Entry) error {
  var err error
  smpo := ibmmq.NewMQSMPO()
  pd := ibmmq.NewMQPD()

  for k, v := range props {
    err = h.SetMP(smpo, k, pd, v)
    if err != nil {
      l.Errorf("Failed to set message property '%s' value '%v': %v", k, v, err)
      return err
    }
  }

  return nil
}
