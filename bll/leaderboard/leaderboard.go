package leaderboard

import (
	"context"

	"github.com/lucbarr/leaderboard-manager/dal"
	models "github.com/lucbarr/leaderboard-manager/models/leaderboard"
	"golang.org/x/xerrors"
)

type BLL interface {
	SetScore(ctx context.Context, playerID, leaderboardID string, score int) error
	GetHighestScores(ctx context.Context, leaderboardID string, limit int) ([]*models.Score, error)
	GetPlayerScoreAndRank(ctx context.Context, playerID, leaderboardID string) (int, int, error)
}

func NewBLL(dal *dal.All) *bll {
	return &bll{
		dal: dal,
	}
}

type bll struct {
	dal *dal.All
}

var _ BLL = (*bll)(nil)

func (bll *bll) SetScore(ctx context.Context, playerID, leaderboardID string, score int) error {
	return bll.dal.Leaderboard.SetScore(ctx, playerID, leaderboardID, score)
}

func (bll *bll) GetHighestScores(ctx context.Context, leaderboardID string, limit int) ([]*models.Score, error) {
	return bll.dal.Leaderboard.GetHighestScores(ctx, leaderboardID, limit)
}

func (bll *bll) GetPlayerScoreAndRank(ctx context.Context, playerID, leaderboardID string) (score int, rank int, err error) {
	playerScore, err := bll.dal.Leaderboard.GetPlayerScore(ctx, playerID, leaderboardID)
	if err != nil {
		return 0, 0, xerrors.Errorf("getting player score: %w", err)
	}
	return playerScore.Score, playerScore.Rank, nil
}
