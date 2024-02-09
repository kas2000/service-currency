package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	command "github.com/kas2000/commandlib"
	"github.com/kas2000/logger"
	"github.com/kas2000/http"
	"github.com/kas2000/service-currency/currency"
	"github.com/urfave/cli/v2"
	_ "github.com/jackc/pgx/v4/stdlib"
	"os"
	"time"
)

var (
	port       = ""
	systemName = "CURRENCY_SERVICE"
	dbUri      = ""
	env        = ""

	flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Usage:       "Load configuration from `FILE`",
			Required:    true,
			Destination: &env,
		},
	}
)

func parseEnv() error {
	err := godotenv.Overload(env)
	if err != nil {
		return err
	}
	port = os.Getenv("PORT")
	if port == "" {
		return errors.New("invalid port.")
	}
	dbUri = os.Getenv("DB_URI")
	if dbUri == "" {
		return errors.New("invalid db uri.")
	}


	return nil
}

func main() {
	app := &cli.App{
		Name:      "Currency Service",
		Usage:     "currency-service",
		UsageText: "go run main.go/currency-service --config FILE",
		Flags:     flags,
		Action:    run,
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func run(*cli.Context) error {
	log, _ := logger.New("debug")

	if err := parseEnv(); err != nil {
		log.Fatal("Error parsing .env file: " + err.Error())
	}

	db, err := sql.Open("pgx", dbUri)
	if err != nil {
		log.Fatal("Unable to connect to database: %v\n"+ err.Error())
	}
	defer db.Close()

	var greeting string
	err = db.QueryRow("select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		log.Fatal("QueryRow failed: %v\n"+ err.Error())
	}

	fmt.Println(greeting)

	serverConfig := http.Config{
		Port:            port,
		ShutdownTimeout: time.Second * 20,
		GracefulTimeout: time.Second * 21,
		ApiVersion:      "v1",
		Timeout:         time.Second * 20,
		Logger:          log,
	}
	server := http.NewServer(serverConfig)

	currencyRepo, err := currency.NewPostgresCurrencyRepo(db)
	if err != nil {
		log.Fatal("couldn't initialize roles repository: " + err.Error())
	}

	service := currency.NewService(log, currencyRepo)

	currencyCh := command.NewCommandHandler(service)
	currencyHttpFactory := currency.NewCurrencyHttpHandler(log, currencyCh, systemName)
	currencyController := currency.NewCurrencyController(&server, currencyHttpFactory, "")
	currencyController.Bind()

	server.ListenAndServe()

	return nil
}
