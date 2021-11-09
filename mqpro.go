package mqpro

import (
  "context"
  "github.com/sirupsen/logrus"
  "sync"
  "time"
)

type Mqpro struct {
  rootCtx               context.Context
  ctx                   context.Context
  ctxCancel             context.CancelFunc
  conns                 []*Mqconn
  connGet               []*Mqconn
  connPut               []*Mqconn
  connBrowse            []*Mqconn
  fnEventInMsg          func(*Msg)    // Обработчик входящих сообщений
  mx                    sync.Mutex    // Подключение / отключение
  delayBeforeDisconnect time.Duration // Задержка перед разрывом соединения
  reconnDelay           time.Duration // Задержка при повторных попытках подключения к MQ
  log                   *logrus.Entry
  cfg                   *Config
}

func New(rootCtx context.Context, l *logrus.Entry) *Mqpro {
  return &Mqpro{
    rootCtx:               rootCtx,
    delayBeforeDisconnect: defDisconnDelay,
    reconnDelay:           defReconnDelay,
    log:                   l,
  }
}

func (p *Mqpro) SetConn(connLi ...*Mqconn) {
  for _, conn := range connLi {
    switch conn.Type() {
    case TypeGet:
      p.connGet = append(p.connGet, conn)
      if p.fnEventInMsg != nil {
        conn.RegisterEventInMsg(p.fnEventInMsg)
      }
    case TypePut:
      p.connPut = append(p.connPut, conn)
    case TypeBrowse:
      p.connBrowse = append(p.connBrowse, conn)

    default:
      p.log.Panic("Unknown connection type")
    }

    p.conns = append(p.conns, conn)
  }
}

func (p *Mqpro) SetLogger(l *logrus.Entry) {
  p.log = l
}

func (p *Mqpro) GetConns() []*Mqconn {
  return p.conns
}
