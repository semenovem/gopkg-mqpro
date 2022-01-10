package mqpro

import (
  "errors"
  "time"
)

const (
  defDisconnDelay  = time.Millisecond * 100 // По умолчанию задержка перед разрывом соединения
  defReconnDelay   = time.Second * 3        // По умолчанию задержка повторных попыток соединения
)

var (
  ErrNoEstablishedConnection = errors.New("ibm mq: no established connections")
  ErrInvalidConfig  = errors.New("ibm mq: не валидная конфигурация")
  ErrNoConfig         = errors.New("ibm mq: no configuration")
  ErrAlreadyConnected = errors.New("ibm mq: connection already established")
)

type state int32

const (
  stateDisconn state = iota
  stateConn
  stateConnecting
  stateErr
)

var stateKey = map[state]string{
  stateDisconn:    "stateDisconn",
  stateConn:       "stateConn",
  stateConnecting: "stateConnecting",
  stateErr:        "stateErr",
}
