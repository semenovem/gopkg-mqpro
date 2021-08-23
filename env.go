package mqpro

import (
  "github.com/caarlos0/env/v6"
)

// configuration for ibm mq
type envCfg struct {
  MQ0Host          string `env:"ENV_MQ_0_HOST"`
  MQ0Port          int    `env:"ENV_MQ_0_PORT"`
  MQ0Mgr           string `env:"ENV_MQ_0_MGR"`
  MQ0Channel       string `env:"ENV_MQ_0_CHANNEL"`
  MQ0PutQueue      string `env:"ENV_MQ_0_PUT_QUEUE"`
  MQ0GetQueue      string `env:"ENV_MQ_0_GET_QUEUE"`
  MQ0BrowseQ       string `env:"ENV_MQ_0_BROWSE_QUEUE"`
  MQ0Header        string `env:"ENV_MQ_0_HEADER"` // тип заголовков (свойств) rfh2
  MQ0App           string `env:"ENV_MQ_0_APP"`
  MQ0User          string `env:"ENV_MQ_0_USER"`
  MQ0Pass          string `env:"ENV_MQ_0_PASS"`
  MQ0Priority      string `env:"ENV_MQ_0_PRIORITY"`
  MQ0Tls           bool   `env:"ENV_MQ_0_TLS"`
  MQ0KeyRepository string `env:"ENV_MQ_0_KEY_REPOSITORY"`
  MQ0MaxMsgLength  int32  `env:"ENV_MQ_0_MAX_MSG_LENGTH"`
  MQ0WaitInterval  int32  `env:"ENV_MQ_0_WAIT_INTERVAL"`
  MQ0DevMode       bool   `env:"ENV_MQ_0_DEV_MODE"`

  MQ1Host          string `env:"ENV_MQ_1_HOST"`
  MQ1Port          int    `env:"ENV_MQ_1_PORT"`
  MQ1Mgr           string `env:"ENV_MQ_1_MGR"`
  MQ1Channel       string `env:"ENV_MQ_1_CHANNEL"`
  MQ1PutQueue      string `env:"ENV_MQ_1_PUT_QUEUE"`
  MQ1GetQueue      string `env:"ENV_MQ_1_GET_QUEUE"`
  MQ1BrowseQ       string `env:"ENV_MQ_1_BROWSE_QUEUE"`
  MQ1Header        string `env:"ENV_MQ_1_HEADER"`
  MQ1App           string `env:"ENV_MQ_1_APP"`
  MQ1User          string `env:"ENV_MQ_1_USER"`
  MQ1Pass          string `env:"ENV_MQ_1_PASS"`
  MQ1Priority      string `env:"ENV_MQ_1_PRIORITY"`
  MQ1Tls           bool   `env:"ENV_MQ_1_TLS"`
  MQ1KeyRepository string `env:"ENV_MQ_1_KEY_REPOSITORY"`
  MQ1MaxMsgLength  int32  `env:"ENV_MQ_1_MAX_MSG_LENGTH"`
  MQ1WaitInterval  int32  `env:"ENV_MQ_1_WAIT_INTERVAL"`
  MQ1DevMode       bool   `env:"ENV_MQ_1_DEV_MODE"`
}

// UseDefEnv для настройки использовать стандартные названия переменных
func (p *Mqpro) UseDefEnv() {
  p.SetConn(p.getConnFromEnv()...)
}

func (p *Mqpro) getConnFromEnv() []*Mqconn {
  var cfg = &envCfg{}

  if err := env.Parse(cfg); err != nil {
    p.log.Error(err)
  }

  p.log.Debugf("mqpro configuration: %+v", *cfg)

  connLi := make([]*Mqconn, 0)

  if cfg.MQ0Host != "" {
    if cfg.MQ0PutQueue != "" {
      conn := NewMqconn(TypePut, p.log, &Cfg{
        Host:          cfg.MQ0Host,
        Port:          cfg.MQ0Port,
        MgrName:       cfg.MQ0Mgr,
        ChannelName:   cfg.MQ0Channel,
        QueueName:     cfg.MQ0PutQueue,
        Header:        cfg.MQ0Header,
        AppName:       cfg.MQ0App,
        User:          cfg.MQ0User,
        Pass:          cfg.MQ0Pass,
        Priority:      cfg.MQ0Priority,
        Tls:           cfg.MQ0Tls,
        KeyRepository: cfg.MQ0KeyRepository,
        MaxMsgLength:  cfg.MQ0MaxMsgLength,
        DevMode:       cfg.MQ0DevMode,
      })
      connLi = append(connLi, conn)
    }

    if cfg.MQ0GetQueue != "" {
      conn := NewMqconn(TypeGet, p.log, &Cfg{
        Host:          cfg.MQ0Host,
        Port:          cfg.MQ0Port,
        MgrName:       cfg.MQ0Mgr,
        ChannelName:   cfg.MQ0Channel,
        Header:        cfg.MQ0Header,
        QueueName:     cfg.MQ0GetQueue,
        AppName:       cfg.MQ0App,
        User:          cfg.MQ0User,
        Pass:          cfg.MQ0Pass,
        Priority:      cfg.MQ0Priority,
        Tls:           cfg.MQ0Tls,
        KeyRepository: cfg.MQ0KeyRepository,
        MaxMsgLength:  cfg.MQ0MaxMsgLength,
        DevMode:       cfg.MQ0DevMode,
      })
      connLi = append(connLi, conn)
    }

    if cfg.MQ0BrowseQ != "" {
      conn := NewMqconn(TypeBrowse, p.log, &Cfg{
        Host:          cfg.MQ0Host,
        Port:          cfg.MQ0Port,
        MgrName:       cfg.MQ0Mgr,
        ChannelName:   cfg.MQ0Channel,
        QueueName:     cfg.MQ0BrowseQ,
        Header:        cfg.MQ0Header,
        AppName:       cfg.MQ0App,
        User:          cfg.MQ0User,
        Pass:          cfg.MQ0Pass,
        Priority:      cfg.MQ0Priority,
        Tls:           cfg.MQ0Tls,
        KeyRepository: cfg.MQ0KeyRepository,
        MaxMsgLength:  cfg.MQ0MaxMsgLength,
        DevMode:       cfg.MQ0DevMode,
      })
      connLi = append(connLi, conn)
    }
  }

  if cfg.MQ1Host != "" {
    if cfg.MQ1PutQueue != "" {
      conn := NewMqconn(TypePut, p.log, &Cfg{
        Host:          cfg.MQ1Host,
        Port:          cfg.MQ1Port,
        MgrName:       cfg.MQ1Mgr,
        ChannelName:   cfg.MQ1Channel,
        QueueName:     cfg.MQ1PutQueue,
        Header:        cfg.MQ1Header,
        AppName:       cfg.MQ1App,
        User:          cfg.MQ1User,
        Pass:          cfg.MQ1Pass,
        Priority:      cfg.MQ1Priority,
        Tls:           cfg.MQ1Tls,
        KeyRepository: cfg.MQ1KeyRepository,
        MaxMsgLength:  cfg.MQ1MaxMsgLength,
        DevMode:       cfg.MQ0DevMode,
      })
      connLi = append(connLi, conn)
    }

    if cfg.MQ1GetQueue != "" {
      conn := NewMqconn(TypeGet, p.log, &Cfg{
        Host:          cfg.MQ1Host,
        Port:          cfg.MQ1Port,
        MgrName:       cfg.MQ1Mgr,
        ChannelName:   cfg.MQ1Channel,
        QueueName:     cfg.MQ1GetQueue,
        Header:        cfg.MQ1Header,
        AppName:       cfg.MQ1App,
        User:          cfg.MQ1User,
        Pass:          cfg.MQ1Pass,
        Priority:      cfg.MQ1Priority,
        Tls:           cfg.MQ1Tls,
        KeyRepository: cfg.MQ1KeyRepository,
        MaxMsgLength:  cfg.MQ1MaxMsgLength,
        DevMode:       cfg.MQ0DevMode,
      })
      connLi = append(connLi, conn)
    }

    if cfg.MQ1BrowseQ != "" {
      conn := NewMqconn(TypeBrowse, p.log, &Cfg{
        Host:          cfg.MQ1Host,
        Port:          cfg.MQ1Port,
        MgrName:       cfg.MQ1Mgr,
        ChannelName:   cfg.MQ1Channel,
        QueueName:     cfg.MQ1BrowseQ,
        Header:        cfg.MQ1Header,
        AppName:       cfg.MQ1App,
        User:          cfg.MQ1User,
        Pass:          cfg.MQ1Pass,
        Priority:      cfg.MQ1Priority,
        Tls:           cfg.MQ1Tls,
        KeyRepository: cfg.MQ1KeyRepository,
        MaxMsgLength:  cfg.MQ1MaxMsgLength,
        DevMode:       cfg.MQ0DevMode,
      })
      connLi = append(connLi, conn)
    }
  }

  return connLi
}
