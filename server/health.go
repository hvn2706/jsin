package server

import (
	"jsin/logger"
	"net/http"
)

func ready(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("ok"))
	if err != nil {
		logger.Errorf("===== Write failed: %+v", err.Error())
		return
	}
	logger.Infof("===== Ready ok")
}

func liveness(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("live ok"))
	if err != nil {
		logger.Errorf("===== Write failed: %+v", err.Error())
		return
	}
	logger.Infof("===== Liveness ok")
}
