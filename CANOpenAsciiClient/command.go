package canopenasciiclient

import "time"

type dataType string

// Command 代表一条 CANOpen 指令
type Command struct {
	CommandID  int // 唯一标识符
	NodeID     int
	Index      int
	SubIndex   int
	DataType   dataType
	Content    string        // 指令内容，如 "[1] r 1000 0"
	Timeout    time.Duration // 命令超时
	MaxRetries int           // 最大重试次数
}

// Response 代表解析后的响应
type Response struct {
	CommandID string // 对应的命令 CommandID
	Value     string // 成功时的值（如 SDO 读取结果）
	IsError   bool   // 是否为错误响应
	ErrorCode string // 错误代码（如 "100"）
	ErrorMsg  string // 错误消息（如 "Invalid parameter"）
	Raw       string // 原始响应
}

// CommandResult 代表指令执行结果
type CommandResult struct {
	Command  Command
	Response Response
	Error    error
}
