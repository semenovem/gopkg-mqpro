package main

import (
  "fmt"
  "github.com/semenovem/mqm/v2/queue"
  "net/http"
)

// Подписка на входящие сообщения
// curl host:port/sub
func onRegisterInMsg(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Включено получение сообщений из очереди")
  subscr()

  printCfg()
  _, _ = fmt.Fprintf(w, "[sub] Ok\n")
}

// Отписаться
// curl host:port/unsub
func offRegisterInMsg(w http.ResponseWriter, _ *http.Request) {
  if cfg.Mirror {
    fmt.Println("Отключено получение сообщений из очереди")
    _, _ = fmt.Fprintf(w, "[unsub] ERROR. use curl host:port/off-mirror\n")
    return
  }

  fmt.Println("Отключено получение входящих сообщений")
  unsubscr()

  printCfg()
  _, _ = fmt.Fprintf(w, "[unsub] Ok\n")
}

// Подписаться на входящие сообщения
func subscr() {
  cfg.Subscribe = true
  if !mqQueGet.IsSubscribed() {
    mqQueGet.RegisterInMsg(handlerInMsg)
  }
}

// Отписаться
func unsubscr() {
  cfg.Subscribe = false
  if mqQueGet.IsSubscribed() {
    mqQueGet.UnregisterInMsg()
  }
}

// Обработчик входящих сообщений
func handlerInMsg(m *queue.Msg) {
  fmt.Printf("Обработчик входящих сообщений. Mirror = %t", cfg.Mirror)
  fmt.Println(">>>>> ", formatMsgId(m.MsgId))
  if cfg.Mirror {
    mirror(m)
  }
}
