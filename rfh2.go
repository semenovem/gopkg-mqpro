package mqpro

import (
  "encoding/binary"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

const (
  StructId     = "RFH "
  EncodingUTF8 = 1208
  space8       = "        "
)

var endian = binary.LittleEndian

type rfh2Cfg struct {
  StructId       string
  Version        int32
  Encoding       int32
  CodedCharSetId int32
  Format         string
  Flags          int32
  NameValueCCSID int32
}

func newRfh2Cfg() *rfh2Cfg {
  return &rfh2Cfg{
    StructId:       StructId,
    Version:        ibmmq.MQRFH_VERSION_2,
    Encoding:       ibmmq.MQENC_NATIVE,
    CodedCharSetId: ibmmq.MQCCSI_INHERIT,
    Format:         ibmmq.MQFMT_STRING,
    Flags:          ibmmq.MQRFH_NONE,
    NameValueCCSID: EncodingUTF8,
  }
}

type MQRFH2 struct {
  StructId       string
  Version        int32
  StrucLength    int32
  Encoding       int32
  CodedCharSetId int32
  Format         string
  Flags          int32
  NameValueCCSID int32
  NameValues     []map[string]interface{}
  RawXml         [][]byte // Не обработанное содержимое NameValueData
}
