package main

import (
	"github.com/urfave/cli"
	"os"
	"shoppinglistserver/api"
	"shoppinglistserver/constants"
	"shoppinglistserver/log"
	"shoppinglistserver/storage"
)

const (
	API_PORT int = 5454
)

func initLogger(ctx *cli.Context) {
	logLevel := ctx.GlobalString("loglevel")
	log.InitLogger(logLevel)
}

func initStorage(ctx *cli.Context) {
	dbPath := ctx.GlobalString("dbPath")
	if err := storage.InitStorage(dbPath); err != nil {
		log.Logger.Fatalf("Could not create the storage: %+v", err)
	}
}

func start(port int) error {
	return api.Run(port)
}

func main() {
	app := cli.NewApp()
	app.Name = constants.APP_NAME
	app.Version = constants.VERSION
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "loglevel",
			Usage:  "Log Level [TRACE, DEBUG, INFO, WARN, ERROR]",
			EnvVar: "LOGLEVEL",
			Value:  "INFO",
		},
		cli.StringFlag{
			Name:   "dbPath",
			EnvVar: "DB_PATH",
			Value:  "./shopping.sqlite",
		},
		cli.IntFlag{
			Name:   "port",
			EnvVar: "PORT",
			Value:  API_PORT,
		},
	}
	app.Action = func(c *cli.Context) error {
		initLogger(c)
		initStorage(c)
		port := c.Int("port")
		return start(port)
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Logger.Fatal(err)
	}
}
