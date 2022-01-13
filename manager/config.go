package manager

import (
  "fmt"
  "github.com/semenovem/gopkg_mqpro/v2/queue"
)

type Config struct {
  Host          string
  Port          int32
  Manager       string
  Channel       string
  App           string
  User          string
  Pass          string
  Tls           bool
  KeyRepository string
  MaxMsgLength  int32
}

func (m *Mqpro) IsConfigured() bool {
  return m.host != "" && m.port > 0 && m.manager != "" && m.channel != ""
}

func (m *Mqpro) Cfg(c *Config) error {
  m.mx.Lock()
  defer m.mx.Unlock()

  fatal := false

  if c.Host == "" {
    m.log.Error("пустое значение Host")
    fatal = true
  }
  if c.Port == 0 {
    m.log.Error("пустое значение Port")
    fatal = true
  }
  if c.Manager == "" {
    m.log.Error("пустое значение Manager")
    fatal = true
  }
  if c.Channel == "" {
    m.log.Error("пустое значение Channel")
    fatal = true
  }

  if fatal {
    return ErrInvalidConfig
  }

  m.host = c.Host
  m.port = c.Port
  m.manager = c.Manager
  m.channel = c.Channel
  m.app = c.App
  m.user = c.User
  m.pass = c.Pass
  m.tls = c.Tls
  m.keyRepository = c.KeyRepository
  m.maxMsgLen = c.MaxMsgLength

  return nil
}

func (m *Mqpro) PrintSetCli(p string) {
  queue.PrintSetCli(m.getSet(), p)
}

func (m *Mqpro) getSet() []map[string]string {
  return []map[string]string{
    {"host": m.host},
    {"port": fmt.Sprintf("%d", m.port)},
    {"manager": m.manager},
    {"channel": m.channel},
    {"app": m.app},
    {"user": m.user},
    {"pass": m.pass},
    {"tls": fmt.Sprintf("%t", m.tls)},
    {"keyRepository": m.keyRepository},
    {"maxMsgLen": fmt.Sprintf("%d", m.maxMsgLen)},
  }
}
