package queue

import "github.com/ibm-messaging/mq-golang/v5/ibmmq"

func (q *Queue) IsConnected() bool {
	return q.state == stateConn
}

func (q *Queue) isWarnConn(err error) {
	if err != nil {
		mqret := err.(*ibmmq.MQReturn)
		if mqret == nil || mqret.MQRC != ibmmq.MQRC_CONNECTION_BROKEN {
			q.log.Warn(err)
		}
	}
}

func (q *Queue) isWarn(err error) {
	if err != nil {
		q.log.Warn(err)
	}
}

func (q *Queue) permString() []string {
	a := make([]string, len(q.perm))
	for i, v := range q.perm {
		a[i] = permKey[v]
	}
	return a
}
