package _struct

import "net/http"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AppCredentials struct {
	AppKey    string
	AppSecret string
}

type Token struct {
	AccessToken string `json:"access_token"`
	Expiry      string `json:"expires_in"`
}

type GenericResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"code"`
}

type NetworkRequest struct {
	AccessToken string `json:"access_token"`
	Expiry      string `json:"expiry_in"`
	AppKey      string `json:"app_key"`
	AppSecret   string `json:"app_secret"`
	client      HttpClient
}

type StkRequest struct {
	ShortCode        string  `json:"BusinessShortCode"`
	Password         string  `json:"Password"`
	Timestamp        string  `json:"Timestamp"`
	TransactionType  string  `json:"TransactionType"`
	Amount           float64 `json:"Amount"`
	Sender           string  `json:"PartyA"`
	Receiver         string  `json:"PartyB"`
	PhoneNumber      string  `json:"PhoneNumber"`
	CallBackURL      string  `json:"CallBackURL"`
	AccountReference string  `json:"AccountReference"`
	TransactionDesc  string  `json:"TransactionDesc"`
}
type Response struct {
	MerchantRequestId       string `json:"MerchantRequestId"`
	CheckoutRequestId       string `json:"CheckoutRequestId"`
	OriginatorCoversationID string `json:"OriginatorCoversationID"`
	ResponseCode            string `json:"ResponseCode"`
	ResponseDescription     string `json:"ResponseDescription"`
	CustomerMessage         string `json:"CustomerMessage"`
	ResultCode              string `json:"ResultCode"`
	ResultDesc              string `json:"ResultDesc"`
}

type C2B struct {
	ShortCode     string  `json:"ShortCode"`
	CommandID     string  `json:"CommandID"`
	Amount        float64 `json:"Amount"`
	Msisdn        string  `json:"Msisdn"`
	BillRefNumber string  `json:"BillRefNumber"`
}

type StkPushQuery struct {
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	CheckoutRequestID string `json:"CheckoutRequestID"`
}

type BadRequestResponse struct {
	ErrorId      string `json:"requestId"`
	ResponseCode string `json:"responseCode"`
	ResponseDesc string `json:"responseDesc"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

type C2BRegisterUrlBody struct {
	ShortCode       string `json:"shortCode"`
	ResponseType    string `json:"ResponseType"`
	ConfirmationUrl string `json:"ConfirmationUrl"`
	ValidationUrl   string `json:"ValidationUrl"`
}
