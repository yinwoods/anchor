package cmd

import (
	"math/rand"
	"time"

	docker "github.com/docker/docker/client"
)

// Client wraps docker client
var Client *docker.Client

func init() {
	var err error
	Client, err = docker.NewEnvClient()
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())
}

func randRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func randString(str ...string) string {
	return str[rand.Intn(len(str))]
}
