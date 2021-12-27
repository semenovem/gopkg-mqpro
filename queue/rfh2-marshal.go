package queue

import (
  "fmt"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "reflect"
  "strconv"
)

func (c *Conn) Rfh2Marshal(p map[string]interface{}) ([]byte, error) {
  var bufPros []byte

  if p != nil && len(p) > 0 {
    var err error

    if c.rfh2RootTag != "" {
      p = map[string]interface{}{
        c.rfh2RootTag: p,
      }
    }

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

    b, err := toXml(n, v)
    if err != nil {
      return nil, err
    }

    b = append(b, make([]byte, tailFour(len(b)))...)
    buf = append(buf, 0, 0, 0, 0)
    endian.PutUint32(buf[ofs:], uint32(len(b)))

    buf = append(buf, b...)
  }

  return buf, nil
}

func toXml(k string, v interface{}) ([]byte, error) {
  var b []byte

  switch reflect.TypeOf(v).Kind() {
  case reflect.Ptr:
    bb, err := toXml(k, reflect.Indirect(reflect.ValueOf(v)).Interface())
    if err != nil {
      return nil, err
    }
    return bb, nil

  case reflect.Bool:
    b = []byte(strconv.FormatBool(v.(bool)))
  case reflect.Int:
    b = []byte(strconv.FormatInt(int64(v.(int)), 10))
  case reflect.Int8:
    b = []byte(strconv.FormatUint(uint64(v.(int8)), 10))
  case reflect.Int16:
    b = []byte(strconv.FormatUint(uint64(v.(int16)), 10))
  case reflect.Int32:
    b = []byte(strconv.FormatUint(uint64(v.(int32)), 10))
  case reflect.Int64:
    b = []byte(strconv.FormatUint(uint64(v.(int64)), 10))
  case reflect.Uint:
    b = []byte(strconv.FormatUint(uint64(v.(uint)), 10))
  case reflect.Uint8:
    b = []byte(strconv.FormatUint(uint64(v.(uint8)), 10))
  case reflect.Uint16:
    b = []byte(strconv.FormatUint(uint64(v.(uint16)), 10))
  case reflect.Uint32:
    b = []byte(strconv.FormatUint(uint64(v.(uint32)), 10))
  case reflect.Uint64:
    b = []byte(strconv.FormatUint(v.(uint64), 10))
  case reflect.Float32:
    b = []byte(strconv.FormatFloat(float64(v.(float32)), 'f', -1, 32))
  case reflect.Float64:
    b = []byte(strconv.FormatFloat(v.(float64), 'f', -1, 64))
  case reflect.Complex64:
    // TODO
    fmt.Printf("In progress")
  case reflect.Complex128:
    // TODO
    fmt.Printf("In progress")
  case reflect.String:
    b = []byte(v.(string))

  case reflect.Array, reflect.Slice:
    a := reflect.ValueOf(v)
    l := a.Len()
    for i := 0; i < l; i++ {
      xml, err := toXml(k, a.Index(i).Interface())
      if err != nil {
        return nil, err
      }
      b = append(b, xml...)
    }

    return b, nil

  case reflect.Map:
    m := reflect.ValueOf(v)
    iter := m.MapRange()

    for iter.Next() {
      key := iter.Key().Interface()
      if reflect.TypeOf(key).Kind() != reflect.String {
        return nil, ErrXml
      }

      xml, err := toXml(key.(string), iter.Value().Interface())
      if err != nil {
        return nil, err
      }

      b = append(b, xml...)
    }
  default:
    return nil, ErrXmlInconvertible
  }

  res := []byte("<" + k + ">")
  res = append(res, b...)
  res = append(res, []byte("</"+k+">")...)

  return res, nil
}
