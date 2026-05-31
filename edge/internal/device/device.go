// Package device 管理门店网络打印机/标签机及其状态采集。
package device

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Status 打印机运行状态。
type Status struct {
	Online   bool
	PaperOut bool
	Message  string
}

// Printer 抽象一台门店打印机（ESC/POS 收据机，P2 扩展 TSPL 标签机）。
type Printer interface {
	ID() int64
	// Print 将已渲染好的字节流发送到打印机。
	Print(ctx context.Context, payload []byte) error
	// Probe 探测在线状态。
	Probe(ctx context.Context) Status
}

// NetworkPrinter 通过 TCP(默认 9100) 直连的网络打印机。
type NetworkPrinter struct {
	id   int64
	addr string // ip:port
}

// NewNetworkPrinter 创建网络打印机；addr 形如 "192.168.1.51:9100"。
func NewNetworkPrinter(id int64, addr string) *NetworkPrinter {
	return &NetworkPrinter{id: id, addr: addr}
}

func (p *NetworkPrinter) ID() int64 { return p.id }

func (p *NetworkPrinter) Print(ctx context.Context, payload []byte) error {
	d := net.Dialer{Timeout: 3 * time.Second}
	conn, err := d.DialContext(ctx, "tcp", p.addr)
	if err != nil {
		return fmt.Errorf("dial printer %s: %w", p.addr, err)
	}
	defer conn.Close()
	if _, err := conn.Write(payload); err != nil {
		return fmt.Errorf("write printer %s: %w", p.addr, err)
	}
	return nil
}

func (p *NetworkPrinter) Probe(ctx context.Context) Status {
	d := net.Dialer{Timeout: 2 * time.Second}
	conn, err := d.DialContext(ctx, "tcp", p.addr)
	if err != nil {
		return Status{Online: false, Message: err.Error()}
	}
	_ = conn.Close()
	return Status{Online: true}
}
