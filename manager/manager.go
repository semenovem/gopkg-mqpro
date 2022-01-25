package manager

import (
  "context"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "github.com/sirupsen/logrus"
  "sync"
  "time"
)

type Mqpro struct {
  log            *logrus.Entry
  rootCtx        context.Context
  ctx            context.Context
  ctxCanc        context.CancelFunc
  mx             sync.Mutex                      // Подключение / отключение
  disconnDelay   time.Duration                   // Задержка перед разрывом соединения
  reconnDelay    time.Duration                   // Таймаут попыток подключения к MQ
  mgr            *ibmmq.MQQueueManager           // Менеджер очереди
  state          state                           // Состояние подключения к менеджеру IBMMQ
  chState        chan state                      // Канал изменения состояния
  chRegisterConn chan chan *ibmmq.MQQueueManager // Ожидание подключения к ibmmq

  host                              string
  port                              int32
  manager, channel, app, user, pass string
  tls                               bool
  keyRepository                     string
  maxMsgLen                         int32
}

func New(rootCtx context.Context, l *logrus.Entry) *Mqpro {
  o := &Mqpro{
    rootCtx:        rootCtx,
    disconnDelay:   defDisconnDelay,
    reconnDelay:    defReconnDelay,
    log:            l,
    state:          stateDisconn,
    chState:        make(chan state, 10),
    chRegisterConn: make(chan chan *ibmmq.MQQueueManager, 10),
  }

  go o.workerState()
  go o.workerRegisterConn()

  return o
}

func (m *Mqpro) Reconnect() {
  m.stateErr()
}
