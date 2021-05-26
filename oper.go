package mqpro

import (
  "errors"
)

// IsConnect есть ли соединение к менеджерам IBM MQ
func (p *Mqpro) IsConnect() bool {
  if len(p.conns) == 0 {
    return false
  }

  a := len(p.connPut) == 0 || p.isConnect(p.connPut)
  b := len(p.connGet) == 0 || p.isConnect(p.connGet)
  c := len(p.connBrowse) == 0 || p.isConnect(p.connBrowse)

  return a && b && c
}

func (p *Mqpro) isConnect(li []*Mqconn) bool {
  for _, conn := range li {
    if conn.IsConnected() {
      return true
    }
  }
  return false
}

func (p *Mqpro) Ready() error {
  if len(p.conns) == 0 {
    return errors.New("no data to connect to IBM MQ")
  }

  if p.IsConnect() {
    return nil
  }

  return errors.New("no connection to IBM MQ")
}
