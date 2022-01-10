package queue

import (
	"bytes"
	"fmt"
)

// CoreSet Данные подключения
type CoreSet struct {
	DevMode            bool
	Header             Header // Тип заголовков
	Rfh2CodedCharSetId int32  // Тип кодирования
	Rfh2RootTag        string // Корневой тег
}

type CfgQueue struct {
	Name   string      // Название очереди
	Access []permQueue // Разрешения на очередь
}

func (q *Queue) Set(cfg *CoreSet) {
	q.mx.Lock()
	defer q.mx.Unlock()

	q.h = cfg.Header
	q.devMode = cfg.DevMode
	q.rfh2RootTag = cfg.Rfh2RootTag

	if q.h == HeaderRfh2 {
		q.rfh2 = newRfh2Cfg()
		if cfg.Rfh2CodedCharSetId != 0 {
			q.rfh2.CodedCharSetId = cfg.Rfh2CodedCharSetId
		}
	}
}

func (q *Queue) CfgQueue(s string) error {
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

func (q *Queue) SetDevMode(v bool) {
	q.devMode = v
}

// PrintCfg
// Deprecated
func (q *Queue) PrintCfg() {
	var buf = bytes.NewBufferString("")
	f := func(s string, i ...interface{}) {
		buf.WriteString(fmt.Sprintf(s, i...))
	}

	f("Environment variable values:\n")
	f("DevMode:           = %t\n", q.devMode)
	f("Header             = %s\n", q.h)
	//f("Rfh2CodedCharSetId = %s\n", cfg.Rfh2CodedCharSetId)
	f("Rfh2RootTag        = %t\n", q.rfh2RootTag)

	fmt.Println(buf.String())
}
