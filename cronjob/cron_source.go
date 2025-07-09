package cronjob

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
)

const (
	TaskPeriodPrefix = "cronjob:task_period:schedule"
)

type PeriodicTaskConfigContainer map[string]string

type CronConfigRepo interface {
	GetCron(ctx context.Context) (PeriodicTaskConfigContainer, error)
	AddCron(ctx context.Context, containers PeriodicTaskConfigContainer) error
	UpdateCron(ctx context.Context, containers PeriodicTaskConfigContainer) error
	DelCron(ctx context.Context, topics []string) error
}

type CronConfigProvider struct {
	repo CronConfigRepo
}

func NewCronConfigProvider(repo CronConfigRepo) *CronConfigProvider {
	return &CronConfigProvider{
		repo: repo,
	}
}

func (p *CronConfigProvider) GetConfigs() ([]*asynq.PeriodicTaskConfig, error) {

	var (
		ctx = context.Background()
	)

	c, err := p.repo.GetCron(ctx)
	if err != nil {
		fmt.Printf("获取任务周期失败: value[%v] err[%v]", c, err)
		return nil, err
	}

	var configs []*asynq.PeriodicTaskConfig
	for topic, cron := range c {
		configs = append(configs, &asynq.PeriodicTaskConfig{Cronspec: cron, Task: asynq.NewTask(topic, nil)})
	}
	return configs, nil
}

func (p *CronConfigProvider) AddCron(ctx context.Context, kv PeriodicTaskConfigContainer) error {
	return p.repo.AddCron(ctx, kv)

}

func (p *CronConfigProvider) UpdateCron(ctx context.Context, kv PeriodicTaskConfigContainer) error {
	return p.repo.UpdateCron(ctx, kv)

}

func (p *CronConfigProvider) DelCron(ctx context.Context, topic []string) error {
	return p.repo.DelCron(ctx, topic)

}

func (p *CronConfigProvider) GetCron(ctx context.Context) (PeriodicTaskConfigContainer, error) {
	return p.repo.GetCron(ctx)
}
