package main

import (
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"net/http"
	"strings"
)

var (
	pubsub *redis.PubSub
)

func setupLogging() {
	switch strings.ToLower(loggingConfig.level) {
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// serve serves the web socket connection and creates the new user
func serve(room string, res http.ResponseWriter, req *http.Request) error {
	log.Debug("attempting to upgrade connection")

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(res, req, nil)

	if err != nil {
		log.Infof("failed upgrading connection: %s", err)
		return err
	}
	log.Info("successfully upgrade connection")


	_user := &user{room: room, connection: conn, messages: make(chan *message)}
	_user.assignId()

	username := req.URL.Query().Get("username")
	if username != "" {
		_user.setUsername(username)
	}
	_user.addToUsers()

	log.Infof("successfully added user to room: %s / %s - @%s", room, _user.ID, _user.Username)

	go _user.pull()
	go _user.push()

	_user.sendSelf()

	return nil
}

// redisPull receives new messages from channel
func redisPull() {
	pubsub = r.Subscribe(redisConfig.channel)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

	// Go channel which receives messages.
	ch := pubsub.Channel()

	for {
		select {
		case msg := <-ch:
			jsonError, message := messageFromJson([]byte(msg.Payload))
			if jsonError != nil {
				log.Infof("failed encoding message from redis: %s", msg.Payload)
			}
			for _, user := range users {
				if user.room == message.Room && message.UserId != user.ID {
					user.messages <- message
				}
			}
		}
	}
}

// publish publishes new messages to redis
func publish(msg *message) {
	jsonData, jsonError := msg.toJson(false)
	if jsonError != nil {
		log.Infof("failed encoding message when publishing on redis: %s", msg)
	}
	err := r.Publish(redisConfig.channel, jsonData).Err()
	if err != nil {
		log.Infof("Failed publishing new message")
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
				Value:       "127.0.0.1:6379",
				Usage:       "Redis server address",
				EnvVars:     []string{"CHATTER_REDIS_ADDR"},
				Destination: &redisConfig.uri,
			},
			&cli.StringFlag{
				Name:        "logging-level",
				Value:       "debug",
				Usage:       "Set logging level",
				EnvVars:     []string{"CHATTER_LOG_LEVEL"},
				Destination: &loggingConfig.level,
			},
			&cli.StringFlag{
				Name:        "ws-address",
				Value:       ":8009",
				Usage:       "Set web-socket server bind address",
				EnvVars:     []string{"CHATTER_WS_BIND_ADDRESS"},
				Destination: &websocketConfig.addr,
			},
		},
		Action: func(c *cli.Context) error {
			http.HandleFunc("/room/", func(res http.ResponseWriter, req *http.Request) {

				room := strings.TrimPrefix(req.URL.Path, "/room/")

				serveError := serve(room, res, req)
				if serveError != nil {
					log.Infof("Failed serving the connection")
				}
			})

			go redisPull()

			log.Infof("web-socket server is running at %s", websocketConfig.addr)
			err := http.ListenAndServe(websocketConfig.addr, nil)
			if err != nil {
				log.Infof("Failed running the web-socket server")
				return err
			}

			return nil
		},
	}
}
