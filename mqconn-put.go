package mqpro

import (
	"fmt"
	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
	"github.com/sirupsen/logrus"
)

// Put отправка сообщения в очередь
func (c *Mqconn) Put(msg *Msg) ([]byte, error) {
	l := c.log.WithField("method", "Put")

	if msg.CorrelId != nil {
		l = c.log.WithField("correlId", fmt.Sprintf("%x", msg.CorrelId))
	}

	if !c.IsConnected() {
		return nil, ErrNoConnection
	}

	d, err := c.put(msg, l)

	if err != nil {
		if HasConnBroken(err) {
			c.reqError()
			err = ErrConnBroken
		} else {
			err = ErrPutMsg
		}
	}

	return d, err
}

func (c *Mqconn) put(msg *Msg, l *logrus.Entry) ([]byte, error) {
	c.mxPut.Lock()
	defer c.mxPut.Unlock()

	l.Trace("Start")

	cmho := ibmmq.NewMQCMHO()
	putMsgHandle, err := c.mgr.CrtMH(cmho)
	if err != nil {
		return nil, err
	}

	err = setProps(&putMsgHandle, msg.Props, l)
	if err != nil {
		return nil, err
	}

	putmqmd := ibmmq.NewMQMD()
	pmo := ibmmq.NewMQPMO()

	if msg.CorrelId != nil {
		putmqmd.CorrelId = msg.CorrelId
	}

	pmo.Options = ibmmq.MQPMO_NO_SYNCPOINT
	pmo.OriginalMsgHandle = putMsgHandle
	putmqmd.Format = ibmmq.MQFMT_STRING

	var d []byte
	if msg.Payload == nil {
		d = make([]byte, 0)
	} else {
		d = msg.Payload
	}

	err = c.que.Put(putmqmd, pmo, d)
	if err != nil {
		return nil, err
	}

	l.Debugf("Success. MsgId: %x", putmqmd.MsgId)

	return putmqmd.MsgId, nil
}
