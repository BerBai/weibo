package main

import (
	"github.com/berbai/weibo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
	"time"
)

var (
	logger = log.New(os.Stdout, "weibo: ", log.Ldate|log.Lmicroseconds)
)

type App struct {
	cli      *cli.App
	client   *weibo.Client
	database *weibo.Database
	full     bool
	tz       string
	spec     string
	userid   string
	page     int
	sleep    int
}

func (app *App) Run() error {
	app.client = &weibo.Client{}
	app.database = &weibo.Database{}
	app.cli = &cli.App{
		Usage: "weibo collector and monitor",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "cookie",
				Aliases:     []string{"c"},
				Value:       "SUB=_2AkMTKKlJf8NxqwFRmP8RzWLkbY10zwrEieKldFiSJRMxHRl-yT9kqlM8tRB6OKiHpmrIgcUy6YQdWlF4Q9LVcDAvvpWG;",
				Usage:       "client cookie",
				Destination: &app.client.Cookie,
				EnvVars:     []string{"WEIBO_COLLECTOR_COOKIE"},
			},
			&cli.StringFlag{
				Name:        "userid",
				Aliases:     []string{"u"},
				Value:       "1223178222",
				Usage:       "userid list",
				Destination: &app.userid,
				EnvVars:     []string{"WEIBO_COLLECTOR_USERID"},
			},
			&cli.IntFlag{
				Name:        "page",
				Aliases:     []string{"p"},
				Value:       1,
				Usage:       "start page",
				Destination: &app.page,
				EnvVars:     []string{"WEIBO_COLLECTOR_PAGE"},
			},
			&cli.IntFlag{
				Name:        "sleep",
				Aliases:     []string{"s"},
				Value:       10,
				Usage:       "sleep seconds between page",
				Destination: &app.sleep,
				EnvVars:     []string{"WEIBO_COLLECTOR_SLEEP"},
			},
			&cli.StringFlag{
				Name:        "dn",
				Aliases:     []string{"d"},
				Value:       "mysql",
				Usage:       "database driver name",
				Destination: &app.database.DN,
				EnvVars:     []string{"WEIBO_COLLECTOR_DN"},
			},
			&cli.StringFlag{
				Name:        "dsn",
				Value:       "root:root@/weibo",
				Usage:       "database source name",
				Destination: &app.database.DSN,
				EnvVars:     []string{"WEIBO_COLLECTOR_DSN"},
			},
			&cli.BoolFlag{
				Name:        "full",
				Aliases:     []string{"f"},
				Value:       false,
				Usage:       "full collect",
				Destination: &app.full,
				EnvVars:     []string{"WEIBO_COLLECTOR_FULl"},
			},
			&cli.StringFlag{
				Name:        "cron",
				Value:       "*/1 * * * *",
				Usage:       "cron spec",
				Destination: &app.spec,
				EnvVars:     []string{"WEIBO_COLLECTOR_SPEC"},
			},
			&cli.StringFlag{
				Name:        "tz",
				Aliases:     []string{"t"},
				Value:       "Local",
				Usage:       "time zone",
				Destination: &app.tz,
				EnvVars:     []string{"WEIBO_COLLECTOR_TZ"},
			},
		},
		Action: app.run,
	}
	print(app.cli)
	return app.cli.Run(os.Args)
}

func (app *App) run(c *cli.Context) error {
	if err := app.database.Migrate(); err != nil {
		return err
	}
	if app.full {
		logger.Printf("full collecting.")
		if _, err := app.collect(true); err != nil {
			return err
		}
		logger.Printf("full collecting finished.")
	}
	return app.cron()
}

func (app *App) cron() error {
	logger.Printf("monitoring.")
	c := cron.New(
		cron.WithLocation(location(app.tz)),
		cron.WithLogger(cron.VerbosePrintfLogger(logger)),
		cron.WithChain(cron.SkipIfStillRunning(cron.VerbosePrintfLogger(logger))),
	)
	c.AddFunc(app.spec, app.monitoring)
	c.Run()

	return nil
}

func (app *App) collect(full bool) ([]*weibo.Mblog, error) {
	page := 1
	if full {
		page = 98
	}

	var mblogs []*weibo.Mblog
	for _, userid := range strings.Split(app.userid, ",") {
		for i := app.page; i <= page; i++ {
			logger.Printf("collecting. userid=%s, page=%d", userid, i)
			_mblogs, err := app.client.GetMblogs(userid, i, true)
			if err != nil {
				return nil, err
			}

			for _, mblog := range _mblogs {
				if has, err := app.database.HasMblog(mblog); err != nil {
					return nil, err
				} else if has {
					continue
				}

				if err := app.database.AddMblog(mblog); err != nil {
					return nil, err
				}
				mblogs = append(mblogs, mblog)
			}
			logger.Printf("sleep %d seconds", app.sleep)
			time.Sleep(time.Duration(app.sleep) * time.Second)
		}
	}
	return mblogs, nil
}

func (app *App) monitoring() {
	if mblogs, err := app.collect(false); err != nil {
		logger.Printf("monitoring, err='%s'\n", err)
	} else {
		if len(mblogs) > 0 {
			logger.Printf("monitoring found new weibo mblog.")
			app.notification(mblogs)
		}
	}
}

func (app *App) notification(mblogs []*weibo.Mblog) {
	logger.Printf("send notification.")
}

func location(tz string) *time.Location {
	if loc, err := time.LoadLocation(tz); err != nil {
		return time.Local
	} else {
		return loc
	}
}

func main() {
	if err := (&App{}).Run(); err != nil {
		logger.Printf("run, err='%s'\n", err)
	}
}
