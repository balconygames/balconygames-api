package handlers

import (
	"net/http"

	httpreq "gitlab.com/balconygames/analytics/pkg/http"
)

type serverSigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type serverSigninResponse struct {
	ID string `json:"id"`
}

func (h *Handler) ServerSignin(w http.ResponseWriter, r *http.Request) {
	httpreq.NotImplemented(w)
}

func (h *Handler) ServerSignout(w http.ResponseWriter, r *http.Request) {
	httpreq.NotImplemented(w)
}

type serverSignupRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`

	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type serverSignupResponse struct {
	ID string
}

func (h *Handler) ServerSignup(w http.ResponseWriter, r *http.Request) {
	httpreq.NotImplemented(w)
}
