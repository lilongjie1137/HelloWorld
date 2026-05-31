package printq

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/lilongjie1137/HelloWorld/common/types"
	"github.com/lilongjie1137/HelloWorld/edge/internal/device"
)

// flakyPrinter 前 failTimes 次返回错误，之后成功。
type flakyPrinter struct {
	id        int64
	failTimes int32
	calls     int32
}

func (p *flakyPrinter) ID() int64 { return p.id }
func (p *flakyPrinter) Print(_ context.Context, _ []byte) error {
	n := atomic.AddInt32(&p.calls, 1)
	if n <= p.failTimes {
		return errors.New("printer offline")
	}
	return nil
}
func (p *flakyPrinter) Probe(_ context.Context) device.Status { return device.Status{Online: true} }

func resolverFor(p device.Printer) PrinterResolver {
	return func(id int64) (device.Printer, bool) {
		if id == p.ID() {
			return p, true
		}
		return nil, false
	}
}

func TestEnqueue_Idempotent(t *testing.T) {
	q := New(resolverFor(&flakyPrinter{id: 1}), 3, 0)
	job := Job{Key: "O-SH001-20260530-0001:KITCHEN", PrinterID: 1, Type: types.JobKitchen}
	if !q.Enqueue(job) {
		t.Fatal("first enqueue should succeed")
	}
	if q.Enqueue(job) {
		t.Fatal("duplicate job_key must be dropped")
	}
}

func TestProcess_RetryThenSucceed(t *testing.T) {
	p := &flakyPrinter{id: 1, failTimes: 2}
	q := New(resolverFor(p), 3, 0)
	res := q.process(context.Background(), Job{Key: "k", PrinterID: 1, Type: types.JobBar})
	if !res.Printed {
		t.Fatalf("expected printed after retries, got err=%v", res.Err)
	}
	if res.Retries != 2 {
		t.Errorf("retries = %d, want 2", res.Retries)
	}
}

func TestProcess_ExhaustRetries(t *testing.T) {
	p := &flakyPrinter{id: 1, failTimes: 99}
	q := New(resolverFor(p), 3, 0)
	res := q.process(context.Background(), Job{Key: "k", PrinterID: 1})
	if res.Printed {
		t.Fatal("should fail after exhausting retries")
	}
	if res.Retries != 3 {
		t.Errorf("retries = %d, want 3", res.Retries)
	}
}

func TestProcess_PrinterNotFound(t *testing.T) {
	q := New(resolverFor(&flakyPrinter{id: 1}), 3, 0)
	res := q.process(context.Background(), Job{Key: "k", PrinterID: 999})
	if res.Err == nil {
		t.Fatal("expected error for missing printer")
	}
}
