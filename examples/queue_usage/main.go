package main

import (
	"flag"
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
	opt := flag.String("opt", "", "operate usage")
	flag.Parse()
	switch *opt {
	case "normal":
		normalUsage()
	case "stop1":
		stopUsage1()
	case "stop2":
		stopUsage2()
	default:
		yu.Warnf("opt should be normal, stop1, stop2")
	}
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
		yu.Logf("arg1 idx: %d", a.idx)
		a.w.Done()
	case *arg2:
		yu.Logf("arg2 idx: %d", a.idx)
		a.w.Done()
	default:
		yu.Errf("deal args type is unknown")
	}
}

func stopUsage1() {
	qm := yu.NewQueueManager()
	qm.AddQueue(yu.NewQueue(1024, 10, deal2))
	qm.StartQueue()
	yu.Logf("start")
	time.Sleep(10 * time.Millisecond)
	qm.StopQueue()
	yu.Logf("stop")
}

func stopUsage2() {
	qm := yu.NewQueueManager()
	qm.AddQueue(yu.NewQueue(1024, 10, deal2))
	qm.StartQueue()
	yu.Logf("start")
	go func() {
		for i := 0; i < 1000000; i++ {
			qm.PushQueue(0, arg3{idx: i}, false)
		}
	}()
	time.Sleep(10 * time.Millisecond)
	qm.StopQueue()
	yu.Logf("stop")
}

func deal2(args interface{}) {
	a, ok := args.(arg3)
	if !ok {
		yu.Errf("args is not arg3")
		return
	}
	yu.Logf("arg3 idx: %d", a.idx)
}
