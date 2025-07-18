package CANOpenAscii

import (
	"bufio"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"
)

// CO_CANOpen_309_3

// regexp
var (
	// 标准错误响应：! code message
	errorRe = regexp.MustCompile(`! (\d+) (.+)`)
	// 非标准错误响应：[NodeId] ERROR:code #message，支持十进制和十六进制
	nonStandardErrorRe = regexp.MustCompile(`\[(\d+)] ERROR:(0x[0-9a-fA-F]+|\d+) #(.+)`)
	// 成功响应正则表达式
	//	[NodeID] OK
	okRe = regexp.MustCompile(`\[(\d+)] OK`)
	//	[1] r = 1234
	cmdValueRe = regexp.MustCompile(`\[(\d+)] \w+ = (-?\d+|.+)`)
	//
	cmdOkRe = regexp.MustCompile(`\[(\d+)] \w+ OK`)
	//	[1] 01 00 00 00
	hexValuesRe = regexp.MustCompile(`\[(\d+)] ((?:[0-9a-fA-F]{2}\s*)+)`)
	//	[1]	0
	simpleValueRe = regexp.MustCompile(`\[(\d+)] (-?\d+)`)
)

// CANOpenClient 代表 CANOpen 客户端
type CANOpenClient struct {
	conn          net.Conn
	connected     bool
	responses     chan Response
	mutex         sync.Mutex
	host          string
	port          int
	timeout       time.Duration
	retryInterval time.Duration
	wg            sync.WaitGroup
	responseMu    sync.Mutex
}

// NewCANOpenClient 创建新的 CANOpen 客户端
func NewCANOpenClient(host string, port int, timeout, retryInterval time.Duration) *CANOpenClient {
	return &CANOpenClient{
		host:          host,
		port:          port,
		timeout:       timeout,
		retryInterval: retryInterval,
		responses:     make(chan Response, 100), // 响应缓冲区
	}
}

// Connect 连接到 CANOpen 网关
func (c *CANOpenClient) Connect() error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", c.host, c.port), c.timeout)
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}
	c.conn = conn
	c.connected = true
	log.Infof("[CANOpen] Connected to CANOpen gateway at %s:%d\n", c.host, c.port)

	// 启动接收协程
	c.wg.Add(1)
	go c.receive()
	return nil
}

// Disconnect 断开连接
func (c *CANOpenClient) Disconnect() {
	if c.conn != nil {
		c.conn.Close()
		c.connected = false
	}
	close(c.responses)
	c.wg.Wait()
	log.Infof("[CANOpen] Disconnected from CANOpen gateway")
}

// receive 异步接收并解析网关响应
func (c *CANOpenClient) receive() {
	defer c.wg.Done()
	reader := bufio.NewReader(c.conn)
	for c.connected {
		data, err := reader.ReadString('\n')
		if err != nil {
			if c.connected {
				log.Errorf("[CANOpen] Receive error: %v\n", err)
				c.connected = false
			}
			break
		}
		data = strings.TrimRight(data, "\r\n")
		resp := c.parseResponse(data)
		c.responseMu.Lock()
		c.responses <- resp
		c.responseMu.Unlock()
	}
}

// parseResponse 解析 CiA 309-3 响应
func (c *CANOpenClient) parseResponse(data string) Response {
	resp := Response{Raw: data}

	if matches := nonStandardErrorRe.FindStringSubmatch(data); matches != nil {
		resp.IsError = true
		resp.CommandID = matches[1]
		resp.ErrorCode = matches[2]
		resp.ErrorMsg = matches[3]
	} else if matches := errorRe.FindStringSubmatch(data); matches != nil {
		resp.IsError = true
		resp.ErrorCode = matches[1]
		resp.ErrorMsg = matches[2]
	} else if matches := okRe.FindStringSubmatch(data); matches != nil {
		resp.CommandID = matches[1]
		resp.Value = "OK"
	} else if matches := cmdOkRe.FindStringSubmatch(data); matches != nil {
		resp.CommandID = matches[1]
		resp.Value = "OK"
	} else if matches := cmdValueRe.FindStringSubmatch(data); matches != nil {
		resp.CommandID = matches[1]
		resp.Value = matches[2]
	} else if matches := hexValuesRe.FindStringSubmatch(data); matches != nil {
		resp.CommandID = matches[1]
		resp.Value = matches[2]
	} else if matches := simpleValueRe.FindStringSubmatch(data); matches != nil {
		resp.CommandID = matches[1]
		resp.Value = matches[2]
	}
	return resp
}

// SendCommand 发送单条命令，支持重试
func (c *CANOpenClient) SendCommand(cmd Command) (CommandResult, error) {
	if !c.connected {
		errConn := fmt.Errorf("not connected to gateway")
		return CommandResult{Command: cmd, Error: errConn}, errConn
	}

	var resp Response
	var err error
	attempts := 0

	for attempts <= cmd.MaxRetries {
		attempts++
		// 添加 <CR><LF>
		command := fmt.Sprintf("%s\r\n", cmd.Content)
		_, err = c.conn.Write([]byte(command))
		if err != nil {
			sendErr := fmt.Errorf("send error (attempt %d): %v", attempts, err)
			return CommandResult{Command: cmd, Error: sendErr}, sendErr
		}

		// 等待响应
		select {
		case resp = <-c.responses:
			if !resp.IsError {
				return CommandResult{Command: cmd, Response: resp}, nil
			}
			err = fmt.Errorf("response error: %s %s", resp.ErrorCode, resp.ErrorMsg)
		case <-time.After(cmd.Timeout):
			err = fmt.Errorf("timeout after %v", cmd.Timeout)
		}

		if attempts <= cmd.MaxRetries {
			log.Infof("[CANOpen] Retrying command %d (%d/%d) due to: %v\n", cmd.CommandID, attempts, cmd.MaxRetries, err)
			time.Sleep(c.retryInterval) // 重试间隔
		}
	}
	errRetry := fmt.Errorf("failed after %d attempts: %v", attempts-1, err)
	return CommandResult{Command: cmd, Response: resp, Error: errRetry}, errRetry
}

// SendMultipleCommands 并发发送多条命令
func (c *CANOpenClient) SendMultipleCommands(commands []Command) []CommandResult {
	var results []CommandResult
	for _, cmd := range commands {
		result, err := c.SendCommand(cmd)
		results = append(results, result)
		if err != nil {
			log.Errorf("[CANOpen] [Break] Command %d failed: %v\n", cmd.CommandID, err)
			break
		}
	}

	return results
}
