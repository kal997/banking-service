package dto

import (
	"net/http"
	"testing"
)

func Test_should_return_error_when_transaction_type_is_not_deposit_or_withdrawal(t *testing.T) {
	//AAA

	// Arrange
	transaction := TransactionRequest{Amount: 50, TransactionType: "invalid"}

	// Act
	appErr := transaction.Validate()

	// Assert
	if appErr.Message != "Transaction type can only be deposit or withdrawal" {
		t.Error("Invalid message while testing transaction type")
	}

	if appErr.Code != http.StatusUnprocessableEntity {
		t.Error("Invalid code while testing transaction type")
	}
}

func Test_should_return_error_when_amount_is_less_than_zero(t *testing.T) {
	//AAA

	// Arrange
	transaction := TransactionRequest{Amount: -5, TransactionType: "deposit"}

	// Act
	appErr := transaction.Validate()

	// Assert
	if appErr.Message != "Amount cannot be less than zero" {
		t.Error("Invalid message while validating amount")
	}

	if appErr.Code != http.StatusUnprocessableEntity {
		t.Error("Invalid code while validating amount")
	}
}


func Test_should_return_nil_when_transaction_is_valid(t *testing.T) {
	//AAA

	// Arrange
	transaction := TransactionRequest{Amount: 500, TransactionType: "deposit"}

	// Act
	appErr := transaction.Validate()

	// Assert
	if appErr != nil {
		t.Error("Invalid AppErr while validating amount")
	}

	
}


