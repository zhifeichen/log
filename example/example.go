package main

import (
	"fmt"
	"time"

	"github.com/zhifeichen/log"
)

func main() {
	logger := log.New(log.NewOptions(
		log.Filename("example.txt"),
		log.MaxSize(60),
		log.Level("trace"),
	))
	now := time.Now()
	for range [3]struct{}{} {
		for range [500000]struct{}{} {
			logger.Tracef("C:Usersccworkspacegomyapplogexample>benchmarck %s\n", "test")
			logger.Debugf("C:Usersccworkspacegomyapplogexample>benchmarck %s\n", "test")
			logger.Infof("C:sersccworkspacegomyapplogexample>benchmarck %s\n", "test")
			logger.Errorf("C:Usersccworkspacegomyapplogexample>benchmarck %s\n", "test")
		}
	}
	logger.Info("done")
	logger.Flush()
	fmt.Printf("%d\n", time.Since(now).Milliseconds())
}
