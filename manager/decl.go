package manager

import (
  "errors"
  "time"
)

const (
  defDisconnDelay = time.Millisecond * 100 // Задержка перед разрывом соединения
  defReconnDelay  = time.Second * 3        // Задержка повторных попыток соединения
)

var (
  ErrNoEstablishedConnection = errors.New("ibm mq: no established connections")
  ErrInvalidConfig           = errors.New("ibm mq: invalid configuration")
  ErrNotConfigured           = errors.New("ibm mq: not configured")
  ErrAlreadyConnected        = errors.New("ibm mq: connection already established")
)

type state int32

const (
  stateDisconn state = iota
  stateConn
  stateConnecting
  stateErr
)

var stateMapByKey = map[state]string{
  stateDisconn:    "stateDisconn",
  stateConn:       "stateConn",
  stateConnecting: "stateConnecting",
  stateErr:        "stateErr",
}
