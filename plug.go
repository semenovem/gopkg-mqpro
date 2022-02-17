package mqm

import (
  "context"
  "github.com/semenovem/mqm/v2/queue"
  "github.com/sirupsen/logrus"
)

type plug struct {
  log *logrus.Entry
}

func (p plug) Put(context.Context, *queue.Msg) error {
  p.log.Panic(ErrNotInitialised)
  return nil
}

func (p plug) Get(context.Context, *queue.Msg) error {
  p.log.Panic(ErrNotInitialised)
  return nil
}

func (p plug) GetByCorrelId(context.Context, []byte) (*queue.Msg, error) {
  p.log.Panic(ErrNotInitialised)
  return nil, nil
}

func (p plug) GetByMsgId(context.Context, []byte) (*queue.Msg, error) {
  p.log.Panic(ErrNotInitialised)
  return nil, nil
}

func (p plug) Browse(context.Context) (<-chan *queue.Msg, error) {
  p.log.Panic(ErrNotInitialised)
  return nil, nil
}

func (p plug) Alias() string {
  p.log.Panic(ErrNotInitialised)
  return ""
}

func (p plug) CfgByStr(string) error {
  p.log.Panic(ErrNotInitialised)
  return nil
}

func (p plug) IsConfigured() bool {
  p.log.Panic(ErrNotInitialised)
  return false
}

func (p plug) Open() error {
  p.log.Panic(ErrNotInitialised)
  return nil
}

func (p plug) Close() error {
  p.log.Panic(ErrNotInitialised)
  return nil
}

func (p plug) UpdateBaseCfg() {
  p.log.Panic(ErrNotInitialised)
}

func (p plug) IsSubscribed() bool {
  p.log.Panic(ErrNotInitialised)
  return false
}

func (p plug) RegisterInMsg(func(*queue.Msg)) {
  p.log.Panic(ErrNotInitialised)
}

func (p plug) UnregisterInMsg() {
  p.log.Panic(ErrNotInitialised)
}

func (p plug) Ready() bool {
  p.log.Panic(ErrNotInitialised)
  return false
}
