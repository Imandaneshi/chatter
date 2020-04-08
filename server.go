package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"strings"
)

func setupLogging() {
	switch strings.ToLower(loggingConfig.level) {
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func Server() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Starts the chatter web-socket server",
		Before: func(c *cli.Context) error {
			setupLogging()

			redisError := setupRedis()
			if redisError != nil {
				return redisError
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "redis-addr",
				Value:       "redis:6379",
				Usage:       "Redis server address",
				EnvVars:     []string{"CHATTER_REDIS_ADDR"},
				Destination: &redisConfig.uri,
			},
			&cli.StringFlag{
				Name:        "logging-level",
				Value:       "info",
				Usage:       "Set logging level",
				EnvVars:     []string{"CHATTER_LOG_LEVEL"},
				Destination: &loggingConfig.level,
			},
		},
		Action: func(c *cli.Context) error {
			log.Info("go rules")
			return nil
		},
	}
}
