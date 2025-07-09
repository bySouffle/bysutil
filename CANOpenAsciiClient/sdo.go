package canopenasciiclient

import (
	"fmt"
	"time"
)

//Datatypes:
//b                  # Boolean.
//i8, i16, i32, i64  # Signed integers.
//u8, u16, u32, u64  # Unsigned integers.
//x8, x16, x32, x64  # Unsigned integers, displayed as hexadecimal, non-standard.
//r32, r64           # Real numbers.
//t, td              # Time of day, time difference.
//vs                 # Visible string (between double quotes if multi-word).
//os, us             # Octet, unicode string, (mime-base64 (RFC2045) based, line).
//d                  # domain (mime-base64 (RFC2045) based, one line).
//hex                # Hexagonal data, optionally space separated, non-standard.

const (
	i8  = "i8"
	i16 = "i16"
	i32 = "i32"
	i64 = "i64"

	u8  = "u8"
	u16 = "u16"
	u32 = "u32"
	u64 = "u64"

	x8  = "x8"
	x16 = "x16"
	x32 = "x32"
	x64 = "x64"

	r32 = "r32"
	r64 = "r64"

	t  = "t"
	td = "td"

	vs = "vs"

	os = "os"
	us = "us"

	d = "d"

	hex = "hex"
)

//Command strings start with '"["<sequence>"]"' followed by:
//[<node>] r[ead] <index> <subindex> [<datatype>]        # SDO upload.
//[<node>] w[rite] <index> <subindex> <datatype> <value> # SDO download.

// ReadSDO 生成读取 SDO 的命令
func ReadSDO(cmdID int, nodeID, index, subindex int, dataType dataType, timeout time.Duration, retries int) Command {
	return Command{
		CommandID:  cmdID,
		NodeID:     nodeID,
		Index:      index,
		SubIndex:   subindex,
		DataType:   dataType,
		Content:    fmt.Sprintf("[%d] %d r %#x %x %s ", cmdID, nodeID, index, subindex, dataType),
		Timeout:    timeout,
		MaxRetries: retries,
	}
}

// WriteSDO 生成写入 SDO 的命令
func WriteSDO(cmdID int, nodeID, index, subindex int, value, dataType dataType, timeout time.Duration, retries int) Command {
	return Command{
		CommandID:  cmdID,
		Content:    fmt.Sprintf("[%d] %d w %#x %x %s %s", cmdID, nodeID, index, subindex, dataType, value),
		Timeout:    timeout,
		MaxRetries: retries,
	}
}
