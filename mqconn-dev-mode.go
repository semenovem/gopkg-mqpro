package mqpro

import "fmt"

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
func logMsg(msg *Msg, b []byte, l int, n string) {
  fmt.Println("\n--------------------------------")
  fmt.Printf("[DevMode] Сообщение: %s\n", n)
  fmt.Printf("[DevMode] original len    = %d\n", l)
  fmt.Printf("[DevMode] original byte   = %+v\n", b)
  fmt.Printf("[DevMode] original string = %s\n", b)
  fmt.Printf("[DevMode] Payload str = %s\n", msg.Payload)
  fmt.Printf("[DevMode] Props       = %+v\n", msg.Props)
  fmt.Printf("[DevMode] CorrelId    = %x\n", msg.CorrelId)
  fmt.Printf("[DevMode] MsgId       = %x\n", msg.MsgId)
  fmt.Printf("[DevMode] Time        = %s\n", msg.Time)

  if len(msg.MQRFH2) == 0 {
    fmt.Printf("[DevMode] MQRFH2      = %x\n", msg.MQRFH2)
  } else {
    for i, h := range msg.MQRFH2 {
      fmt.Printf("[DevMode] MQRFH2[%d].StructId       = %s\n", i, h.StructId)
      fmt.Printf("[DevMode] MQRFH2[%d].Version        = %d\n", i, h.Version)
      fmt.Printf("[DevMode] MQRFH2[%d].StructLength   = %d\n", i, h.StructLength)
      fmt.Printf("[DevMode] MQRFH2[%d].Encoding       = %d\n", i, h.Encoding)
      fmt.Printf("[DevMode] MQRFH2[%d].CodedCharSetId = %d\n", i, h.CodedCharSetId)
      fmt.Printf("[DevMode] MQRFH2[%d].Format         = %s\n", i, h.Format)
      fmt.Printf("[DevMode] MQRFH2[%d].Flags          = %d\n", i, h.Flags)
      fmt.Printf("[DevMode] MQRFH2[%d].NameValueCCSID = %d\n", i, h.NameValueCCSID)
      fmt.Printf("[DevMode] MQRFH2[%d].NameValues     = %+v\n", i, h.NameValues)

      for ii, raw := range h.RawXml {
        fmt.Printf("[DevMode] MQRFH2[%d].RawXml[%d] = %+v\n", i, ii, raw)
        fmt.Printf("[DevMode] MQRFH2[%d].RawXml[%d] = %s\n", i, ii, raw)
      }
    }
  }
}
