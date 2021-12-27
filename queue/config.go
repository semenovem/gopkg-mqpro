package queue

import (
  "bytes"
  "fmt"
  "time"
)

// Cfg Данные подключения
type Cfg struct {
  DevMode            bool
  Host               string
  Port               int
  Manager            string
  Channel            string
  Queue              string // Название очереди
  App                string
  User               string
  Pass               string
  Header             string // Тип заголовков [prop | rfh2]
  Rfh2CodedCharSetId int32  // Тип кодирования
  Rfh2RootTag        string
  MaxMsgLength       int32
  Tls                bool
  KeyRepository      string
  CertificateLabel   string
  ReconnectDelay     time.Duration // Повтор попытки подключения
  RetryOper          time.Duration // Повтор операции
}

func (c *Conn) Cfg(cfg *Cfg) {
  c.mx.Lock()
  defer c.mx.Unlock()

  c.cfg = cfg

  if cfg.Rfh2RootTag == "" {
    c.log.Debugf("Не указан корневой тег, используем значение по умолчанию %s",
      defRootTagHeader)
    cfg.Rfh2RootTag = defRootTagHeader
  }

  if cfg.Queue == "" {
    c.log.Warnf("Не указано название очереди")
  }

  if cfg.MaxMsgLength == 0 {
    cfg.MaxMsgLength = defMaxMsgLength
  }

  m := map[string]interface{}{
    "conn": fmt.Sprintf("%s|%s", c.endpoint(), cfg.Manager),
  }
  c.log = c.log.WithFields(m)

  if cfg.ReconnectDelay != 0 {
    c.reconnectDelay = cfg.ReconnectDelay
  }

  if cfg.Header != "" {
    h, err := parseHeaderType(cfg.Header)
    if err == nil {
      c.h = h
    } else {
      c.log.Warnf("Передано не валидное значение типа заголовков. "+
        "Используем значение по умолчанию: %s", headerVal[defHeader])
    }

    if c.h == headerRfh2 {
      c.rfh2 = newRfh2Cfg()
      if cfg.Rfh2CodedCharSetId != 0 {
        c.rfh2.CodedCharSetId = cfg.Rfh2CodedCharSetId
      }
    }
  }
}

func (c *Conn) PrintCfg() {
  var buf = bytes.NewBufferString("")
  f := func(s string, i ...interface{}) {
    buf.WriteString(fmt.Sprintf(s, i...))
  }

  cfg := c.cfg

  f("Environment variable values:\n")
  f("DevMode:           = %t\n", cfg.DevMode)
  f("Host               = %s\n", cfg.Host)
  f("Port               = %d\n", cfg.Port)
  f("Manager            = %s\n", cfg.Manager)
  f("Channel            = %s\n", cfg.Channel)
  f("Queue              = %s\n", cfg.Queue)
  f("App                = %s\n", cfg.App)
  f("User               = %s\n", cfg.User)
  f("Pass               = %s\n", cfg.Pass)
  f("Header             = %s\n", cfg.Header)
  f("Rfh2CodedCharSetId = %s\n", cfg.Rfh2CodedCharSetId)
  f("Rfh2RootTag        = %t\n", cfg.Rfh2RootTag)
  f("MaxMsgLength       = %s\n", cfg.MaxMsgLength)
  f("Tls                = %d\n", cfg.Tls)
  f("KeyRepository      = %d\n", cfg.KeyRepository)
  f("CertificateLabel   = %d\n", cfg.CertificateLabel)

  fmt.Println(buf.String())
}
