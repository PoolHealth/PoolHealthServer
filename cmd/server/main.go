package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	gql "github.com/99designs/gqlgen/graphql"
	"github.com/kelseyhightower/envconfig"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"

	"github.com/PoolHealth/PoolHealthServer/internal/adapters/sheets"
	"github.com/PoolHealth/PoolHealthServer/internal/repo/influx"
	repoRedis "github.com/PoolHealth/PoolHealthServer/internal/repo/redis"
	"github.com/PoolHealth/PoolHealthServer/internal/services/actionsmanager"
	"github.com/PoolHealth/PoolHealthServer/internal/services/additiveshistory"
	auth2 "github.com/PoolHealth/PoolHealthServer/internal/services/auth"
	"github.com/PoolHealth/PoolHealthServer/internal/services/estimator"
	"github.com/PoolHealth/PoolHealthServer/internal/services/measurementhistory"
	"github.com/PoolHealth/PoolHealthServer/internal/services/poolmanager"
	"github.com/PoolHealth/PoolHealthServer/internal/services/poolsettingsmanager"
	"github.com/PoolHealth/PoolHealthServer/internal/services/sheetsmigrator"
	googleOAuth "github.com/PoolHealth/PoolHealthServer/internal/transport/google_oauth_handler"
	"github.com/PoolHealth/PoolHealthServer/internal/transport/server"
	"github.com/PoolHealth/PoolHealthServer/pkg/api/v1/graphql"
	"github.com/PoolHealth/PoolHealthServer/pkg/log"
)

var version = "dev"

const (
	pkgKey = "pkg"
)

type configuration struct {
	LoggerLevel  logrus.Level `envconfig:"LOG_LEVEL" default:"info"`
	LogToEcs     bool         `envconfig:"LOG_TO_ECS" default:"false"`
	RedisAddress string       `envconfig:"REDIS_ADDRESS" default:"localhost:6379"`

	InfluxDBAddress string `envconfig:"INFLUXDB_ADDRESS" default:"http://localhost:8086"`
	InfluxDBToken   string `envconfig:"INFLUXDB_TOKEN" default:""`
	InfluxDBOrg     string `envconfig:"INFLUXDB_ORG" default:"poolhealth"`
	InfluxDBBucket  string `envconfig:"INFLUXDB_BUCKET" default:"history"`

	SheetsCredentialsPath string `envconfig:"SHEETS_CREDENTIALS_PATH" default:".google/"`
}

func main() {
	printVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *printVersion {
		fmt.Println(version)
		return
	}

	cfg := new(configuration)
	if err := envconfig.Process("", cfg); err != nil {
		panic(err)
	}

	zeroLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	defer func() {
		cancel()
	}()

	if !cfg.LogToEcs {
		zeroLogger = zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()
	}

	logger := log.NewZerologLogger(&zeroLogger)

	db := redis.NewClient(&redis.Options{Addr: cfg.RedisAddress})
	repo := repoRedis.NewDB(db, logger.WithField(pkgKey, "repo"))

	if err := repo.Migrate(ctx); err != nil {
		panic(err)
	}

	poolManager := poolmanager.NewManager(repo, logger.WithField(pkgKey, "poolmanager"))

	authCfg := auth2.Config()
	authM := auth2.NewAuth(authCfg, repo, logger.WithField(pkgKey, "auth"))

	if err := authM.Start(); err != nil {
		panic(err)
	}

	idb, err := influx.New(cfg.InfluxDBAddress, cfg.InfluxDBToken, cfg.InfluxDBOrg, cfg.InfluxDBBucket, logger.WithField(pkgKey, "influx"))
	if err != nil {
		panic(err)
	}

	mhistory := measurementhistory.NewMeasurementHistory(idb, logger.WithField(pkgKey, "measurementhistory"))
	ahistory := additiveshistory.NewAdditivesHistory(idb, logger.WithField(pkgKey, "additiveshistory"))

	est := estimator.NewEstimator(repo, idb, logger.WithField(pkgKey, "estimator"))

	actionsManager := actionsmanager.NewManager(idb, logger.WithField(pkgKey, "actions"))

	poolSettingsManager := poolsettingsmanager.NewPoolSettingsManager(repo, logger.WithField(pkgKey, "poolsettingsmanager"))

	sheetAdapter := sheets.New(logger.WithField(pkgKey, "sheet"), cfg.SheetsCredentialsPath)
	if err = sheetAdapter.Start(context.Background()); err != nil {
		panic(err)
	}

	migrateMnager := sheetsmigrator.NewMigrator(
		sheetAdapter,
		poolManager,
		idb,
		logger.WithField(pkgKey, "sheetsmigrator"),
	)

	resolvers := graphql.NewResolver(
		logger.WithField(pkgKey, "graphql"),
		poolManager, mhistory, ahistory, est, actionsManager, poolSettingsManager,
		authM, migrateMnager,
	)

	goauthHandler := googleOAuth.NewHandler(sheetAdapter)

	s := server.NewServer(
		resolvers,
		[]gql.HandlerExtension{authM.Middleware()},
		authM.WsInitFunc,
		goauthHandler,
		logger.WithField(pkgKey, "server"),
	)

	s.InitV1Api()

	if err := s.Run(ctx); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("server closed")
		} else {
			panic(err)
		}
	}
}
