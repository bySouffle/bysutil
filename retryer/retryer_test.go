package retryer

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMaxRetries(t *testing.T) {
	retryVal := 0

	err := Retry(func() (RetrySignal, error) {
		retryVal++
		return FuncFailure, nil
	}, 8, 100*time.Millisecond, 100*time.Millisecond)
	assert.Error(t, err) // Exceeding max-retries is an error.

	assert.Equal(t, 8, retryVal)
}

func TestFailRetry(t *testing.T) {
	retryVal := 0
	err := Retry(func() (RetrySignal, error) {
		retryVal++
		println(retryVal)
		if retryVal < 7 {
			//return FuncFailure, fmt.Errorf("fail retry") //	if return error, will break retry
			return FuncFailure, nil //	if return error, will break retry
		}
		if retryVal == 7 {
			return FuncError, fmt.Errorf("fail retry") //	if return error, will break retry
		}
		return FuncComplete, nil
	}, 7, 100*time.Millisecond, 100*time.Millisecond)

	if err != nil {
		println(retryVal, err.Error())
	}
}
