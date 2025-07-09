package canopenasciiclient

import (
	"fmt"
	"testing"
	"time"
)

func TestNewCANOpenClient(t *testing.T) {
	client := NewCANOpenClient("localhost", 6543, 5*time.Second)
	err := client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	if client.connected {
		// 准备多条指令
		commands := []Command{
			ReadSDO(1, 4, 0x1000, 0, u8, 200*time.Millisecond, 2),
			ReadSDO(1, 4, 0x1000, 0, u16, 200*time.Millisecond, 2), // 读取 SDO，超时 200ms，重试 2 次 // 读取 SDO，超时 200ms，重试 2 次
			ReadSDO(1, 4, 0x1000, 0, u32, 200*time.Millisecond, 2), // 读取 SDO，超时 200ms，重试 2 次 // 读取 SDO，超时 200ms，重试 2 次

			WriteSDO(2, 4, 0x6040, 0, "8", "u16", 300*time.Millisecond, 1), // 写入 SDO，超时 300ms，重试 1 次
			SetNMTState(3, 4, MNT_Start, 100*time.Millisecond, 0),          // 设置 NMT，超时 100ms，无重试
		}

		// 并发发送指令
		results := client.SendMultipleCommands(commands)

		// 打印结果
		for _, result := range results {
			if result.Error != nil {
				fmt.Printf("Command %d failed: %v\n", result.Command.CommandID, result.Error)
			} else {
				resp := result.Response
				if resp.IsError {
					fmt.Printf("Command %d error response: %s %s\n", result.Command.CommandID, resp.ErrorCode, resp.ErrorMsg)
				} else {
					fmt.Printf("Command %d success: %s\n", result.Command.CommandID, resp.Value)
				}
			}
		}

		// 断开连接
		client.Disconnect()
	}
}
