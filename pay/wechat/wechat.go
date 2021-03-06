package wechat

import (
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"

	"github.com/douglarek/apikit"
	"github.com/fatih/structs"
)

const (
	orderURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	queryURL = "https://api.mch.weixin.qq.com/pay/orderquery"
)

// Wechat ...
type Wechat struct {
	client *apikit.Client
}

// New makes a wechat ...
func New(httpClient *http.Client) *Wechat {
	c := apikit.NewClient(httpClient)
	c.SetHeader(apikit.H{"Content-Type": apikit.MediaXML})
	return &Wechat{client: c}
}

// Req ...
type Req struct {
	XMLName  xml.Name `xml:"xml" json:"-"`
	AppID    string   `xml:"appid" structs:"appid" json:"appId"`
	MchID    string   `xml:"mch_id" structs:"mch_id" json:"partnerId"`
	NonceStr string   `xml:"nonce_str" structs:"nonce_str" json:"nonceStr"`
	Sign     string   `xml:"sign" structs:"sign" json:"sign"`
}

// OrderReq ...
type OrderReq struct {
	Req
	DeviceInfo     string `xml:"device_info" structs:"device_info"`
	Body           string `xml:"body" structs:"body"`
	Detail         string `xml:"detail" structs:"detail"`
	Attach         string `xml:"attach" structs:"attach"`
	OutTradeNo     string `xml:"out_trade_no" structs:"out_trade_no"`
	FeeType        string `xml:"fee_type" structs:"fee_type"`
	TotalFee       int    `xml:"total_fee" structs:"total_fee"`
	SpbillCreateIP string `xml:"spbill_create_ip" structs:"spbill_create_ip"`
	TimeStart      string `xml:"time_start" structs:"time_start"`
	TimeExpire     string `xml:"time_expire" structs:"time_expire"`
	GoodsTag       string `xml:"goods_tag" structs:"goods_tag"`
	NotifyURL      string `xml:"notify_url" structs:"notify_url"`
	TradeType      string `xml:"trade_type" structs:"trade_type"`
	ProductID      string `xml:"product_id" structs:"product_id"`
	LimitPay       string `xml:"limit_pay" structs:"limit_pay"`
	OpenID         string `xml:"openid" structs:"openid"`
}

// Resp ...
type Resp struct {
	ReturnCode string `structs:"return_code" xml:"return_code" json:"returnCode"`
	ReturnMsg  string `structs:"return_msg" xml:"return_msg" json:"returnMsg"`
	ResultCode string `structs:"result_code" xml:"result_code" json:"resultCode"`
	ErrCode    string `structs:"err_code" xml:"err_code" json:"errCode"`
	ErrCodeDes string `structs:"err_code_des" xml:"err_code_des" json:"errCodeDes"`
}

// OrderResp ...
type OrderResp struct {
	Resp
	Req
	DeviceInfo string `xml:"device_info" structs:"device_info"`
	TradeType  string `xml:"trade_type"`
	PrepayID   string `xml:"prepay_id"`
	CodeURL    string `xml:"code_url"`
}

// QueryReq ...
type QueryReq struct {
	Req
	TransactionID string `xml:"transaction_id" structs:"transaction_id"`
	OutTradeNo    string `xml:"out_trade_no" structs:"out_trade_no"`
}

// QueryResp ...
type QueryResp struct {
	Resp
	Req
	DeviceInfo     string `structs:"device_info" xml:"device_info" json:"deviceInfo"`
	OpenID         string `structs:"openid" xml:"openid" json:"openId"`
	IsSubscribe    string `structs:"is_subscribe" xml:"is_subscribe" json:"isSubscribe"`
	TradeType      string `structs:"trade_type" xml:"trade_type" json:"tradeType"`
	TradeState     string `structs:"trade_state" xml:"trade_state" json:"tradeState"`
	BankType       string `structs:"bank_type" xml:"bank_type" json:"bankType"`
	TotalFee       string `structs:"total_fee" xml:"total_fee" json:"totalFee"`
	FeeType        string `structs:"fee_type" xml:"fee_type" json:"feeType"`
	CashFee        int    `structs:"cash_fee" xml:"cash_fee" json:"cashFee"`
	CashFeeType    string `structs:"cash_fee_type" xml:"cash_fee_type" json:"cashFeeType"`
	CouponFee      int    `structs:"coupon_fee" xml:"coupon_fee" json:"couponFee"`
	CouponCount    int    `structs:"coupon_count" xml:"coupon_count" json:"couponCount"`
	TransactionID  string `structs:"transaction_id" xml:"transaction_id" json:"transactionId"`
	OutTradeNo     string `structs:"out_trade_no" xml:"out_trade_no" json:"outTradeNo"`
	Attach         string `structs:"attach" xml:"attach" json:"attach"`
	TimeEnd        string `structs:"time_end" xml:"time_end" json:"timeEnd"`
	TradeStateDesc string `structs:"trade_state_desc" xml:"trade_state_desc" json:"tradeStateDesc"`
}

// NotifyReq ...
type NotifyReq struct {
	QueryResp
}

// NotifyResp ...
type NotifyResp struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}

// AppReq packages needed params for client.
type AppReq struct {
	AppID     string `structs:"appid" json:"appid"`
	PartnerID string `structs:"partnerid" json:"partnerid"`
	PrepayID  string `structs:"prepayid" json:"prepayid"`
	Package   string `structs:"package" json:"package"`
	NonceStr  string `structs:"noncestr" json:"noncestr"`
	TimeStamp string `structs:"timestamp" json:"timestamp"`
	Sign      string `structs:"sign" json:"sign"`
}

// Sign ...
func (w *Wechat) Sign(s interface{}, secret string) string {
	m := apikit.Params(structs.Map(s))
	delete(m, "sign")
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for _, k := range keys {
		buf.WriteString(fmt.Sprintf("%s=%s", k, m[k]))
		buf.WriteString("&")
	}
	buf.WriteString("key=" + secret)
	return fmt.Sprintf("%X", md5.Sum(buf.Bytes()))
}

// Order ...
func (w *Wechat) Order(r *OrderReq) (*OrderResp, error) {
	req, err := w.client.NewRequest("POST", orderURL, r)
	if err != nil {
		return nil, err
	}
	res := new(OrderResp)
	if _, err := w.client.Do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Query ...
func (w *Wechat) Query(r *QueryReq) (*QueryResp, error) {
	req, err := w.client.NewRequest("POST", queryURL, r)
	if err != nil {
		return nil, err
	}
	res := new(QueryResp)
	if _, err := w.client.Do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}
