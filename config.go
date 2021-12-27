package mqpro

import (
  "bytes"
  "fmt"
  "github.com/caarlos0/env/v6"
  "github.com/pkg/errors"
  "github.com/sirupsen/logrus"
  "gopkg.in/yaml.v3"
  "io/ioutil"
)

type QueCfg struct {
  Alias  string `yaml:"alias"`
  Browse string `yaml:"browse"`
  Put    string `yaml:"put"`
  Get    string `yaml:"get"`
}

type Queues struct {
  Alias string   `yaml:"alias"`
  Name  string   `yaml:"name"`
}

type Config struct {
  DevMode        bool     `yaml:"devMode" env:"ENV_MQPRO_DEV_MODE"`
  Host           string   `yaml:"host" env:"ENV_MQPRO_HOST"`
  Port           int      `yaml:"port" env:"ENV_MQPRO_PORT"`
  Manager        string   `yaml:"manager" env:"ENV_MQPRO_MANAGER"`
  Channel        string   `yaml:"channel" env:"ENV_MQPRO_CHANNEL"`
  App            string   `yaml:"app" env:"ENV_MQPRO_APP"`
  User           string   `yaml:"user" env:"ENV_MQPRO_USER"`
  Pass           string   `yaml:"pass" env:"ENV_MQPRO_PASS"`
  BrowseQueue    string   `yaml:"browseQueue" env:"ENV_MQPRO_BROWSE_QUEUE"`
  PutQueue       string   `yaml:"putQueue" env:"ENV_MQPRO_PUT_QUEUE"`
  GetQueue       string   `yaml:"getQueue" env:"ENV_MQPRO_GET_QUEUE"`
  Header         string   `yaml:"header" env:"ENV_MQPRO_HEADER"`
  CodedCharSetId int32    `yaml:"codedCharSetId" env:"ENV_MQPRO_CODED_CHAR_SET_ID"`
  RootTag        string   `yaml:"rootTag" env:"ENV_MQPRO_ROOT_TAG"`
  OffRootTag     bool     `yaml:"offRootTag" env:"ENV_MQPRO_OFF_ROOT_TAG"`
  Tls            bool     `yaml:"tls" env:"ENV_MQPRO_TLS"`
  KeyRepository  string   `yaml:"keyRepository" env:"ENV_MQPRO_KEY_REPOSITORY"`
  MaxMsgLength   int32    `yaml:"maxMsgLength" env:"ENV_MQPRO_MAX_MSG_LENGTH"`
  MultiQueues    []QueCfg `yaml:"multiQueues"`
  Queues         []Queues `yaml:"queues"`
}

func ParseConfig(f string) (*Config, error) {
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

// UseDefEnv2 получение конфигурации из стандартных переменных окружения
func UseDefEnv2(l *logrus.Entry) *Config {
  var c = &Config{}
  if err := env.Parse(c); err != nil {
    l.Error(err)
  }
  return c
}

func (p *Mqpro) Cfg(c *Config) {
  p.mx.Lock()
  defer p.mx.Unlock()

  p.cfg = c

  if c.OffRootTag {
    c.RootTag = ""
  } else {
    if c.RootTag == "" {
      p.log.Debugf("Не указан корневой тег, используем значение по умолчанию %s",
        defRootTagHeader)

      c.RootTag = defRootTagHeader
    }
  }

  if (c.GetQueue != "" || c.PutQueue != "" || c.BrowseQueue != "") && len(c.MultiQueues) != 0 {
    p.log.Warnf("Указаны данные очередей в single использовании и multi наборе. " +
      "Будут использоваться очереди только из single набора")
  }
}

func (p *Mqpro) PrintCfg() {
  var buf = bytes.NewBufferString("")
  f := func(s string, i ...interface{}) {
    buf.WriteString(fmt.Sprintf(s, i...))
  }

  c := p.cfg

  f("Environment variable values:\n")
  f("ENV_MQPRO_DEV_MODE:    = %t\n", c.DevMode)
  f("ENV_MQPRO_HOST         = %s\n", c.Host)
  f("ENV_MQPRO_PORT         = %d\n", c.Port)
  f("ENV_MQPRO_MANAGER      = %s\n", c.Manager)
  f("ENV_MQPRO_CHANNEL      = %s\n", c.Channel)
  f("ENV_MQPRO_APP          = %s\n", c.App)
  f("ENV_MQPRO_USER         = %s\n", c.User)
  f("ENV_MQPRO_PASS         = %s\n", c.Pass)
  f("ENV_MQPRO_BROWSE       = %s\n", c.BrowseQueue)
  f("ENV_MQPRO_PUT          = %s\n", c.PutQueue)
  f("ENV_MQPRO_GET          = %s\n", c.GetQueue)
  f("ENV_MQPRO_HEADER       = %s\n", c.Header)
  f("ENV_MQPRO_ROOT_TAG     = %s\n", c.RootTag)
  f("ENV_MQPRO_OFF_ROOT_TAG = %t\n", c.OffRootTag)
  f("ENV_MQPRO_TLS          = %t\n", c.Tls)
  f("ENV_MQPRO_KEY_REPOSITORY = %s\n", c.KeyRepository)
  f("ENV_MQPRO_MAX_MSG_LENGTH = %d\n", c.MaxMsgLength)

  if len(c.MultiQueues) != 0 {
    for _, v := range c.MultiQueues {
      f("TRANSPORT::ALIAS   = %s\n", v.Alias)
      f("TRANSPORT::GET     = %s\n", v.Get)
      f("TRANSPORT::PUT     = %s\n", v.Put)
      f("TRANSPORT::BROWSER = %s\n", v.Browse)
      f("TRANSPORT::-------------------\n")
    }
  }

  fmt.Println(buf.String())
}
