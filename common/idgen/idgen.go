// Package idgen 生成业务单号。
//
// 单号规则：{业务前缀}{门店码}{yyyyMMdd}{当日自增序号(4位起)}
// 例：堂食订单 O-SH001-20260530-0007。
// 序号在「门店 + 业务 + 自然日」维度内单调递增；离线场景由边缘网关用门店本地序列生成，
// 携带门店码前缀天然避免跨店冲突，上行云端按 order_no 幂等去重。
package idgen

import (
	"fmt"
	"sync"
	"time"
)

// Biz 业务前缀。
type Biz string

const (
	BizOrder Biz = "O" // 订单
	BizBill  Biz = "B" // 账单
	BizPrint Biz = "P" // 打印任务
)

// Sequencer 按 (业务,门店,日期) 维度产出单调递增序号。
// 进程内线程安全；分布式/多实例场景应替换为 Redis INCR 或边缘本地序列实现 Source。
type Sequencer struct {
	mu  sync.Mutex
	seq map[string]int64
}

// NewSequencer 创建内存序列器（开发/边缘单机用）。
func NewSequencer() *Sequencer {
	return &Sequencer{seq: make(map[string]int64)}
}

// Next 返回下一个单号。day 用于按自然日归零，便于测试注入。
func (s *Sequencer) Next(biz Biz, storeCode string, day time.Time) string {
	dateStr := day.Format("20060102")
	key := fmt.Sprintf("%s:%s:%s", biz, storeCode, dateStr)

	s.mu.Lock()
	s.seq[key]++
	n := s.seq[key]
	s.mu.Unlock()

	return Format(biz, storeCode, day, n)
}

// Format 按统一规则拼装单号（纯函数，便于云端用外部序列来源时复用）。
func Format(biz Biz, storeCode string, day time.Time, seq int64) string {
	return fmt.Sprintf("%s-%s-%s-%04d", biz, storeCode, day.Format("20060102"), seq)
}
