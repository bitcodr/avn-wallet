package helper

import (
	"encoding/json"
	"github.com/amiraliio/avn-wallet/config"
	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack/v4"
	"log"
	"net/http"
	"strings"
)

//the below response message model implemented according
//to the OASIS Standard incorporating Approved standard
//referenced by link http://docs.oasis-open.org/odata/odata-json-format/v4.0/errata02/os/odata-json-format-v4.0-errata02-os-complete.html#_Toc403940655

type ResponseModel struct {
	Success    bool             `json:"success" msgpack:"success"`
	Error      *ErrorModel      `json:"error" msgpack:"error"`
	Data       interface{}      `json:"data" msgpack:"data"`
	Pagination *PaginationModel `json:"pagination" msgpack:"pagination"`
}

type PaginationModel struct {
	Page  int `json:"page" msgpack:"page"`
	Limit int `json:"limit" msgpack:"limit"`
}

type ErrorModel struct {
	Code       int            `json:"code" msgpack:"code"`
	Message    string         `json:"message" msgpack:"message"`
	Details    []*ErrorDetail `json:"details" msgpack:"details"`
	ErrorTrace *ErrorTrace    `json:"errorTrace" msgpack:"errorTrace"`
}

type ErrorTrace struct {
	Trace   string `json:"trace" msgpack:"trace"`
	Context string `json:"context" msgpack:"context"`
}

type ErrorDetail struct {
	Code    string `json:"code" msgpack:"code"`
	Target  string `json:"target" msgpack:"target"`
	Message string `json:"message" msgpack:"message"`
}

func jsonSerializer(response *ResponseModel) ([]byte, error) {
	raw, err := json.Marshal(response)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Encode")
	}
	return raw, nil
}

func msgpackSerializer(response *ResponseModel) ([]byte, error) {
	raw, err := msgpack.Marshal(response)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Encode")
	}
	return raw, nil
}

func serialize(res http.ResponseWriter, contentType string, httpCode int, response *ResponseModel) {
	var err error
	var raw []byte
	switch contentType {
	case "application/json":
		raw, err = jsonSerializer(response)
	case "application/x-msgpack":
		raw, err = msgpackSerializer(response)
	default:
		contentType = "application/json"
		raw, err = jsonSerializer(response)
	}
	if err != nil {
		log.Fatal(err)
	}
	res.Header().Set("Content-Type", contentType)
	res.WriteHeader(httpCode)
	if _, err := res.Write(raw); err != nil {
		log.Println(err)
	}
}

func ResponseError(res http.ResponseWriter, err error, httpCode int, contentType, internalCode, detailMessage string) {
	errorMessage := new(ErrorModel)
	errorMessage.Code = httpCode
	errorMessage.Message = http.StatusText(httpCode)
	body := new(ErrorDetail)
	body.Code = internalCode
	if len(detailMessage) > 0 {
		target := strings.Fields(detailMessage)
		body.Target = strings.ToLower(target[0])
	}
	body.Message = detailMessage
	errorMessage.Details = append(errorMessage.Details, body)
	if config.AppConfig.GetBool("APP.DEBUG") && err != nil {
		innerError := new(ErrorTrace)
		innerError.Trace = err.Error()
		errorMessage.ErrorTrace = innerError
	}
	response := new(ResponseModel)
	response.Success = false
	response.Error = errorMessage
	response.Data = nil
	response.Pagination = nil
	serialize(res, contentType, httpCode, response)
}

func ResponseOk(res http.ResponseWriter, httpCode int, contentType string, data interface{}) {
	response := new(ResponseModel)
	response.Success = true
	response.Error = nil
	response.Data = data
	response.Pagination = nil
	serialize(res, contentType, httpCode, response)
}
