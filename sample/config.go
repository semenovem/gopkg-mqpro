package main

import (
  "fmt"
  "github.com/caarlos0/env"
  "github.com/sirupsen/logrus"
)

var cfg = &appConfig{}

type appConfig struct {
  ApiPort int `env:"ENV_API_PORT"`

  LogLevel string `env:"ENV_LOG_LEVEL"`

  // При старте подписаться на входящие сообщения
  SimpleSubscriber bool `env:"ENV_SIMPLE_SUBSCRIBER"`

  // В ответ на входящее сообщение - отправить полученное, установив CorrelID
  Mirror bool `env:"ENV_MIRROR"`
}

func init() {
  fatal := false

  if err := env.Parse(cfg); err != nil {
    fmt.Println("ERROR: ", err)
  }

  if cfg.ApiPort == 0 {
    fmt.Println("WARN: не установлен порт api ENV_API_PORT")
  }

  lev, err := logrus.ParseLevel(cfg.LogLevel)
  if err != nil {
    fmt.Println("WARN: не установлен уровень логирования ENV_LOG_LEVEL")
    fatal = true
  }

  l := logrus.New()
  l.SetLevel(lev)
  ibmmq.SetLogger(logrus.NewEntry(l).WithField("pkg", "mqpro"))

  fmt.Println("Список переменных окружения")
  fmt.Println("ENV_API_PORT           =  ", cfg.ApiPort)
  fmt.Println("ENV_LOG_LEVEL          =  ", cfg.LogLevel)
  fmt.Println("ENV_SIMPLE_SUBSCRIBER  =  ", cfg.SimpleSubscriber)
  fmt.Println("ENV_MIRROR             =  ", cfg.Mirror)

  if fatal {
    panic("")
  }
}
