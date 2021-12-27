package mqpro

import (
  "errors"
  "strings"
  "time"
)

const (
  defDisconnDelay  = time.Millisecond * 500 // По умолчанию задержка перед разрывом соединения
  defReconnDelay   = time.Second * 3        // По умолчанию задержка повторных попыток соединения
  defRootTagHeader = "usr"                  // Корневой тэг для заголовков rhf2 формата
)

var (
  ErrNoEstablishedConnection = errors.New("ibm mq: no established connections")
  ErrNoConnection            = errors.New("ibm mq: no connections")
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
  ErrAliasExist         = errors.New("ibm mq: the queue already exists")
)

type TypeConn int
type stateConn int
type reqStateConn int
type queueOper int
type header int

const (
  HeaderBase header = iota
  HeaderRfh2
  HeaderBaseTxt = "prop"
  HeaderRfh2Txt = "rfh2"
)

var headerMap = map[string]header{
  HeaderBaseTxt: HeaderBase,
  HeaderRfh2Txt: HeaderRfh2,
}

func parseHeaderType(n string) (header, error) {
  h, ok := headerMap[strings.ToLower(n)]
  if ok {
    return h, nil
  }
  return 0, errHeaderParseType
}

const (
  TypePut TypeConn = iota + 1
  TypeGet
  TypeBrowse
  defReconnectDelay  = time.Second * 3
  defMaxMsgLength    = 1024 * 1024 * 100
  defMsgWaitInterval = time.Millisecond * 100
)

const (
  stateDisconnect stateConn = iota
  stateConnect
  stateErr
)

const (
  reqConnect reqStateConn = iota
  reqReconnect
  reqDisconnect
)

const (
  operGet queueOper = iota
  operGetByMsgId
  operGetByCorrelId
  operBrowseFirst
  operBrowseNext
)

var typeConnTxt = map[TypeConn]string{
  TypePut:    "put",
  TypeGet:    "get",
  TypeBrowse: "browse",
}
