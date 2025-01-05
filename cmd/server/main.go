package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	gql "github.com/99designs/gqlgen/graphql"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/kelseyhightower/envconfig"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"

	"github.com/PoolHealth/PoolHealthServer/internal/auth"
	"github.com/PoolHealth/PoolHealthServer/internal/poolmanager"
	repoRedis "github.com/PoolHealth/PoolHealthServer/internal/repo/redis"
	"github.com/PoolHealth/PoolHealthServer/internal/server"
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
	RedisAddress string       `default:"localhost:6379"`
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

	authCfg := auth.Config()
	authM := auth.NewAuth(authCfg, repo, logger.WithField(pkgKey, "auth"))

	if err := authM.Start(); err != nil {
		panic(err)
	}

	resolvers := graphql.NewResolver(
		logger.WithField(pkgKey, "graphql"), poolManager, authM,
	)

	s := server.NewServer(
		resolvers,
		[]gql.HandlerExtension{authM.Middleware()},
		authM.WsInitFunc,
		logger.WithField(pkgKey, "server"),
	)

	s.InitV1Api()

	if err := s.Run(ctx); err != nil {
		panic(err)
	}
}
