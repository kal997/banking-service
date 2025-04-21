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
		return nil, errs.NewUnexpectedError("unexpected Database error")
	}

	id, err := res.LastInsertId()
	if err != nil {
		logger.Error("Error while getting last inserted id for new account " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected Database error")

	}
	acc.AccountId = strconv.FormatInt(id, 10)
	return &acc, nil

}

func (d AccountRepositoryDb) FindBy(accountId string) (*Account, *errs.AppError) {
	sqlGet := "SELECT account_id, customer_id, opening_date, account_type, amount, status from accounts WHERE account_id = ?"
	var acc Account
	err := d.client.Get(&acc, sqlGet, accountId)
	if err != nil {
		logger.Error("Error whle fetching account information: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	return &acc, nil

}

func (d AccountRepositoryDb) SaveTransaction(transaction Transaction) (*Transaction, *errs.AppError) {

	// starting database transaction block
	tx, err := d.client.Begin()
	if err != nil {
		logger.Error("Error while starting a new transaction for bank account transaction: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	sqlNewTransaction := "INSERT INTO transactions (account_id, amount, transaction_type, transaction_date) VALUES (?, ?, ?, ?)"

	// Inserting the transaction into transactions table
	result, err := tx.Exec(sqlNewTransaction, transaction.AccountId, transaction.Amount, transaction.TransactionType, transaction.TransactionDate)

	if err != nil {
		logger.Error("Error while inserting a transaction " + err.Error())

		// abort the transaction since inserting transaction record failed
		tx.Rollback()

		logger.Error(err.Error())

		return nil, errs.NewUnexpectedError(err.Error())
	}

	// based on the transaction type, withdrawl or deposit, the query will change, so we nned to check
	if transaction.IsWithdrawal() {

		_, err = tx.Exec(`UPDATE accounts SET amount = amount - ? where account_id = ?`, transaction.Amount, transaction.AccountId)
	} else {
		_, err = tx.Exec(`UPDATE accounts SET amount = amount + ? where account_id = ?`, transaction.Amount, transaction.AccountId)
	}

	if err != nil {

		logger.Error("Error while inserting a transaction " + err.Error())

		tx.Rollback()
		// abort the transaction since updating account balance failed
		logger.Error("Error while saving transaction: " + err.Error())
		return nil, errs.NewUnexpectedError(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Error("Error while commiting transaction for bank account: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	// getting the latest transactionId
	transactionId, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error while getting transaction id " + err.Error())
		return nil, errs.NewUnexpectedError(err.Error())
	}

	// getting the updated account balance to insert it into the returned transaction
	updatedAcc, appErr := d.FindBy(transaction.AccountId)

	if err != nil {
		return nil, appErr
	}

	transaction.TransactionId = strconv.FormatInt(transactionId, 10)
	transaction.Amount = updatedAcc.Amount

	return &transaction, nil
}

func NewAccountRepositoryDb(dbClient *sqlx.DB) AccountRepositoryDb {
	return AccountRepositoryDb{dbClient}
}
