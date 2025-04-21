package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kal997/banking/dto"
	"github.com/kal997/banking/logger"
	"github.com/kal997/banking/service"
)

type AccountHandler struct {
	service service.AccountService
}

func (h *AccountHandler) MakeTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customer_id := vars["customer_id"]
	account_id := vars["account_id"]

	var userReq dto.TransactionRequest
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {

		// amount and transaction type are comming froom the user in the body
		// date will be inserted at domain level
		userReq.AccountId = account_id
		userReq.CustomerId = customer_id

		userResponse, err := h.service.MakeTransaction(userReq)
		if err != nil {
			logger.Error(err.Message)
			writeResponse(w, http.StatusInternalServerError, err.Message)
		} else {
			writeResponse(w, http.StatusCreated, userResponse)
		}

	}
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
