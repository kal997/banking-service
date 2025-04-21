package app

import (
	"encoding/json"
	"net/http"

	"github.com/kal997/banking/dto"
	"github.com/kal997/banking/service"
)

type AuthHandler struct {
	service service.AuthService
}

func (ah AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	token, appErr := ah.service.Login(loginRequest)
	if appErr != nil {
		writeResponse(w, appErr.Code, appErr.AsMessage())
		return
	}

	writeResponse(w, http.StatusOK, token)
}
