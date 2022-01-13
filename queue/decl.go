package queue

import (
  "errors"
  "strings"
  "time"
)

const (
  defDisconnDelay   = time.Millisecond * 500 // По умолчанию задержка перед разрывом соединения
  defReconnDelay    = time.Second * 3        // По умолчанию задержка повторных попыток соединения
  DefRootTagHeader  = "usr"                  // Корневой тэг для заголовков rhf2 формата
  defReconnectDelay = time.Second * 3        // Тамаут повторных попыток
  defDelayClose     = time.Millisecond * 50  // Ожидание при закрытии очереди
  DefHeader         = HeaderBase
)

var (
  ErrNotOpen         = errors.New("ibm mq: queue is not open")
  ErrBusySubsc       = errors.New("ibm mq: the queue object is busy subscribing")
  ErrInterrupted     = errors.New("ibm mq: operation interrupted")
  ErrAlreadyOpen     = errors.New("ibm mq: queue already open")
  ErrConnBroken      = errors.New("ibm mq: conn: connection broken")
  ErrPutMsg          = errors.New("ibm mq: failed to put message")
  ErrGetMsg          = errors.New("ibm mq: failed to get message")
  errHeaderParseType = errors.New("ibm mq: header type parsing error")
  ErrFormatRFH2      = errors.New("ibm mq: rfh2: error decoding header")
  ErrParseRfh2       = errors.New("ibm mq: rfh2: error parse value")
  ErrRegisterInMsg   = errors.New("ibm mq: the handler is already assigned")
  ErrNotGetOpen      = errors.New("ibm mq: no rights to view messages")
  ErrXml             = errors.New("ibm mq: error when converting headers to xml. " +
    "Permissible: 'map[string]interface{}'")
  ErrXmlInconvertible     = errors.New("ibm mq: non-convertible data format")
  ErrNotConfigured        = errors.New("ibm mq: not configured")
  ErrManagerNotConfigured = errors.New("ibm mq: manager not configured")
)

const (
  msgErrPropCreation = "creating a message property: %s"
  msgErrPropDeletion = "deleting message properties: %s"
  msgErrPropGetting  = "getting message property: %s"
)

type state int32
type queueOper int
type Header int
type permQueue int32

const (
  HeaderBase Header = iota
  HeaderRfh2
)

var headerMapByVal = map[string]Header{
  HeaderMapByKey[HeaderBase]: HeaderBase,
  HeaderMapByKey[HeaderRfh2]: HeaderRfh2,
}

var HeaderMapByKey = map[Header]string{
  HeaderBase: "prop",
  HeaderRfh2: "rfh2",
}

func ParseHeader(n string) (Header, error) {
  h, ok := headerMapByVal[strings.ToLower(n)]
  if ok {
    return h, nil
  }
  return 0, errHeaderParseType
}

const (
  stateClosed state = iota
  stateOpen
  stateConnecting
  stateErr
)

var stateMapByKey = map[state]string{
  stateClosed:     "stateClosed",
  stateOpen:       "stateOpen",
  stateConnecting: "stateConnecting",
  stateErr:        "stateConnecting",
}

var stateMapByVal = map[string]state{
  stateMapByKey[stateClosed]:     stateClosed,
  stateMapByKey[stateOpen]:       stateOpen,
  stateMapByKey[stateConnecting]: stateConnecting,
  stateMapByKey[stateErr]:        stateErr,
}

const (
  operGet queueOper = iota
  operGetByMsgId
  operGetByCorrelId
  operBrowseFirst
  operBrowseNext
)

const (
  permGet permQueue = iota
  permBrowse
  permPut
)

var permMapByKey = map[permQueue]string{
  permGet:    "get",
  permBrowse: "browse",
  permPut:    "put",
}

var permMapByVal = map[string]permQueue{
  permMapByKey[permGet]:    permGet,
  permMapByKey[permBrowse]: permBrowse,
  permMapByKey[permPut]:    permPut,
}
