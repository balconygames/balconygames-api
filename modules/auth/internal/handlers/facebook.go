package handlers

import (
	"net/http"

	httpreq "gitlab.com/balconygames/analytics/pkg/http"
)

// FacebookMiddleware should have callback implementation to support web sign in
func (h Handler) FacebookMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// todo: search for creds per app_id
		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FacebookLogin should have callback implementation to support web sign in
func (h Handler) FacebookLogin(w http.ResponseWriter, r *http.Request) {
	httpreq.NotImplemented(w)
}

// FacebookCallback should have callback implementation to support web sign in
func (h Handler) FacebookCallback(w http.ResponseWriter, r *http.Request) {
	httpreq.NotImplemented(w)
}

// GoogleMiddleware should have callback implementation to support web sign in
func (h Handler) GoogleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// todo: search for creds per app_id
		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GoogleLogin should have callback implementation to support web sign in
func (h Handler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	httpreq.NotImplemented(w)
}

// GoogleCallback should have callback implementation to support web sign in
func (h Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	httpreq.NotImplemented(w)
}

// TwitterMiddleware should have callback implementation to support web sign in
func (h Handler) TwitterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// todo: search for creds per app_id
		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// TwitterLogin should have callback implementation to support web sign in
func (h Handler) TwitterLogin(w http.ResponseWriter, r *http.Request) {
	httpreq.NotImplemented(w)
}

// TwitterCallback should have callback implementation to support web sign in
func (h Handler) TwitterCallback(w http.ResponseWriter, r *http.Request) {
	httpreq.NotImplemented(w)
}
