package queue

import (
  "bytes"
  "encoding/binary"
  "encoding/xml"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "strings"
)

// Rfh2Unmarshal получение заголовков rfh2
func (c *Conn) Rfh2Unmarshal(b []byte) ([]*MQRFH2, error) {
  var (
    tot []*MQRFH2
    rfh *MQRFH2
    err error
    ofs int32
  )

  for ofs < int32(len(b)) {
    rfh, err = rfh2ParseHeader(b[ofs:], endian)
    if err != nil {
      if err != ErrFormatRFH2 {
        return nil, err
      }
      // Из-за того, что мы не знаем порядок кодирования байтов в rfh2 заголовке, пробуем оба варианта
      rfh, err = rfh2ParseHeader(b[ofs:], endian2)
      if err != nil {
        return nil, err
      }
    }

    if rfh == nil {
      break
    }
    tot = append(tot, rfh)
    ofs += rfh.StructLength
  }

  return tot, nil
}

func rfh2ParseHeader(b []byte, ord binary.ByteOrder) (*MQRFH2, error) {
  if len(b) < 4 {
    return nil, nil
  }

  if !bytes.Equal([]byte(StructId), b[:4]) {
    return nil, nil
  }

  if int32(len(b)) < ibmmq.MQRFH_STRUC_LENGTH_FIXED_2 {
    return nil, ErrFormatRFH2
  }

  h := &MQRFH2{}
  var err error

  h.StructId = string(b[:4])
  h.Version = int32(ord.Uint32(b[4:8]))
  h.StructLength = int32(ord.Uint32(b[8:12]))
  h.Encoding = int32(ord.Uint32(b[12:16]))
  h.CodedCharSetId = int32(ord.Uint32(b[16:20]))
  h.Format = strings.TrimRight(string(b[20:28]), " ")
  h.Flags = int32(ord.Uint32(b[28:32]))
  h.NameValueCCSID = int32(ord.Uint32(b[32:36]))

  if int32(len(b)) < h.StructLength {
    return nil, ErrFormatRFH2
  }
  err = rfh2ParseData(b[36:h.StructLength], h, ord)
  if err != nil {
    return nil, err
  }

  return h, nil
}

// Обработка пар NameValueLength NameValueData
// https://www.ibm.com/docs/en/ibm-mq/9.0?topic=mqrfh2-namevaluelength-mqlong
func rfh2ParseData(buf []byte, rfh *MQRFH2, ord binary.ByteOrder) error {
  ofs := 0

  for ofs+4 < len(buf) {
    l := int(ord.Uint32(buf[ofs : ofs+4]))
    ofs += 4

    if len(buf) < l+ofs {
      return ErrFormatRFH2
    }

    b := bytes.TrimRight(buf[ofs:l+ofs], "\x00")
    ofs += l
    m, err := rfh2ParseXml(b)
    if err != nil {
      return ErrParseRfh2
    }

    rfh.RawXml = append(rfh.RawXml, b)
    rfh.NameValues = append(rfh.NameValues, m)
  }

  return nil
}

func rfh2ParseXml(buf []byte) (map[string]interface{}, error) {
  m := rfh2Xml{}
  err := xml.Unmarshal(buf, &m)
  if err != nil {
    return nil, err
  }

  return m.m, nil
}

type rfh2Xml struct {
  m map[string]interface{}
}

type gg struct {
  root *gg
  a    map[string]interface{}
}

func (c *rfh2Xml) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
  root := make(map[string]interface{})
  root[start.Name.Local] = nil
  path := []string{start.Name.Local}

loop:
  for {
    t, _ := d.Token()

    key := path[len(path)-1]
    mnt := root

    for _, k := range path[:len(path)-1] {
      switch mnt[k].(type) {
      case map[string]interface{}:
        mnt = mnt[k].(map[string]interface{})
      }
    }

    switch tt := t.(type) {
    case xml.StartElement:
      //fmt.Printf(">>> StartElement: <%s>:<%s>\n", key, tt.Name.Local)

      switch mnt[key].(type) {
      case map[string]interface{}:
      default:
        m := make(map[string]interface{})
        mnt[key] = m
        m[tt.Name.Local] = nil
      }

      m := mnt[key].(map[string]interface{})
      m[tt.Name.Local] = nil

      path = append(path, tt.Name.Local)

    case xml.EndElement:
      //fmt.Printf(">>> EndElement:   <%s>\n", key)

      if tt.Name == start.Name {
        break loop
      }
      path = path[:len(path)-1]

    case xml.CharData:
      //fmt.Printf(">>> CharData:     <%s>  [%s] \n", key, tt.Copy())

      switch mnt[key].(type) {
      case map[string]interface{}:
      default:
        mnt[key] = string(tt.Copy())
      }
    }
  }

  c.m = root

  return nil
}
