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

func (m *walletRepo) Get(cellphone uint64) (*model.Wallet, error) {
	db := m.app.DB()
	defer db.Close()
	wallet := new(model.Wallet)
	if err := db.QueryRow("call getWallet(?)", cellphone).Scan(&wallet.Charge); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (w *walletRepo) Insert(wallet *model.Wallet) (*model.Wallet, error) {
	return nil, nil
}
