package main

import (
  "bytes"
  "fmt"
  "github.com/caarlos0/env/v6"
  mqpro "github.com/semenovem/gopkg_mqpro/v2"
  "github.com/sirupsen/logrus"
  "strings"
)

var cfg = &appConfig{}

type appConfig struct {
  ApiPort  int    `env:"ENV_API_PORT"`  // Порт api управления
  LogLevel string `env:"ENV_LOG_LEVEL"` // Уровень логирования приложения

  // mq
  MqLogLevel      string `env:"ENV_MQPRO_LOG_LEVEL"` // Уровень логирования ibmmq
  MqQueOper1Put   string `env:"ENV_MQPRO_QUEUE_OPER1_PUT"`
  MqQueOper1Get   string `env:"ENV_MQPRO_QUEUE_OPER1_GET"`

  MqQueOper2Put string `env:"ENV_MQPRO_QUEUE_OPER2_PUT"`
  MqQueOper2Get string `env:"ENV_MQPRO_QUEUE_OPER2_GET"`

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
    log.Error(err)
    fatal = true
  }

  if cfg.ApiPort == 0 {
    log.Warn("Не установлен порт api ENV_API_PORT")
    fatal = true
  }

  lev, err := logrus.ParseLevel(cfg.LogLevel)
  if err != nil {
    log.Warn("Не установлен уровень логирования приложения ENV_LOG_LEVEL. <%s>\n", err)
    fatal = true
    lev = logrus.TraceLevel
  }

  if lev >= logrus.InfoLevel {
    cfg.logInfo = true
  }

  // mq
  cfgIbmmq, err := mqpro.ParseDefaultEnv()
  if err != nil {
    log.Warn("Ошибка получения значений переменных окружения для настройки ibmmq")
    fatal = true
  }

  err = ibmmq.Cfg(cfgIbmmq)
  if err != nil {
    log.Warn(err)
    fatal = true
  }

  lev, err = logrus.ParseLevel(cfg.MqLogLevel)
  if err == nil {
    logIbmmq.Logger.SetLevel(lev)
  } else {
    log.Warn("Не установлен уровень логирования ibmmq ENV_MQ_LOG_LEVEL. <%s>\n", err)
    fatal = true
    lev = logrus.TraceLevel
  }

  err = ibmmqOper1Get.CfgByStr(cfg.MqQueOper1Get)
  if err != nil {
    fatal = true
    log.Warn(err)
  }

  err = ibmmqOper1Put.CfgByStr(cfg.MqQueOper1Put)
  if err != nil {
    fatal = true
    log.Warn(err)
  }

  ibmmq.PrintDefaultEnv()
  ibmmq.PrintSetCli("mgr")
  ibmmqOper1Get.PrintSetCli("queue/" + ibmmqOper1Get.Alias())

  if fatal {
    panic("При подготовке конфигурации есть фатальные ошибки. Подробности в логах")
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

  fmt.Println(buf.String())
}
