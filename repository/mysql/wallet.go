//Package mysql ...
package mysql

import (
	"github.com/amiraliio/avn-wallet/config"
	"github.com/amiraliio/avn-wallet/domain/model"
	"github.com/amiraliio/avn-wallet/domain/service"
)

type walletRepo struct {
	app *config.App
}

func NewMysqlWalletRepository(app *config.App) service.WalletRepository {
	return &walletRepo{
		app,
	}
}

func (w *walletRepo) Get(cellphone uint64) (*model.Wallet, error) {
	db := w.app.DB()
	defer db.Close()
	wallet := new(model.Wallet)
	if err := db.QueryRow("call getWallet(?)", cellphone).Scan(&wallet.Charge); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (w *walletRepo) Insert(wallet *model.Wallet) (*model.Wallet, error) {
	db := w.app.DB()
	defer db.Close()
	walletModel := new(model.Wallet)
	if err := db.QueryRow("call insertWallet(?,?)", wallet.Charge, wallet.User.Cellphone).Scan(&walletModel.Charge); err != nil {
		return nil, err
	}
	return walletModel, nil
}

func (w *walletRepo) Transactions(cellphone uint64) ([]*model.Transaction, error) {
	db := w.app.DB()
	defer db.Close()
	row, err := db.Query("call getTransactions(?)", cellphone)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var transactions []*model.Transaction
	for row.Next() {
		transaction := new(model.Transaction)
		if err := row.Scan(&transaction.ID, &transaction.Balance, &transaction.Type, &transaction.CreatedAt, &transaction.Cause, &transaction.Description); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
