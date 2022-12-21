package mixpay

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// get by mixpay result traceId

type Data struct {
	Status           string `json:"status"`
	QuoteAmount      string `json:"quoteAmount"`
	QuoteSymbol      string `json:"quoteSymbol"`
	PaymentAmount    string `json:"paymentAmount"`
	PaymentSymbol    string `json:"paymentSymbol"`
	Payee            string `json:"payee"`
	PayeeMixinNumber string `json:"payeeMixinNumber"`
	PayeeAvatarURL   string `json:"payeeAvatarUrl"`
	Txid             string `json:"txid"`
	BlockExplorerURL string `json:"blockExplorerUrl"`
	Date             int64  `json:"date"`
	SurplusAmount    string `json:"surplusAmount"`
	SurplusStatus    string `json:"surplusStatus"`
	Confirmations    int64  `json:"confirmations"`
	PayableAmount    string `json:"payableAmount"`
	FailureCode      string `json:"failureCode"`
	FailureReason    string `json:"failureReason"`
	ReturnTo         string `json:"returnTo"`
}

type MixpayResult struct {
	Code        int64  `json:"code"`
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	Data        Data   `json:"data"`
	TimestampMS int64  `json:"timestampMs"`
}

// type MixpayRequest struct {
// 	TraceId string `json:"traceId"`
// }

func GetMixpayResult(orderId string, payeeId string) (MixpayResult, error) {
	mixpayAPIURL := fmt.Sprintf("https://api.mixpay.me/v1/payments_result?orderId=" + orderId + "&payeeId=" + payeeId)
	req, err := http.NewRequest(http.MethodGet, mixpayAPIURL, nil)
	if err != nil {
		log.Println(err)
		return MixpayResult{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return MixpayResult{}, err
	}
	defer res.Body.Close()

	mixpayResult := MixpayResult{}
	if err = json.NewDecoder(res.Body).Decode(&mixpayResult); err != nil {
		return mixpayResult, err
	}

	if mixpayResult.Data.Status != "success" {
		return mixpayResult, err
	}

	fmt.Println(mixpayResult)

	return mixpayResult, nil
}
