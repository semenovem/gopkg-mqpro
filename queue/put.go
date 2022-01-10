package queue

import (
	"context"
	"fmt"
	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
	"github.com/sirupsen/logrus"
)

// Put отправка сообщения в очередь
func (q *Queue) Put(ctx context.Context, msg *Msg) ([]byte, error) {
	// TODO использовать контекст

	l := q.log.WithField("method", "Put")

	if msg.CorrelId != nil {
		l = l.WithField("correlId", fmt.Sprintf("%x", msg.CorrelId))
	}

	d, err := q.put(msg, l)
	if err == ErrConnBroken {
		q.stateError()
	}

	return d, err
}

func (q *Queue) put(msg *Msg, l *logrus.Entry) ([]byte, error) {
	if !q.IsConnected() {
		return nil, ErrNoConnection
	}

	q.mxMsg.Lock()
	defer q.mxMsg.Unlock()

	l.Trace("Start")

	var payload []byte
	if msg.Payload == nil {
		payload = make([]byte, 0)
	} else {
		payload = msg.Payload
	}

	putmqmd := ibmmq.NewMQMD()
	pmo := ibmmq.NewMQPMO()
	cmho := ibmmq.NewMQCMHO()

	pmo.Options = ibmmq.MQPMO_NO_SYNCPOINT

	if msg.CorrelId != nil {
		putmqmd.CorrelId = msg.CorrelId
	}

	var devMsg Msg

	if q.devMode {
		devMsg = *msg
		f := devMode(&devMsg, payload, "put")
		defer func() {
			f()
		}()
	}

	switch q.h {
	case HeaderRfh2:
		putmqmd.Format = ibmmq.MQFMT_RF_HEADER_2
		hd, err := q.Rfh2Marshal(msg.Props)
		if err != nil {
			l.Error("Не удалось подготовить сообщение с заголовками rfh2: ", err)
			return nil, err
		}
		payload = append(hd, payload...)

		if q.devMode {
			devMsg.MQRFH2, err = q.Rfh2Unmarshal(hd)
			if err != nil {
				return nil, err
			}
			devMsg.Payload = payload
		}

	default:
		putmqmd.Format = ibmmq.MQFMT_STRING

		putMsgHandle, err := q.mgr.CrtMH(cmho)
		if err != nil {
			l.Errorf("Ошибка создания объекта свойств сообщения: %s", err)

			if IsConnBroken(err) {
				err = ErrConnBroken
			} else {
				err = ErrPutMsg
			}

			return nil, err
		}

		err = setProps(&putMsgHandle, msg.Props, l)
		if err != nil {
			return nil, ErrPutMsg
		}
		pmo.OriginalMsgHandle = putMsgHandle
	}

	err := q.que.Put(putmqmd, pmo, payload)
	if err != nil {
		l.Error("Ошибка отправки сообщения: ", err)

		if IsConnBroken(err) {
			err = ErrConnBroken
		} else {
			err = ErrPutMsg
		}

		return nil, err
	}

	l.Debugf("Success. MsgId: %x", putmqmd.MsgId)

	if q.devMode {
		devMsg.Time = putmqmd.PutDateTime
		devMsg.MsgId = putmqmd.MsgId
	}

	return putmqmd.MsgId, nil
}
