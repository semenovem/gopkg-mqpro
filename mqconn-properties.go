package mqpro

import (
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "github.com/sirupsen/logrus"
)

func Properties(getMsgHandle ibmmq.MQMessageHandle, l *logrus.Entry) (map[string]interface{}, error) {
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
        l.Error("Ошибка получения свойств сообщения: ", err)
        return nil, err
      }
      break
    }
    props[name] = value
  }
  return props, nil
}

func DltMh(mh ibmmq.MQMessageHandle, l *logrus.Entry) {
  dmho := ibmmq.NewMQDMHO()
  err := mh.DltMH(dmho)
  if err != nil {
    l.Warn("Ошибка удаления объекта свойств сообщения: ", err)
  }
}
