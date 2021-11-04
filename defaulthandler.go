package httprouter

import "net/http"

type DefaultHandler struct {
	hasHandler bool
	handler    http.Handler
}

func AlternativeHandler(h http.Handler) DefaultHandler {
	return DefaultHandler{true, h}
}

var _ http.Handler = DefaultHandler{}

func (d DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !d.hasHandler {
		w.WriteHeader(404)
		w.Write([]byte("Not found"))
	} else {
		d.handler.ServeHTTP(w, r)
	}
}
