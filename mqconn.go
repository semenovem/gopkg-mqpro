package mqpro

import (
	"errors"
	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Mqconn struct {
	cfg            Cfg
	log            *logrus.Entry
	typeConn       TypeConn              // тип подключения
	mgr            *ibmmq.MQQueueManager // Менеджер очереди
	que            *ibmmq.MQObject       // Объект открытой очереди
	mx             sync.Mutex
	stateConn      stateConn
	chMgr          chan reqStateConn
	fnInMsg        func(*Msg)               // подписка на входящие сообщения
	ctlo           *ibmmq.MQCTLO            // объект подписки ibmmq
	fnsConn        map[uint32]chan struct{} // подписки на установку соединения
	fnsDisconn     map[uint32]chan struct{} // подписки на закрытие соединения
	ind            uint32                   // простой атомарный счетчик
	reconnectDelay time.Duration            // таймаут попыток повторного подключения

	// менеджер imbmq одновременно может отправлять/принимать одно сообщение
	mxPut    sync.Mutex
	mxGet    sync.Mutex
	mxBrowse sync.Mutex
}

// Cfg Данные подключения
type Cfg struct {
	Host        string
	Port        int
	MgrName     string
	ChannelName string
	QueueName   string
	AppName     string
	User        string
	Pass        string
	Priority    string
}

type TypeConn int
type stateConn int
type reqStateConn int
type queueOper int

const (
	TypePut TypeConn = iota + 1
	TypeGet
	TypeBrowse
	defReconnectDelay = time.Second * 3
)

const (
	stateDisconnect stateConn = iota
	stateConnect
	stateErr
)

const (
	reqConnect reqStateConn = iota
	reqReconnect
	reqDisconnect
)

const (
	operGet queueOper = iota
	operGetByMsgId
	operGetByCorrelId
	operPut
)

var typeConnTxt = map[TypeConn]string{
	TypePut:    "TypePut",
	TypeGet:    "TypeGet",
	TypeBrowse: "TypeBrowse",
}

type Msg struct {
	MsgId    []byte
	CorrelId []byte
	Payload  []byte
	Props    map[string]interface{}
}

var (
	ErrConnBroken = errors.New("ibm mq conn: connection broken")
	ErrPutMsg     = errors.New("ibm mq: failed to put message")
	ErrGetMsg     = errors.New("ibm mq: failed to get message")
	ErrBrowseMsg  = errors.New("ibm mq: failed to browse message")
)
