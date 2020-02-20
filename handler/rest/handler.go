//Package rest ...
package rest

import (
	"net/http"
	"strconv"

	"github.com/amiraliio/avn-wallet/config"
	"github.com/amiraliio/avn-wallet/domain/model"
	"github.com/amiraliio/avn-wallet/domain/service"
	"github.com/amiraliio/avn-wallet/helper"
	"github.com/amiraliio/avn-wallet/serializer/json"
	"github.com/amiraliio/avn-wallet/serializer/msgpack"
	"github.com/gorilla/mux"
)

type WalletHandler interface {
	Get(res http.ResponseWriter, req *http.Request)
	Insert(res http.ResponseWriter, req *http.Request)
	Transactions(res http.ResponseWriter, req *http.Request)
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
	acceptHeader := req.Header.Get("Accept")
	promotionCode := req.FormValue("promotionCode")
	if promotionCode == "" {
		helper.ResponseError(res, nil, http.StatusUnprocessableEntity, acceptHeader, "W-1003", config.LangConfig.GetString("MESSAGES.PROMOTION_CODE_IS_REQUIRED"))
		return
	}
	//TODO get promotion code is verified from promotion server with grpc
	walletModel := new(model.Wallet)
	wallet, err := w.walletService.Insert(walletModel)
	if err != nil {
		helper.ResponseError(res, err, http.StatusNotFound, acceptHeader, "W-1004", config.LangConfig.GetString("MESSAGES.DATA_NOT_FOUND"))
		return
	}
	//TODO send an event to promotion server who get the promotion
	//TODO waite for acknowledge from broker
	helper.ResponseOk(res, http.StatusOK, acceptHeader, wallet)
}

func (w *walletHandler) Transactions(res http.ResponseWriter, req *http.Request) {
	acceptHeader := req.Header.Get("Accept")
	params := mux.Vars(req)
	if params == nil {
		helper.ResponseError(res, nil, http.StatusUnprocessableEntity, acceptHeader, "W-1010", config.LangConfig.GetString("MESSAGES.PARAM_EMPTY"))
		return
	}
	cellphone, err := strconv.ParseUint(params["cellphone"], 10, 64)
	if err != nil {
		helper.ResponseError(res, err, http.StatusInternalServerError, acceptHeader, "W-1011", config.LangConfig.GetString("MESSAGES.PARSE_CELLPHONE"))
		return
	}
	transactions, err := w.walletService.Transactions(cellphone)
	if err != nil {
		helper.ResponseError(res, err, http.StatusNotFound, acceptHeader, "W-1012", config.LangConfig.GetString("MESSAGES.DATA_NOT_FOUND"))
		return
	}
	helper.ResponseOk(res, http.StatusOK, acceptHeader, transactions)
}
