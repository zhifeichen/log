package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/zhifeichen/log"
)

func main() {
	// ========== 测试包级函数 (向后兼容) ==========
	fmt.Println("=== 包级函数测试 ===")
	log.Init(log.NewOptions(
		log.Filename("example.txt"),
		log.MaxSize(60),
		log.Level("trace"),
	))

	log.Info("start...")
	log.Debug("debug message")
	log.Warn("warn message")
	log.Error("error message")
	log.Trace("trace message")

	wg := sync.WaitGroup{}
	now := time.Now()
	for j := range [5]struct{}{} {
		wg.Add(1)
		go func(jj int) {
			for i := range [100]struct{}{} {
				log.Tracef("benchmark %s_%d_%d\n", "test", jj, i)
				log.Debugf("benchmark %s_%d_%d\n", "test", jj, i)
				log.Infof("benchmark %s_%d_%d\n", "test", jj, i)
				log.Errorf("benchmark %s_%d_%d\n", "test", jj, i)
			}
			wg.Done()
		}(j)
	}
	wg.Wait()
	log.Info("done")
	fmt.Printf("elapsed: %dms\n", time.Since(now).Milliseconds())

	log.Flush()

	// ========== 测试实例化 Logger ==========
	fmt.Println("=== 实例化 Logger 测试 ===")
	l := log.New(log.NewOptions(
		log.Filename("instance.log"),
		log.Level("debug"),
		log.MaxSize(10),
	))

	l.Info("instance info message")
	l.Debug("instance debug message")
	l.Warn("instance warn message")

	// ========== 测试 Discard/ResumeWriter ==========
	fmt.Println("=== Discard/ResumeWriter 测试 ===")
	l.Discard()
	l.Info("this should NOT appear in log") // 不会写入
	l.ResumeWriter()
	l.Info("this should appear in log") // 正常写入

	// ========== 测试 Writer ==========
	fmt.Println("=== Writer 测试 ===")
	w := l.Writer()
	w.Write([]byte("write via io.Writer\n"))

	// ========== 测试 Context/TraceID ==========
	fmt.Println("=== Context/TraceID 测试 ===")
	ctx := l.SetContext(context.Background())
	logWithCtx := l.WithContext(ctx)
	logWithCtx.Info("message with traceID from context")

	logWithTraceID := l.WithTraceID(12345)
	logWithTraceID.Info("message with explicit traceID")

	newCtx := l.GetNewContext(ctx)
	logWithNewCtx := l.WithContext(newCtx)
	logWithNewCtx.Info("message with new context")

	// ========== 测试包级 Context 函数 ==========
	fmt.Println("=== 包级 Context 测试 ===")
	pkgCtx := log.SetContext(context.Background())
	pkgLogger := log.WithContext(pkgCtx)
	pkgLogger.Info("package-level with context")

	pkgLogger2 := log.WithTraceID(99999)
	pkgLogger2.Info("package-level with traceID")

	// ========== 测试 LoggingRoundTripper ==========
	fmt.Println("=== RoundTripper 测试 ===")
	rt := &log.LoggingRoundTripper{
		Transport: http.DefaultTransport,
		Logger:    l,
		LogBodies: false,
	}
	_ = rt // 编译验证通过即可，实际 HTTP 请求会依赖外部网络

	l.Flush()
	fmt.Println("all tests done")
}
