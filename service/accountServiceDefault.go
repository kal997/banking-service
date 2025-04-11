package service

import (
	"time"

	"github.com/kal997/banking/domain"
	"github.com/kal997/banking/dto"
	"github.com/kal997/banking/errs"
)

type DefaultAccountService struct {
	repo domain.AccountRepository
}

func (s DefaultAccountService) NewAccount(req dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError) {
	
	// validate the request
	appErr := req.Validate()
	if appErr != nil{
		return nil, appErr
	}
	
	// in request: transform from dto to domain
	a := domain.Account{
		AccountId:   "",
		CustomerId:  req.CustomerId,
		OpeningDate: time.Now().Format("2006-01-02 15:04:05"),
		AccountType: req.AccountType,
		Amount:      req.Amount,
		Status:      "1",
	}
	newAccount, err := s.repo.Save(a)
	if err != nil {
		return nil, err
	}

	// in response; transform from domain to dto
	responseDto := newAccount.ToNewAccountResponseDto()
	return &responseDto, nil
}


func NewAccountService(repo domain.AccountRepository) DefaultAccountService{
	return DefaultAccountService{repo}
}