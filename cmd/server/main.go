package main

import (
	"os"
	"shoppinglistserver/api"
	"shoppinglistserver/constants"
	"shoppinglistserver/log"
	"shoppinglistserver/storage"

	"github.com/urfave/cli"
)

const (
	API_PORT int = 5454

	LOG_LEVEL_FLAG       = "loglevel"
	PORT_FLAG            = "port"
	DB_URL_FLAG          = "dbUrl"
	SECRET_ENDPOINT_FLAG = "secretEndpoint"
	SECRET_BEARER_FLAG   = "secretBearer"
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

func start(port int, secretEndpoint string, secretBearer string) error {
	return api.Run(port, secretEndpoint, secretBearer)
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
	}
	app.Action = func(c *cli.Context) error {
		initLogger(c)
		initStorage(c)
		port := c.Int(PORT_FLAG)
		secretEndpoint := c.GlobalString(SECRET_ENDPOINT_FLAG)
		secretBearer := c.GlobalString(SECRET_BEARER_FLAG)
		return start(port, secretEndpoint, secretBearer)
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Logger.Errorf("Error running the application: %+v", err)
		os.Exit(1)
	}
}
