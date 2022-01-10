package queue

import (
  "context"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "github.com/sirupsen/logrus"
  "sync"
  "time"
)

type Queue struct {
  rootCtx   context.Context
  ctx       context.Context
  ctxCanc   context.CancelFunc
  log       *logrus.Entry
  manager   manager
  queueName string
  perm      []permQueue     // Разрешения очереди
  que       *ibmmq.MQObject // Объект открытой очереди

  // Deprecated
  mgr *ibmmq.MQQueueManager // Менеджер ibmmq

  state           state         // Состояние
  chState         chan state    // Канал изменения состояния
  ctlo            *ibmmq.MQCTLO // Объект подписки ibmmq
  fnInMsg         func(*Msg)    // Подписка на входящие сообщения
  reconnectDelay  time.Duration // Таймаут попыток повторного подключения
  msgWaitInterval time.Duration // Ожидание сообщения
  delayClose      time.Duration // Ожидание закрытия
  devMode         bool          // Режим разработки расширенное логирование
  h               Header        // Тип заголовков
  rfh2            *rfh2Cfg      // Данные для заголовков RFH2
  rfh2RootTag     string        // Корневой тег тела сообщения

  mx    sync.Mutex
  mxMsg sync.Mutex

  hndInMsg       func(*Msg)
  chRegisterConn chan chan *ibmmq.MQObject
}

func New(ctx context.Context, l *logrus.Entry, m manager) *Queue {
  q := &Queue{
    rootCtx:         ctx,
    reconnectDelay:  defReconnectDelay,
    state:           stateDisconn,
    msgWaitInterval: defMsgWaitInterval,
    log:             l,
    h:               DefHeader,
    chState:         make(chan state),
    chRegisterConn:  make(chan chan *ibmmq.MQObject, 100),
    manager:         m,
    delayClose:      defDelayClose,
  }

  go q.workerState()
  go q.workerRegisterConn()

  go func() {
    <-ctx.Done()
    close(q.chState)
    //  TODO освободить ресурсы
  }()

  return q
}
