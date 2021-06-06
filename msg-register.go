package mqpro

// RegisterEvenInMsg Добавляет обработчик входящих сообщений
func (p *Mqpro) RegisterEvenInMsg(fn func(*Msg)) {
  if p.fnEventInMsg != nil {
    p.log.Panic("Subscription already exists")
  }

  p.fnEventInMsg = fn

  for _, conn := range p.connGet {
    conn.RegisterEventInMsg(p.fnEventInMsg)
  }
}

// UnregisterEvenInMsg Удалит подписку на входящие сообщения
func (p *Mqpro) UnregisterEvenInMsg() {
  p.fnEventInMsg = nil
  for _, conn := range p.connGet {
    conn.UnregisterInMsg()
  }
}
