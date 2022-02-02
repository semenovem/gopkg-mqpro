package mqm

import "github.com/semenovem/mqm/v2/queue"

type Channel struct {
  get *queue.Queue
  put *queue.Queue
}


