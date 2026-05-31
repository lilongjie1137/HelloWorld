// Package offline 负责弱网/断网时在本地落库订单与收银流水，公网恢复后幂等上行。
//
// V1.0 生产实现建议用 BoltDB / SQLite（单机、无依赖）；此处给出接口与内存实现占位，
// 便于上层(sync)联调与单元测试。
package offline

import (
	"context"
	"errors"
	"sync"
)

// OrderRecord 离线订单/收银快照（含单调递增 OpSeq，供云端幂等去重）。
type OrderRecord struct {
	OrderNo string
	OpSeq   int64
	Payload []byte // 序列化后的订单+明细+收银流水
	Synced  bool
}

// Store 离线存储抽象。
type Store interface {
	SaveOrder(ctx context.Context, rec OrderRecord) error
	PendingUploads(ctx context.Context) ([]OrderRecord, error)
	MarkSynced(ctx context.Context, orderNo string) error
}

// ErrNotFound 记录不存在。
var ErrNotFound = errors.New("offline: record not found")

// MemStore 内存实现（开发/测试用）。
type MemStore struct {
	mu   sync.Mutex
	recs map[string]OrderRecord
}

// NewMemStore 创建内存存储。
func NewMemStore() *MemStore {
	return &MemStore{recs: make(map[string]OrderRecord)}
}

func (s *MemStore) SaveOrder(_ context.Context, rec OrderRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 幂等：同 orderNo 仅保留更高 OpSeq。
	if old, ok := s.recs[rec.OrderNo]; ok && old.OpSeq >= rec.OpSeq {
		return nil
	}
	s.recs[rec.OrderNo] = rec
	return nil
}

func (s *MemStore) PendingUploads(_ context.Context) ([]OrderRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []OrderRecord
	for _, r := range s.recs {
		if !r.Synced {
			out = append(out, r)
		}
	}
	return out, nil
}

func (s *MemStore) MarkSynced(_ context.Context, orderNo string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.recs[orderNo]
	if !ok {
		return ErrNotFound
	}
	r.Synced = true
	s.recs[orderNo] = r
	return nil
}
