//Package rest ...
package rest

import (
	"context"
	jsonEncoder "encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/amiraliio/avn-grpc-promotion-proto/proto"
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
	contentTypeHeader := req.Header.Get("Content-Type")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		helper.ResponseError(res, err, http.StatusNotFound, contentTypeHeader, "W-1004", config.LangConfig.GetString("MESSAGES.DATA_NOT_FOUND"))
		return
	}
	defer req.Body.Close()
	chargeRequest := new(model.ChargeRequest)
	if err := jsonEncoder.Unmarshal(body, chargeRequest); err != nil {
		helper.ResponseError(res, nil, http.StatusUnprocessableEntity, contentTypeHeader, "W-1003", config.LangConfig.GetString("MESSAGES.PROMOTION_CODE_IS_REQUIRED"))
		return
	}
	charge, err := w.getVerifyClient(chargeRequest.PromotionCode)
	if err != nil {
		helper.ResponseError(res, err, http.StatusNotFound, contentTypeHeader, "W-1004", config.LangConfig.GetString("MESSAGES.DATA_NOT_FOUND"))
		return
	}
	walletModel := new(model.Wallet)
	walletModel.Charge = charge
	userModel := new(model.User)
	userModel.Cellphone = chargeRequest.Cellphone
	walletModel.User = userModel
	wallet, err := w.walletService.Insert(walletModel)
	if err != nil {
		helper.ResponseError(res, err, http.StatusNotFound, contentTypeHeader, "W-1005", config.LangConfig.GetString("MESSAGES.DATA_NOT_FOUND"))
		return
	}
	nats, err := config.NATSClient()
	if err != nil {
		helper.ResponseError(res, err, http.StatusNotFound, contentTypeHeader, "W-1006", config.LangConfig.GetString("MESSAGES.DATA_NOT_FOUND"))
		return
	}
	defer nats.Close()
	chargeRequest.Fullname = wallet.User.FirstName + " " + wallet.User.LastName
	if err := nats.Publish("promotion."+strconv.FormatUint(chargeRequest.Cellphone, 10), chargeRequest); err != nil {
		helper.ResponseError(res, err, http.StatusNotFound, contentTypeHeader, "W-1007", config.LangConfig.GetString("MESSAGES.DATA_NOT_FOUND"))
		return
	}
	helper.ResponseOk(res, http.StatusOK, contentTypeHeader, wallet)
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

func (w *walletHandler) getVerifyClient(promotionCode string) (float64, error) {
	conn, err := config.GRPCConnection(config.AppConfig.GetString("APP.PROMOTION_GRPC_SERVER"))
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	client := proto.NewPromotionClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	requestModel := new(proto.Request)
	requestModel.PromotionCode = promotionCode
	response, err := client.Verify(ctx, requestModel)
	if err != nil {
		return 0, err
	}
	return response.GetCharge(), nil
}
