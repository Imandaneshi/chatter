package main

import "time"

type websocketConf struct {
	addr string
}

// database holds the information about mongo
type redisConf struct {
	uri string
	channel string
}

// logging holds the information about logrus
type loggingConf struct {
	level string
}

var (
	redisConfig   = &redisConf{channel:"chatter"}
	loggingConfig = &loggingConf{}
	websocketConfig = &websocketConf{}
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 1024
)