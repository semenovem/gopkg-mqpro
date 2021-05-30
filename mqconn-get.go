package mqpro

import (
	"fmt"
	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
	"github.com/sirupsen/logrus"
)

func (c *Mqconn) Get() (*Msg, bool, error) {
	l := c.log.WithField("method", "Get")

	msg, ok, err := c.get(operGet, nil, l)
	if err != nil {
		if HasConnBroken(err) {
			c.reqError()
			err = ErrConnBroken
		} else {
			err = ErrGetMsg
		}
	}
	return msg, ok, err
}

// GetByCorrelId Извлекает сообщение из очереди по его CorrelId
func (c *Mqconn) GetByCorrelId(correlId []byte) (*Msg, bool, error) {
	l := c.log.WithFields(map[string]interface{}{
		"correlId": fmt.Sprintf("%x", correlId),
		"method":   "GetByCorrelId",
	})

	msg, ok, err := c.get(operGetByCorrelId, correlId, l)

	if err != nil {
		if HasConnBroken(err) {
			c.reqError()
			err = ErrConnBroken
		} else {
			err = ErrGetMsg
		}
	}

	return msg, ok, err
}

// GetByMsgId Извлекает сообщение из очереди по его MsgId
func (c *Mqconn) GetByMsgId(msgId []byte) (*Msg, bool, error) {
	l := c.log.WithFields(map[string]interface{}{
		"msgId":  fmt.Sprintf("%x", msgId),
		"method": "GetByCorrelId",
	})

	msg, ok, err := c.get(operGetByMsgId, msgId, l)

	if err != nil {
		if HasConnBroken(err) {
			c.reqError()
			err = ErrConnBroken
		} else {
			err = ErrGetMsg
		}
	}

	return msg, ok, err
}

// получение сообщения
func (c *Mqconn) get(oper queueOper, id []byte, l *logrus.Entry) (*Msg, bool, error) {
	if !c.IsConnected() {
		return nil, false, ErrNoConnection
	}

	c.mxGet.Lock()
	defer c.mxGet.Unlock()

	l.Trace("Start")

	var datalen int
	var err error

	getmqmd := ibmmq.NewMQMD()
	gmo := ibmmq.NewMQGMO()
	getmqmd.Format = ibmmq.MQFMT_STRING
	gmo.Options = ibmmq.MQGMO_NO_SYNCPOINT

	gmo.Options |= ibmmq.MQGMO_WAIT
	gmo.WaitInterval = int32(3 * 1000) // The WaitInterval is in milliseconds

	switch oper {
	case operGet:
	case operGetByMsgId:
		gmo.MatchOptions = ibmmq.MQMO_MATCH_MSG_ID
		getmqmd.MsgId = id
	case operGetByCorrelId:
		gmo.MatchOptions = ibmmq.MQMO_MATCH_CORREL_ID
		getmqmd.CorrelId = id
	default:
		l.Panicf("Unknown operation. queueOper = %v", oper)
	}

	cmho := ibmmq.NewMQCMHO()

	getMsgHandle, err := c.mgr.CrtMH(cmho)
	if err != nil {
		l.Error("Ошибка создания объекта свойств сообщения: ", err)
		return nil, false, err
	}
	defer func() {
		err := dltMh(getMsgHandle)
		if err != nil {
			l.Warnf("Ошибка удаления объекта свойств сообщения: %s", err.Error())
		}
	}()

	gmo.MsgHandle = getMsgHandle
	gmo.Options |= ibmmq.MQGMO_PROPERTIES_IN_HANDLE

	buffer := make([]byte, 0, 1024)

	for i := 0; i < 2; i++ {
		buffer, datalen, err = c.que.GetSlice(getmqmd, gmo, buffer)
		if err != nil {
			mqret := err.(*ibmmq.MQReturn)

			if mqret.MQRC == ibmmq.MQRC_TRUNCATED_MSG_FAILED {
				buffer = make([]byte, 0, datalen)
				continue
			}

			if mqret.MQRC == ibmmq.MQRC_NO_MSG_AVAILABLE {
				err = nil
			} else {
				l.Error("Ошибка получения сообщения: ", err, "  len: ", datalen)
			}

			return nil, false, err
		}

		break
	}

	props, err := properties(getMsgHandle)
	if err != nil {
		l.Error("Ошибка получения свойств сообщения: ", err)
		return nil, false, err
	}

	l.Debug("Success")

	ret := &Msg{
		Payload:  buffer,
		Props:    props,
		CorrelId: getmqmd.CorrelId,
		MsgId:    getmqmd.MsgId,
	}

	return ret, true, nil
}
