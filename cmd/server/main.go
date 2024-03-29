package main

import (
	"errors"
	"fmt"
	"os"
	"shoppinglistserver/api"
	"shoppinglistserver/constants"
	"shoppinglistserver/log"
	"shoppinglistserver/storage"
	"shoppinglistserver/utils"

	"github.com/urfave/cli"
)

const (
	API_PORT int = 5454

	LOG_LEVEL_FLAG        = "loglevel"
	PORT_FLAG             = "port"
	DB_URL_FLAG           = "dbUrl"
	SECRET_ENDPOINT_FLAG  = "secretEndpoint"
	SECRET_BEARER_FLAG    = "secretBearer"
	REPLACEMENTS_FLAG     = "replacements"
	TODOIST_APP_ID_FLAG   = "todoistAppId"
	TODOIST_ENDPOINT_FLAG = "todoistEndpoint"
)

func initLogger(ctx *cli.Context) {
	logLevel := ctx.GlobalString(LOG_LEVEL_FLAG)
	log.InitLogger(logLevel)
}

func initStorage(ctx *cli.Context) {
	dbPath := ctx.GlobalString(DB_URL_FLAG)
	if err := storage.InitStorage(dbPath); err != nil {
		log.Logger.Fatalf("Could not create the storage: %+v", err)
	}
}

func getTodoistConfig(ctx *cli.Context) api.TodoistConfig {
	config := api.TodoistConfig{}
	appId := ctx.GlobalString(TODOIST_APP_ID_FLAG)
	endpoint := ctx.GlobalString(TODOIST_ENDPOINT_FLAG)

	if appId != "" && endpoint != "" {
		config.Enabled = true
		config.AppId = appId
		config.Endpoint = endpoint
	}

	return config
}

func start(config *api.ApiConfig) error {
	return api.Run(config)
}

func main() {
	app := cli.NewApp()
	app.Name = constants.APP_NAME
	app.Version = constants.VERSION
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   LOG_LEVEL_FLAG,
			Usage:  "Log Level [TRACE, DEBUG, INFO, WARN, ERROR]",
			EnvVar: "LOGLEVEL",
			Value:  "INFO",
		},
		cli.StringFlag{
			Name:     DB_URL_FLAG,
			Usage:    "Database URL connection (either sqlite3://PATH_TO_DB or PostgreSQL connection string)",
			EnvVar:   "DB_URL",
			Required: true,
		},
		cli.StringFlag{
			Name:     SECRET_ENDPOINT_FLAG,
			Usage:    "Secret endpoint without authorization (useful for 3rd party integrations)",
			EnvVar:   "SECRET_ENDPOINT",
			Value:    "",
			Required: false,
		},
		cli.StringFlag{
			Name:     SECRET_BEARER_FLAG,
			Usage:    "Secret bearer authorization in order to secure your server (Header: Authorization: Bearer YOUR_SECRET)",
			EnvVar:   "SECRET_BEARER",
			Value:    "",
			Required: false,
		},
		cli.IntFlag{
			Name:   PORT_FLAG,
			Usage:  "Port where the server will listen",
			EnvVar: "PORT",
			Value:  API_PORT,
		},
		cli.StringFlag{
			Name:     REPLACEMENTS_FLAG,
			Usage:    "String for configuring any needed replacements for fixing third-party spelling mistages (Format: source=dest,c=d)",
			EnvVar:   "REPLACEMENTS",
			Value:    "",
			Required: false,
		},
		cli.StringFlag{
			Name:     TODOIST_APP_ID_FLAG,
			Usage:    "String containing the Todoist APP ID for enabling todoist integration",
			EnvVar:   "TODOIST_APP_ID",
			Value:    "",
			Required: false,
		},
		cli.StringFlag{
			Name:     TODOIST_ENDPOINT_FLAG,
			Usage:    "String containing the Todoist endpoint for enabling todoist integration",
			EnvVar:   "TODOIST_ENDPOINT",
			Value:    "",
			Required: false,
		},
	}
	app.Action = func(c *cli.Context) error {
		initLogger(c)
		initStorage(c)

		replacements := c.GlobalString(REPLACEMENTS_FLAG)
		err := utils.SetReplacements(replacements)
		if err != nil {
			return errors.New(fmt.Sprintf("startup error: invalid replacements: %+v", err))
		}

		config := &api.ApiConfig{
			Port:           c.Int(PORT_FLAG),
			SecretEndpoint: c.GlobalString(SECRET_ENDPOINT_FLAG),
			SecretBearer:   c.GlobalString(SECRET_BEARER_FLAG),
			TodoistConfig:  getTodoistConfig(c),
		}
		return start(config)
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Logger.Errorf("Error running the application: %+v", err)
		os.Exit(1)
	}
}
