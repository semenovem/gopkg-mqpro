package queue

import (
  "bytes"
  "encoding/xml"
  "fmt"
  "github.com/stretchr/testify/assert"
  "testing"
)

func TestRfh2Unmarshal_parsePayload(t *testing.T) {
  intOrd := endian
  txt1 := "<first>value_first3</first>"
  tag1 := []byte(txt1)
  l1 := len(tag1) + tailFour(len(tag1))
  b1 := make([]byte, 4+l1)
  intOrd.PutUint32(b1[:4], uint32(l1))
  copy(b1[4:], tag1)

  txt2 := "<second11>value_second</second11>"
  tag2 := []byte(txt2)
  l2 := len(tag2) + tailFour(len(tag2))
  b2 := make([]byte, 4+l2)
  intOrd.PutUint32(b2[:4], uint32(l2))
  copy(b2[4:], tag2)
  buf := append([]byte{}, b1...)
  buf = append(buf, b2...)

  h := &MQRFH2{}
  err := rfh2ParseData(buf, h, intOrd)
  assert.NoError(t, err)

  assert.Len(t, h.NameValues, 2)
  assert.Len(t, h.RawXml, 2)

  assert.True(t, bytes.Equal(tag1, h.RawXml[0]))
  assert.True(t, bytes.Equal(tag2, h.RawXml[1]))

  assert.Equal(t, "map[first:value_first3]", fmt.Sprintf("%v", h.NameValues[0]))
  assert.Equal(t, "map[second11:value_second]", fmt.Sprintf("%v", h.NameValues[1]))
}

func TestRfh2Unmarshal_parseNameValue(t *testing.T) {
  b := []byte("<first>value_first3</first>")
  b = append(b, make([]byte, len(b)%4)...)

  m, err := rfh2ParseXml(b)
  assert.NoError(t, err)
  assert.Equal(t, "map[first:value_first3]", fmt.Sprintf("%v", m))
}

func TestRfh2Unmarshal_rfh2Xml(t *testing.T) {
  b := []byte("<usr><first>value_first1</first>  <second>34523<third>333</third></second></usr>")

  m := rfh2Xml{}
  err := xml.Unmarshal(b, &m)
  assert.NoError(t, err)

  assert.Equal(t, "map[usr:map[first:value_first1 second:map[third:333]]]",
    fmt.Sprintf("%v", m.m))
}
