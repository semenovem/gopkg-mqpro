package mqpro

// RegisterEventInMsg Добавляет обработчик входящих сообщений
func (p *Mqpro) RegisterEventInMsg(fn func(*Msg)) error {
  if p.fnEventInMsg != nil {
    p.log.Error("Subscription already exists")
    return ErrRegisterEventInMsg
  }

  p.fnEventInMsg = fn

  for _, conn := range p.connGet {
    conn.RegisterEventInMsg(p.fnEventInMsg)
  }

  return nil
}

// UnregisterEventInMsg Удалит подписку на входящие сообщения
func (p *Mqpro) UnregisterEventInMsg() {
  p.fnEventInMsg = nil
  for _, conn := range p.connGet {
    conn.UnregisterInMsg()
  }
}
