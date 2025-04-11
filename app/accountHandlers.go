package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kal997/banking/dto"
	"github.com/kal997/banking/service"
)

type AccountHandler struct {
	service service.AccountService
}

func (h *AccountHandler) NewAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerId := vars["customer_id"]

	var request dto.NewAccountRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())

	} else {
		request.CustomerId = customerId
		response, appErr := h.service.NewAccount(request)
		if appErr != nil {
			writeResponse(w, appErr.Code, appErr.AsMessage())
		} else {
			writeResponse(w, http.StatusCreated, response)
		}

	}

}
