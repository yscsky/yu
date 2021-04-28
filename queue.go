package yu

import (
	"sync"
)

type dealFunc func(interface{})

// Queue 队列
type Queue struct {
	num  int
	deal dealFunc
	que  chan interface{}
	wg   *sync.WaitGroup
	stop bool
	stch chan struct{}
}

// NewQueue 新建队列
func NewQueue(size, num int, deal dealFunc) *Queue {
	return &Queue{
		num:  num,
		deal: deal,
		que:  make(chan interface{}, size),
		wg:   new(sync.WaitGroup),
		stop: false,
		stch: make(chan struct{}),
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
	q.tryStop()
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

func (q *Queue) isStop() bool {
	select {
	case <-q.stch:
		q.stop = true
	default:
	}
	return q.stop
}

func (q *Queue) submit(i interface{}) {
	if q.isStop() {
		return
	}
	q.que <- i
}

func (q *Queue) tryStop() {
	select {
	case <-q.que:
		q.stch <- struct{}{}
	default:
	}
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
func (qm *QueueManager) PushQueue(i int, args interface{}, async bool) {
	if async {
		qm.ques[i].SubmitAsync(args)
		return
	}
	qm.ques[i].SubmitSync(args)
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
