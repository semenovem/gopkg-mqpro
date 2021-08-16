package mqpro

import (
  "fmt"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "strconv"
)

func (c *Mqconn) Rfh2Marshal(p map[string]interface{}) ([]byte, error) {
  var bufPros []byte

  if p != nil && len(p) > 0 {
    var err error
    bufPros, err = toHeadersBinary(p)
    if err != nil {
      return nil, ErrFormatRFH2
    }
  }

  l := ibmmq.MQRFH_STRUC_LENGTH_FIXED_2 + int32(len(bufPros))
  l += tailFour32(l)
  b := make([]byte, l)

  copy(b[:4], c.rfh2.StructId)
  endian.PutUint32(b[4:], uint32(c.rfh2.Version))
  endian.PutUint32(b[8:], uint32(l))
  endian.PutUint32(b[12:], uint32(c.rfh2.Encoding))
  endian.PutUint32(b[16:], uint32(c.rfh2.CodedCharSetId))
  copy(b[20:28], []byte(c.rfh2.Format + space8)[:8])
  endian.PutUint32(b[28:], uint32(c.rfh2.Flags))
  endian.PutUint32(b[32:], uint32(c.rfh2.NameValueCCSID))

  if len(bufPros) > 0 {
    copy(b[36:], bufPros)
  }

  return b, nil
}

func toHeadersBinary(m map[string]interface{}) ([]byte, error) {
  buf := make([]byte, 0)
  var ofs int

  for n, v := range m {
    ofs = len(buf)

    b := make([]byte, 0)
    b = append(b, []byte("<"+n+">")...)
    b = append(b, toXmlVal(v)...)
    b = append(b, []byte("</"+n+">")...)
    b = append(b, make([]byte, tailFour(len(b)))...)

    buf = append(buf, 0, 0, 0, 0)
    endian.PutUint32(buf[ofs:], uint32(len(b)))

    buf = append(buf, b...)
  }

  return buf, nil
}

// TODO сделать поддержку вложенности, либо использовать xml либу
func toXmlVal(val interface{}) []byte {
  var b []byte

  switch v := val.(type) {
  case []byte:
    b = v
  case string:
    b = []byte(v)
  case bool:
    b = []byte(strconv.FormatBool(v))
  case uint:
    b = []byte(strconv.FormatUint(uint64(v), 10))
  case uint8:
    b = []byte(strconv.FormatUint(uint64(v), 10))
  case uint16:
    b = []byte(strconv.FormatUint(uint64(v), 10))
  case uint32:
    b = []byte(strconv.FormatUint(uint64(v), 10))
  case uint64:
    b = []byte(strconv.FormatUint(v, 10))
  case int:
    b = []byte(strconv.FormatInt(int64(v), 10))
  case int8:
    b = []byte(strconv.FormatInt(int64(v), 10))
  case int16:
    b = []byte(strconv.FormatInt(int64(v), 10))
  case int32:
    b = []byte(strconv.FormatInt(int64(v), 10))
  case int64:
    b = []byte(strconv.FormatInt(v, 10))
  case float32:
    b = []byte(strconv.FormatFloat(float64(v), 'f', -1, 32))
  case float64:
    b = []byte(strconv.FormatFloat(v, 'f', -1, 64))
  default:
    b = []byte(fmt.Sprintf("%v", val))
  }

  return b
}
