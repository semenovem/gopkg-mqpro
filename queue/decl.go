package queue

import (
  "errors"
  "strings"
  "time"
)

const (
  defDisconnDelay    = time.Millisecond * 500 // По умолчанию задержка перед разрывом соединения
  defReconnDelay     = time.Second * 3        // По умолчанию задержка повторных попыток соединения
  defRootTagHeader   = "usr"                  // Корневой тэг для заголовков rhf2 формата
  defReconnectDelay  = time.Second * 3
  defMaxMsgLength    = 1024 * 1024 * 100
  defMsgWaitInterval = time.Millisecond * 100
  defHeader          = headerBase
)

var (
  ErrNoEstablishedConnection = errors.New("ibm mq: no established connections")
  ErrNoConnection            = errors.New("ibm mq: no connections")
  ErrHasConnection           = errors.New("ibm mq: connection is already established")
  ErrNoData                  = errors.New("ibm mq: no data to connect to IBM MQ")
  ErrConnBroken              = errors.New("ibm mq conn: connection broken")
  ErrPutMsg                  = errors.New("ibm mq: failed to put message")
  ErrGetMsg                  = errors.New("ibm mq: failed to get message")
  ErrBrowseMsg               = errors.New("ibm mq: failed to browse message")
  ErrPropsNoField            = errors.New("ibm mq: property is missing")
  errMsgNoField              = "ibm mq: property '%s' is missing"
  errMsgFieldTypeTxt         = "ibm mq: invalid field type '%s'. Got '%T'"
  errHeaderParseType         = errors.New("ibm mq: header type parsing error")
  ErrFormatRFH2              = errors.New("ibm mq rfh2: error decoding header")
  ErrParseRfh2               = errors.New("ibm mq rfh2: error parse value")
  ErrRegisterEventInMsg      = errors.New("ibm mq: inbound message handler registration error")
  ErrXml                     = errors.New("ibm mq: error when converting headers to xml. " +
    "Permissible: 'map[string]interface{}'")
  ErrXmlInconvertible = errors.New("ibm mq: Non-convertible data format")
  ErrConfigPathEmpty  = errors.New("ibm mq: configuration file path not specified")
  ErrNoConfig         = errors.New("ibm mq: no configuration")
)

type state int
type queueOper int
type header int

const (
  headerBase header = iota
  headerRfh2
)

var headerKey = map[string]header{
  headerVal[headerBase]: headerBase,
  headerVal[headerRfh2]: headerRfh2,
}

var headerVal = map[header]string{
  headerBase: "prop",
  headerRfh2: "rfh2",
}

func parseHeaderType(n string) (header, error) {
  h, ok := headerKey[strings.ToLower(n)]
  if ok {
    return h, nil
  }
  return 0, errHeaderParseType
}

const (
  stateDisconn state = iota
  stateConn
  stateErr
)

const (
  operGet queueOper = iota
  operGetByMsgId
  operGetByCorrelId
  operBrowseFirst
  operBrowseNext
)
