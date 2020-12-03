package leaderboard

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lucbarr/leaderboard-manager/bll/leaderboard"
	"github.com/lucbarr/leaderboard-manager/handler"
	models "github.com/lucbarr/leaderboard-manager/models/leaderboard"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Handler struct {
	logger logrus.FieldLogger
	bll    leaderboard.BLL
	config *viper.Viper
}

var _ handler.Handler = (*Handler)(nil)

func (h *Handler) Configure(router *mux.Router) {
	router.HandleFunc("/leaderboard/{leaderboardID}/{playerID}", h.setScore).Methods("POST")
	router.HandleFunc("/leaderboard/{leaderboardID}", h.getTopScores).Methods("GET")
	router.HandleFunc("/leaderboard/{leaderboardID}/{playerID}", h.getPlayerScore).Methods("GET")
}

func NewHandler(logger logrus.FieldLogger, leaderboardBLL leaderboard.BLL, config *viper.Viper) *Handler {
	return &Handler{
		logger: logger,
		bll:    leaderboardBLL,
		config: config,
	}
}

type SetScoreRequest struct {
	Score int `json:"score"`
}

func (h *Handler) setScore(w http.ResponseWriter, r *http.Request) {
	req := &SetScoreRequest{}
	vars := mux.Vars(r)
	handler.ParseRequest(h.logger, w, r, req)

	score := req.Score
	leaderboardID := vars["leaderboardID"]
	playerID := vars["playerID"]

	ctx := r.Context()
	err := h.bll.SetScore(ctx, playerID, leaderboardID, score)
	if err != nil {
		handler.WriteJSONResponse(w, &handler.BaseResponse{Code: handler.CodeInternalError, Msg: "error"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	handler.WriteJSONResponse(w, &handler.BaseResponse{Msg: "success"})
}

type GetTopScores struct {
	handler.BaseResponse
	Scores []*Score `json:"scores"`
}

type Score struct {
	PlayerID string `json:"playerID"`
	Score    int    `json:"score"`
}

func (h *Handler) getTopScores(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	limitString := queryValues.Get("limit")

	var limit int
	if limitString == "" {
		limit = h.config.GetInt("leaderboard.topPlayers.limit")
	} else {
		var err error
		limit, err = strconv.Atoi(limitString)
		if err != nil {
			handler.WriteJSONResponse(w, &handler.BaseResponse{Code: handler.CodeBadRequest, Msg: "success"})
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	vars := mux.Vars(r)
	leaderboardID := vars["leaderboardID"]

	ctx := r.Context()
	scores, err := h.bll.GetHighestScores(ctx, leaderboardID, limit)
	if err != nil {
		handler.WriteJSONResponse(w, &handler.BaseResponse{Code: handler.CodeInternalError, Msg: "error"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	retScores := buildScores(scores)

	resp := &GetTopScores{
		BaseResponse: handler.BaseResponse{
			Msg: "success",
		},
		Scores: retScores,
	}

	handler.WriteJSONResponse(w, resp)
}

type GetPlayerScoreResponse struct {
	handler.BaseResponse
	Score int `json:"score"`
	Rank  int `json:"rank"`
}

func (h *Handler) getPlayerScore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	leaderboardID := vars["leaderboardID"]
	playerID := vars["playerID"]

	ctx := r.Context()
	score, rank, err := h.bll.GetPlayerScoreAndRank(ctx, playerID, leaderboardID)
	if err != nil {
		handler.WriteJSONResponse(w, &handler.BaseResponse{Code: handler.CodeInternalError, Msg: "error"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &GetPlayerScoreResponse{
		BaseResponse: handler.BaseResponse{
			Msg: "success",
		},
		Score: score,
		Rank:  rank,
	}

	handler.WriteJSONResponse(w, resp)
}

func buildScores(modelsScores []*models.Score) []*Score {
	scores := make([]*Score, len(modelsScores))
	for i, modelsScore := range modelsScores {
		scores[i] = &Score{
			PlayerID: modelsScore.PlayerID,
			Score:    modelsScore.Score,
		}
	}

	return scores
}
