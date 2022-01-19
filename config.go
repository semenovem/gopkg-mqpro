package mqpro

import (
  "bytes"
  "fmt"
  "github.com/caarlos0/env/v6"
  "github.com/pkg/errors"
  "github.com/semenovem/gopkg_mqpro/v2/manager"
  "github.com/semenovem/gopkg_mqpro/v2/queue"
  "github.com/sirupsen/logrus"
  "gopkg.in/yaml.v3"
  "io/ioutil"
  "os"
  "strings"
  "unicode/utf8"
)

type Config struct {
  DevMode            bool   `env:"ENV_MQPRO_DEV_MODE" yaml:"devMode"`
  LogLev             string `env:"ENV_MQPRO_LOG_LEVEL" yaml:"logLev"`
  Host               string `env:"ENV_MQPRO_HOST" yaml:"host"`
  Port               int32  `env:"ENV_MQPRO_PORT" yaml:"port"`
  Manager            string `env:"ENV_MQPRO_MANAGER" yaml:"manager"`
  Channel            string `env:"ENV_MQPRO_CHANNEL" yaml:"channel"`
  App                string `env:"ENV_MQPRO_APP" yaml:"app"`
  User               string `env:"ENV_MQPRO_USER" yaml:"user"`
  Pass               string `env:"ENV_MQPRO_PASS" yaml:"pass"`
  Header             string `env:"ENV_MQPRO_HEADER" yaml:"header"`
  Rfh2CodedCharSetId int32  `env:"ENV_MQPRO_RFH2_CODE_CHAR_SET_ID" yaml:"rfh2CodedCharSetId"`
  Rfh2RootTag        string `env:"ENV_MQPRO_RFH2_ROOT_TAG" yaml:"rfh2RootTag"`
  Rfh2OffRootTag     bool   `env:"ENV_MQPRO_RFH2_OFF_ROOT_TAG" yaml:"rfh2OffRootTag"`
  Tls                bool   `env:"ENV_MQPRO_TLS" yaml:"tls"`
  KeyRepository      string `env:"ENV_MQPRO_KEY_REPOSITORY" yaml:"keyRepository"`
  MaxMsgLength       int32  `env:"ENV_MQPRO_MESSAGE_LENGTH" yaml:"maxMsgLength"`
  Queues             []*Queues
}

type Queues struct {
  Alias string `yaml:"alias"`
  Name  string `yaml:"name"`
}

func (m *Mqpro) isConfigured() bool {
  return m.managerCfg != nil && m.queueCfg != nil &&
    m.managerCfg.Host != "" && m.managerCfg.Port > 0 &&
    m.managerCfg.Manager != "" && m.managerCfg.Channel != ""
}

func (m *Mqpro) Cfg(c *Config) error {
  m.mx.Lock()
  defer m.mx.Unlock()

  fatal := false

  queCfg := queue.BaseConfig{
    DevMode:            c.DevMode,
    Rfh2CodedCharSetId: c.Rfh2CodedCharSetId,
  }

  if c.Header == "" {
    m.log.Warnf("Не передан тип заголовков. "+
      "Значение по умолчанию: Header={%s}", queue.HeaderMapByKey[queue.DefHeader])
    queCfg.Header = queue.DefHeader
  } else {
    h, err := queue.ParseHeader(c.Header)
    if err != nil {
      m.log.Errorf("не валидное значение Header = %s", c.Header)
      fatal = true
    }
    queCfg.Header = h
  }

  if c.Rfh2OffRootTag {
    queCfg.Rfh2RootTag = ""
  } else {
    if c.Rfh2RootTag == "" {
      m.log.Warnf("Не установлено значение корневого тега. "+
        "Значение по умолчанию: Rfh2OffRootTag={%s}", queue.DefRootTagHeader)
      queCfg.Rfh2RootTag = queue.DefRootTagHeader
    } else {
      queCfg.Rfh2RootTag = c.Rfh2RootTag
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

  if c.LogLev != "" {
    lev, err := logrus.ParseLevel(c.LogLev)
    if err != nil {
      m.log.Errorf("значение уровня логгирования LogLev={%s} не валидно", c.LogLev)
      return err
    }
    m.log.Logger.SetLevel(lev)
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

  m.queueCfg = &queCfg
  return nil
}

// CfgYaml конфигурирование файлом
func (m *Mqpro) CfgYaml(file string) error {
  c, err := ParseCfgYaml(file)
  if err != nil {
    return err
  }
  if err = m.Cfg(c); err != nil {
    return err
  }
  if c.Queues == nil {
    return nil
  }
  for _, qc := range c.Queues {
    q := m.GetQueueByAlias(qc.Alias)
    if q == nil {
      m.log.Warnf("Очередь с алиасом {%s} не существует", qc.Alias)
      continue
    }
    err = q.CfgByStr(qc.Name)
    if err != nil {
      return err
    }
  }
  return nil
}

// CfgEnv конфигурирование переменными окружения
func (m *Mqpro) CfgEnv() error {
  c, err := ParseDefaultEnv()
  if err != nil {
    return err
  }
  return m.Cfg(c)
}

// ParseDefaultEnv использует значения переменных окружения по умолчанию
func ParseDefaultEnv() (*Config, error) {
  c := new(Config)
  return c, env.Parse(c)
}

func ParseCfgYaml(f string) (*Config, error) {
  if f == "" {
    return nil, ErrConfigPathEmpty
  }
  byt, err := ioutil.ReadFile(f)
  if err != nil {
    return nil, errors.Wrap(err, "Error reading configuration file")
  }
  c := new(Config)
  err = yaml.Unmarshal(byt, c)
  if err != nil {
    return nil, errors.Wrapf(err, "Configuration file parsing error '%s'", f)
  }
  return c, nil
}

func (m *Mqpro) SetDevMode(v bool) {
  m.queueCfg.DevMode = v
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

    {"queueCfg/DevMode": fmt.Sprintf("%t", m.queueCfg.DevMode)},
    {"queueCfg/Header": fmt.Sprintf("%s", queue.HeaderMapByKey[m.queueCfg.Header])},
    {"queueCfg/Rfh2CodedCharSetId": fmt.Sprintf("%d", m.queueCfg.Rfh2CodedCharSetId)},
    {"queueCfg/Rfh2RootTag": m.queueCfg.Rfh2RootTag},
  }
}

// PrintDefaultEnv распечатать содержимое переменных окружения
func (m *Mqpro) PrintDefaultEnv() {
  var (
    buf    = bytes.NewBufferString("Standard environment variables:\n")
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
