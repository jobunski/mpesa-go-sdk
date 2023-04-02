package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

/**
URL
*/

const (
	STKPUSH        = "mpesa/stkpush/v1/processrequest"
	STKQUERY       = "mpesa/stkpushquery/v1/query"
	C2BURL         = "mpesa/c2b/v1/simulate"
	C2BREGISTERURL = "mpesa/c2b/v1/registerurl"

	url      = "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
	password = "MTc0Mzc5YmZiMjc5ZjlhYTliZGJjZjE1OGU5N2RkNzFhNDY3Y2QyZTBjODkzMDU5YjEwZjc4ZTZiNzJhZGExZWQyYzkxOTIwMTkxMjEzMTA1NzEz"
)

/**
COMMANDIDS
*/

const (
	PAYBILL = "CustomerPayBillOnline"
)

/*
*
ENV VARIABLES
*/
var CALLBACKURL, APPKEY, APPSECRET, BASEURL string

const shortCode = "174379"
const C2BTransactionType = "CustomerPayBillOnline"

type mpesaAction string

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

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func AssignConfigsToVariables(appKey, appSecret, baseUrl, callbackUrl string) {

	CALLBACKURL = callbackUrl
	APPKEY = appKey
	APPSECRET = appSecret
	BASEURL = baseUrl
}

func (t *Token) GetAccessToken2(ctx context.Context) (string, error) {
	client := http.Client{}
	var appCredentials = &AppCredentials{APPKEY, APPSECRET}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return "", fmt.Errorf("error creating GenerateToken Request: %v", err)
	}

	req.SetBasicAuth(appCredentials.AppKey, appCredentials.AppSecret)
	req.Header.Add("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error while sending token Generation request: %v", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("StatusCode: %d. StatusMessage: %s\n", resp.StatusCode, resp.Status)
	}
	respErr := json.NewDecoder(resp.Body).Decode(&t)

	if respErr != nil {
		return "", fmt.Errorf("error while unmarshalling token Generation request: %v", err)
	}
	return t.AccessToken, nil
}

/**
All the Http Request are made from this point,
Parameters required are the url,action and body
Figured out all the request are made same way hence create one function
that does exactly that

Returns a Response Body taken by another object for processing
*/

func makeHttpRequest(ctx context.Context, method, url string, action mpesaAction, body interface{}) (*http.Response, error) {
	client := http.Client{}
	var tokenBody Token
	reqBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("error while marshalling body for Send Request:")
		log.Fatal(err)
		return nil, err
	}

	token, err := tokenBody.GetAccessToken2(ctx)
	if err != nil {
		fmt.Printf("Could not Generate Token due to specified error: %s\n", err.Error())
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Could not Create Request with Provided body")
		return nil, fmt.Errorf("mpesa: error creating %v request. Error is: - %v", action, err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error from the Http Request Send")
		return nil, fmt.Errorf("mpesa: error after  %v request. Error is: - %v", action, err)
	}

	return resp, nil

}

func CustomResponse(resp *http.Response, err error) (response *GenericResponse) {
	var respErr error

	if err != nil {
		fmt.Print("Error from Response received")
		return &GenericResponse{err.Error(), 500}
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	customResponse := &GenericResponse{"Success", 200}
	if resp.StatusCode == 200 {
		var successResponse Response
		respErr = json.NewDecoder(resp.Body).Decode(&successResponse)
		fmt.Printf("SuccResponse: %v\n", successResponse)

	} else {
		var errorResponse BadRequestResponse
		respErr = json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Printf("ErrResponse: %v\n", errorResponse)
		customResponse.Status = "Error occurred"
		customResponse.StatusCode = 415
	}

	if respErr != nil {
		fmt.Printf("mpesa: error unmarshalling stk push response Request. Response Status: %d . Error: %v\n",
			resp.StatusCode, respErr)
	}

	return customResponse

}

func StkPushRequest(ctx context.Context, businessShortCode, timestamp, partyA, partyB, phoneNumber, accountReference, transDescription string,
	amount float64) (response *GenericResponse) {
	stkPushBody := &StkRequest{
		businessShortCode,
		password,
		timestamp,
		C2BTransactionType,
		amount,
		partyA,
		partyB,
		phoneNumber,
		CALLBACKURL,
		accountReference,
		transDescription,
	}

	resp, err := makeHttpRequest(ctx, http.MethodPost, BASEURL+STKPUSH, "STK PUSH", stkPushBody)
	return CustomResponse(resp, err)

}

func StkPushQueryRequest(ctx context.Context, shortCode, timestamp, requestId string) (response *GenericResponse) {
	stkPushQueryBody := &StkPushQuery{
		shortCode,
		password,
		timestamp,
		requestId,
	}

	resp, err := makeHttpRequest(ctx, http.MethodPost, BASEURL+STKQUERY, "STK PUSH QUERY", stkPushQueryBody)
	return CustomResponse(resp, err)
}

func CustomerToBusiness(ctx context.Context, shortcode, msisdn, billReferenceNumber string, amount float64) (response *GenericResponse) {
	c2BBody := &C2B{shortcode,
		PAYBILL,
		amount,
		msisdn,
		billReferenceNumber}

	resp, err := makeHttpRequest(ctx, http.MethodPost, BASEURL+C2BURL, "C2B Simulate Transactions", c2BBody)
	return CustomResponse(resp, err)

}

func C2BRegisterUrl(ctx context.Context, shortcode, reasonType, confirmationUrl, validationUrl string) (response *GenericResponse) {
	c2bRegisterUrlBody := &C2BRegisterUrlBody{
		shortcode,
		reasonType,
		confirmationUrl,
		validationUrl}

	resp, err := makeHttpRequest(ctx, http.MethodPost, BASEURL+C2BREGISTERURL, "C2B REGISTER URL", c2bRegisterUrlBody)
	return CustomResponse(resp, err)

}

func TransactionStatus(ctx context.Context, transactionId, originalConversationId, creditParty, remarks, occasion string,
	identifierType int) {

}
