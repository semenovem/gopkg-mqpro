package mqm

import (
  "bytes"
  "fmt"
  "github.com/caarlos0/env/v6"
  "github.com/pkg/errors"
  "github.com/semenovem/mqm/v2/manager"
  "github.com/semenovem/mqm/v2/queue"
  "github.com/sirupsen/logrus"
  "gopkg.in/yaml.v3"
  "io/ioutil"
  "os"
  "strings"
  "unicode/utf8"
)

type Config struct {
  DevMode            bool   `env:"ENV_MQM_DEV_MODE" yaml:"devMode"`
  LogLev             string `env:"ENV_MQM_LOG_LEVEL" yaml:"logLev"`
  Host               string `env:"ENV_MQM_HOST" yaml:"host"`
  Port               int32  `env:"ENV_MQM_PORT" yaml:"port"`
  Manager            string `env:"ENV_MQM_MANAGER" yaml:"manager"`
  Channel            string `env:"ENV_MQM_CHANNEL" yaml:"channel"`
  App                string `env:"ENV_MQM_APP" yaml:"app"`
  User               string `env:"ENV_MQM_USER" yaml:"user"`
  Pass               string `env:"ENV_MQM_PASS" yaml:"pass"`
  Header             string `env:"ENV_MQM_HEADER" yaml:"header"`
  Rfh2CodedCharSetId int32  `env:"ENV_MQM_RFH2_CODE_CHAR_SET_ID" yaml:"rfh2CodedCharSetId"`
  Rfh2RootTag        string `env:"ENV_MQM_RFH2_ROOT_TAG" yaml:"rfh2RootTag"`
  Rfh2OffRootTag     bool   `env:"ENV_MQM_RFH2_OFF_ROOT_TAG" yaml:"rfh2OffRootTag"`
  Tls                bool   `env:"ENV_MQM_TLS" yaml:"tls"`
  KeyRepository      string `env:"ENV_MQM_KEY_REPOSITORY" yaml:"keyRepository"`
  MaxMsgLength       int32  `env:"ENV_MQM_MESSAGE_LENGTH" yaml:"maxMsgLength"`

  Queues []QueCfg  `yaml:"queues"`
  Pipes  []PipeCfg `yaml:"pipes"`
}

type QueCfg struct {
  Alias string `yaml:"alias"`
  Name  string `yaml:"name"`
}

type PipeCfg struct {
  Alias string `yaml:"alias"`
  Put   string `yaml:"put"`
  Get   string `yaml:"get"`
}

func (m *Mqm) isConfigured() bool {
  return m.managerCfg != nil && m.queueCfg != nil &&
    m.managerCfg.Host != "" && m.managerCfg.Port > 0 &&
    m.managerCfg.Manager != "" && m.managerCfg.Channel != ""
}

func (m *Mqm) Cfg(c *Config) error {
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
        "Значение по умолчанию: Rfh2OffRootTag=%s", queue.DefRootTagHeader)
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

  if c.Queues != nil {
    for _, qc := range c.Queues {
      q := m.GetQueueByAlias(qc.Alias)
      if q == nil {
        m.log.Warnf("Очередь с алиасом {%s} не существует", qc.Alias)
        continue
      }
      err := q.CfgByStr(qc.Name)
      if err != nil {
        fatal = true
        m.log.Error(err)
      }
    }
  }

  if c.Pipes != nil {
    for _, p := range c.Pipes {
      err := m.cfgPipe(&p)
      if err != nil {
        fatal = true
        m.log.Error(err)
      }
    }
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
func (m *Mqm) CfgYaml(file string) error {
  c, err := ParseCfgYaml(file)
  if err != nil {
    return err
  }
  return m.Cfg(c)
}

// CfgEnv конфигурирование переменными окружения
func (m *Mqm) CfgEnv() error {
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

func (m *Mqm) SetDevMode(v bool) {
  m.queueCfg.DevMode = v

  // TODO конфигурация передается по ссылке, обновление дочерним объектам будет доступно сразу
  // TODO нужно проверить !
  //for _, q := range m.queues {
  //  q.UpdateBaseCfg()
  //}
}

func (m *Mqm) PrintSetCli(p string) {
  queue.PrintSetCli(m.getSet(), p)
}

func (m *Mqm) getSet() []map[string]string {
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
func (m *Mqm) PrintDefaultEnv() {
  var (
    buf    = bytes.NewBufferString("Standard environment variables:\n")
    k, v   string
    max, l int
  )

  li := []string{
    "ENV_MQM_DEV_MODE",
    "ENV_MQM_HOST",
    "ENV_MQM_PORT",
    "ENV_MQM_MANAGER",
    "ENV_MQM_CHANNEL",
    "ENV_MQM_APP",
    "ENV_MQM_USER",
    "ENV_MQM_PASS",
    "ENV_MQM_HEADER",
    "ENV_MQM_RFH2_CODE_CHAR_SET_ID",
    "ENV_MQM_RFH2_ROOT_TAG",
    "ENV_MQM_RFH2_OFF_ROOT_TAG",
    "ENV_MQM_TLS",
    "ENV_MQM_KEY_REPOSITORY",
    "ENV_MQM_MESSAGE_LENGTH",
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
