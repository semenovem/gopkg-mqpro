package mqpro

import "github.com/sirupsen/logrus"

func MqconnNew(tc TypeConn, l *logrus.Entry, c *Cfg) *Mqconn {
	o := &Mqconn{
		cfg:            *c,
		fnsConn:        map[uint32]chan struct{}{},
		fnsDisconn:     map[uint32]chan struct{}{},
		reconnectDelay: defReconnectDelay,
		stateConn:      stateDisconnect,
	}

	m := map[string]interface{}{
		"hostPort": o.endpoint(),
		"manager":  c.MgrName,
		"channel":  c.ChannelName,
		"queue":    c.QueueName,
		"type":     typeConnTxt[tc],
	}

	o.log = l.WithFields(m)

	if tc != TypePut && tc != TypeGet && tc != TypeBrowse {
		o.log.Panic("Unknown connection type")
	}

	o.typeConn = tc

	return o
}
