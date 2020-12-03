package middleware

import (
	"net/http"

	"github.com/lucbarr/leaderboard-manager/handler"
	"github.com/spf13/viper"
)

type basicAuthHandler struct {
	config *viper.Viper
}

func NewBasicAuthHandler(config *viper.Viper) *basicAuthHandler {
	return &basicAuthHandler{
		config: config,
	}
}

func (b *basicAuthHandler) BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || !b.isUserValid(user, pass) {
			handler.WriteJSONResponse(w, &handler.BaseResponse{Code: "LB-001", Msg: "Unauthorized"})
			w.WriteHeader(http.StatusUnauthorized)
		}

		next.ServeHTTP(w, r)
	})
}

func (b *basicAuthHandler) isUserValid(user, pass string) bool {
	return user == b.config.GetString("api.auth.user") && pass == b.config.GetString("api.auth.pass")
}
