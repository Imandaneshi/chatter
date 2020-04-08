package main

// logging holds the information about logrus
type logging struct {
	Level string
}

// database holds the information about mongo
type redisConf struct {
	uri string
}

type loggingConf struct {
	level string
}

var (
	redisConfig   = &redisConf{}
	loggingConfig = &loggingConf{}
)
