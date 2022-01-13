package mqpro

import (
  "errors"
  "time"
)

const (
  defDisconnDelay = time.Millisecond * 100 // Задержка перед разрывом соединения
)

var (
  ErrNoEstablishedConnection = errors.New("ibm mq: no established connections")
  ErrInvalidConfig           = errors.New("ibm mq: не валидная конфигурация")
  ErrNoConfig                = errors.New("ibm mq: no configuration")
  ErrAlreadyConnected        = errors.New("ibm mq: connection already established")
)
