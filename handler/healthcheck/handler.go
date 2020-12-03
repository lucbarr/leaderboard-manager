package healthcheck

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Handler is an HTTP handler made for debugging only
type Handler struct {
	logger logrus.FieldLogger
}

// NewHandler returns debug HTTP handler
func NewHandler(logger logrus.FieldLogger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) Configure(router *mux.Router) {
	router.HandleFunc("/healthcheck", h.healthcheck).Methods("GET")
}

func (d *Handler) healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("working"))
}
