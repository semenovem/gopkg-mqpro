package mqm

import (
  "context"
  "github.com/semenovem/mqm/v2/manager"
  "github.com/semenovem/mqm/v2/queue"
  "github.com/sirupsen/logrus"
  "sync"
  "time"
)

type Mqm struct {
  rootCtx      context.Context
  ctx          context.Context
  ctxCanc      context.CancelFunc
  mx           sync.Mutex // Подключение / отключение
  log          *logrus.Entry
  disconnDelay time.Duration
  isConnected  bool
  queueCfg     *queue.BaseConfig // Конфиг ibmmq очереди
  managerCfg   *manager.Config   // Конфиг ibmmq менеджера
  queues       []*queue.Queue
  managers     []*manager.Manager
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
func (m *Mqm) NewQueue(a string) *queue.Queue {
  l := m.log.WithField("_t", "queue")
  logMag := m.log.WithField("_t", "manager")

  man := manager.New(m.rootCtx, logMag)

  q := queue.New(m.rootCtx, l, man, m, a)
  m.queues = append(m.queues, q)
  m.managers = append(m.managers, man)

  return q
}

func (m *Mqm) GetQueueByAlias(a string) *queue.Queue {
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
