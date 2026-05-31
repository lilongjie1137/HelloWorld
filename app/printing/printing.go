// Command printing 是打印路由服务（rpc + Kafka consumer）入口骨架。
//
// 设计：消费 trade 服务发出的 order.created 事件 → 用 domain/printing 按出品部门类型拆单
// → 查询门店 printer_dept 绑定 → 生成 print_job(job_key 幂等) → 下发门店边缘网关。
// 本文件为装配骨架，Kafka/RPC/DB 接入为 TODO。
package main

import (
	"log"

	"github.com/lilongjie1137/HelloWorld/common/types"
	"github.com/lilongjie1137/HelloWorld/domain/printing"
)

func main() {
	// TODO: 接入 Kafka 消费 order.created；从 DB 读取门店 printer_dept 绑定。
	// 下方为拆单路由的演示装配，证明领域逻辑可用。
	demo := []printing.Item{
		{OrderItemID: 1, Name: "鱼香肉丝", DeptType: types.DeptKitchen},
		{OrderItemID: 2, Name: "珍珠奶茶", DeptType: types.DeptBar},
	}
	bindings := printing.Bindings{types.DeptKitchen: 1, types.DeptBar: 2}
	jobs, unrouted := printing.BuildJobs("O-SH001-20260530-0001", demo, bindings)
	for _, j := range jobs {
		log.Printf("print_job key=%s printer=%d type=%s items=%d", j.Key, j.PrinterID, j.JobType, len(j.Items))
	}
	if len(unrouted) > 0 {
		log.Printf("WARN unrouted depts: %v", unrouted)
	}
	log.Println("printing service skeleton: replace demo with Kafka consumer + DB bindings")
}
