package app

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kal997/banking-lib/errs"
	"github.com/kal997/banking/dto"
	"github.com/kal997/banking/service"
)

type AuthMiddleware struct {
	service service.AuthService
}

func (ah AuthMiddleware) Verify(w http.ResponseWriter, r *http.Request) {

	urlParams := make(map[string]string)

	// converting from query to map type
	for k, _ := range r.URL.Query() {
		urlParams[k] = r.URL.Query().Get(k)
	}

	if urlParams["token"] != "" {
		appErr := ah.service.Verify(urlParams)

		if appErr != nil {
			writeResponse(w, appErr.Code, appErr.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, authorizedResponse())
		}
	} else {

		writeResponse(w, http.StatusForbidden, errs.NewAuthorizationError("Missing token"))

	}
}
func (ah AuthMiddleware) Login(w http.ResponseWriter, r *http.Request) {
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

func (ah AuthMiddleware) authorizationHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			currentRoute := mux.CurrentRoute(r)
			currentRouteVars := mux.Vars(r)
			authHeader := r.Header.Get("Authorization")

			if authHeader != "" {
				token := getTokenFromHeader(authHeader)

				// send git method to the verify endpoint
				isAuthorized := ah.service.IsAuthorized(token, currentRoute.GetName(), currentRouteVars)
				if isAuthorized {
					next.ServeHTTP(w, r)
				} else {
					writeResponse(w, http.StatusForbidden, "Unauthorized")
				}

			} else {
				// check if the request is login, skip verification
				if currentRoute.GetName() == "Login" || currentRoute.GetName() == "Verify" {
					next.ServeHTTP(w, r)
				} else {

					writeResponse(w, http.StatusUnauthorized, "Missing Token")
				}
			}

		})
	}
}

func getTokenFromHeader(header string) string {
	/*
	   token is coming in the format as below
	   "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50cyI6W.yI5NTQ3MCIsIjk1NDcyIiw"
	*/
	splitToken := strings.Split(header, "Bearer")
	if len(splitToken) == 2 {
		return strings.TrimSpace(splitToken[1])
	}
	return ""
}

func authorizedResponse() map[string]bool {
	return map[string]bool{"isAuthorized": true}
}
