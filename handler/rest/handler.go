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
	"github.com/pkg/errors"
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
	params := mux.Vars(req)
	if params == nil {
		http.Error(res, config.LangConfig.GetString("MESSAGES.PARAM_EMPTY"), http.StatusUnprocessableEntity)
		return
	}
	cellphone, err := strconv.ParseUint(params["cellphone"], 10, 32)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	wallet, err := w.walletService.Get(uint(cellphone))
	if err != nil {
		if errors.Cause(err) == service.ErrWalletNotFound {
			http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	requestAccept := req.Header.Get("Accept")
	walletSerializer, err := w.serializer(requestAccept).Encode(wallet)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	helper.SetupResponse(res, requestAccept, walletSerializer, http.StatusOK)
}

func (w *walletHandler) Insert(res http.ResponseWriter, req *http.Request) {

}
