package service

import (
	"github.com/kal997/banking/dto"
	"github.com/kal997/banking/errs"
)

type AccountService interface {
	NewAccount(dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError)
	MakeTransaction(request dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppError)
}



