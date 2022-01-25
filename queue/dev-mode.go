package queue

import (
  "bytes"
  "fmt"
)

const devModeMaxBufferSize = 100
const devModeMaxStrSize = 300

// Вывод информации о ibmmq сообщении
func logMsg(m *Msg, origBuf []byte, n string) {
  var (
    buf    = bytes.NewBufferString("")
    extend = false
    f      = func(s string, i ...interface{}) { buf.WriteString(fmt.Sprintf(s, i...)) }
    f2     = func(b []byte, l int) []byte {
      if len(b) > l {
        return b[:l]
      }
      return b
    }
  )

  f("\n--------------------------------\n")
  f("[MQPRO-DevMode] Сообщение: %s\n", n)

  if origBuf != nil {
    f("[MQPRO-DevMode] origin(len)  = %d\n", len(origBuf))
    f("[MQPRO-DevMode] origin(str)  = %s\n", f2(origBuf, devModeMaxStrSize))
    if extend {
      f("[MQPRO-DevMode] origin(byt)  = %+v\n", f2(origBuf, devModeMaxBufferSize))
    }
  }

  f("[MQPRO-DevMode] Payload(len) = %d\n", len(m.Payload))
  f("[MQPRO-DevMode] Payload(str) = %s\n", f2(m.Payload, devModeMaxStrSize))
  if extend {
    f("[MQPRO-DevMode] Payload(byt) = %v\n", f2(m.Payload, devModeMaxBufferSize))
  }

  f("[MQPRO-DevMode] Props        = %+v\n", m.Props)
  f("[MQPRO-DevMode] CorrelId     = %s\n", logMsgCorr(m.CorrelId))
  f("[MQPRO-DevMode] MsgId        = %x\n", m.MsgId)

  tt := ""
  if m.Time.Year() != 1 {
    tt = m.Time.String()
  }
  f("[MQPRO-DevMode] Time         = %s\n", tt)

  if len(m.MQRFH2) == 0 {
    f("[MQPRO-DevMode] MQRFH2  = []\n")
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
        f("[MQPRO-DevMode] MQRFH2[%d].RawXml[%d](str) = %s\n", i, ii, f2(raw, devModeMaxStrSize))
        if extend {
          f("[MQPRO-DevMode] MQRFH2[%d].RawXml[%d] byt = %+v\n", i, ii, f2(raw, devModeMaxBufferSize))
        }
      }
    }
  }
  fmt.Println(buf.String())
}

func logMsgCorr(b []byte) string {
  var v byte
  for _, v = range b {
    if v != 0 {
      return fmt.Sprintf("%x", b)
    }
  }
  return ""
}
