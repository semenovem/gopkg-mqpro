package mqpro

import (
  "fmt"
)

type tProp int

const (
  tUnknown tProp = iota
  tInt
  tInt8
  tInt16
  tInt32
  tInt64
  tUint
  tUint8
  tUint16
  tUint32
  tUint64
  tByte
  tRune
  tString
  tBool
  tFloat32
  tFloat64
  tComplex64
  tComplex128
)

func (m *Msg) Int(f string) (int, error) {
  v, err := m.getter(f, tInt)
  if err != nil {
    return 0, err
  }
  return v.(int), nil
}

func (m *Msg) Int8(f string) (int8, error) {
  v, err := m.getter(f, tInt8)
  if err != nil {
    return 0, err
  }
  return v.(int8), nil
}

func (m *Msg) Int16(f string) (int16, error) {
  v, err := m.getter(f, tInt16)
  if err != nil {
    return 0, err
  }
  return v.(int16), nil
}

func (m *Msg) Int32(f string) (int32, error) {
  v, err := m.getter(f, tInt32)
  if err != nil {
    return 0, err
  }
  return v.(int32), nil
}

func (m *Msg) Int64(f string) (int64, error) {
  v, err := m.getter(f, tInt64)
  if err != nil {
    return 0, err
  }
  return v.(int64), nil
}

func (m *Msg) Uint(f string) (uint, error) {
  v, err := m.getter(f, tUint)
  if err != nil {
    return 0, err
  }
  return v.(uint), nil
}

func (m *Msg) Uint8(f string) (uint8, error) {
  v, err := m.getter(f, tUint8)
  if err != nil {
    return 0, err
  }
  return v.(uint8), nil
}

func (m *Msg) Uint16(f string) (uint16, error) {
  v, err := m.getter(f, tUint16)
  if err != nil {
    return 0, err
  }
  return v.(uint16), nil
}

func (m *Msg) Uint32(f string) (uint32, error) {
  v, err := m.getter(f, tUint32)
  if err != nil {
    return 0, err
  }
  return v.(uint32), nil
}

func (m *Msg) Uint64(f string) (uint64, error) {
  v, err := m.getter(f, tUint64)
  if err != nil {
    return 0, err
  }
  return v.(uint64), nil
}

func (m *Msg) Float32(f string) (float32, error) {
  v, err := m.getter(f, tFloat32)
  if err != nil {
    return 0, err
  }
  return v.(float32), nil
}

func (m *Msg) Float64(f string) (float64, error) {
  v, err := m.getter(f, tFloat64)
  if err != nil {
    return 0, err
  }
  return v.(float64), nil
}

func (m *Msg) Complex64(f string) (complex64, error) {
  v, err := m.getter(f, tComplex64)
  if err != nil {
    return 0, err
  }
  return v.(complex64), nil
}

func (m *Msg) Complex128(f string) (complex128, error) {
  v, err := m.getter(f, tComplex128)
  if err != nil {
    return 0, err
  }
  return v.(complex128), nil
}

func (m *Msg) String(f string) (string, error) {
  v, err := m.getter(f, tString)
  if err != nil {
    return "", err
  }
  return v.(string), nil
}

func (m *Msg) Bool(f string) (bool, error) {
  v, err := m.getter(f, tBool)
  if err != nil {
    return false, err
  }
  return v.(bool), nil
}

func (m *Msg) Rune(f string) (rune, error) {
  v, err := m.getter(f, tRune)
  if err != nil {
    return 0, err
  }
  return v.(rune), nil
}

func (m *Msg) Byte(f string) (byte, error) {
  v, err := m.getter(f, tByte)
  if err != nil {
    return 0, err
  }
  return v.(byte), nil
}

func (m *Msg) getter(f string, expect tProp) (interface{}, error) {
  o, ok := m.Props[f]
  if !ok {
    return false, fmt.Errorf(errMsgNoField, f)
  }

  re := tUnknown

  switch o.(type) {
  case int:
    re = tInt
  case int8:
    re = tInt8
  case int16:
    re = tInt16
  case int32:
    re = tInt32
  case int64:
    re = tInt64
  case uint:
    re = tUint
  case uint8:
    re = tUint8
  case uint16:
    re = tUint16
  case uint32:
    re = tUint32
  case uint64:
    re = tUint64
  case string:
    re = tString
  case float32:
    re = tFloat32
  case float64:
    re = tFloat64
  case complex64:
    re = tComplex64
  case complex128:
    re = tComplex128
  case bool:
    re = tBool
  }

  if re != expect {
    return nil, fmt.Errorf(errMsgFieldTypeTxt, f, o)
  }

  return o, nil
}
