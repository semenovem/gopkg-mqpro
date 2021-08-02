package mqpro

import (
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "github.com/sirupsen/logrus"
  "sync"
  "time"
)

type Mqconn struct {
  cfg             *Cfg
  log             *logrus.Entry
  typeConn        TypeConn              // тип подключения
  mgr             *ibmmq.MQQueueManager // Менеджер очереди
  que             *ibmmq.MQObject       // Объект открытой очереди
  mx              sync.Mutex
  stateConn       stateConn
  chMgr           chan reqStateConn
  fnInMsg         func(*Msg)               // подписка на входящие сообщения
  ctlo            *ibmmq.MQCTLO            // объект подписки ibmmq
  fnsConn         map[uint32]chan struct{} // подписки на установку соединения
  fnsDisconn      map[uint32]chan struct{} // подписки на закрытие соединения
  ind             uint32                   // простой атомарный счетчик
  reconnectDelay  time.Duration            // таймаут попыток повторного подключения
  msgWaitInterval time.Duration            // Ожидание сообщения

  // менеджер imbmq одновременно может отправлять/принимать одно сообщение
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
  QueueName        string
  AppName          string
  User             string
  Pass             string
  Priority         string
  MaxMsgLength     int32
  Tls              bool
  KeyRepository    string
  CertificateLabel string
}

type TypeConn int
type stateConn int
type reqStateConn int
type queueOper int

const (
  TypePut TypeConn = iota + 1
  TypeGet
  TypeBrowse
  defReconnectDelay  = time.Second * 3
  defMaxMsgLength    = 1024 * 1024 * 100
  defMsgWaitInterval = time.Millisecond * 100
)

const (
  stateDisconnect stateConn = iota
  stateConnect
  stateErr
)

const (
  reqConnect reqStateConn = iota
  reqReconnect
  reqDisconnect
)

const (
  operGet queueOper = iota
  operGetByMsgId
  operGetByCorrelId
  operBrowseFirst
  operBrowseNext
)

var typeConnTxt = map[TypeConn]string{
  TypePut:    "put",
  TypeGet:    "get",
  TypeBrowse: "browse",
}

// Msg отправляемое / получаемое сообщение
type Msg struct {
  MsgId    []byte
  CorrelId []byte
  Payload  []byte
  Props map[string]interface{}
  Time  time.Time
}
