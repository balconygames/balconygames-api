package handlers

import (
	"net/http"

	httpreq "gitlab.com/balconygames/analytics/pkg/http"
)

// ServerSignin used to use to make the dashboard application
func (h *Handler) ServerSignin(w http.ResponseWriter, r *http.Request) {
	httpreq.NotImplemented(w)
}

// ServerSignout log out the user
func (h *Handler) ServerSignout(w http.ResponseWriter, r *http.Request) {
	httpreq.NotImplemented(w)
}

// ServerSignup creates the new user
func (h *Handler) ServerSignup(w http.ResponseWriter, r *http.Request) {
	httpreq.NotImplemented(w)
}
