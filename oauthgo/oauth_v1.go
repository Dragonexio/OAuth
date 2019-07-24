package oauthgo

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"net/http"
	"strings"
)

type OAuthV1 struct {
	OAuth
}

func NewOAuthV1(appId, host, accessKey, secretKey, signKey string) (oAuth *OAuthV1) {
	return &OAuthV1{
		OAuth: NewDefaultOAuth(appId, host, accessKey, secretKey, signKey),
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
		PayToken               string `json:"pay_token"`
	}
}

func (d *OAuthV1) DoLogin(ctx context.Context, code, state, device string, scopes []int) (r *DoLoginResponse, hResp *http.Response, err error) {
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
	hResp, err = d.Do(ctx, req, r)
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
		PayToken               string `json:"pay_token"`
	}
}

func (d *OAuthV1) RefreshToken(ctx context.Context, accessToken, refreshToken string) (r *RefreshTokenResponse, hResp *http.Response, err error) {
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
	hResp, err = d.Do(ctx, req, r)
	return
}

type LogoutTokenResponse struct {
	BaseResponse
	Data struct{}
}

func (d *OAuthV1) LogoutToken(ctx context.Context, accessToken string) (r *LogoutTokenResponse, hResp *http.Response, err error) {
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
	hResp, err = d.Do(ctx, req, r)
	return
}

type QueryUserDetailResponse struct {
	BaseResponse
	Data struct{}
}

func (d *OAuthV1) QueryUserDetail(ctx context.Context, openId string) (r *QueryUserDetailResponse, hResp *http.Response, err error) {
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
	hResp, err = d.Do(ctx, req, r)
	return
}

type PreUser2AppResponse struct {
	BaseResponse
	Data struct {
		PayUrl string `json:"pay_url"`
	}
}

func (d *OAuthV1) PreUser2App(ctx context.Context, tradeNo, coinCode, volume, scene, desc, device, state, redirectUrl, domain string, specifyDragonExUid int) (r *PreUser2AppResponse, hResp *http.Response, err error) {
	var (
		path   = "/api/v1/pay/user2app/pre/"
		method = http.MethodPost
		values = map[string]interface{}{
			"trade_no":             tradeNo,
			"coin_code":            coinCode,
			"volume":               volume,
			"scene":                scene,
			"desc":                 desc,
			"device":               device,
			"state":                state,
			"redirect_url":         redirectUrl,
			"domain":               domain,
			"specify_dragonex_uid": specifyDragonExUid,
		}
		headers = http.Header{}
	)
	req, err := d.NewRequest(method, path, values, headers)
	if err != nil {
		return
	}

	r = new(PreUser2AppResponse)
	hResp, err = d.Do(ctx, req, r)
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
	OrderType   int             `json:"order_type"`
}

type DoApp2UserResponse struct {
	BaseResponse
	Data *OrderDetail
}

func (d *OAuthV1) doApp2UserByOpenIdOrDragonExUid(ctx context.Context, openId, tradeNo, coinCode, volume, scene, desc, device string, dragonExUid int64) (r *DoApp2UserResponse, hResp *http.Response, err error) {
	var (
		path   = "/api/v1/pay/app2user/do/"
		method = http.MethodPost
		values = map[string]interface{}{
			"uid":       dragonExUid,
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
	hResp, err = d.Do(ctx, req, r)
	return
}

func (d *OAuthV1) DoApp2UserByOpenId(ctx context.Context, openId, tradeNo, coinCode, volume, scene, desc, device string) (r *DoApp2UserResponse, hResp *http.Response, err error) {
	return d.doApp2UserByOpenIdOrDragonExUid(ctx, openId, tradeNo, coinCode, volume, scene, desc, device, 0)
}

func (d *OAuthV1) DoApp2UserByDragonExUid(ctx context.Context, dragonExUid int64, tradeNo, coinCode, volume, scene, desc, device string) (r *DoApp2UserResponse, hResp *http.Response, err error) {
	return d.doApp2UserByOpenIdOrDragonExUid(ctx, "", tradeNo, coinCode, volume, scene, desc, device, dragonExUid)
}

func (d *OAuthV1) openDoApp2UserByOpenIdOrDragonExUid(ctx context.Context, openId, tradeNo, coinCode, volume, scene, desc, device string, dragonExUid int64) (r *DoApp2UserResponse, hResp *http.Response, err error) {
	var (
		path   = "/api/v1/open/pay/app2user/do/"
		method = http.MethodPost
		values = map[string]interface{}{
			"uid":       dragonExUid,
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
	hResp, err = d.Do(ctx, req, r)
	return
}

func (d *OAuthV1) OpenDoApp2UserByOpenId(ctx context.Context, openId, tradeNo, coinCode, volume, scene, desc, device string) (r *DoApp2UserResponse, hResp *http.Response, err error) {
	return d.openDoApp2UserByOpenIdOrDragonExUid(ctx, openId, tradeNo, coinCode, volume, scene, desc, device, 0)
}

func (d *OAuthV1) OpenDoApp2UserByDragonExUid(ctx context.Context, dragonExUid int64, tradeNo, coinCode, volume, scene, desc, device string) (r *DoApp2UserResponse, hResp *http.Response, err error) {
	return d.openDoApp2UserByOpenIdOrDragonExUid(ctx, "", tradeNo, coinCode, volume, scene, desc, device, dragonExUid)
}

func (d *OAuthV1) openDoUser2AppByOpenIdOrDragonExUid(ctx context.Context, openId, tradeNo, coinCode, volume, scene, desc, device, payToken string, dragonExUid int64) (r *DoApp2UserResponse, hResp *http.Response, err error) {
	var (
		path   = "/api/v1/open/pay/user2app/do/"
		method = http.MethodPost
		values = map[string]interface{}{
			"uid":       dragonExUid,
			"open_id":   openId,
			"trade_no":  tradeNo,
			"coin_code": coinCode,
			"volume":    volume,
			"scene":     scene,
			"desc":      desc,
			"device":    device,
			"pay_token": payToken,
		}
		headers = http.Header{}
	)
	req, err := d.NewRequest(method, path, values, headers)
	if err != nil {
		return
	}

	r = new(DoApp2UserResponse)
	hResp, err = d.Do(ctx, req, r)
	return
}

func (d *OAuthV1) OpenDoUser2AppByOpenId(ctx context.Context, openId, tradeNo, coinCode, volume, scene, desc, device, payToken string) (r *DoApp2UserResponse, hResp *http.Response, err error) {
	return d.openDoUser2AppByOpenIdOrDragonExUid(ctx, openId, tradeNo, coinCode, volume, scene, desc, device, payToken, 0)
}

func (d *OAuthV1) OpenDoUser2AppByDragonExUid(ctx context.Context, dragonExUid int64, tradeNo, coinCode, volume, scene, desc, device, payToken string) (r *DoApp2UserResponse, hResp *http.Response, err error) {
	return d.openDoUser2AppByOpenIdOrDragonExUid(ctx, "", tradeNo, coinCode, volume, scene, desc, device, payToken, dragonExUid)
}

type QueryOrderDetailResponse struct {
	BaseResponse
	Data *OrderDetail
}

func (d *OAuthV1) QueryOrderDetail(ctx context.Context, tradeNo string) (r *QueryOrderDetailResponse, hResp *http.Response, err error) {
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
	hResp, err = d.Do(ctx, req, r)
	return
}

type ListOrdersResponse struct {
	BaseResponse
	Data struct {
		List  []*OrderDetail `json:"list"`
		Total int64          `json:"total"`
	}
}

func (d *OAuthV1) ListOrdersByDragonExUid(ctx context.Context, dragonExUid int64, coinCode string, direction, startTime, endTime, orderType, offset, limit int64) (r *ListOrdersResponse, hResp *http.Response, err error) {
	var (
		path   = "/api/v1/pay/order/history/"
		method = http.MethodPost
		values = map[string]interface{}{
			"uid":        dragonExUid,
			"coin_code":  coinCode,
			"direction":  direction,
			"start_time": startTime,
			"end_time":   endTime,
			"order_type": orderType,
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
	hResp, err = d.Do(ctx, req, r)
	return
}

type RedoPayCallbackResponse struct {
	BaseResponse
	Data struct{}
}

func (d *OAuthV1) RedoPayCallback(ctx context.Context, tradeNo string) (r *RedoPayCallbackResponse, hResp *http.Response, err error) {
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
	hResp, err = d.Do(ctx, req, r)
	return
}

type UserCoinDetail struct {
	Uid    int64           `json:"uid"`
	CoinId int64           `json:"coin_id"`
	Code   string          `json:"code"`
	Total  decimal.Decimal `json:"total"`
	Frozen decimal.Decimal `json:"frozen"`
}

type ListUserCoinsResponse struct {
	BaseResponse
	Data []*UserCoinDetail
}

func (d *OAuthV1) listUserCoins(ctx context.Context, openId string, dragonExUid int64) (r *ListUserCoinsResponse, hResp *http.Response, err error) {
	var (
		path   = "/api/v1/user/coin/list/"
		method = http.MethodPost
		values = map[string]interface{}{
			"open_id":      openId,
			"dragonex_uid": dragonExUid,
		}
		headers = http.Header{}
	)
	req, err := d.NewRequest(method, path, values, headers)
	if err != nil {
		return
	}

	r = new(ListUserCoinsResponse)
	hResp, err = d.Do(ctx, req, r)
	return
}

func (d *OAuthV1) ListUserCoinsByOpenId(ctx context.Context, openId string) (r *ListUserCoinsResponse, hResp *http.Response, err error) {
	return d.listUserCoins(ctx, openId, 0)
}

func (d *OAuthV1) ListUserCoinsByDragonExUid(ctx context.Context, dragonExUid int64) (r *ListUserCoinsResponse, hResp *http.Response, err error) {
	return d.listUserCoins(ctx, "", dragonExUid)
}

func (d *OAuthV1) listUserOpenCoins(ctx context.Context, openId string, dragonExUid int64) (r *ListUserCoinsResponse, hResp *http.Response, err error) {
	var (
		path   = "/api/v1/open/user/coin/list/"
		method = http.MethodPost
		values = map[string]interface{}{
			"open_id":      openId,
			"dragonex_uid": dragonExUid,
		}
		headers = http.Header{}
	)
	req, err := d.NewRequest(method, path, values, headers)
	if err != nil {
		return
	}

	r = new(ListUserCoinsResponse)
	hResp, err = d.Do(ctx, req, r)
	return
}

func (d *OAuthV1) ListUserOpenCoinsByOpenId(ctx context.Context, openId string) (r *ListUserCoinsResponse, hResp *http.Response, err error) {
	return d.listUserOpenCoins(ctx, openId, 0)
}

func (d *OAuthV1) ListUserOpenCoinsByDragonExUid(ctx context.Context, dragonExUid int64) (r *ListUserCoinsResponse, hResp *http.Response, err error) {
	return d.listUserOpenCoins(ctx, "", dragonExUid)
}

type QueryUserCoinResponse struct {
	BaseResponse
	Data *UserCoinDetail
}

func (d *OAuthV1) queryUserCoin(ctx context.Context, openId string, dragonExUid, CoinId int64, coinCode string) (r *QueryUserCoinResponse, hResp *http.Response, err error) {
	var (
		path   = "/api/v1/user/coin/detail/"
		method = http.MethodPost
		values = map[string]interface{}{
			"open_id":      openId,
			"dragonex_uid": dragonExUid,
			"coin_id":      CoinId,
			"coin_code":    coinCode,
		}
		headers = http.Header{}
	)
	req, err := d.NewRequest(method, path, values, headers)
	if err != nil {
		return
	}

	r = new(QueryUserCoinResponse)
	hResp, err = d.Do(ctx, req, r)
	return
}

func (d *OAuthV1) QueryUserCoinByOpenIdCoinId(ctx context.Context, openId string, CoinId int64) (r *QueryUserCoinResponse, hResp *http.Response, err error) {
	return d.queryUserCoin(ctx, openId, 0, CoinId, "")
}

func (d *OAuthV1) QueryUserCoinByOpenIdCoinCode(ctx context.Context, openId, CoinCode string) (r *QueryUserCoinResponse, hResp *http.Response, err error) {
	return d.queryUserCoin(ctx, openId, 0, 0, CoinCode)
}

func (d *OAuthV1) QueryUserCoinByDragonExUidCoinId(ctx context.Context, dragonExUid, CoinId int64) (r *QueryUserCoinResponse, hResp *http.Response, err error) {
	return d.queryUserCoin(ctx, "", dragonExUid, CoinId, "")
}

func (d *OAuthV1) QueryUserCoinByDragonExUidCoinCode(ctx context.Context, dragonExUid int64, CoinCode string) (r *QueryUserCoinResponse, hResp *http.Response, err error) {
	return d.queryUserCoin(ctx, "", dragonExUid, 0, CoinCode)
}

func (d *OAuthV1) queryUserOpenCoin(ctx context.Context, openId string, dragonExUid, CoinId int64, coinCode string) (r *QueryUserCoinResponse, hResp *http.Response, err error) {
	var (
		path   = "/api/v1/open/user/coin/detail/"
		method = http.MethodPost
		values = map[string]interface{}{
			"open_id":      openId,
			"dragonex_uid": dragonExUid,
			"coin_id":      CoinId,
			"coin_code":    coinCode,
		}
		headers = http.Header{}
	)
	req, err := d.NewRequest(method, path, values, headers)
	if err != nil {
		return
	}

	r = new(QueryUserCoinResponse)
	hResp, err = d.Do(ctx, req, r)
	return
}

func (d *OAuthV1) QueryUserOpenCoinByOpenIdCoinId(ctx context.Context, openId string, CoinId int64) (r *QueryUserCoinResponse, hResp *http.Response, err error) {
	return d.queryUserOpenCoin(ctx, openId, 0, CoinId, "")
}

func (d *OAuthV1) QueryUserOpenCoinByOpenIdCoinCode(ctx context.Context, openId, CoinCode string) (r *QueryUserCoinResponse, hResp *http.Response, err error) {
	return d.queryUserOpenCoin(ctx, openId, 0, 0, CoinCode)
}

func (d *OAuthV1) QueryUserOpenCoinByDragonExUidCoinId(ctx context.Context, dragonExUid, CoinId int64) (r *QueryUserCoinResponse, hResp *http.Response, err error) {
	return d.queryUserOpenCoin(ctx, "", dragonExUid, CoinId, "")
}

func (d *OAuthV1) QueryUserOpenCoinByDragonExUidCoinCode(ctx context.Context, dragonExUid int64, CoinCode string) (r *QueryUserCoinResponse, hResp *http.Response, err error) {
	return d.queryUserOpenCoin(ctx, "", dragonExUid, 0, CoinCode)
}
