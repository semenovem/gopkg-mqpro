package main

import (
  "github.com/semenovem/gopkg_mqpro/v2/queue"
)

func hndIncomingMsg(msg *queue.Msg) {

  logMsgIn(msg)

}

