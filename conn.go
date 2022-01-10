package mqpro

import (
  "context"
  "fmt"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "github.com/semenovem/gopkg_mqpro/v2/queue"
  "sync"
  "time"
)

func (m *Mqpro) Connect() error {
  m.log.Debug("Request to establish connection to IBM MQ...")

  m.mx.Lock()
  defer m.mx.Unlock()

  if !m.isConfigured() {
    m.log.Warn(ErrNoConfig)
    return ErrNoConfig
  }

  if !m.IsDisconn() {
    m.log.Warn(ErrAlreadyConnected)
    return ErrAlreadyConnected
  }

  m.ctx, m.ctxCanc = context.WithCancel(m.rootCtx)

  m.stateConn()

  select {
  case <-m.ctx.Done():
    return nil
  case <-m.RegisterConn():
  }
  m.log.Info("Установлено подключение к менеджеру и очередям ibmmq")

  // Открытие очередей
  select {
  case <-m.ctx.Done():
  case err := <-m.openQues():
    if err != nil {
      return err
    }
  }

  return nil
}

func (m *Mqpro) openQues() <-chan error {
  var (
  	ch = make(chan error)
    wg = sync.WaitGroup{}
    err1 error
  )

  wg.Add(len(m.queues))

  // Запуск открытия очередей
  for _, q := range m.queues {
    go func(q *queue.Queue) {
      defer wg.Done()

      err := q.Open()
      if err != nil {
        err1 = err
      }
    }(q)
  }

  go func() {
    defer close(ch)
    wg.Wait()
    if err1 != nil {
      ch <- err1
      return
    }

    // Ожидание открытия очередей

  }()

  return ch
}

func (m *Mqpro) connect() error {
  cd := ibmmq.NewMQCD()
  cno := ibmmq.NewMQCNO()
  csp := ibmmq.NewMQCSP()

  cd.ChannelName = m.channel
  cd.ConnectionName = m.endpoint()
  cd.MaxMsgLength = m.maxMsgLen

  // TODO попробовать mutual authentication
  //cd.CertificateLabel

  cno.SecurityParms = csp
  cno.ClientConn = cd
  cno.Options = ibmmq.MQCNO_CLIENT_BINDING
  cno.ApplName = m.app

  if m.tls {
    sco := ibmmq.NewMQSCO()
    sco.KeyRepository = m.keyRepository

    cno.SSLConfig = sco

    cd.SSLCipherSpec = "ANY_TLS12"
    cd.SSLClientAuth = ibmmq.MQSCA_OPTIONAL
  }

  if m.user == "" {
    csp.AuthenticationType = ibmmq.MQCSP_AUTH_NONE
  } else {
    csp.AuthenticationType = ibmmq.MQCSP_AUTH_USER_ID_AND_PWD
    csp.UserId = m.user
    csp.Password = m.pass
  }

  mgr, err := ibmmq.Connx(m.manager, cno)
  if err != nil {
    return err
  }
  m.mgr = &mgr

  m.log.WithFields(map[string]interface{}{
    "endpoint": cd.ConnectionName,
    "manager":  m.manager,
  }).Info("Connected to ibmmq manager")

  return nil
}

func (m *Mqpro) endpoint() string {
  return fmt.Sprintf("%s(%d)", m.host, m.port)
}

func (m *Mqpro) Disconnect() error {
  m.log.Debug("Request to disconnect from IBM MQ...")

  if m.IsDisconn() {
    m.log.Warn(ErrNoEstablishedConnection)
    return ErrNoEstablishedConnection
  }

  m.ctxCanc()
  m.stateDisconn()

  m.mx.Lock()
  defer m.mx.Unlock()

  select {
  case <-m.rootCtx.Done():
  case <-time.After(m.disconnDelay):
  }

  m.log.Info("Connection dropped")

  return nil
}

func (m *Mqpro) disconnect() {
  mgr := m.mgr
  if mgr != nil {
    m.mgr = nil
    err := mgr.Disc()
    if err != nil {
      m.log.WithField("mod", "disconnect").Warn(err)
    }
  }
}
