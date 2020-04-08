package main

import (
	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
)

var (
	r *redis.Client
)

func setupRedis() error {
	r = redis.NewClient(&redis.Options{
		Addr:     redisConfig.uri,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	log.Debug("pinging redis: %s", redisConfig.uri)
	_, err := r.Ping().Result()
	if err != nil {
		log.Infof("pinging redis failed: %s", redisConfig.uri)
		return err
	}
	log.Infof("successfully pinged redis: %s", redisConfig.uri)
	return nil
}
