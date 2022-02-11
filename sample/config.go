package main

import (
  "bytes"
  "fmt"
  "github.com/caarlos0/env/v6"
  "github.com/semenovem/mqm/v2"
  "strings"
)

var cfg = &appConfig{}

type appConfig struct {
  ApiPort  int    `env:"ENV_API_PORT"`  // Порт api управления
  LogLevel string `env:"ENV_LOG_LEVEL"` // Уровень логирования приложения

  // mq
  MqYamlCfgFile string `env:"ENV_MQM_YAML_CFG_FILE"`

  // При старте подписаться на входящие сообщения
  Subscribe bool `env:"ENV_SIMPLE_SUBSCRIBER"`

  // В ответ на входящее сообщение - отправить полученное, установив CorrelID
  Mirror bool `env:"ENV_MIRROR"`
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

  // mq. Файл конфигурации
  c, err := mqm.ParseCfgYaml(cfg.MqYamlCfgFile)
  if err != nil {
    log.Errorf("Ошибка парсинга файла конфигурации MQM '%s', %v\n",
      cfg.MqYamlCfgFile, err)
    fatal = true
  } else {
    err = mq.Cfg(c)
    if err != nil {
      log.Error("Ошибка при установке конфигурации MQM: ", err)
      fatal = true
    }
  }

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
  f("ENV_SIMPLE_SUBSCRIBER  = %t\n", cfg.Subscribe)
  f("ENV_MIRROR             = %t\n", cfg.Mirror)

  fmt.Println(buf.String())

  mq.PrintSetCli("mqm")
}
