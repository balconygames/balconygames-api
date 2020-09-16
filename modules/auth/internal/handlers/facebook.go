package handlers

import (
	"net/http"

	httpreq "gitlab.com/balconygames/analytics/pkg/http"
)

func (h Handler) FacebookMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// todo: search for creds per app_id
		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h Handler) FacebookLogin(w http.ResponseWriter, r *http.Request) {
	httpreq.NotImplemented(w)
}

func (h Handler) FacebookCallback(w http.ResponseWriter, r *http.Request) {
	httpreq.NotImplemented(w)
}

func (h Handler) GoogleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// todo: search for creds per app_id
		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h Handler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	httpreq.NotImplemented(w)
}

func (h Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	httpreq.NotImplemented(w)
}

func (h Handler) TwitterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// todo: search for creds per app_id
		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h Handler) TwitterLogin(w http.ResponseWriter, r *http.Request) {
	httpreq.NotImplemented(w)
}

func (h Handler) TwitterCallback(w http.ResponseWriter, r *http.Request) {
	httpreq.NotImplemented(w)
}
