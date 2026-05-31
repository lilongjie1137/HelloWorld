// Package sync 负责边缘网关与云端的双向同步：
//   - 下行：拉取云端打印任务并入本地队列；
//   - 上行：网络恢复后按 orderNo+opSeq 幂等上报离线订单，并回报打印回执。
package sync

import (
	"context"
	"time"

	"github.com/lilongjie1137/HelloWorld/edge/internal/offline"
	"github.com/lilongjie1137/HelloWorld/edge/internal/printq"
)

// CloudClient 与云端交互的抽象（HTTP/WSS 实现注入）。
type CloudClient interface {
	// FetchPrintJobs 拉取待执行打印任务。
	FetchPrintJobs(ctx context.Context, storeID int64) ([]printq.Job, error)
	// UploadOrder 幂等上报离线订单（orderNo+opSeq 去重由云端保证）。
	UploadOrder(ctx context.Context, rec offline.OrderRecord) error
	// AckPrint 回报打印结果。
	AckPrint(ctx context.Context, res printq.Result) error
}

// Syncer 周期性同步器。
type Syncer struct {
	storeID  int64
	cloud    CloudClient
	store    offline.Store
	queue    *printq.Queue
	interval time.Duration
}

// New 创建同步器。
func New(storeID int64, cloud CloudClient, store offline.Store, queue *printq.Queue, interval time.Duration) *Syncer {
	if interval <= 0 {
		interval = 3 * time.Second
	}
	return &Syncer{storeID: storeID, cloud: cloud, store: store, queue: queue, interval: interval}
}

// Run 启动同步循环 + 回执上报，直到 ctx 取消。
func (s *Syncer) Run(ctx context.Context) {
	go s.forwardAcks(ctx)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.pullJobs(ctx)
			s.pushOffline(ctx)
		}
	}
}

// pullJobs 下行：拉取打印任务入本地队列（队列自带幂等防重打）。
func (s *Syncer) pullJobs(ctx context.Context) {
	jobs, err := s.cloud.FetchPrintJobs(ctx, s.storeID)
	if err != nil {
		return // 弱网：下次重试
	}
	for _, j := range jobs {
		s.queue.Enqueue(j)
	}
}

// pushOffline 上行：网络恢复后幂等上报本地离线订单。
func (s *Syncer) pushOffline(ctx context.Context) {
	pending, err := s.store.PendingUploads(ctx)
	if err != nil {
		return
	}
	for _, rec := range pending {
		if err := s.cloud.UploadOrder(ctx, rec); err == nil {
			_ = s.store.MarkSynced(ctx, rec.OrderNo)
		}
	}
}

// forwardAcks 把打印回执上报云端。
func (s *Syncer) forwardAcks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case res := <-s.queue.Acks():
			_ = s.cloud.AckPrint(ctx, res)
		}
	}
}
