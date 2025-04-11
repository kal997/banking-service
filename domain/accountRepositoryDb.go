package domain

import (
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/kal997/banking/errs"
	"github.com/kal997/banking/logger"
)

type AccountRepositoryDb struct {
	client *sqlx.DB
}

func (d AccountRepositoryDb) Save(acc Account) (*Account, *errs.AppError) {

	saveAcc := "INSERT INTO accounts (customer_id, opening_date, account_type, amount, status) VALUES (?, ?, ?, ?, ?)"
	res, err := d.client.Exec(saveAcc, acc.CustomerId, acc.OpeningDate, acc.AccountType, acc.Amount, acc.Status)
	if err != nil {
		logger.Error("Error while creating a new account" + err.Error())
		return nil, errs.NewNotExpectedError("unexpected Database error")
	}

	id, err := res.LastInsertId()
	if err != nil {
		logger.Error("Error while getting last inserted id for new account " + err.Error())
		return nil, errs.NewNotExpectedError("unexpected Database error")

	}
	acc.AccountId = strconv.FormatInt(id, 10)
	return &acc, nil

}

func NewAccountRepositoryDb(dbClient *sqlx.DB) AccountRepositoryDb {
	return AccountRepositoryDb{dbClient}
}
