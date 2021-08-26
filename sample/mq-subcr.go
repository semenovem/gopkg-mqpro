package main

import (
  "fmt"
  mqpro "github.com/semenovem/gopkg_mqpro"
  "net/http"
)

// Подписка на входящие сообщения
// curl host:port/sub
func onRegisterInMsg(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Включено получение сообщений из очереди")
  subscr()

  _, _ = fmt.Fprintf(w, "[subcribe] Ok\n")
}

// Отписаться
// curl host:port/unsub
func offRegisterInMsg(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Отключено получение входящих сообщений")
  unsubscr()

  _, _ = fmt.Fprintf(w, "[unsubcribe] Ok\n")
}

// Обработчик входящих сообщений
func handlerInMsg(msg *mqpro.Msg) {
  fmt.Println("Вызван обработчик входящих сообщений")
  logMsg(msg)
}

func subscr()  {
  ibmmq.RegisterEvenInMsg(handlerInMsg)
}

func unsubscr()  {
  ibmmq.UnregisterEvenInMsg()
}
