package mqm

import (
  "context"
  "errors"
  "github.com/semenovem/mqm/v2/manager"
  "github.com/semenovem/mqm/v2/queue"
  "github.com/sirupsen/logrus"
  "sync"
  "time"
)

type Mqm struct {
  rootCtx      context.Context
  ctx          context.Context
  ctxEsc       context.CancelFunc
  mx           sync.Mutex // Подключение / отключение
  log          *logrus.Entry
  disconnDelay time.Duration
  isConnected  bool
  queueCfg     *queue.BaseConfig // Конфиг ibmmq очереди
  managerCfg   *manager.Config   // Конфиг ibmmq менеджера
  queues       []Queue
  managers     []*manager.Manager
  pipes        []*Pipe
}

func New(rootCtx context.Context, l *logrus.Entry) *Mqm {
  o := &Mqm{
    rootCtx:      rootCtx,
    log:          l,
    disconnDelay: defDisconnDelay,
  }
  return o
}

// NewQueue Объект очереди
func (m *Mqm) NewQueue(a string) Queue {
  var (
    l    = m.log.WithField("_t", "queue")
    lMag = m.log.WithField("_t", "manager")
    man  = manager.New(m.rootCtx, lMag)
    q    = queue.New(m.rootCtx, l, man, m, a)
  )
  m.queues = append(m.queues, q)
  m.managers = append(m.managers, man)
  return q
}

func (m *Mqm) GetQueueByAlias(a string) Queue {
  for _, q := range m.queues {
    if q.Alias() == a {
      return q
    }
  }
  return nil
}

func (m *Mqm) GetBaseCfg() *queue.BaseConfig {
  return m.queueCfg
}

func (m *Mqm) Ready() bool {
  return m.isConnected
}

const (
  defDisconnDelay = time.Millisecond * 100 // Задержка перед разрывом соединения
)

var (
  ErrNoConnection     = errors.New("ibm mq: no established connections")
  ErrInvalidConfig    = errors.New("ibm mq: invalid configuration")
  ErrNoConfig         = errors.New("ibm mq: no configuration")
  ErrAlreadyConnected = errors.New("ibm mq: connection already established")
  ErrConfigPathEmpty  = errors.New("ibm mq: configuration file path not specified")
  ErrChannelCfgNotSup = errors.New("ibm mq: channel configuration not supported")
  ErrNotFoundPipe     = errors.New("ibm mq: not found pipe by alias")
  ErrNotInitialised   = errors.New("ibm mq: queue not initialized")
)

type Queue interface {
  Put(ctx context.Context, msg *queue.Msg) error
  Get(ctx context.Context, msg *queue.Msg) error
  GetByCorrelId(ctx context.Context, correlId []byte) (*queue.Msg, error)
  GetByMsgId(ctx context.Context, msgId []byte) (*queue.Msg, error)
  Browse(ctx context.Context) (<-chan *queue.Msg, error)
  Alias() string
  CfgByStr(cfg string) error
  IsConfigured() bool
  Open() error
  Close() error
  UpdateBaseCfg()
  IsSubscribed() bool
  RegisterInMsg(hnd func(*queue.Msg))
  UnregisterInMsg()
}
