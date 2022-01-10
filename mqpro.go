package mqpro

import (
  "context"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "github.com/semenovem/gopkg_mqpro/v2/queue"
  "github.com/sirupsen/logrus"
  "sync"
  "time"
)

type Mqpro struct {
  rootCtx        context.Context
  ctx            context.Context
  ctxCanc        context.CancelFunc
  mx             sync.Mutex    // Подключение / отключение
  disconnDelay   time.Duration // Задержка перед разрывом соединения
  reconnDelay    time.Duration // Задержка при повторных попытках подключения к MQ
  log            *logrus.Entry
  coreSet        *queue.CoreSet
  mgr            *ibmmq.MQQueueManager           // Менеджер очереди
  state          state                           // Состояние подключения к менеджеру IBMMQ
  chState        chan state                      // Канал изменения состояния
  chRegisterConn chan chan *ibmmq.MQQueueManager // Ожидание подключения к ibmmq
  queues         []*queue.Queue                  // Очереди

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
    chState:        make(chan state, 100),
    chRegisterConn: make(chan chan *ibmmq.MQQueueManager, 100),
    coreSet:        &queue.CoreSet{},
  }

  go o.workerState()
  go o.workerRegisterConn()

  return o
}

// Queue Объект очереди
func (m *Mqpro) Queue(a string) *queue.Queue {
  l := m.log.WithField("que", a)

  q := queue.New(m.rootCtx, l, m)
  m.queues = append(m.queues, q)

  return q
}
