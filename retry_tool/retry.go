package retry_tool

import (
	"context"
	"time"
)

func Retry(ctx context.Context, f func(ctx context.Context) error, retryTimes int) error {

	if retryTimes == 0 {
		retryTimes = 1
	}
	var err error
	for i := 0; i < retryTimes; i++ {
		err = f(ctx)
		if err == nil { // 成功则直接结束
			return nil
		}
		time.Sleep(time.Millisecond * 100)
	}
	return err
}

func RetryWithInterval(ctx context.Context, f func(ctx context.Context) error, retryTimes int, interval time.Duration) error {

	if retryTimes == 0 {
		retryTimes = 1
	}
	var err error
	for i := 0; i < retryTimes; i++ {
		err = f(ctx)
		if err == nil { // 成功则直接结束
			return nil
		}
		time.Sleep(interval)
	}
	return err
}
