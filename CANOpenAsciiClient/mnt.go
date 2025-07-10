package CANOpenAscii

import (
	"fmt"
	"time"
)

//[<node>] start                   # NMT Start node.
//[<node>] stop                    # NMT Stop node.
//[<node>] preop[erational]        # NMT Set node to pre-operational.
//[<node>] reset node              # NMT Reset node.
//[<node>] reset comm[unication]   # NMT Reset communication.

type MNTOpt string

const (
	MNT_Start     = "start"
	MNT_Stop      = "stop"
	MNT_Preop     = "preop"
	MNT_ResetNode = "reset node"
	MNT_ResetComm = "reset comm"
)

// SetNMTState 生成设置 NMT 状态的命令
func SetNMTState(cmdID int, nodeID int, state MNTOpt, timeout time.Duration, retries int) Command {
	return Command{
		CommandID:  cmdID,
		Content:    fmt.Sprintf("[%d] %d %s", cmdID, nodeID, state),
		Timeout:    timeout,
		MaxRetries: retries,
	}
}
