package leaderboard

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
	models "github.com/lucbarr/leaderboard-manager/models/leaderboard"
	"github.com/spf13/viper"
)

type DAL interface {
	SetScore(ctx context.Context, playerID, leaderboardID string, score int) error
	GetHighestScores(ctx context.Context, leaderboardID string, limit int) ([]*models.Score, error)
	GetPlayerScore(ctx context.Context, playerID, leaderboardID string) (*models.Score, error)
}

type leaderboard struct {
	redisClient *redis.Client
	config      *viper.Viper
}

func NewDAL(redisClient *redis.Client, config *viper.Viper) *leaderboard {
	return &leaderboard{
		redisClient: redisClient,
		config:      config,
	}
}

var _ DAL = (*leaderboard)(nil)

func (l *leaderboard) buildLeaderboardKey(leaderboardID string) string {
	return fmt.Sprintf("%s::%s", l.config.GetString("redis.prefix"), leaderboardID)
}

const setScoreScript string = `
local s = redis.call('ZSCORE', KEYS[1], ARGV[2])
if not s or tonumber(s) < tonumber(ARGV[1]) then
  redis.call('ZADD', KEYS[1], ARGV[1], ARGV[2])
end
`

func (l *leaderboard) SetScore(ctx context.Context, playerID, leaderboardID string, score int) error {
	redisClient := l.redisClient.WithContext(ctx)

	leaderboardKey := l.buildLeaderboardKey(leaderboardID)

	cmd := redisClient.Eval(setScoreScript, []string{leaderboardKey}, score, playerID)

	if cmd.Err() == redis.Nil {
		return nil
	}

	return cmd.Err()
}
func (l *leaderboard) GetHighestScores(ctx context.Context, leaderboardID string, limit int) ([]*models.Score, error) {
	redisClient := l.redisClient.WithContext(ctx)

	leaderboardKey := l.buildLeaderboardKey(leaderboardID)

	cmd := redisClient.ZRevRangeByScoreWithScores(leaderboardKey, redis.ZRangeBy{
		Max:   "+inf",
		Min:   "-inf",
		Count: int64(limit),
	})

	results, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	scores := make([]*models.Score, len(results))
	for i, result := range results {
		playerID := result.Member.(string)
		scores[i] = &models.Score{
			LeaderboardID: leaderboardID,
			PlayerID:      playerID,
			Score:         int(result.Score),
		}
	}

	return scores, nil
}
func (l *leaderboard) GetPlayerScore(ctx context.Context, playerID, leaderboardID string) (*models.Score, error) {
	redisClient := l.redisClient.WithContext(ctx)

	leaderboardKey := l.buildLeaderboardKey(leaderboardID)

	txPipeline := redisClient.TxPipeline()

	scoreCmd := txPipeline.ZScore(leaderboardKey, playerID)
	rankCmd := txPipeline.ZRevRank(leaderboardKey, playerID)

	_, err := txPipeline.Exec()
	if err != nil {
		return nil, err
	}

	score, err := scoreCmd.Result()
	if err != nil {
		return nil, err
	}

	rank, err := rankCmd.Result()
	if err != nil {
		return nil, err
	}

	return &models.Score{
		LeaderboardID: leaderboardID,
		PlayerID:      playerID,
		Score:         int(score),
		Rank:          int(rank),
	}, nil
}
