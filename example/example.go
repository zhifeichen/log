package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/zhifeichen/log"
)

func main() {
	// logger := log.New(log.NewOptions(
	// 	log.Filename("example.txt"),
	// 	log.MaxSize(60),
	// 	log.Level("trace"),
	// ))
	log.Init(log.NewOptions(
		log.Filename("example.txt"),
		log.MaxSize(60),
		log.Level("trace"),
	))
	wg := sync.WaitGroup{}
	now := time.Now()
	// logger.Info("start...")
	log.Info("start...")
	for j := range [5]struct{}{} {
		wg.Add(1)
		go func(jj int) {
			for i := range [300000]struct{}{} {
				// logger.Tracef("C:Usersccworkspacegomyapplogexample>benchmarck %s_%d_%d\n", "test", jj, i)
				// logger.Debugf("C:Usersccworkspacegomyapplogexample>benchmarck %s_%d_%d\n", "test", jj, i)
				// logger.Infof("C:sersccworkspacegomyapplogexample>benchmarck %s_%d_%d\n", "test", jj, i)
				// logger.Errorf("C:Usersccworkspacegomyapplogexample>benchmarck %s_%d_%d\n", "test", jj, i)
				log.Tracef("C:Usersccworkspacegomyapplogexample>benchmarck %s_%d_%d\n", "test", jj, i)
				log.Debugf("C:Usersccworkspacegomyapplogexample>benchmarck %s_%d_%d\n", "test", jj, i)
				log.Infof("C:sersccworkspacegomyapplogexample>benchmarck %s_%d_%d\n", "test", jj, i)
				log.Errorf("C:Usersccworkspacegomyapplogexample>benchmarck %s_%d_%d\n", "test", jj, i)
			}
			wg.Done()
		}(j)
	}
	wg.Wait()
	log.Info("done")
	// logger.Flush()
	fmt.Printf("%d\n", time.Since(now).Milliseconds())
	time.Sleep(time.Second)
	log.Info("done2")
	time.Sleep(time.Second * 2)
}
