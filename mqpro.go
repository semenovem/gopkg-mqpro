package mqpro

import (
  "context"
  "github.com/semenovem/gopkg_mqpro/v2/manager"
  "github.com/semenovem/gopkg_mqpro/v2/queue"
  "github.com/sirupsen/logrus"
  "sync"
  "time"
)

type Mqpro struct {
  rootCtx      context.Context
  ctx          context.Context
  ctxCanc      context.CancelFunc
  mx           sync.Mutex // Подключение / отключение
  log          *logrus.Entry
  disconnDelay time.Duration

  isConnected bool

  // TODO rename to `queueCfg`
  coreSet    *queue.CoreSet  // Конфиг ibmmq очереди
  managerCfg *manager.Config // Конфиг ibmmq менеджера

  queues   []*queue.Queue   // Очереди
  managers []*manager.Mqpro // Очереди
}

func New(rootCtx context.Context, l *logrus.Entry) *Mqpro {
  o := &Mqpro{
    rootCtx:      rootCtx,
    log:          l,
    disconnDelay: defDisconnDelay,
  }

  return o
}

// Queue Объект очереди
func (m *Mqpro) Queue(a string) *queue.Queue {
  l := m.log.WithField("_t", "queue")
  logMag := m.log.WithField("_t", "manager")

  man := manager.New(m.rootCtx, logMag)

  q := queue.New(m.rootCtx, l, man, m, a)
  m.queues = append(m.queues, q)
  m.managers = append(m.managers, man)

  return q
}
