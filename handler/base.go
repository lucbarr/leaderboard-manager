package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

func ParseRequest(logger logrus.FieldLogger, w http.ResponseWriter, r *http.Request, req interface{}) {
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.WithError(err).Error("Failed reading request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(data, req)
	if err != nil {
		logger.WithError(err).Error("Failed parsing request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

type BaseResponse struct {
	Code StatusCode `json:"code"`
	Msg  string     `json:"msg"`
}

func WriteJSONResponse(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	w.Write(data)
	return nil
}
