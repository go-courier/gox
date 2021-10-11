package renderer

import (
	"sync"
	"time"

	"github.com/go-courier/gox/pkg/browser"
)

type commitQueue struct {
	queue    chan func()
	queueBuf []func()
	rw       sync.RWMutex
}

func (q *commitQueue) push(fn func()) {
	q.rw.Lock()
	q.queueBuf = append(q.queueBuf, fn)
	q.rw.Unlock()
}

func (q *commitQueue) runIfExists() {
	q.rw.Lock()
	func(queue []func()) {
		actionQueueBufferEach(queue, bufSize, func(buf []func()) {
			browser.RequestAnimationFrame(func() {
				for i := range buf {
					buf[i]()
				}
			})
		})
	}(q.queueBuf)
	q.queueBuf = make([]func(), 0)
	q.rw.Unlock()
}

const bufSize = 2048

func actionQueueBufferEach(queue []func(), bufSize int, each func(queue []func())) {
	n := len(queue)
	partN := n / bufSize
	for i := 0; i < partN+1; i++ {
		left := i * bufSize
		if left > n-1 {
			break
		}
		right := (i + 1) * bufSize
		if right > n {
			right = n
		}
		each(queue[i*bufSize : right])
	}
}

func (q *commitQueue) ForceCommit() {
	q.runIfExists()
}

func (q *commitQueue) Close() error {
	close(q.queue)
	return nil
}

func (q *commitQueue) Start() {
	q.queue = make(chan func())

	go func() {
	LOOP:
		for {
			select {
			case fn, ok := <-q.queue:
				if !ok {
					break LOOP
				}
				q.push(fn)
			case <-time.After(10 * time.Millisecond):
				q.runIfExists()
			}
		}
	}()
}

func (q *commitQueue) Dispatch(fn func()) {
	q.queue <- fn
}
