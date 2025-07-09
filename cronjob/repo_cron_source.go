package cronjob

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type ConfigRepo struct {
	redis *redis.Client
}

func (c *ConfigRepo) GetCron(ctx context.Context) (PeriodicTaskConfigContainer, error) {
	b, err := c.redis.Get(ctx, TaskPeriodPrefix).Bytes()
	if err != nil {
		return nil, err
	}
	container := PeriodicTaskConfigContainer{}
	return container, json.Unmarshal(b, &container)

}

func (c *ConfigRepo) AddCron(ctx context.Context, containers PeriodicTaskConfigContainer) error {
	b, err := c.redis.Get(ctx, TaskPeriodPrefix).Bytes()
	if err != nil {
		return err
	}

	container := PeriodicTaskConfigContainer{}
	if err := json.Unmarshal(b, &container); err != nil {
		fmt.Printf("解析任务周期失败:%v", err)
		return err
	}

	for k, v := range containers {
		container[k] = v
	}

	confB, encErr := json.Marshal(container)
	if encErr != nil {
		return err
	}

	return c.redis.Set(ctx, TaskPeriodPrefix, confB, 0).Err()
}

func (c *ConfigRepo) UpdateCron(ctx context.Context, containers PeriodicTaskConfigContainer) error {
	b, err := c.redis.Get(ctx, TaskPeriodPrefix).Bytes()
	if err != nil {
		return err
	}

	container := PeriodicTaskConfigContainer{}
	if err := json.Unmarshal(b, &container); err != nil {
		fmt.Printf("解析任务周期失败:%v", err)
		return err
	}

	for k, v := range containers {
		container[k] = v
	}

	confB, encErr := json.Marshal(container)
	if encErr != nil {
		return err
	}

	return c.redis.Set(ctx, TaskPeriodPrefix, confB, 0).Err()

}

func (c *ConfigRepo) DelCron(ctx context.Context, topics []string) error {
	b, err := c.redis.Get(ctx, TaskPeriodPrefix).Bytes()
	if err != nil {
		return err
	}

	container := PeriodicTaskConfigContainer{}
	if err := json.Unmarshal(b, &container); err != nil {
		fmt.Printf("解析任务周期失败:%v", err)
		return err
	}
	for _, v := range topics {
		delete(container, v)
	}
	confB, encErr := json.Marshal(container)
	if encErr != nil {
		return err
	}

	return c.redis.Set(ctx, TaskPeriodPrefix, confB, 0).Err()
}

func NewCronConfigRepo(r *redis.Client, TaskPeriod map[string]string) CronConfigRepo {
	var (
		ctx = context.Background()
	)
	b, err := r.Get(ctx, TaskPeriodPrefix).Bytes()
	if err != nil || string(b) == "null" {
		initConf, err := json.Marshal(TaskPeriod)
		if err != nil {
			panic("[server] NewCronConfigProvider init config error")
		}
		if err := r.Set(ctx, TaskPeriodPrefix, initConf, 0).Err(); err != nil {
			panic("[server] NewCronConfigProvider Set config error")
		}
	}

	return &ConfigRepo{
		redis: r,
	}
}
