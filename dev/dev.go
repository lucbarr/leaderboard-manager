// +build test unit integration

package dev

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-redis/redis"
	"github.com/lucbarr/leaderboard-manager/dal"
	"github.com/spf13/viper"
)

var (
	testConfig *viper.Viper
	testRedis  *redis.Client

	cachedDAL *dal.All
)

func init() {
	var err error
	testConfig, err = readConfigs("../../cfg/default.yaml")
	if err != nil {
		panic(err)
	}

	testRedis, err = configureRedis(testConfig)
	if err != nil {
		panic(err)
	}
}

func GetDAL(t *testing.T) (*dal.All, func()) {
	if cachedDAL != nil {
		return cachedDAL, nil
	}

	dal := dal.NewAll(testRedis, testConfig)

	dropAllDBs := func() {
		flushRedis(testRedis)
	}

	return dal, dropAllDBs
}

func flushRedis(redis *redis.Client) {
	redis.FlushAll()
}

func readConfigs(configPath string) (*viper.Viper, error) {
	cfg := viper.New()

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	cfgData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	cfg.SetConfigType("yaml")
	cfg.ReadConfig(bytes.NewBuffer(cfgData))

	return cfg, nil
}

func configureRedis(cfg *viper.Viper) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.GetString("redis.host"),
		Password: cfg.GetString("redis.password"),
	})

	cmd := redisClient.Ping()
	if err := cmd.Err(); err != nil {
		return nil, err
	}

	return redisClient, nil
}
