package queue

import (
  "context"
  "fmt"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "github.com/sirupsen/logrus"
  "sync"
  "time"
)

type Conn struct {
  cfg             *Cfg
  rootCtx         context.Context
  ctx             context.Context
  ctxCns          context.CancelFunc
  log             *logrus.Entry
  que             *ibmmq.MQObject          // Объект открытой очереди
  mgr             *ibmmq.MQQueueManager    // Менеджер ibmmq
  state           state                    // Состояние
  chState         chan state               // Канал изменения состояния
  ctlo            *ibmmq.MQCTLO            // Объект подписки ibmmq
  fnInMsg         func(*Msg)               // Подписка на входящие сообщения
  fnsConn         map[uint32]chan struct{} // Подписки на установку соединения
  fnsDisconn      map[uint32]chan struct{} // Подписки на закрытие соединения
  reconnectDelay  time.Duration            // Таймаут попыток повторного подключения
  msgWaitInterval time.Duration            // Ожидание сообщения
  devMode         bool                     // Режим разработки расширенное логирование
  h               header                   // тип заголовков
  rfh2            *rfh2Cfg                 // Данные для заголовков RFH2
  rfh2RootTag     string                   // Корневой тег тела сообщения

  mx    sync.Mutex
  mxMsg sync.Mutex

  subConnMgr chan<- chan *ibmmq.MQQueueManager
}

func New(ctx context.Context, l *logrus.Entry) *Conn {
  o := &Conn{
    rootCtx:         ctx,
    fnsConn:         map[uint32]chan struct{}{},
    fnsDisconn:      map[uint32]chan struct{}{},
    reconnectDelay:  defReconnectDelay,
    state:           stateDisconn,
    msgWaitInterval: defMsgWaitInterval,
    log:             l,
    h:               defHeader,
    chState:         make(chan state),
  }

  go o.workerState()

  go func() {
    <-ctx.Done()
    close(o.chState)
    //  TODO освободить ресурсы
  }()

  return o
}

func (c *Conn) endpoint() string {
  return fmt.Sprintf("%s(%d)", c.cfg.Host, c.cfg.Port)
}

func (c *Conn) Connect() error {
  c.mx.Lock()
  defer c.mx.Unlock()

  c.log.Trace("Request to connect...")

  if c.cfg == nil {
    return ErrNoConfig
  }
  if c.state != stateDisconn {
    return ErrHasConnection
  }
  c.ctx, c.ctxCns = context.WithCancel(c.rootCtx)
  c.reqConn()

  select {
  case <-c.ctx.Done():
  case <-c.regConnMgr():
  }

  return nil
}

func (c *Conn) Disconnect() {
  c.reqDisconn()
}
