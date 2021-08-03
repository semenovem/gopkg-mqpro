package mqpro

import (
  "bytes"
  "encoding/binary"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "io"
  "strconv"
)

const (
  StrucId      = "RFH "
  EncodingUTF8 = 1208
)

type MQRFH2 struct {
  StrucId        string
  Version        int32
  StrucLength    int32
  Encoding       int32
  CodedCharSetId int32
  Format         string
  Flags          int32
  NameValueCCSID int32
  NameValues     H
}

type H map[string]interface{}

func NewMQRFH2() *MQRFH2 {
  // cmqc.h:5021
  return &MQRFH2{
    StrucId:        StrucId,
    Version:        ibmmq.MQRFH_VERSION_2,
    StrucLength:    ibmmq.MQRFH_STRUC_LENGTH_FIXED_2,
    Encoding:       ibmmq.MQENC_NATIVE,
    CodedCharSetId: ibmmq.MQCCSI_INHERIT,
    Format:         ibmmq.MQFMT_NONE,
    Flags:          ibmmq.MQRFH_NONE,
    NameValueCCSID: EncodingUTF8,
    NameValues:     make(H),
  }
}

func (hdr *MQRFH2) MarshalBinary() ([]byte, error) {
  const space8 = "        "
  endian := hdr.getEndian()
  buf := &bytes.Buffer{}
  nameValues := &bytes.Buffer{}
  if hdr.NameValues != nil && len(hdr.NameValues) > 0 {
    nameValues = bytes.NewBuffer(hdr.NameValues.Bytes(endian))
  }
  _, _ = buf.Write([]byte(hdr.StrucId))
  _ = binary.Write(buf, endian, uint32(hdr.Version))
  _ = binary.Write(buf, endian, uint32(hdr.StrucLength+int32(nameValues.Len())))
  _ = binary.Write(buf, endian, uint32(hdr.Encoding))
  _ = binary.Write(buf, endian, uint32(hdr.CodedCharSetId))
  _, _ = buf.Write([]byte((hdr.Format + space8)[0:8]))
  _ = binary.Write(buf, endian, uint32(hdr.Flags))
  _ = binary.Write(buf, endian, uint32(hdr.NameValueCCSID))
  _, _ = io.Copy(buf, nameValues)
  return buf.Bytes(), nil
}

func HeaderOffset(payload []byte) int {
  offset := uint32(0)
  var endian binary.ByteOrder
  if ibmmq.MQENC_NATIVE%2 == 0 {
    endian = binary.LittleEndian
  } else {
    endian = binary.BigEndian
  }
  for string(payload[offset:4]) == StrucId {
    offset += endian.Uint32(payload[offset+8 : offset+12])
  }
  return int(offset)
}

func (hdr *MQRFH2) getEndian() binary.ByteOrder {
  if hdr.Encoding%2 == 0 {
    return binary.LittleEndian
  } else {
    return binary.BigEndian
  }
}

func (h H) Bytes(order binary.ByteOrder) []byte {
  w := &bytes.Buffer{}
  buf := &bytes.Buffer{}
  for folderName, folder := range h {
    buf.Reset()
    buf.Write([]byte("<" + folderName + ">"))
    h.writeFolder(folder, buf)
    buf.Write([]byte("</" + folderName + ">"))
    _ = binary.Write(w, order, uint32(buf.Len()))
    _, _ = io.Copy(w, buf)
  }
  return w.Bytes()
}

func (h H) writeFolder(folder interface{}, out io.Writer) {
  switch v := folder.(type) {
  case bool:
    _, _ = out.Write([]byte(strconv.FormatBool(v)))
  case uint:
    _, _ = out.Write([]byte(strconv.FormatUint(uint64(v), 10)))
  case uint8:
    _, _ = out.Write([]byte(strconv.FormatUint(uint64(v), 10)))
  case uint16:
    _, _ = out.Write([]byte(strconv.FormatUint(uint64(v), 10)))
  case uint32:
    _, _ = out.Write([]byte(strconv.FormatUint(uint64(v), 10)))
  case uint64:
    _, _ = out.Write([]byte(strconv.FormatUint(v, 10)))
  case int:
    _, _ = out.Write([]byte(strconv.FormatInt(int64(v), 10)))
  case int8:
    _, _ = out.Write([]byte(strconv.FormatInt(int64(v), 10)))
  case int16:
    _, _ = out.Write([]byte(strconv.FormatInt(int64(v), 10)))
  case int32:
    _, _ = out.Write([]byte(strconv.FormatInt(int64(v), 10)))
  case int64:
    _, _ = out.Write([]byte(strconv.FormatInt(v, 10)))
  case float32:
    _, _ = out.Write([]byte(strconv.FormatFloat(float64(v), 'b', -1, 10)))
  case float64:
    _, _ = out.Write([]byte(strconv.FormatFloat(v, 'b', -1, 10)))
  case map[string]interface{}:
    for name, value := range v {
      _, _ = out.Write([]byte("<" + name + ">"))
      h.writeFolder(value, out)
      _, _ = out.Write([]byte("</" + name + ">"))
    }
  }
}
