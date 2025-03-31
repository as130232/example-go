package job

import (
	"example-go/cmd"
	"example-go/common/global"
	"github.com/robfig/cron/v3"
)

func Init(app *cmd.App) *cron.Cron {
	c := cron.New(cron.WithSeconds())

	if global.ServerConfig.AppEnv != "local" {
		//_, err := c.AddFunc("0/30 * * * * *", NewMetricJob().Execute)
		//if err != nil {
		//	panic(err)
		//}
	}

	c.Start()

	return c
}
