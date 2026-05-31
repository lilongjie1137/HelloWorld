// Package printq 实现门店本地打印队列：job_key 幂等防重打 + 失败重试 + 状态回执。
package printq

import (
	"context"
	"sync"
	"time"

	"github.com/lilongjie1137/HelloWorld/common/types"
	"github.com/lilongjie1137/HelloWorld/edge/internal/device"
)

// Job 一条本地打印任务。
type Job struct {
	Key       string             // 幂等键（云端 print_job.job_key）
	PrinterID int64              // 目标打印机
	Type      types.PrintJobType // KITCHEN/BAR/LABEL/CHECKOUT/REFUND
	Payload   []byte             // 已渲染字节流
}

// Result 打印回执。
type Result struct {
	Key     string
	Printed bool
	Retries int
	Err     error
}

// PrinterResolver 按 printerID 取得打印机实例。
type PrinterResolver func(printerID int64) (device.Printer, bool)

// Queue 本地打印队列。
type Queue struct {
	resolver PrinterResolver
	maxRetry int
	retryGap time.Duration

	mu   sync.Mutex
	seen map[string]struct{} // 已接收过的 job_key（幂等）
	jobs chan Job
	acks chan Result
}

// New 创建队列。
func New(resolver PrinterResolver, maxRetry int, retryGap time.Duration) *Queue {
	if maxRetry <= 0 {
		maxRetry = 3
	}
	return &Queue{
		resolver: resolver,
		maxRetry: maxRetry,
		retryGap: retryGap,
		seen:     make(map[string]struct{}),
		jobs:     make(chan Job, 256),
		acks:     make(chan Result, 256),
	}
}

// Acks 返回回执通道（供 sync 模块上报云端）。
func (q *Queue) Acks() <-chan Result { return q.acks }

// Enqueue 入队；重复 job_key 直接丢弃（防重打），返回是否真正入队。
func (q *Queue) Enqueue(j Job) bool {
	q.mu.Lock()
	if _, dup := q.seen[j.Key]; dup {
		q.mu.Unlock()
		return false
	}
	q.seen[j.Key] = struct{}{}
	q.mu.Unlock()

	q.jobs <- j
	return true
}

// Run 启动 worker，处理任务直到 ctx 取消。
func (q *Queue) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case j := <-q.jobs:
			q.acks <- q.process(ctx, j)
		}
	}
}

func (q *Queue) process(ctx context.Context, j Job) Result {
	printer, ok := q.resolver(j.PrinterID)
	if !ok {
		return Result{Key: j.Key, Err: errPrinterNotFound}
	}
	var lastErr error
	for attempt := 0; attempt < q.maxRetry; attempt++ {
		if attempt > 0 && q.retryGap > 0 {
			select {
			case <-ctx.Done():
				return Result{Key: j.Key, Retries: attempt, Err: ctx.Err()}
			case <-time.After(q.retryGap):
			}
		}
		if err := printer.Print(ctx, j.Payload); err == nil {
			return Result{Key: j.Key, Printed: true, Retries: attempt}
		} else {
			lastErr = err
		}
	}
	return Result{Key: j.Key, Retries: q.maxRetry, Err: lastErr}
}
