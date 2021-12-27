package main

import (
  "bytes"
  "fmt"
  "github.com/caarlos0/env"
  mqpro "github.com/semenovem/gopkg_mqpro/v2"
  "github.com/sirupsen/logrus"
  "strings"
)

var cfg = &appConfig{}

type appConfig struct {
  ApiPort    int    `env:"ENV_API_PORT"`        // Порт api управления
  LogLevel   string `env:"ENV_LOG_LEVEL"`       // Уровень логирования приложения
  MqLogLevel string `env:"ENV_MQPRO_LOG_LEVEL"` // Уровень логирования ibmmq
  MqConfig   string `env:"ENV_MQPRO_CONFIG"`    // Файл конфигурации IBMMQ

  // При старте подписаться на входящие сообщения
  SimpleSubscriber bool `env:"ENV_SIMPLE_SUBSCRIBER"`

  // В ответ на входящее сообщение - отправить полученное, установив CorrelID
  Mirror bool `env:"ENV_MIRROR"`

  logInfo bool // Логировать содержимое полученных/отправленных сообщений
}

func init() {
  var (
    err   error
    fatal bool
  )

  if err := env.Parse(cfg); err != nil {
    fmt.Println("ERROR: ", err)
  }

  if cfg.ApiPort == 0 {
    fmt.Println("WARN: не установлен порт api ENV_API_PORT")
  }

  lev, err := logrus.ParseLevel(cfg.LogLevel)
  if err != nil {
    fmt.Printf("WARN: не установлен уровень логирования приложения ENV_LOG_LEVEL. <%s>\n", err)
    fatal = true
    lev = logrus.TraceLevel
  }

  if lev >= logrus.InfoLevel {
    cfg.logInfo = true
  }

  lev, err = logrus.ParseLevel(cfg.MqLogLevel)
  if err != nil {
    fmt.Printf("WARN: не установлен уровень логирования ibmmq ENV_MQ_LOG_LEVEL. <%s>\n", err)
    fatal = true
    lev = logrus.TraceLevel
  }

  l := logrus.New()
  l.SetLevel(lev)
  ibmmq.SetLogger(logrus.NewEntry(l).WithField("pkg", "mqpro"))

  if cfg.MqConfig != "" {
    c, err := mqpro.ParseConfig(cfg.MqConfig)
    if err == nil {
      ibmmq.Cfg(c)
    } else {
      fmt.Println("ERROR: не удалось получить конфигурацию из файла")
      fatal = true
    }
  } else {
    c := mqpro.UseDefEnv2(log)
    ibmmq.Cfg(c)
  }

  ibmmq.SetConn()

  // Вывод конфигурации
  //ibmmq.PrintCfg()
  //printCfg()

  if fatal {
    panic("")
  }
}

func printCfg() {
  var buf = bytes.NewBufferString("")
  f := func(s string, i ...interface{}) {
    buf.WriteString(fmt.Sprintf(s, i...))
  }

  f("Список переменных окружения и настроек:\n")
  f("ENV_API_PORT           = %d\n", cfg.ApiPort)
  f("ENV_LOG_LEVEL          = %s\n", strings.ToUpper(cfg.LogLevel))
  f("ENV_MQPRO_LOG_LEVEL    = %s\n", strings.ToUpper(cfg.MqLogLevel))
  f("ENV_SIMPLE_SUBSCRIBER  = %t\n", cfg.SimpleSubscriber)
  f("ENV_MIRROR             = %t\n", cfg.Mirror)
  f("cfg.logInfo            = %t\n", cfg.logInfo)

  conns := ibmmq.GetConns()
  if len(conns) > 0 {
    f("mqpro.DevMode          = %t\n", conns[0].DevMode)
  }

  fmt.Println(buf.String())
}
