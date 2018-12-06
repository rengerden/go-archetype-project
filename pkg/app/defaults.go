package app

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	defaultPgDsn      = "postgres://docker:docker@localhost:5432/app"
	defaultRedisDsn   = "foobared@localhost:6379"
	defaultApiAddress = ":8888"

	QueueWorker1    = "worker#1"
	QueueWorker2    = "worker#2"
	QueueRcvTimeout = time.Second

	defaultBadWords = "fee,nee,cruul,leent"
)

type Config struct {
	PgDsn     string
	RedisPass string
	RedisAddr string
	ApiAddr   string
	BadWords  []string
}

func GetConfig() (Config, error) {
	dsnPg := flag.String("pg", defaultPgDsn, "e.g.: postgres://lgn:pwd@localhost:5432/app")
	redis := flag.String("redis", defaultRedisDsn, "e.g.: pwd@localhost:5432")
	apiAddr := flag.String("api", defaultApiAddress, "e.g.: :8888")
	words := flag.String("words", defaultBadWords, "bad words e.g.: "+defaultBadWords)
	flag.Parse()

	fmt.Println(*apiAddr)
	c := Config{
		PgDsn:   *dsnPg,
		ApiAddr: *apiAddr,
	}
	var redisDsn = regexp.MustCompile(`(.+)@(([\w|\.|\d]+):\d+)`)
	if !redisDsn.MatchString(*redis) {
		return c, errors.New("wrong Redis dsn")
	}
	res := redisDsn.FindAllSubmatch([]byte(*redis), -1)
	c.RedisPass = string(res[0][1])
	c.RedisAddr = string(res[0][2])

	var reWord = regexp.MustCompile(`\W+`)
	for _, word := range reWord.Split(*words, -1) {
		c.BadWords = append(c.BadWords, strings.ToLower(word))
	}
	return c, nil
}