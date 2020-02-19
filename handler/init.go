//Package wallet ...
package handler

import (
	"net/http"

	"github.com/amiraliio/avn-wallet/config"
	"github.com/amiraliio/avn-wallet/domain/service"
	"github.com/amiraliio/avn-wallet/handler/rest"
	"github.com/amiraliio/avn-wallet/repository/mysql"
	"github.com/gorilla/mux"
)

const (
	REST_GET_WALLET    = "REST_GET_WALLET"
	REST_INSERT_WALLET = "REST_INSERT_WALLET"
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

	router.HandleFunc("wallet/{cellphone}", walletRestHandler.Get).Methods(http.MethodGet).Name(REST_GET_WALLET)
	router.HandleFunc("wallet/{cellphone}", walletRestHandler.Insert).Methods(http.MethodPost).Name(REST_INSERT_WALLET)
	router.HandleFunc("wallet/{cellphone}/transactions", walletRestHandler.Insert).Methods(http.MethodPost).Name(REST_INSERT_WALLET)

}

func GRPC(app *config.App) {
	//implement grpc handler route here
}