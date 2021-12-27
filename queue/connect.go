package queue

import (
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

// Подключение к менеджеру
func (c *Conn) connToMgr() (*ibmmq.MQQueueManager, error) {
  cd := ibmmq.NewMQCD()
  cno := ibmmq.NewMQCNO()
  csp := ibmmq.NewMQCSP()

  cd.ChannelName = c.cfg.Channel
  cd.ConnectionName = c.endpoint()
  cd.MaxMsgLength = c.cfg.MaxMsgLength

  // TODO попробовать mutual authentication
  //cd.CertificateLabel

  cno.SecurityParms = csp
  cno.ClientConn = cd
  cno.Options = ibmmq.MQCNO_CLIENT_BINDING
  cno.ApplName = c.cfg.App

  if c.cfg.Tls {
    sco := ibmmq.NewMQSCO()
    sco.KeyRepository = c.cfg.KeyRepository

    cno.SSLConfig = sco

    cd.SSLCipherSpec = "ANY_TLS12"
    cd.SSLClientAuth = ibmmq.MQSCA_OPTIONAL
  }

  if c.cfg.User == "" {
    csp.AuthenticationType = ibmmq.MQCSP_AUTH_NONE
  } else {
    csp.AuthenticationType = ibmmq.MQCSP_AUTH_USER_ID_AND_PWD
    csp.UserId = c.cfg.User
    csp.Password = c.cfg.Pass
  }

  mgr, err := ibmmq.Connx(c.cfg.Manager, cno)
  if err != nil {
    return nil, err
  }
  return &mgr, nil
}




//// Открытие очереди отправки
//func (c *Conn) openQuePut() error {
//  mqod := ibmmq.NewMQOD()
//  mqod.ObjectType = ibmmq.MQOT_Q
//  mqod.ObjectName = c.cfg.QueueName
//  que, err := c.mgr.Open(mqod, ibmmq.MQOO_OUTPUT)
//  if err != nil {
//    return err
//  }
//  c.que = &que
//  return nil
//}
//
//// Открыть очередь получения
//func (c *Conn) openQueGet() error {
//  mqod := ibmmq.NewMQOD()
//  mqod.ObjectType = ibmmq.MQOT_Q
//  mqod.ObjectName = c.cfg.QueueName
//  que, err := c.mgr.Open(mqod, ibmmq.MQOO_INPUT_SHARED)
//  if err != nil {
//    return err
//  }
//  c.que = &que
//  return nil
//}
//
//func (c *Conn) openBrowse() error {
//  mqod := ibmmq.NewMQOD()
//  mqod.ObjectType = ibmmq.MQOT_Q
//  mqod.ObjectName = c.cfg.QueueName
//  que, err := c.mgr.Open(mqod, ibmmq.MQOO_BROWSE)
//  if err != nil {
//    return err
//  }
//  c.que = &que
//  return nil
//}
