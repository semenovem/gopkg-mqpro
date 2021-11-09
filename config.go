package mqpro

import (
  "github.com/pkg/errors"
  "gopkg.in/yaml.v3"
  "io/ioutil"
)

type TransportCfg struct {
  Alias  string `yaml:"alias"`
  Browse string `yaml:"browse"`
  Put    string `yaml:"put"`
  Get    string `yaml:"get"`
}

type Config struct {
  LogMode   string         `yaml:"logMode"`
  LogCli    bool           `yaml:"logCli"`
  LogLevel  string         `yaml:"logLevel"`
  DevMode   bool           `yaml:"devMode"`
  Host      string         `yaml:"host"`
  Port      int            `yaml:"port"`
  Manager   string         `yaml:"manager"`
  Channel   string         `yaml:"channel"`
  App       string         `yaml:"app"`
  User      string         `yaml:"user"`
  Pass      int            `yaml:"pass"`
  Header    string         `yaml:"header"`
  RootTag   string         `yaml:"rootTag"`
  Transport []TransportCfg `yaml:"transport"`
}

func (p *Mqpro) ParseConfig(f string) error {
  if f == "" {
    return ErrConfigPathEmpty
  }
  byt, err := ioutil.ReadFile(f)
  if err != nil {
    return errors.Wrap(err, "Error reading configuration file")
  }

  c := new(Config)
  err = yaml.Unmarshal(byt, c)
  if err != nil {
    return errors.Wrapf(err, "Configuration file parsing error '%s'", f)
  }

  p.Cfg(c)

  return nil
}

func (p *Mqpro) Cfg(c *Config) {
  p.cfg = c
}
