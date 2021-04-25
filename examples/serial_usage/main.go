package main

import (
	"log"
	"sync"

	"github.com/yscsky/yu"
)

func main() {
	serial := yu.NewSerial(100)
	serial.Start()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			log.Println("go 1 get:", serial.Get())
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			log.Println("go 2 get:", serial.Get())
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			log.Println("go 3 get:", serial.Get())
		}
		wg.Done()
	}()
	wg.Wait()
}
