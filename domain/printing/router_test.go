package printing

import (
	"testing"

	"github.com/lilongjie1137/HelloWorld/common/types"
)

func TestBuildJobs_SplitKitchenAndBar(t *testing.T) {
	items := []Item{
		{OrderItemID: 1, Name: "鱼香肉丝", DeptType: types.DeptKitchen},
		{OrderItemID: 2, Name: "珍珠奶茶", DeptType: types.DeptBar},
		{OrderItemID: 3, Name: "宫保鸡丁", DeptType: types.DeptKitchen},
	}
	bindings := Bindings{types.DeptKitchen: 101, types.DeptBar: 102}

	jobs, unrouted := BuildJobs("O-SH001-20260530-0001", items, bindings)
	if len(unrouted) != 0 {
		t.Fatalf("unrouted = %v, want none", unrouted)
	}
	if len(jobs) != 2 {
		t.Fatalf("jobs = %d, want 2", len(jobs))
	}

	byDept := map[types.ProductionDeptType]Job{}
	for _, j := range jobs {
		byDept[j.DeptType] = j
	}
	k := byDept[types.DeptKitchen]
	if k.PrinterID != 101 || k.JobType != types.JobKitchen || len(k.Items) != 2 {
		t.Errorf("kitchen job wrong: %+v", k)
	}
	if k.Key != "O-SH001-20260530-0001:KITCHEN" {
		t.Errorf("kitchen job_key = %q", k.Key)
	}
	b := byDept[types.DeptBar]
	if b.PrinterID != 102 || b.JobType != types.JobBar || len(b.Items) != 1 {
		t.Errorf("bar job wrong: %+v", b)
	}
}

func TestBuildJobs_UnboundDeptReported(t *testing.T) {
	items := []Item{{OrderItemID: 1, Name: "凉拌黄瓜", DeptType: types.DeptCold}}
	jobs, unrouted := BuildJobs("O1", items, Bindings{types.DeptKitchen: 101})
	if len(jobs) != 0 {
		t.Errorf("jobs = %d, want 0", len(jobs))
	}
	if len(unrouted) != 1 || unrouted[0] != types.DeptCold {
		t.Errorf("unrouted = %v, want [COLD]", unrouted)
	}
}
