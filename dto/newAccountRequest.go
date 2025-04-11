package dto

import (
	"strings"

	"github.com/kal997/banking/errs"
)

type NewAccountRequest struct {
	CustomerId  string  `json:"customer_id"`
	AccountType string  `json:"account_type"`
	Amount      float64 `json:"amount"`
}

// thhe request nows how to validate itself, thats why we put validate here
func (r NewAccountRequest) Validate() *errs.AppError {
	if r.Amount < 5000 {
		return errs.NewValidationError("Amount should be greater than or equal to 5000")
	}

	if strings.ToLower(r.AccountType) != "saving" && strings.ToLower(r.AccountType) != "checking" {
		return errs.NewValidationError("Acount type should be saving or checking")
	}

	return nil
}
