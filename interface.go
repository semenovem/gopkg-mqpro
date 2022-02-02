package mqm

import (
  "context"
  "github.com/semenovem/mqm/v2/queue"
)

type Queue interface {
  Put(ctx context.Context, msg *queue.Msg) error
  Get(ctx context.Context) (*queue.Msg, error)
  GetByCorrelId(ctx context.Context, correlId []byte) (*queue.Msg, error)
  GetByMsgId(ctx context.Context, msgId []byte) (*queue.Msg, error)
}
