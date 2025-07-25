package log

import (
	"github.com/bySouffle/bysutil/utils"
)

// GetLogFile get log file absolute path
func GetLogFile(filename string, suffix string) string {
	return utils.ConcatString(logDir, "/", hostname, "/", filename, suffix)
}
