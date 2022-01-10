package queue

import (
  "errors"
  "strings"
  "time"
)

const (
  defDisconnDelay    = time.Millisecond * 500 // По умолчанию задержка перед разрывом соединения
  defReconnDelay     = time.Second * 3        // По умолчанию задержка повторных попыток соединения
  DefRootTagHeader   = "usr"                  // Корневой тэг для заголовков rhf2 формата
  defReconnectDelay  = time.Second * 3
  defDelayClose      = time.Millisecond * 50 // Ожидание при закрытии очереди
  DefHeader          = HeaderBase
)

var (
  ErrClosed             = errors.New("ibm mq: queue is closed")
  ErrNoConnection       = errors.New("ibm mq: no connections")
  ErrInterrupted       = errors.New("ibm mq: operation interrupted")
  ErrAlreadyOpen        = errors.New("ibm mq: queue already open")
  ErrConnBroken         = errors.New("ibm mq conn: connection broken")
  ErrPutMsg             = errors.New("ibm mq: failed to put message")
  ErrGetMsg             = errors.New("ibm mq: failed to get message")
  ErrBrowseMsg          = errors.New("ibm mq: failed to browse message")
  ErrPropsNoField       = errors.New("ibm mq: property is missing")
  errMsgNoField         = "ibm mq: property '%s' is missing"
  errMsgFieldTypeTxt    = "ibm mq: invalid field type '%s'. Got '%T'"
  errHeaderParseType    = errors.New("ibm mq: header type parsing error")
  ErrFormatRFH2         = errors.New("ibm mq rfh2: error decoding header")
  ErrParseRfh2          = errors.New("ibm mq rfh2: error parse value")
  ErrRegisterEventInMsg = errors.New("ibm mq: the handler is already assigned")
  ErrXml                = errors.New("ibm mq: error when converting headers to xml. " +
    "Permissible: 'map[string]interface{}'")
  ErrXmlInconvertible = errors.New("ibm mq: Non-convertible data format")
  ErrNoConfig         = errors.New("ibm mq: не задана конфигурация")
)

type state int32
type queueOper int
type Header int
type permQueue int32

const (
  HeaderBase Header = iota
  HeaderRfh2
)

var headerKey = map[string]Header{
  headerVal[HeaderBase]: HeaderBase,
  headerVal[HeaderRfh2]: HeaderRfh2,
}

var headerVal = map[Header]string{
  HeaderBase: "prop",
  HeaderRfh2: "rfh2",
}

func ParseHeader(n string) (Header, error) {
  h, ok := headerKey[strings.ToLower(n)]
  if ok {
    return h, nil
  }
  return 0, errHeaderParseType
}

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
  stateErr:        "stateConnecting",
}

var stateVal = map[string]state{
  stateKey[stateDisconn]:    stateDisconn,
  stateKey[stateConn]:       stateConn,
  stateKey[stateConnecting]: stateConnecting,
  stateKey[stateErr]:        stateErr,
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

var permKey = map[permQueue]string{
  permGet:    "get",
  permBrowse: "browse",
  permPut:    "put",
}

var permVal = map[string]permQueue{
  permKey[permGet]:    permGet,
  permKey[permBrowse]: permBrowse,
  permKey[permPut]:    permPut,
}
