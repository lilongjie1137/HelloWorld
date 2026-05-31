// Package printing 实现拆单打印路由：按 order_item 的出品部门类型分组，
// 落到该门店绑定的打印机，生成幂等的打印任务（鱼香肉丝→后厨，珍珠奶茶→吧台，互不干扰）。
package printing

import (
	"fmt"
	"sort"

	"github.com/lilongjie1137/HelloWorld/common/types"
)

// Item 参与拆单的订单明细（仅需出品部门类型）。
type Item struct {
	OrderItemID int64
	Name        string
	DeptType    types.ProductionDeptType
}

// Bindings 门店内出品部门类型 → 打印机 ID。
type Bindings map[types.ProductionDeptType]int64

// Job 生成的打印任务。
type Job struct {
	Key       string // 幂等键：{orderNo}:{deptType}
	PrinterID int64
	DeptType  types.ProductionDeptType
	JobType   types.PrintJobType
	Items     []Item
}

var deptToJobType = map[types.ProductionDeptType]types.PrintJobType{
	types.DeptKitchen: types.JobKitchen,
	types.DeptBar:     types.JobBar,
	types.DeptCold:    types.JobKitchen,
}

// BuildJobs 按出品部门类型把订单明细拆成多张打印任务。
// 同一部门类型的明细合并为一张单；未绑定打印机的部门会被跳过并在 unrouted 中返回。
func BuildJobs(orderNo string, items []Item, bindings Bindings) (jobs []Job, unrouted []types.ProductionDeptType) {
	grouped := map[types.ProductionDeptType][]Item{}
	for _, it := range items {
		grouped[it.DeptType] = append(grouped[it.DeptType], it)
	}

	// 稳定输出顺序，便于测试与打印顺序一致。
	depts := make([]types.ProductionDeptType, 0, len(grouped))
	for d := range grouped {
		depts = append(depts, d)
	}
	sort.Slice(depts, func(i, j int) bool { return depts[i] < depts[j] })

	for _, dept := range depts {
		printerID, ok := bindings[dept]
		if !ok {
			unrouted = append(unrouted, dept)
			continue
		}
		jobType, ok := deptToJobType[dept]
		if !ok {
			jobType = types.JobKitchen
		}
		jobs = append(jobs, Job{
			Key:       fmt.Sprintf("%s:%s", orderNo, dept),
			PrinterID: printerID,
			DeptType:  dept,
			JobType:   jobType,
			Items:     grouped[dept],
		})
	}
	return jobs, unrouted
}
