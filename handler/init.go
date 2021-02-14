//Package wallet ...
package handler

import (
	"net/http"

	"github.com/bitcodr/avn-wallet/config"
	"github.com/bitcodr/avn-wallet/domain/service"
	"github.com/bitcodr/avn-wallet/handler/rest"
	"github.com/bitcodr/avn-wallet/repository/mysql"
	"github.com/gorilla/mux"
)

const (
	REST_GET_WALLET          = "REST_GET_WALLET"
	REST_INSERT_WALLET       = "REST_INSERT_WALLET"
	REST_WALLET_TRANSACTIONS = "REST_WALLET_TRANSACTIONS"
)

func chooseWalletRepo(connection string, app *config.App) service.WalletRepository {
	switch connection {
	case "mysql":
		return mysql.NewMysqlWalletRepository(app)
	default:
		return nil
	}
}

func HTTP(app *config.App, router *mux.Router) {

	walletRepo := chooseWalletRepo("mysql", app)

	walletService := service.NewWalletService(walletRepo)

	walletRestHandler := rest.NewRestWalletHandler(walletService)

	router.HandleFunc("/wallet/{cellphone}", walletRestHandler.Get).Methods(http.MethodGet).Name(REST_GET_WALLET)
	router.HandleFunc("/wallet", walletRestHandler.Insert).Methods(http.MethodPost).Name(REST_INSERT_WALLET)
	router.HandleFunc("/wallet/{cellphone}/transactions", walletRestHandler.Transactions).Methods(http.MethodGet).Name(REST_WALLET_TRANSACTIONS)

}

func GRPC(app *config.App) {
	//implement grpc handler here
}
