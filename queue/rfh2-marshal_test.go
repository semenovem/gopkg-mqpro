package queue

import (
	"bytes"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testLog() *logrus.Entry {
	l := logrus.New()
	l.SetLevel(logrus.WarnLevel)
	return logrus.NewEntry(l)
}

func testRfh2Conn(tag string) *Queue {
	cfg := &CoreSet{
		Header:      headerVal[HeaderRfh2],
		Rfh2RootTag: tag,
	}

	c := New(context.Background(), testLog(), nil)
	c.Set(cfg)
	return c
}

func TestRfh2Marshal_Marshal(t *testing.T) {
	p := map[string]interface{}{
		"first":  "value_first",
		"second": 123,
	}
	c := testRfh2Conn("")

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

func TestRfh2Marshal_toXml(t *testing.T) {
	f := func(k string, v interface{}) string {
		b, _ := toXml(k, v)
		return string(b)
	}

	assert.Equal(t, "<tag>string</tag>", f("tag", "string"))
	assert.Equal(t, "<tag>true</tag>", f("tag", true))
	assert.Equal(t, "<tag>100</tag>", f("tag", int8(100)))
	assert.Equal(t, "<tag>100</tag>", f("tag", int16(100)))
	assert.Equal(t, "<tag>100</tag>", f("tag", int32(100)))
	assert.Equal(t, "<tag>100</tag>", f("tag", int64(100)))
	assert.Equal(t, "<tag>100</tag>", f("tag", uint(100)))
	assert.Equal(t, "<tag>100</tag>", f("tag", uint8(100)))
	assert.Equal(t, "<tag>100</tag>", f("tag", uint16(100)))
	assert.Equal(t, "<tag>100</tag>", f("tag", uint32(100)))
	assert.Equal(t, "<tag>100</tag>", f("tag", uint64(100)))
	assert.Equal(t, "<tag>100</tag>", f("tag", byte(100)))
	assert.Equal(t, "<tag>100</tag>", f("tag", rune(100)))
	assert.Equal(t, "<tag>100.111</tag>", f("tag", float32(100.111)))
	assert.Equal(t, "<tag>100.111555</tag>", f("tag", 100.111555))

	m := map[string]interface{}{
		"first": "val-first",
	}

	assert.Equal(t,
		"<tag><first>val-first</first></tag>",
		f("tag", m))

	_, err := toXml("tag", f)
	assert.EqualError(t, err, ErrXmlInconvertible.Error())

	m["first"] = map[string]string{
		"second": "val-second",
	}

	assert.Equal(t,
		"<tag><first><second>val-second</second></first></tag>",
		f("tag", m))

	li := []string{
		"valfirst", "valsecond",
	}

	assert.Equal(t,
		"<tag>valfirst</tag><tag>valsecond</tag>",
		f("tag", li))
}
