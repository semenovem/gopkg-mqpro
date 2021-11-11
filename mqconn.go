package mqpro

import (
  "fmt"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "github.com/sirupsen/logrus"
  "sync"
  "time"
)

type Mqconn struct {
  cfg             *Cfg
  log             *logrus.Entry
  typeConn        TypeConn              // Тип подключения
  mgr             *ibmmq.MQQueueManager // Менеджер очереди
  que             *ibmmq.MQObject       // Объект открытой очереди
  h               header                // тип заголовков
  mx              sync.Mutex
  stateConn       stateConn
  chMgr           chan reqStateConn
  fnInMsg         func(*Msg)               // Подписка на входящие сообщения
  ctlo            *ibmmq.MQCTLO            // Объект подписки ibmmq
  fnsConn         map[uint32]chan struct{} // Подписки на установку соединения
  fnsDisconn      map[uint32]chan struct{} // Подписки на закрытие соединения
  ind             uint32                   // Счетчик
  reconnectDelay  time.Duration            // Таймаут попыток повторного подключения
  msgWaitInterval time.Duration            // Ожидание сообщения
  rfh2            *rfh2Cfg                 // Данные для заголовков RFH2
  DevMode         bool                     // Режим разработки расширенное логирование
  rfh2RootTag     string                   // В какой тег оборачивать заголовки

  // менеджер imbmq одновременно может отправлять/принимать одно сообщение
  // TODO - использовать только один мьютекс
  mxOper sync.Mutex

  mxPut    sync.Mutex
  mxGet    sync.Mutex
  mxBrowse sync.Mutex
}

// Cfg Данные подключения
type Cfg struct {
  Host             string
  Port             int
  MgrName          string
  ChannelName      string
  QueueName        string // название очереди
  Header           string // Тип заголовков [prop | rfh2]
  AppName          string
  User             string
  Pass             string
  Priority         string
  MaxMsgLength     int32
  Tls              bool
  KeyRepository    string
  CertificateLabel string
  DevMode          bool
  RootTag          string
}

func NewMqconn(tc TypeConn, l *logrus.Entry, c *Cfg) *Mqconn {
  o := &Mqconn{
    cfg:             c,
    fnsConn:         map[uint32]chan struct{}{},
    fnsDisconn:      map[uint32]chan struct{}{},
    reconnectDelay:  defReconnectDelay,
    stateConn:       stateDisconnect,
    msgWaitInterval: defMsgWaitInterval,
    DevMode:         c.DevMode,
    rfh2RootTag:     c.RootTag,
  }

  m := map[string]interface{}{
    "conn": fmt.Sprintf("%s|%s|%s|%s",
      o.endpoint(), c.MgrName, c.QueueName, typeConnTxt[tc]),
  }

  o.log = l.WithFields(m)

  if tc != TypePut && tc != TypeGet && tc != TypeBrowse {
    o.log.Panic("Unknown connection type")
  }

  o.typeConn = tc

  if c.MaxMsgLength == 0 {
    c.MaxMsgLength = defMaxMsgLength
  }

  if c.Header != "" {
    h, err := parseHeaderType(c.Header)
    if err != nil {
      o.log.Panic(errHeaderParseType)
    }
    o.h = h
    o.rfh2 = newRfh2Cfg()
  }

  return o
}
