package main

import "github.com/zhifeichen/log"

func main() {
	logger := log.New(log.NewOptions(
		log.Filename("example.log"),
	))
	logger.Trace("trace")
	logger.Tracef("tracef\n")
	logger.Debug("debug")
	logger.Debugf("debugf\n")
	logger.Info("info")
	logger.Infof("info\n")
	logger.Warn("warn")
	logger.Warnf("warn\n")
	logger.Error("error")
	logger.Errorf("errorf")

	log.Init(log.NewOptions())
	log.Trace("trace")
	log.Tracef("tracef\n")
	log.Debug("debug")
	log.Debugf("debugf\n")
	log.Info("info")
	log.Infof("info\n")
	log.Warn("warn")
	log.Warnf("warn\n")
	log.Error("error")
	log.Errorf("errorf")
}
