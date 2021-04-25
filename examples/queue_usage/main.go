package main

import (
	"log"
	"sync"
	"time"

	"github.com/yscsky/yu"
)

type arg1 struct {
	idx int
	w   *sync.WaitGroup
}

type arg2 struct {
	idx int
	w   *sync.WaitGroup
}

type arg3 struct {
	idx int
}

func main() {
	normalUsage()
	timeoutUsage()
	stopUsage()
}

func normalUsage() {
	queue := yu.NewQueue(1024, 10, deal)
	queue.Start()
	defer queue.Stop()

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		if i%2 == 0 {
			queue.SubmitSync(&arg2{idx: i, w: wg})
		} else {
			queue.SubmitSync(&arg1{idx: i, w: wg})
		}
	}
	wg.Wait()
}

func deal(args interface{}) {
	switch a := args.(type) {
	case *arg1:
		log.Println("[info] - arg1 idx:", a.idx)
		a.w.Done()
	case *arg2:
		log.Println("[info] - arg2 idx:", a.idx)
		a.w.Done()
	default:
		log.Println("[err] - deal args type is unknown")
	}
}

func timeoutUsage() {
	q := yu.NewQueueWithTimeout(10, 10, deal2, dealOut, 2*time.Second)
	q.Start()
	for i := 0; i < 100; i++ {
		if i > 90 {
			q.Stop()
			break
		}
		q.SubmitSync(arg3{idx: i})
	}
	time.Sleep(10 * time.Second)
}

func deal2(args interface{}) {
	a, ok := args.(arg3)
	if !ok {
		yu.Errf("args is not arg3")
		return
	}
	yu.Logf("idx: %d", a.idx)
	time.Sleep(3 * time.Second)
}

func dealOut(args interface{}) {
	a, ok := args.(arg3)
	if !ok {
		yu.Errf("args is not arg3")
		return
	}
	yu.Warnf("idx: %d timeout", a.idx)
}

func stopUsage() {
	queue := yu.NewQueue(10, 10, dealStop)
	queue.Start()
	go func() {
		for i := 0; i < 100; i++ {
			queue.SubmitSync(arg3{idx: i})
		}
	}()
	time.Sleep(1 * time.Millisecond)
	queue.Stop()
	yu.Logf("stop")
}

func dealStop(args interface{}) {
	a, ok := args.(arg3)
	if !ok {
		yu.Errf("args is not arg3")
		return
	}
	time.Sleep(3 * time.Second)
	yu.Logf("dealStop idx: %d", a.idx)
}
