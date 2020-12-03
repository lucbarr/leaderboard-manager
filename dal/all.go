package dal

import (
	"github.com/go-redis/redis"
	"github.com/lucbarr/leaderboard-manager/dal/leaderboard"
	"github.com/spf13/viper"
)

type All struct {
	Leaderboard leaderboard.DAL
}

func NewAll(redisClient *redis.Client, config *viper.Viper) *All {
	return &All{
		Leaderboard: leaderboard.NewDAL(redisClient, config),
	}
}
