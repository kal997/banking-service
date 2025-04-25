package domain

import (
	"time"

	"github.com/kal997/banking-lib/errs"
	"github.com/kal997/banking/dto"
)

const dbTSLayout = "2006-01-02 15:04:05"

type Account struct {
	AccountId   string  `db:"account_id"`
	CustomerId  string  `db:"customer_id"`
	OpeningDate string  `db:"opening_date"`
	AccountType string  `db:"account_type"`
	Amount      float64 `db:"amount"`
	Status      string  `db:"status"`
}

//go:generate mockgen -destination=../mocks/domain/mockAccountRepository.go -package=domain github.com/kal997/banking/domain AccountRepository
type AccountRepository interface {
	Save(Account) (*Account, *errs.AppError)
	FindBy(string) (*Account, *errs.AppError)
	SaveTransaction(transaction Transaction) (*Transaction, *errs.AppError)
}

func (a Account) CanWithdraw(amount float64) bool {
	if a.Amount < amount {
		return false
	}
	return true
}

func (a Account) ToNewAccountResponseDto() *dto.NewAccountResponse {
	return &dto.NewAccountResponse{AccountId: a.AccountId}
}

func NewAccount(customerId, accountType string, amount float64) Account {
	return Account{
		CustomerId:  customerId,
		Amount:      amount,
		OpeningDate: time.Now().Format(dbTSLayout),
		AccountType: accountType,
		Status:      "1",
	}
}
