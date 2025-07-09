package cronjob

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/types/known/durationpb"
	"testing"
	"time"
)

func TestCron(t *testing.T) {

	ctx := context.Background()
	cronjobServer := FNewCronServer()
	go cronjobServer.Start(ctx)
	cronConfigRepo := NewCronConfigRepo(FNewRedis(), map[string]string{"test": "@every 5s"})
	cronConfigProvider := FNewCronConfigProvider(cronConfigRepo)
	periodicManager := FNewCronManager(log.DefaultLogger, cronConfigProvider)
	go periodicManager.Start(ctx)
	time.Sleep(10 * time.Second)

	go periodicManager.Stop(ctx)
	cronjobServer.Stop(ctx)
}

func FNewCronServer() *Server {
	cc := &CronConfig{
		Addr:         "127.0.0.1:16379",
		Db:           6,
		Password:     "root",
		DialTimeout:  durationpb.New(time.Millisecond * 500),
		ReadTimeout:  durationpb.New(time.Millisecond * 500),
		WriteTimeout: durationpb.New(time.Millisecond * 500),
		MinIdleConn:  200,
		PoolSize:     100,
		PoolTimeout:  durationpb.New(time.Second * 240),
		Concurrency:  10,
	}

	cron := NewCronServer(cc, log.DefaultLogger)

	cb1 := func(context.Context, *asynq.Task) error {
		fmt.Print("cb1=>", time.Now())
		err := errors.New("")
		if err != nil {
			//	Notice SkipRetry
			return fmt.Errorf(": %v: %w", err, asynq.SkipRetry)
		}
		return nil
	}

	RegisterCronHandler(cron, NewHandler("test", cb1))
	return cron
}

func FNewCronConfigProvider(repo CronConfigRepo) *CronConfigProvider {
	return &CronConfigProvider{
		repo: repo,
	}
}

func FNewCronManager(logger log.Logger, source *CronConfigProvider) *PeriodicManager {
	cc := &CronConfig{
		Addr:         "127.0.0.1:16379",
		Db:           6,
		Password:     "root",
		DialTimeout:  durationpb.New(time.Millisecond * 500),
		ReadTimeout:  durationpb.New(time.Millisecond * 500),
		WriteTimeout: durationpb.New(time.Millisecond * 500),
		MinIdleConn:  200,
		PoolSize:     100,
		PoolTimeout:  durationpb.New(time.Second * 240),
		Concurrency:  10,
	}
	return NewPeriodicManager(cc, logger, source)
}

func FNewRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:16379",
		DB:           6,
		Password:     "root",
		DialTimeout:  time.Millisecond * 500,
		ReadTimeout:  time.Millisecond * 500,
		WriteTimeout: time.Millisecond * 500,
	})
	return rdb
}
