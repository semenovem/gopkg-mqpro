package mqpro

import (
  "bytes"
  "fmt"
)

const devModeMaxBufferSize = 200

func devMode(m *Msg, b []byte, n string) func(...[]*Msg) {
  var (
    byt []byte
  )
  l := len(b)

  if len(b) < devModeMaxBufferSize {
    byt = append(byt, b...)
  } else {
    byt = append(byt, b[:devModeMaxBufferSize]...)
  }

  return func(mm ...[]*Msg) {
    logMsg(m, byt, l, n)
  }
}

// Вывод информации о ibmmq сообщении
func logMsg(m *Msg, b []byte, l int, n string) {
  var buf = bytes.NewBufferString("")
  f := func(s string, i ...interface{}) {
    buf.WriteString(fmt.Sprintf(s, i...))
  }

  f("\n--------------------------------\n")
  f("[MQPRO-DevMode] Сообщение: %s\n", n)
  f("[MQPRO-DevMode] original len    = %d\n", l)
  f("[MQPRO-DevMode] original byte   = %+v\n", b)
  f("[MQPRO-DevMode] original string = %s\n", b)
  f("[MQPRO-DevMode] Payload str = %s\n", m.Payload)
  f("[MQPRO-DevMode] Props       = %+v\n", m.Props)
  f("[MQPRO-DevMode] CorrelId    = %x\n", m.CorrelId)
  f("[MQPRO-DevMode] MsgId       = %x\n", m.MsgId)
  f("[MQPRO-DevMode] Time        = %s\n", m.Time)

  if len(m.MQRFH2) == 0 {
    f("[MQPRO-DevMode] MQRFH2      = %x\n", m.MQRFH2)
  } else {
    for i, h := range m.MQRFH2 {
      f("[MQPRO-DevMode] MQRFH2[%d].StructId       = %s\n", i, h.StructId)
      f("[MQPRO-DevMode] MQRFH2[%d].Version        = %d\n", i, h.Version)
      f("[MQPRO-DevMode] MQRFH2[%d].StructLength   = %d\n", i, h.StructLength)
      f("[MQPRO-DevMode] MQRFH2[%d].Encoding       = %d\n", i, h.Encoding)
      f("[MQPRO-DevMode] MQRFH2[%d].CodedCharSetId = %d\n", i, h.CodedCharSetId)
      f("[MQPRO-DevMode] MQRFH2[%d].Format         = %s\n", i, h.Format)
      f("[MQPRO-DevMode] MQRFH2[%d].Flags          = %d\n", i, h.Flags)
      f("[MQPRO-DevMode] MQRFH2[%d].NameValueCCSID = %d\n", i, h.NameValueCCSID)
      f("[MQPRO-DevMode] MQRFH2[%d].NameValues     = %+v\n", i, h.NameValues)

      for ii, raw := range h.RawXml {
        f("[MQPRO-DevMode] MQRFH2[%d].RawXml[%d] byt = %+v\n", i, ii, raw)
        f("[MQPRO-DevMode] MQRFH2[%d].RawXml[%d] str = %s\n", i, ii, raw)
      }
    }
  }

  fmt.Println(buf.String())
}
