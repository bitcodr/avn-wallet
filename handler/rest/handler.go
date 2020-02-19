//Package rest ...
package rest

import (
	"net/http"
	"strconv"

	"github.com/amiraliio/avn-wallet/helper"

	"github.com/amiraliio/avn-wallet/config"
	"github.com/amiraliio/avn-wallet/domain/service"
	"github.com/amiraliio/avn-wallet/serializer/json"
	"github.com/amiraliio/avn-wallet/serializer/msgpack"
	"github.com/gorilla/mux"
)

type WalletHandler interface {
	Get(res http.ResponseWriter, req *http.Request)
	Insert(res http.ResponseWriter, req *http.Request)
}

type walletHandler struct {
	walletService service.WalletService
}

func NewRestWalletHandler(walletService service.WalletService) WalletHandler {
	return &walletHandler{
		walletService,
	}
}

func (h *walletHandler) serializer(contentType string) service.WalletSerializer {
	switch contentType {
	case "application/json":
		return &json.Wallet{}
	case "application/x-msgpack":
		return &msgpack.Wallet{}
	default:
		return &json.Wallet{}
	}
}

func (w *walletHandler) Get(res http.ResponseWriter, req *http.Request) {
	acceptHeader := req.Header.Get("Accept")
	params := mux.Vars(req)
	if params == nil {
		helper.ResponseError(res, nil, http.StatusUnprocessableEntity, acceptHeader, "W-1000", config.LangConfig.GetString("MESSAGES.PARAM_EMPTY"))
		return
	}
	cellphone, err := strconv.ParseUint(params["cellphone"], 10, 64)
	if err != nil {
		helper.ResponseError(res, err, http.StatusInternalServerError, acceptHeader, "W-1001", config.LangConfig.GetString("MESSAGES.PARSE_CELLPHONE"))
		return
	}
	wallet, err := w.walletService.Get(cellphone)
	if err != nil {
		helper.ResponseError(res, err, http.StatusNotFound, acceptHeader, "W-1002", config.LangConfig.GetString("MESSAGES.DATA_NOT_FOUND"))
		return
	}
	helper.ResponseOk(res, http.StatusOK, acceptHeader, wallet)
}

func (w *walletHandler) Insert(res http.ResponseWriter, req *http.Request) {

}
