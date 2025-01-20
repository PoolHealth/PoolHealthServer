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
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/kelseyhightower/envconfig"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"

	"github.com/PoolHealth/PoolHealthServer/internal/repo/influx"
	repoRedis "github.com/PoolHealth/PoolHealthServer/internal/repo/redis"
	"github.com/PoolHealth/PoolHealthServer/internal/services/actionsmanager"
	"github.com/PoolHealth/PoolHealthServer/internal/services/additiveshistory"
	auth2 "github.com/PoolHealth/PoolHealthServer/internal/services/auth"
	"github.com/PoolHealth/PoolHealthServer/internal/services/estimator"
	"github.com/PoolHealth/PoolHealthServer/internal/services/measurementhistory"
	"github.com/PoolHealth/PoolHealthServer/internal/services/poolmanager"
	"github.com/PoolHealth/PoolHealthServer/internal/services/poolsettingsmanager"
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

	logrusLogger := logrus.New()
	logrusLogger.SetLevel(cfg.LoggerLevel)

	logrusLogger.SetFormatter(&nested.Formatter{
		FieldsOrder:     []string{pkgKey},
		TimestampFormat: "01-02|15:04:05",
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	defer func() {
		cancel()
	}()

	if cfg.LogToEcs {
		logrusLogger.SetFormatter(&ecslogrus.Formatter{})
	}

	logger := log.NewLogger(logrusLogger)

	db := redis.NewClient(&redis.Options{Addr: cfg.RedisAddress})
	repo := repoRedis.NewDB(db, logger.WithField(pkgKey, "repo"))

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

	resolvers := graphql.NewResolver(
		logger.WithField(pkgKey, "graphql"),
		poolManager, mhistory, ahistory, est, actionsManager, poolSettingsManager,
		authM,
	)

	s := server.NewServer(
		resolvers,
		[]gql.HandlerExtension{authM.Middleware()},
		authM.WsInitFunc,
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
