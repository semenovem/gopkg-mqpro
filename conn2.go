package mqpro

import (
  "context"
  "fmt"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "github.com/sirupsen/logrus"
)

func asdfasfd()  {
  _, _, _ = mqPay.Get(context.Background())
}

var (
  mq = New(context.Background(), &logrus.Entry{})
  mqPay = mq.ConnByAlias("pay")
)


func(p *Mqpro) ConnByAlias(alias string) *Mqconn {


  return &Mqconn{}
}




func (p *Mqpro) Connect2() error {
  p.mx.Lock()
  defer p.mx.Unlock()

  ctx, cancel := context.WithCancel(p.rootCtx)
  p.ctx = ctx
  p.ctxCancel = cancel

  p.log.Info("Request to establish connection to IBM MQ...")

  if p.cfg == nil {
    return ErrNoConfig
  }
  err := p.connect2()
  if err != nil {
    p.log.Error(err)
    return err
  }

  return nil
}

func (p *Mqpro) connect2() error {
  cd := ibmmq.NewMQCD()
  cno := ibmmq.NewMQCNO()
  csp := ibmmq.NewMQCSP()

  cd.ChannelName = p.cfg.Channel
  cd.ConnectionName = p.endpoint()
  cd.MaxMsgLength = p.cfg.MaxMsgLength

  // TODO попробовать mutual authentication
  //cd.CertificateLabel

  cno.SecurityParms = csp
  cno.ClientConn = cd
  cno.Options = ibmmq.MQCNO_CLIENT_BINDING
  cno.ApplName = p.cfg.App

  if p.cfg.Tls {
    sco := ibmmq.NewMQSCO()
    sco.KeyRepository = p.cfg.KeyRepository

    cno.SSLConfig = sco

    cd.SSLCipherSpec = "ANY_TLS12"
    cd.SSLClientAuth = ibmmq.MQSCA_OPTIONAL
  }

  if p.cfg.User == "" {
    csp.AuthenticationType = ibmmq.MQCSP_AUTH_NONE
  } else {
    csp.AuthenticationType = ibmmq.MQCSP_AUTH_USER_ID_AND_PWD
    csp.UserId = p.cfg.User
    csp.Password = p.cfg.Pass
  }

  mgr, err := ibmmq.Connx(p.cfg.Manager, cno)
  if err != nil {
    return err
  }
  p.mgr2 = &mgr
  return nil
}

func (p *Mqpro) endpoint() string {
  return fmt.Sprintf("%s(%d)", p.cfg.Host, p.cfg.Port)
}

