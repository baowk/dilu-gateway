package handler

import (
	"net/http"
)

type ProxyHandler interface {
	GetName() string
	BeforeHander(w http.ResponseWriter, r *http.Request, args ...interface{}) (int, string)
	AfferHandler(w http.ResponseWriter, r *http.Request, args ...interface{}) (int, string)
}
