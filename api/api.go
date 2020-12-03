package api

import (
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/lucbarr/leaderboard-manager/handler"
	"github.com/lucbarr/leaderboard-manager/handler/middleware"
	"github.com/spf13/viper"
)

type API struct {
	handlers *handler.All
	config   *viper.Viper
}

func NewAPI(config *viper.Viper, handlers *handler.All) *API {
	return &API{config: config, handlers: handlers}
}

func (a *API) Start() {
	r := mux.NewRouter()

	allHandlers := reflect.ValueOf(*a.handlers)

	for i := 0; i < allHandlers.NumField(); i++ {
		untypedHandler := allHandlers.Field(i).Interface()
		typedHandler, ok := untypedHandler.(handler.Handler)
		if !ok {
			continue
		}

		typedHandler.Configure(r)
	}

	http.Handle("/", r)
	middlewareHandler := middleware.NewBasicAuthHandler(a.config)
	r.Use(middlewareHandler.BasicAuth)

	hostURL := a.config.GetString("api.url")

	go log.Fatal(http.ListenAndServe(hostURL, nil))
}
