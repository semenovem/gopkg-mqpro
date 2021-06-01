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
  c.log.Info("Opened the queue")

  if err != nil {
    c.log.Errorf("Failed attempt to open queue: %v", err)
    return err
  }

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

// подключение к менеджеру
func (c *Mqconn) _connectMgr() error {
  if c.mgr != nil {
    return nil
  }
  cd := ibmmq.NewMQCD()
  cno := ibmmq.NewMQCNO()
  csp := ibmmq.NewMQCSP()

  cd.ChannelName = c.cfg.ChannelName
  cd.ConnectionName = c.endpoint()

  cno.SecurityParms = csp
  cno.ClientConn = cd
  cno.Options = ibmmq.MQCNO_CLIENT_BINDING
  cno.ApplName = c.cfg.AppName

  // TODO - в настройки
  cd.MaxMsgLength = 104857600

  // -------------------------
  //sco := ibmmq.NewMQSCO()
  //
  //cno.SSLConfig = sco
  //
  //// TLS
  //// The CipherSpec must match what is configured on the corresponding SVRCONN
  //cd.SSLCipherSpec = "TLS_RSA_WITH_AES_128_CBC_SHA256"
  //
  //// The ClientAuth field says whether or not the client needs to present its own certificate
  //// This too must match the SVRCONN definition.
  //cd.SSLClientAuth = ibmmq.MQSCA_OPTIONAL
  //
  //// The keystore contains at least the certificate to verify the qmgr's cert (usually from
  //// a Certificate Authority) and optionally the client's own certificate.
  //// We could also optionally specify which certificate represents the client, based on its label
  //// but don't need to do this when using the MQSCA_OPTIONAL flag.
  //sco.KeyRepository = "./crypto/client-prv-key.pem"
  ////sco.CryptoHardware
  // -------------------------

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

// открытие очереди отправки
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

// открыть очередь получения
func (c *Mqconn) openQueGet() error {
  mqod := ibmmq.NewMQOD()
  mqod.ObjectType = ibmmq.MQOT_Q
  mqod.ObjectName = c.cfg.QueueName
  que, err := c.mgr.Open(mqod, ibmmq.MQOO_INPUT_EXCLUSIVE)
  if err != nil {
    return err
  }
  c.que = &que
  return nil
}

func (c *Mqconn) openBrowse() error {
  mqod := ibmmq.NewMQOD()
  openOptions := ibmmq.MQOO_BROWSE
  mqod.ObjectType = ibmmq.MQOT_Q
  mqod.ObjectName = c.cfg.QueueName
  que, err := c.mgr.Open(mqod, openOptions)
  if err != nil {
    return err
  }
  c.que = &que
  return nil
}
