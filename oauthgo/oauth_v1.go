package oauthgo

import (
	"fmt"
	"github.com/shopspring/decimal"
	"net/http"
	"strings"
)

type OAuthV1 struct {
	OAuth
}

func NewOAuthV1(appId, host, accessKey, secretKey string) (oAuth *OAuthV1) {
	return &OAuthV1{
		OAuth: NewDefaultOAuth(appId, host, accessKey, secretKey),
	}
}

type DoLoginResponse struct {
	BaseResponse
	Data struct {
		AccessToken            string `json:"access_token"`
		AccessTokenExpireTime  int64  `json:"access_token_et"`
		RefreshToken           string `json:"refresh_token"`
		RefreshTokenExpireTime int64  `json:"refresh_token_et"`
		CompanyId              string `json:"company_id"`
		AppId                  string `json:"app_id"`
		OpenId                 string `json:"open_id"`
		UnionId                string `json:"union_id"`
		Scopes                 []int  `json:"scopes"`
	}
}

func (d *OAuthV1) DoLogin(code, state, device string, scopes []int) (r *DoLoginResponse, hResp *http.Response, err error) {
	scopesSlice := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		scopesSlice = append(scopesSlice, fmt.Sprint(scope))
	}

	var (
		path   = "/api/v1/login/do/"
		method = http.MethodPost
		values = map[string]interface{}{
			"app_id": d.GetAppId(),
			"code":   code,
			"state":  state,
			"device": device,
			"scopes": strings.Join(scopesSlice, ","),
		}
		headers = http.Header{}
	)
	req, err := d.NewRequest(method, path, values, headers)
	if err != nil {
		return
	}

	r = new(DoLoginResponse)
	hResp, err = d.Do(req, r)
	return
}

type RefreshTokenResponse struct {
	BaseResponse
	Data struct {
		AccessToken            string `json:"access_token"`
		AccessTokenExpireTime  int64  `json:"access_token_et"`
		RefreshToken           string `json:"refresh_token"`
		RefreshTokenExpireTime int64  `json:"refresh_token_et"`
		Scopes                 []int  `json:"scopes"`
	}
}

func (d *OAuthV1) RefreshToken(accessToken, refreshToken string) (r *RefreshTokenResponse, hResp *http.Response, err error) {
	var (
		path   = "/api/v1/login/refresh/"
		method = http.MethodPost
		values = map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		}
		headers = http.Header{}
	)
	req, err := d.NewRequest(method, path, values, headers)
	if err != nil {
		return
	}

	r = new(RefreshTokenResponse)
	hResp, err = d.Do(req, r)
	return
}

type LogoutTokenResponse struct {
	BaseResponse
	Data struct{}
}

func (d *OAuthV1) LogoutToken(accessToken string) (r *LogoutTokenResponse, hResp *http.Response, err error) {
	var (
		path   = "/api/v1/login/logout/"
		method = http.MethodPost
		values = map[string]interface{}{
			"access_token": accessToken,
		}
		headers = http.Header{}
	)
	req, err := d.NewRequest(method, path, values, headers)
	if err != nil {
		return
	}

	r = new(LogoutTokenResponse)
	hResp, err = d.Do(req, r)
	return
}

type QueryUserDetailResponse struct {
	BaseResponse
	Data struct{}
}

func (d *OAuthV1) QueryUserDetail(openId string) (r *QueryUserDetailResponse, hResp *http.Response, err error) {
	var (
		path   = "/api/v1/user/detail/"
		method = http.MethodPost
		values = map[string]interface{}{
			"open_id": openId,
		}
		headers = http.Header{}
	)
	req, err := d.NewRequest(method, path, values, headers)
	if err != nil {
		return
	}

	r = new(QueryUserDetailResponse)
	hResp, err = d.Do(req, r)
	return
}

type PreUser2AppResponse struct {
	BaseResponse
	Data struct {
		PayUrl string `json:"pay_url"`
	}
}

func (d *OAuthV1) PreUser2App(tradeNo, coinCode, volume, scene, desc, device, state, redirectUrl string) (r *PreUser2AppResponse, hResp *http.Response, err error) {
	var (
		path   = "/api/v1/pay/user2app/pre/"
		method = http.MethodPost
		values = map[string]interface{}{
			"trade_no":     tradeNo,
			"coin_code":    coinCode,
			"volume":       volume,
			"scene":        scene,
			"desc":         desc,
			"device":       device,
			"state":        state,
			"redirect_url": redirectUrl,
		}
		headers = http.Header{}
	)
	req, err := d.NewRequest(method, path, values, headers)
	if err != nil {
		return
	}

	r = new(PreUser2AppResponse)
	hResp, err = d.Do(req, r)
	return
}

type OrderDetail struct {
	Id          int64           `json:"id"`
	TradeNo     string          `json:"trade_no"`
	OpenId      string          `json:"open_id"`
	DragonExUid int64           `json:"uid"`
	CoinCode    string          `json:"coin_code"`
	Volume      decimal.Decimal `json:"volume"`
	Direction   int             `json:"direction"`
	Status      int             `json:"status"`
	CutVolume   decimal.Decimal `json:"cut_volume"`
	ArriveTime  int64           `json:"arrive_time"`
	CreateTime  int64           `json:"create_time"`
}

type DoApp2UserResponse struct {
	BaseResponse
	Data *OrderDetail
}

func (d *OAuthV1) DoApp2UserByOpenId(openId, tradeNo, coinCode, volume, scene, desc, device string) (r *DoApp2UserResponse, hResp *http.Response, err error) {
	var (
		path   = "/api/v1/pay/app2user/do/"
		method = http.MethodPost
		values = map[string]interface{}{
			"open_id":   openId,
			"trade_no":  tradeNo,
			"coin_code": coinCode,
			"volume":    volume,
			"scene":     scene,
			"desc":      desc,
			"device":    device,
		}
		headers = http.Header{}
	)
	req, err := d.NewRequest(method, path, values, headers)
	if err != nil {
		return
	}

	r = new(DoApp2UserResponse)
	hResp, err = d.Do(req, r)
	return
}

func (d *OAuthV1) DoApp2UserByDragonExUid(dragonExUid int64, tradeNo, coinCode, volume, scene, desc, device string) (r *DoApp2UserResponse, hResp *http.Response, err error) {
	var (
		path   = "/api/v1/pay/app2user/do/"
		method = http.MethodPost
		values = map[string]interface{}{
			"uid":       dragonExUid,
			"trade_no":  tradeNo,
			"coin_code": coinCode,
			"volume":    volume,
			"scene":     scene,
			"desc":      desc,
			"device":    device,
		}
		headers = http.Header{}
	)
	req, err := d.NewRequest(method, path, values, headers)
	if err != nil {
		return
	}

	r = new(DoApp2UserResponse)
	hResp, err = d.Do(req, r)
	return
}

type QueryOrderDetailResponse struct {
	BaseResponse
	Data *OrderDetail
}

func (d *OAuthV1) QueryOrderDetail(tradeNo string) (r *QueryOrderDetailResponse, hResp *http.Response, err error) {
	var (
		path   = "/api/v1/pay/order/detail/"
		method = http.MethodPost
		values = map[string]interface{}{
			"trade_no": tradeNo,
		}
		headers = http.Header{}
	)
	req, err := d.NewRequest(method, path, values, headers)
	if err != nil {
		return
	}

	r = new(QueryOrderDetailResponse)
	hResp, err = d.Do(req, r)
	return
}

type ListOrdersResponse struct {
	BaseResponse
	Data struct {
		List  []*OrderDetail `json:"list"`
		Total int64          `json:"total"`
	}
}

func (d *OAuthV1) ListOrdersByDragonExUid(dragonExUid int64, coinCode string, direction, startTime, endTime, offset, limit int64) (r *ListOrdersResponse, hResp *http.Response, err error) {
	var (
		path   = "/api/v1/pay/order/history/"
		method = http.MethodPost
		values = map[string]interface{}{
			"uid":        dragonExUid,
			"coin_code":  coinCode,
			"direction":  direction,
			"start_time": startTime,
			"end_time":   endTime,
			"offset":     offset,
			"limit":      limit,
		}
		headers = http.Header{}
	)
	req, err := d.NewRequest(method, path, values, headers)
	if err != nil {
		return
	}

	r = new(ListOrdersResponse)
	hResp, err = d.Do(req, r)
	return
}

type RedoPayCallbackResponse struct {
	BaseResponse
	Data struct{}
}

func (d *OAuthV1) RedoPayCallback(tradeNo string) (r *RedoPayCallbackResponse, hResp *http.Response, err error) {
	var (
		path   = "/api/v1/pay/callback/redo/"
		method = http.MethodPost
		values = map[string]interface{}{
			"trade_no": tradeNo,
		}
		headers = http.Header{}
	)
	req, err := d.NewRequest(method, path, values, headers)
	if err != nil {
		return
	}

	r = new(RedoPayCallbackResponse)
	hResp, err = d.Do(req, r)
	return
}
