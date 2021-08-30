package mqpro

import (
  "fmt"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "time"
)

func (c *Mqconn) Connect(delay time.Duration) {
  if delay > 0 {
    c.reconnectDelay = delay
  }
  c.mx.Lock()
  defer c.mx.Unlock()
  if c.stateConn == stateConnect {
    return
  }
  c.log.Trace("Request to connect...")
  if c.stateConn == stateDisconnect {
    c.chMgr = make(chan reqStateConn)
    go c.state()
  }
  c.reqConnect()
}

func (c *Mqconn) connect() error {
  c.log.Trace("Connecting to IBM MQ manager...")

  err := c._connectMgr()
  if err != nil {
    c.log.Warnf("Failed connection attempt to IBM MQ manager: %v", err)
    return err
  }
  c.log.Info("Connected to IBM MQ manager")

  // открытие очереди
  c.log.Trace("Opening the queue...")

  switch c.typeConn {
  case TypePut:
    err = c.openQuePut()
  case TypeGet:
    err = c.openQueGet()
  case TypeBrowse:
    err = c.openBrowse()

  default:
    c.log.Panic("Unknown connection type")
  }

  if err != nil {
    c.log.Errorf("Failed attempt to open queue: %v", err)
    return err
  }

  c.log.Info("Opened the queue")

  if c.fnInMsg != nil {
    err := c.registerEventInMsg()
    if err != nil {
      return err
    }
  }

  return nil
}

func (c *Mqconn) endpoint() string {
  return fmt.Sprintf("%s(%d)", c.cfg.Host, c.cfg.Port)
}

// Подключение к менеджеру
func (c *Mqconn) _connectMgr() error {
  if c.mgr != nil {
    return nil
  }
  cd := ibmmq.NewMQCD()
  cno := ibmmq.NewMQCNO()
  csp := ibmmq.NewMQCSP()

  cd.ChannelName = c.cfg.ChannelName
  cd.ConnectionName = c.endpoint()
  cd.MaxMsgLength = c.cfg.MaxMsgLength

  // TODO попробовать mutual authentication
  //cd.CertificateLabel

  cno.SecurityParms = csp
  cno.ClientConn = cd
  cno.Options = ibmmq.MQCNO_CLIENT_BINDING
  cno.ApplName = c.cfg.AppName

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

  mgr, err := ibmmq.Connx(c.cfg.MgrName, cno)
  if err != nil {
    return err
  }
  c.mgr = &mgr
  return nil
}

// Открытие очереди отправки
func (c *Mqconn) openQuePut() error {
  mqod := ibmmq.NewMQOD()
  mqod.ObjectType = ibmmq.MQOT_Q
  mqod.ObjectName = c.cfg.QueueName
  que, err := c.mgr.Open(mqod, ibmmq.MQOO_OUTPUT)
  if err != nil {
    return err
  }
  c.que = &que
  return nil
}

// Открыть очередь получения
func (c *Mqconn) openQueGet() error {
  mqod := ibmmq.NewMQOD()
  mqod.ObjectType = ibmmq.MQOT_Q
  mqod.ObjectName = c.cfg.QueueName
  que, err := c.mgr.Open(mqod, ibmmq.MQOO_INPUT_SHARED)
  if err != nil {
    return err
  }
  c.que = &que
  return nil
}

func (c *Mqconn) openBrowse() error {
  mqod := ibmmq.NewMQOD()
  mqod.ObjectType = ibmmq.MQOT_Q
  mqod.ObjectName = c.cfg.QueueName
  que, err := c.mgr.Open(mqod, ibmmq.MQOO_BROWSE)
  if err != nil {
    return err
  }
  c.que = &que
  return nil
}
