package queue

import (
  "context"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "github.com/sirupsen/logrus"
  "sync"
  "time"
)

type Queue struct {
  rootCtx        context.Context
  ctx            context.Context
  ctxCanc        context.CancelFunc
  log            *logrus.Entry
  base           provider
  manager        manager
  queueName      string
  perm           []permQueue   // Разрешения очереди
  conn           *mqConn       // Данные подключений
  state          state         // Состояние
  chState        chan state    // Канал изменения состояния
  ctlo           *ibmmq.MQCTLO // Объект подписки ibmmq
  reconnectDelay time.Duration // Таймаут попыток повторного подключения
  delayClose     time.Duration // Ожидание закрытия
  devMode        bool          // Режим разработки расширенное логирование
  h              Header        // Тип заголовков
  rfh2           *rfh2Cfg      // Данные для заголовков RFH2
  rfh2RootTag    string        // Корневой тег тела сообщения

  mx      sync.Mutex
  mxMsg   sync.Mutex
  mxSubsc sync.Mutex

  hndInMsg       func(*Msg)        // Обработчик сообщений
  chRegisterOpen chan chan *mqConn // Канал с подписками на открытие очереди
  alias          string
}

// Данные подключений
type mqConn struct {
  q *ibmmq.MQObject
  m *ibmmq.MQQueueManager
}

func New(
  ctx context.Context,
  l *logrus.Entry,
  m manager,
  base provider,
  alias string) *Queue {
  q := &Queue{
    alias:          alias,
    log:            l.WithField("alias", alias),
    rootCtx:        ctx,
    base:           base,
    manager:        m,
    state:          stateClosed,
    h:              DefHeader,
    chState:        make(chan state),
    chRegisterOpen: make(chan chan *mqConn, 100),
    reconnectDelay: defReconnectDelay,
    delayClose:     defDelayClose,
  }

  go q.workerState()
  go q.workerRegisterOpen()

  go func() {
    <-ctx.Done()
    q.close()
  }()

  return q
}

func (q *Queue) convPermToVal() []string {
  a := make([]string, len(q.perm))
  for i, v := range q.perm {
    a[i] = permMapByKey[v]
  }
  return a
}

func (q *Queue) HasPermQueue(p permQueue) bool {
  for _, v := range q.perm {
    if v == p {
      return true
    }
  }
  return false
}

func (q *Queue) isConnErr(err error) bool {
  if err == nil {
    return false
  }
  mqret := err.(*ibmmq.MQReturn)
  return mqret != nil && mqret.MQRC != ibmmq.MQRC_CONNECTION_BROKEN
}

func (q *Queue) isWarn(err error) {
  if err != nil {
    q.log.Warn(err)
  }
}
