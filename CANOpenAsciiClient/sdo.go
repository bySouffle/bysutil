package CANOpenAscii

import (
	"fmt"
	"time"
)

//Datatypes:
//b                  # Boolean.
//I8, i16, i32, i64  # Signed integers.
//U8, U16, U32, U64  # Unsigned integers.
//x8, x16, x32, x64  # Unsigned integers, displayed as hexadecimal, non-standard.
//r32, r64           # Real numbers.
//t, td              # Time of day, time difference.
//vs                 # Visible string (between double quotes if multi-word).
//os, us             # Octet, unicode string, (mime-base64 (RFC2045) based, line).
//d                  # domain (mime-base64 (RFC2045) based, one line).
//hex                # Hexagonal data, optionally space separated, non-standard.

const (
	I8  = "I8"
	I16 = "i16"
	I32 = "i32"
	I64 = "i64"

	U8  = "U8"
	U16 = "U16"
	U32 = "U32"
	U64 = "U64"

	X8  = "x8"
	X16 = "x16"
	X32 = "x32"
	X64 = "x64"

	R32 = "r32"
	R64 = "r64"

	T  = "t"
	TD = "td"

	VS = "vs"

	OS = "os"
	US = "us"

	D = "d"

	HEX = "hex"
)

//Command strings start with '"["<sequence>"]"' followed by:
//[<node>] r[ead] <index> <subindex> [<datatype>]        # SDO upload.
//[<node>] w[rite] <index> <subindex> <datatype> <value> # SDO download.

// ReadSDO 生成读取 SDO 的命令
func ReadSDO(cmdID int, nodeID, index, subindex int, dataType DataType, timeout time.Duration, retries int) Command {
	return Command{
		CommandID:  cmdID,
		NodeID:     nodeID,
		Index:      index,
		SubIndex:   subindex,
		DataType:   dataType,
		Content:    fmt.Sprintf("[%d] %d r %#x %x %s", cmdID, nodeID, index, subindex, dataType),
		Timeout:    timeout,
		MaxRetries: retries,
	}
}

// WriteSDO 生成写入 SDO 的命令
func WriteSDO(cmdID int, nodeID, index, subindex int, value, dataType DataType, timeout time.Duration, retries int) Command {
	return Command{
		CommandID:  cmdID,
		Content:    fmt.Sprintf("[%d] %d w %#x %x %s %s", cmdID, nodeID, index, subindex, dataType, value),
		Timeout:    timeout,
		MaxRetries: retries,
	}
}
