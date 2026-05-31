package idgen

import (
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	day := time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC)
	got := Format(BizOrder, "SH001", day, 7)
	want := "O-SH001-20260530-0007"
	if got != want {
		t.Errorf("Format = %q, want %q", got, want)
	}
}

func TestSequencer_MonotonicPerKey(t *testing.T) {
	s := NewSequencer()
	day := time.Date(2026, 5, 30, 0, 0, 0, 0, time.UTC)

	if got := s.Next(BizOrder, "SH001", day); got != "O-SH001-20260530-0001" {
		t.Errorf("first = %q", got)
	}
	if got := s.Next(BizOrder, "SH001", day); got != "O-SH001-20260530-0002" {
		t.Errorf("second = %q", got)
	}
	// 不同门店独立计数。
	if got := s.Next(BizOrder, "BJ002", day); got != "O-BJ002-20260530-0001" {
		t.Errorf("other store = %q", got)
	}
	// 不同业务独立计数。
	if got := s.Next(BizBill, "SH001", day); got != "B-SH001-20260530-0001" {
		t.Errorf("other biz = %q", got)
	}
}

func TestSequencer_ResetsByDay(t *testing.T) {
	s := NewSequencer()
	d1 := time.Date(2026, 5, 30, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC)
	_ = s.Next(BizOrder, "SH001", d1)
	if got := s.Next(BizOrder, "SH001", d2); got != "O-SH001-20260531-0001" {
		t.Errorf("new day = %q, want reset to 0001", got)
	}
}
