package mqpro

import (
  "bytes"
  "github.com/sirupsen/logrus"
  "github.com/stretchr/testify/assert"
  "testing"
)

func testLog() *logrus.Entry {
  l := logrus.New()
  l.SetLevel(logrus.WarnLevel)
  return logrus.NewEntry(l)
}

func testRfh2Conn() *Mqconn {
  cfg := &Cfg{
    Header: HeaderRfh2Txt,
  }
  return NewMqconn(TypePut, testLog(), cfg)
}

func TestRfh2Marshal_Marshal(t *testing.T) {
  p := map[string]interface{}{
    "first":  "value_first",
    "second": 123,
  }
  c := testRfh2Conn()

  b, err := c.Rfh2Marshal(p)
  assert.NoError(t, err)

  headers, err := c.Rfh2Unmarshal(b)
  assert.NoError(t, err)

  assert.Len(t, headers, 1)
  h := headers[0]

  assert.Len(t, h.NameValues, 2)
  assert.Len(t, h.RawXml, 2)
}

func TestRfh2Marshal_toHeadersBinary(t *testing.T) {
  h := map[string]interface{}{
    "first":  "value_first",
    "second": 123,
  }

  // https://www.ibm.com/docs/en/ibm-mq/9.0?topic=mqrfh2-namevaluelength-mqlong
  tag1 := []byte("<first>value_first</first>")
  l1 := len(tag1) + tailFour(len(tag1))
  b1 := make([]byte, 4+l1)
  endian.PutUint32(b1[:4], uint32(l1))
  copy(b1[4:], tag1)
  tag2 := []byte("<second>123</second>")
  l2 := len(tag2) + tailFour(len(tag2))
  b2 := make([]byte, 4+l2)
  endian.PutUint32(b2[:4], uint32(l2))
  copy(b2[4:], tag2)
  buf := append([]byte{}, b1...)
  buf = append(buf, b2...)

  b, err := toHeadersBinary(h)
  assert.NoError(t, err)
  assert.True(t, bytes.Equal(buf, b))
}

func TestRfh2Marshal_cnvVal(t *testing.T) {
  f := func(s string, b[]byte) bool {
    return bytes.Equal([]byte(s), b)
  }

  assert.True(t, f("string", toXmlVal("string")))
  assert.True(t, f("[]byte", toXmlVal([]byte("[]byte"))))
  assert.True(t, f("true", toXmlVal(true)))

  assert.True(t, f("3232", toXmlVal(uint(3232))))
  assert.True(t, f("8", toXmlVal(uint8(8))))
  assert.True(t, f("16", toXmlVal(uint16(16))))
  assert.True(t, f("32", toXmlVal(uint32(32))))
  assert.True(t, f("64", toXmlVal(uint64(64))))

  assert.True(t, f("3232", toXmlVal(int(3232))))
  assert.True(t, f("8", toXmlVal(int8(8))))
  assert.True(t, f("16", toXmlVal(int16(16))))
  assert.True(t, f("32", toXmlVal(int32(32))))
  assert.True(t, f("64", toXmlVal(int64(64))))

  assert.True(t, f("234.345", toXmlVal(float32(234.345))))
  assert.True(t, f("234.3456666", toXmlVal(234.3456666)))
}
