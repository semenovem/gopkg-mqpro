package main

import (
  "fmt"
  "github.com/semenovem/gopkg_mqpro/v2/queue"
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

func subscr() {
  cfg.SimpleSubscriber = true
  ibmmqOper1In.RegisterInMsg(handlerInMsg)
}

func unsubscr() {
  cfg.SimpleSubscriber = false
  ibmmqOper1In.UnregisterInMsg()
}

// Обработчик входящих сообщений
func handlerInMsg(m *queue.Msg) {
  fmt.Println("Вызван обработчик входящих сообщений")
  fmt.Printf("Режим Mirror = %v", cfg.Mirror)
  logMsgIn(m)

  if cfg.Mirror {
    mirror(m)
  }
}
