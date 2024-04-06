package handler

import (
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

func WriteErrResponseAndLog(rw http.ResponseWriter, logger *logrus.Logger, statusCode int, logMsg string, respMsg string) {
	if logMsg != "" {
		logger.Errorf(logMsg)
	}

	rw.WriteHeader(statusCode)

	if respMsg != "" {
		_, err := rw.Write([]byte(respMsg))
		if err != nil {
			logger.Errorf("error occurred writing response: %s", err)
		}
	}
}

func GetIntParamFromQuery(req *http.Request, key string) (int, error) {
	return strconv.Atoi(req.URL.Query().Get(key))
}

func GetStringParamFromQuery(req *http.Request, key string) (string, error) {
	str := req.URL.Query().Get(key)

	if str == "" {
		return str, ErrNoQueryParamProvided
	}

	return str, nil
}

func GetIntHeaderByKey(req *http.Request, key string) (int, error) {
	str := req.Header.Get(key)
	if str == "" {
		return -1, ErrNoHeaderProvided
	}

	val, err := strconv.Atoi(str)
	if err != nil {
		return -1, ErrInvalidHeaderProvided
	}

	return val, nil
}

func GetStringHeaderByKey(req *http.Request, key string) (string, error) {
	str := req.Header.Get(key)
	if str == "" {
		return str, ErrNoHeaderProvided
	}

	return str, nil
}
