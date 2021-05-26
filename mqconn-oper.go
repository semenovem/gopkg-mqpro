package mqpro

func (c *Mqconn) Type() TypeConn {
  return c.typeConn
}

func (c *Mqconn) IsConnected() bool {
  return c.stateConn == stateConnect
}

func (c *Mqconn) IsDisconnected() bool {
  return c.stateConn == stateDisconnect
}
