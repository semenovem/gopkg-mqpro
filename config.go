package mqpro

import (
  "bytes"
  "fmt"
  "github.com/caarlos0/env/v6"
  "github.com/semenovem/gopkg_mqpro/v2/queue"
)

type Config struct {
  DevMode            bool   `env:"ENV_MQPRO_DEV_MODE"`
  Host               string `env:"ENV_MQPRO_HOST"`
  Port               int32  `env:"ENV_MQPRO_PORT"`
  Manager            string `env:"ENV_MQPRO_MANAGER"`
  Channel            string `env:"ENV_MQPRO_CHANNEL"`
  App                string `env:"ENV_MQPRO_APP"`
  User               string `env:"ENV_MQPRO_USER"`
  Pass               string `env:"ENV_MQPRO_PASS"`
  Header             string `env:"ENV_MQPRO_HEADER"`
  Rfh2CodedCharSetId int32  `env:"ENV_MQPRO_RFH2_CODE_CHAR_SET_ID"`
  Rfh2RootTag        string `env:"ENV_MQPRO_RFH2_ROOT_TAG"`
  Rfh2OffRootTag     bool   `env:"ENV_MQPRO_RFH2_OFF_ROOT_TAG"`
  Tls                bool   `env:"ENV_MQPRO_TLS"`
  KeyRepository      string `env:"ENV_MQPRO_KEY_REPOSITORY"`
  MaxMsgLength       int32  `env:"ENV_MQPRO_MESSAGE_LENGTH"`
}

func (m *Mqpro) isConfigured() bool {
  return m.host != "" && m.port > 0 && m.manager != "" && m.channel != ""
}

func (m *Mqpro) Cfg(c *Config) error {
  m.mx.Lock()
  defer m.mx.Unlock()

  fatal := false

  coreSet := queue.CoreSet{
    DevMode:            c.DevMode,
    Rfh2CodedCharSetId: c.Rfh2CodedCharSetId,
  }

  if c.Header == "" {
    m.log.Warn("Не передан тип заголовков. "+
      "Используем значение по умолчанию {%s}", queue.DefHeader)
    coreSet.Header = queue.DefHeader
  } else {
    h, err := queue.ParseHeader(c.Header)
    if err != nil {
      m.log.Errorf("не валидное значение Header = %s", c.Header)
      fatal = true
    }
    coreSet.Header = h
  }

  if c.Rfh2OffRootTag {
    coreSet.Rfh2RootTag = ""
  } else {
    if c.Rfh2RootTag == "" {
      m.log.Warnf("Не установлено значение корневого тега. "+
        "Значение по умолчанию = {%s}", queue.DefRootTagHeader)
    } else {
      coreSet.Rfh2RootTag = c.Rfh2RootTag
    }
  }

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

  m.coreSet = &coreSet
  return nil
}

// ParseDefaultEnv использует значения переменных окружения по умолчанию
func ParseDefaultEnv() (*Config, error) {
  c := new(Config)
  return c, env.Parse(c)
}

func (m *Mqpro) GetQueueConfig() *queue.CoreSet {
  return m.coreSet
}

func (m *Mqpro) SetDevMode(v bool) {
  m.coreSet.DevMode = v
  for _, q := range m.queues {
    q.Set(m.coreSet)
  }
}

// PrintCfg
// Deprecated
func (m *Mqpro) PrintCfg() {
  var buf = bytes.NewBufferString("")
  f := func(s string, i ...interface{}) {
    buf.WriteString(fmt.Sprintf(s, i...))
  }

  f("Environment variable values:\n")
  f("ENV_MQPRO_DEV_MODE:    = %t\n", m.coreSet.DevMode)
  f("ENV_MQPRO_HOST         = %s\n", m.host)
  f("ENV_MQPRO_PORT         = %d\n", m.port)
  f("ENV_MQPRO_MANAGER      = %s\n", m.manager)
  f("ENV_MQPRO_CHANNEL      = %s\n", m.channel)
  f("ENV_MQPRO_APP          = %s\n", m.app)
  f("ENV_MQPRO_USER         = %s\n", m.user)
  f("ENV_MQPRO_PASS         = %s\n", m.pass)
  f("ENV_MQPRO_HEADER       = %s\n", m.coreSet.Header)
  f("ENV_MQPRO_RFH2_ROOT_TAG= %s\n", m.coreSet.Rfh2RootTag)
  f("ENV_MQPRO_TLS          = %t\n", m.tls)
  f("ENV_MQPRO_KEY_REPOSITORY = %s\n", m.keyRepository)
  f("ENV_MQPRO_MAX_MSG_LENGTH = %d\n", m.maxMsgLen)

  fmt.Println(buf.String())
}
