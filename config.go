package mqpro

import (
  "bytes"
  "fmt"
  "github.com/caarlos0/env/v6"
  "github.com/semenovem/gopkg_mqpro/v2/manager"
  "github.com/semenovem/gopkg_mqpro/v2/queue"
  "os"
  "strings"
  "unicode/utf8"
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
  return m.managerCfg != nil && m.coreSet != nil &&
    m.managerCfg.Host != "" && m.managerCfg.Port > 0 &&
    m.managerCfg.Manager != "" && m.managerCfg.Channel != ""
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

  m.managerCfg = &manager.Config{
    Host:          c.Host,
    Port:          c.Port,
    Manager:       c.Manager,
    Channel:       c.Channel,
    App:           c.App,
    User:          c.User,
    Pass:          c.Pass,
    Tls:           c.Tls,
    KeyRepository: c.KeyRepository,
    MaxMsgLength:  c.MaxMsgLength,
  }
  
  for _, man := range m.managers {
    err := man.Cfg(m.managerCfg)
    if err != nil {
      return err
    }
  }

  m.coreSet = &coreSet
  return nil
}

// ParseDefaultEnv использует значения переменных окружения по умолчанию
func ParseDefaultEnv() (*Config, error) {
  c := new(Config)
  return c, env.Parse(c)
}

func (m *Mqpro) GetCoreSet() *queue.CoreSet {
  return m.coreSet
}

func (m *Mqpro) SetDevMode(v bool) {
  m.coreSet.DevMode = v
  for _, q := range m.queues {
    q.UpdateBaseCfg()
  }
}

func (m *Mqpro) PrintSetCli(p string) {
  queue.PrintSetCli(m.getSet(), p)
}

func (m *Mqpro) getSet() []map[string]string {
  return []map[string]string{
    {"manager/host": m.managerCfg.Host},
    {"manager/port": fmt.Sprintf("%d", m.managerCfg.Port)},
    {"manager/manager": m.managerCfg.Manager},
    {"manager/channel": m.managerCfg.Channel},
    {"manager/app": m.managerCfg.App},
    {"manager/user": m.managerCfg.User},
    {"manager/pass": m.managerCfg.Pass},
    {"manager/tls": fmt.Sprintf("%t", m.managerCfg.Tls)},
    {"manager/keyRepository": m.managerCfg.KeyRepository},
    {"manager/maxMsgLen": fmt.Sprintf("%d", m.managerCfg.MaxMsgLength)},

    {"coreSet/DevMode": fmt.Sprintf("%t", m.coreSet.DevMode)},
    {"coreSet/Header": fmt.Sprintf("%s", queue.HeaderMapByKey[m.coreSet.Header])},
    {"coreSet/Rfh2CodedCharSetId": fmt.Sprintf("%d", m.coreSet.Rfh2CodedCharSetId)},
    {"coreSet/Rfh2RootTag": m.coreSet.Rfh2RootTag},
  }
}

func (m *Mqpro) PrintDefaultEnv() {
  var (
    buf    = bytes.NewBufferString("Default environment:\n")
    k, v   string
    max, l int
  )

  li := []string{
    "ENV_MQPRO_DEV_MODE",
    "ENV_MQPRO_HOST",
    "ENV_MQPRO_PORT",
    "ENV_MQPRO_MANAGER",
    "ENV_MQPRO_CHANNEL",
    "ENV_MQPRO_APP",
    "ENV_MQPRO_USER",
    "ENV_MQPRO_PASS",
    "ENV_MQPRO_HEADER",
    "ENV_MQPRO_RFH2_CODE_CHAR_SET_ID",
    "ENV_MQPRO_RFH2_ROOT_TAG",
    "ENV_MQPRO_RFH2_OFF_ROOT_TAG",
    "ENV_MQPRO_TLS",
    "ENV_MQPRO_KEY_REPOSITORY",
    "ENV_MQPRO_MESSAGE_LENGTH",
  }

  for _, k = range li {
    l = utf8.RuneCountInString(k)
    if max < l {
      max = l
    }
  }

  for _, k = range li {
    v = strings.Repeat(" ", max-utf8.RuneCountInString(k)) + "= " + os.Getenv(k)
    buf.WriteString(fmt.Sprintf("%s%s\n", k, v))
  }

  fmt.Println(buf.String())
}
