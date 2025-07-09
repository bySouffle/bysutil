package mutual

import (
	"fmt"
	"reflect"
	"time"
)

type KV map[string]interface{}

type DuplexChan struct {
	Req         chan any
	Resp        chan any
	reqMsgType  reflect.Type
	respMsgType reflect.Type
}

func NewDuplexChan(reqSize int, respSize int, reqType any, respType any) DuplexChan {
	return DuplexChan{
		Req:         make(chan any, reqSize),
		Resp:        make(chan any, respSize),
		reqMsgType:  reflect.TypeOf(reqType),
		respMsgType: reflect.TypeOf(respType),
	}
}

func (c *DuplexChan) Clean() {
	for len(c.Resp) > 0 {
		<-c.Resp
	}
	for len(c.Req) > 0 {
		<-c.Req
	}
}

func (c *DuplexChan) IsRespType(msg any) bool {
	return reflect.TypeOf(msg) == c.respMsgType
}

func (c *DuplexChan) IsReqType(msg any) bool {
	return reflect.TypeOf(msg) == c.reqMsgType
}

func (c *DuplexChan) RespReceive(duration time.Duration) (any, error) {
	select {
	case msg := <-c.Resp:
		if reflect.TypeOf(msg) == c.respMsgType {
			return msg, nil
		} else {
			return msg, fmt.Errorf("消息格式错误")
		}
	case <-time.After(duration):
		return nil, fmt.Errorf("接收超时")
	}
}

func (c *DuplexChan) RespSend(msg any, duration time.Duration) error {
	if reflect.TypeOf(msg) != c.respMsgType {
		return fmt.Errorf("消息格式错误")
	}
	select {
	case c.Resp <- msg:
		return nil
	case <-time.After(duration):
		return fmt.Errorf("接收超时")
	}
}

func (c *DuplexChan) ReqReceive(duration time.Duration) (any, error) {
	select {
	case msg := <-c.Req:
		if reflect.TypeOf(msg) == c.reqMsgType {
			return msg, nil
		} else {
			return msg, fmt.Errorf("消息格式错误")
		}
	case <-time.After(duration):
		return nil, fmt.Errorf("接收超时")
	}
}

func (c *DuplexChan) ReqSend(msg any, duration time.Duration) error {
	if reflect.TypeOf(msg) != c.reqMsgType {
		return fmt.Errorf("消息格式错误")
	}
	select {
	case c.Req <- msg:
		return nil
	case <-time.After(duration):
		return fmt.Errorf("接收超时")
	}
}
