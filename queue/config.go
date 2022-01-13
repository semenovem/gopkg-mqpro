package queue

import (
  "fmt"
  "strings"
  "time"
)

// BaseConfig Конфигурация
type BaseConfig struct {
  DevMode            bool
  Header             Header // Тип заголовков
  Rfh2CodedCharSetId int32  // Тип кодирования
  Rfh2RootTag        string // Корневой тег
}

func (q *Queue) IsConfigured() bool {
  return q.queueName != "" && len(q.perm) != 0
}

func (q *Queue) UpdateBaseCfg() {
  c := q.base.GetBaseCfg()

  q.h = c.Header
  q.devMode = c.DevMode
  q.rfh2RootTag = c.Rfh2RootTag

  if q.h == HeaderRfh2 {
    q.rfh2 = newRfh2Cfg()
    if c.Rfh2CodedCharSetId != 0 {
      q.rfh2.CodedCharSetId = c.Rfh2CodedCharSetId
    }
  }
}

func (q *Queue) CfgByStr(s string) error {
  if s == "" {
    return nil
  }

  var err error
  q.queueName, q.perm, err = parseQueue(s)
  if err != nil {
    return err
  }
  m := map[string]interface{}{
    "n": q.queueName,
  }
  q.log = q.log.WithFields(m)

  return nil
}

func (q *Queue) Alias() string {
  return q.alias
}

func (q *Queue) SetDevMode(v bool) {
  q.devMode = v
}

func (q *Queue) PrintSetCli(p string) {
  PrintSetCli(q.getSet(), p)
}

func (q *Queue) getSet() []map[string]string {
  q.UpdateBaseCfg()
  m := []map[string]string{
    {"queueName": q.queueName},
    {"perm": strings.Join(q.convPermToVal(), ",")},
    {"state": stateMapByKey[q.state]},
    {"reconnectDelay": fmt.Sprintf("%d sec", q.reconnectDelay/time.Second)},
    {"delayClose": fmt.Sprintf("%d ms", q.delayClose/time.Millisecond)},
    {"devMode": fmt.Sprintf("%t", q.devMode)},
    {"header": HeaderMapByKey[q.h]},
    {"rfh2RootTag": q.rfh2RootTag},
  }

  if q.h == HeaderRfh2 {
    m1 := []map[string]string{
      {"StructId": q.rfh2.StructId},
      {"Version": fmt.Sprintf("%d", q.rfh2.Version)},
      {"Encoding": fmt.Sprintf("%d", q.rfh2.Encoding)},
      {"CodedCharSetId": fmt.Sprintf("%d", q.rfh2.CodedCharSetId)},
      {"Format": q.rfh2.Format},
      {"Flags": fmt.Sprintf("%d", q.rfh2.Flags)},
      {"NameValueCCSID": fmt.Sprintf("%d", q.rfh2.NameValueCCSID)},
      {},
    }

    m = append(m, m1...)
  }

  return m
}
