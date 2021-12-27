package mqpro

import (
  "context"
  "github.com/semenovem/gopkg_mqpro/v2/queue"
  "github.com/sirupsen/logrus"
)

type Factory struct {
  ctx   context.Context
  log   *logrus.Entry
  conns map[string]*queue.Conn
}

func (f *Factory) ByAlias(alias string) *Conn {
  _, ok := f.conns[alias]
  if ok {
    f.log.Panic(ErrAliasExist)
  }
  f.conns[alias] = queue.New(f.ctx, f.log.WithField("alias", alias))
  return &Mqconn{}
}
