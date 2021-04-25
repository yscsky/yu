package yu

import (
	"sync"
	"time"
)

type dealFunc func(interface{})

// Queue 队列
type Queue struct {
	num     int
	deal    dealFunc
	out     time.Duration
	dealOut dealFunc
	que     chan interface{}
	wg      *sync.WaitGroup
	stop    bool
	stch    chan struct{}
}

// NewQueue 新建队列
func NewQueue(size, num int, deal dealFunc) *Queue {
	return &Queue{
		num:     num,
		deal:    deal,
		out:     0,
		dealOut: nil,
		que:     make(chan interface{}, size),
		wg:      new(sync.WaitGroup),
		stop:    false,
		stch:    make(chan struct{}),
	}
}

// NewQueue 新建带提交超时的队列
func NewQueueWithTimeout(size, num int, deal, dealOut dealFunc, out time.Duration) *Queue {
	if out == 0 || dealOut == nil {
		panic("NewQueueWithTimeout must set out and outDeal")
	}
	return &Queue{
		num:     num,
		deal:    deal,
		out:     out,
		dealOut: dealOut,
		que:     make(chan interface{}, size),
		wg:      new(sync.WaitGroup),
		stop:    false,
		stch:    make(chan struct{}),
	}
}

// Start 启动队列
func (q *Queue) Start() {
	for i := 0; i < q.num; i++ {
		q.wg.Add(1)
		go func() {
			defer q.wg.Done()
			for item := range q.que {
				q.deal(item)
			}
		}()
	}
}

// Stop 停止队列
func (q *Queue) Stop() {
	q.stch <- struct{}{}
	close(q.que)
	q.wg.Wait()
}

// SubmitAsync 非阻塞提交队列
func (q *Queue) SubmitAsync(item interface{}) {
	go func(i interface{}) { q.submit(i) }(item)
}

// SubmitSync 阻塞提交队列
func (q *Queue) SubmitSync(item interface{}) {
	q.submit(item)
}

func (q *Queue) submit(item interface{}) {
	if q.isStop() {
		return
	}
	if q.out <= 0 {
		q.que <- item
		return
	}
	t := time.NewTimer(q.out)
	defer t.Stop()
	select {
	case q.que <- item:
	case <-t.C:
		if q.dealOut != nil {
			q.dealOut(item)
		}
	}
}

func (q *Queue) isStop() bool {
	select {
	case <-q.stch:
		q.stop = true
	default:
	}
	return q.stop
}

// QueueManager 队列管理器
type QueueManager struct {
	ques []*Queue
}

// NewQueueManager 新建队列管理器
func NewQueueManager() *QueueManager {
	return &QueueManager{
		ques: make([]*Queue, 0),
	}
}

// AddQueue 管理器中添加队列
func (qm *QueueManager) AddQueue(q *Queue) {
	qm.ques = append(qm.ques, q)
}

// GetQueue 管理器中获取队列
func (qm *QueueManager) GetQueue(i int) *Queue {
	q := qm.ques[i]
	return q
}

// PushQueue 添加任务到管理器队列中
func (qm *QueueManager) PushQueue(i int, args interface{}) {
	qm.ques[i].submit(args)
}

// StartQueue 管理器启动所有队列
func (qm *QueueManager) StartQueue() {
	for _, q := range qm.ques {
		q.Start()
	}
}

// StopQueue 管理器停止所有队列
func (qm *QueueManager) StopQueue() {
	for _, q := range qm.ques {
		q.Stop()
	}
}
