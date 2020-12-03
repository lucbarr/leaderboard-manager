package handler

import (
	"github.com/gorilla/mux"
)

type Handler interface {
	Configure(router *mux.Router)
}

type All struct {
	Debug       Handler
	Leaderboard Handler
}
