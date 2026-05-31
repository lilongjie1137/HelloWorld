// Package config 定义门店边缘网关的运行配置。
package config

// Config 边缘网关配置（门店迷你主机本地 etc/edge.yaml 或环境变量注入）。
type Config struct {
	StoreID      int64
	StoreCode    string
	SN           string // 网关序列号，注册到云端
	CloudBaseURL string // 云端地址（HTTPS/WSS）
	ListenAddr   string // 局域网内给 Android POS 直连
	MaxRetry     int    // 打印任务最大重试次数
	Printers     []Printer
}

// Printer 门店打印机配置。
type Printer struct {
	ID    int64
	Name  string
	IP    string   // ip:9100
	Proto string   // ESC_POS / TSPL
	Type  string   // RECEIPT / LABEL
	Depts []string // 绑定的出品部门类型：KITCHEN/BAR/COLD
}
