package main

import (
  "bytes"
  "fmt"
  "github.com/caarlos0/env/v6"
  "github.com/sirupsen/logrus"
  "strings"
)

var cfg = &appConfig{}

type appConfig struct {
  ApiPort  int    `env:"ENV_API_PORT"`  // Порт api управления
  LogLevel string `env:"ENV_LOG_LEVEL"` // Уровень логирования приложения

  // mq
  MqLogLevel    string `env:"ENV_MQPRO_LOG_LEVEL"` // Уровень логирования ibmmq
  MqQueOper1Put string `env:"ENV_MQPRO_QUEUE_OPER1_PUT"`
  MqQueOper1Get string `env:"ENV_MQPRO_QUEUE_OPER1_GET"`

  MqQueOper2Put string `env:"ENV_MQPRO_QUEUE_OPER2_PUT"`
  MqQueOper2Get string `env:"ENV_MQPRO_QUEUE_OPER2_GET"`

  MqYamlCfgFile string `env:"ENV_MQPRO_YAML_CFG_FILE"`

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
  if cfg.MqYamlCfgFile != "" {
    err = mq.CfgYaml(cfg.MqYamlCfgFile)
    if err != nil {
      log.Warn("Ошибка конфигурации из файла YAML", err)
      fatal = true
    }
  } else {
    err = mq.CfgEnv()
    if err != nil {
      log.Warn("Ошибка конфигурации из переменных окружения с дефолтными значениями", err)
      fatal = true
    }

    // Настройка очередей значениями
    err = mqQueFooGet.CfgByStr(cfg.MqQueOper1Get)
    if err != nil {
      fatal = true
      log.Warn(err)
    }

    err = mqQueFooPut.CfgByStr(cfg.MqQueOper1Put)
    if err != nil {
      fatal = true
      log.Warn(err)
    }
  }

  //mq.PrintDefaultEnv()
  mq.PrintSetCli("mgr")
  //mqQueFooGet.PrintSetCli("queue/" + mqQueFooGet.Alias())

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
