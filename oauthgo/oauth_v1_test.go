package oauthgo

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	apiV1 *OAuthV1
)

func TestMain(m *testing.M) {
	const (
		appId     = "appidfortest"
		host      = "https://oauth.dragonex.io"
		accessKey = "6b7ef684a6645325ae1505d5636cfb12"
		secretKey = "2721a56546de5baf9cfc4b5ceb3cc59f"
		signKey   = "testKey"
	)

	apiV1 = NewOAuthV1(appId, host, accessKey, secretKey, signKey)

	apiV1.After(displayRequestAndResponseMiddleware)
	apiV1.After(CheckResponseMiddleware)

	m.Run()
}

const (
	scene  = "scenefortest"
	desc   = "descrortest"
	device = "devicefortest"
	state  = "statefortest"
)

func TestOAuthV1_DoLogin(t *testing.T) {
	var (
		code   = "59f47ee75a"
		scopes = []int{ScopeLogin}
	)
	resp, _, err := apiV1.DoLogin(context.Background(), code, state, device, scopes)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_RefreshToken(t *testing.T) {
	var (
		accessToken  = "3d6524isnEAoGBx5p3Z0C/3IBXesM7biMwDZKkvhj6lsXgWAb6VSlr/433k6q0hX8IG1WXFvyiZP9TBdMNm7Zwq2WdGdZ3ntdVCqBw4bZI5PMsRsN/L+aR7813rnfCKJA6GQrhlLsZV5l6hRbM/4G0Mmz0m8l+Bs9aAEtXvZ8viyDsUOGO2t5x7k4m7G61Db1CTaCsB+ig0MADb1ssiPY+F4mrTm0p+LR81l8b0aLSA="
		refreshToken = "wRNoiNTk5OTFjMGUzNTg1NWE0OGJiMTE1NzI5NWRkNjNjZTQiLCJhIjoiYXBwaWRmb3J0ZXN0IiwiZXhwIjoxNTY0NTY3NTAyLCJuYmYiOjE1NjM5NjI3MDJ9tbREEzHLsbfMx5GEpMdAyoan4-47_FgIaeHkegSJ_pQM"
	)
	resp, _, err := apiV1.RefreshToken(context.Background(), accessToken, refreshToken)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_LogoutToken(t *testing.T) {
	var (
		accessToken = "3d6524isnEAoGBx5p3Z0C/3IBXesM7biMwDZKkvhj6lsXgWAb6VSlr/433k6q0hX8IG1WXFvyiZP9TBdMNm7Zwq2WdGdZ3ntdVCqBw4bZI5PMsRsN/L+aR7813rnfCKJA6GQrhlLsZV5l6hRbM/4G/4h1cI/04V9BYmd/9fKMijz0y1v2OCcx4Ge7xe8vwxMj+4JpbWIedbWPWFoZV10oLeHaQm+K51rMJ8Onxnid60="
	)
	resp, _, err := apiV1.LogoutToken(context.Background(), accessToken)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_QueryUserDetail(t *testing.T) {
	var (
		openId = "59991c0e35855a48bb1157295dd63ce4"
	)
	resp, _, err := apiV1.QueryUserDetail(context.Background(), openId)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_PreUser2App(t *testing.T) {
	var (
		tradeNo                  = fmt.Sprint(time.Now().Unix())
		coinCode                 = "usdt"
		volume                   = "0.1"
		redirectUrl              = "https://www.google.com"
		domain                   = "dragonex.io"
		specifyDragonExUid int64 = 0
	)
	resp, _, err := apiV1.PreUser2App(context.Background(), tradeNo, coinCode, volume, scene, desc, device, state, redirectUrl, domain, specifyDragonExUid)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_DoApp2UserByOpenId(t *testing.T) {
	var (
		openId   = "59991c0e35855a48bb1157295dd63ce4"
		tradeNo  = fmt.Sprint(time.Now().Unix())
		coinCode = "usdt"
		volume   = "0.1"
	)
	resp, _, err := apiV1.DoApp2UserByOpenId(context.Background(), openId, tradeNo, coinCode, volume, scene, desc, device)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_DoApp2UserByDragonExUid(t *testing.T) {
	var (
		dragonExUid int64 = 1000000
		tradeNo           = fmt.Sprint(time.Now().Unix())
		coinCode          = "usdt"
		volume            = "0.1"
	)
	resp, _, err := apiV1.DoApp2UserByDragonExUid(context.Background(), dragonExUid, tradeNo, coinCode, volume, scene, desc, device)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_Refund(t *testing.T) {
	var (
		oriTradeNo    = ""
		refundTradeNo = ""
		refundRate    = "1"
	)
	resp, _, err := apiV1.rReturn(context.Background(), oriTradeNo, refundTradeNo, refundRate, scene, desc, device)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_OpenDoUser2AppByOpenId(t *testing.T) {
	var (
		openId   = "59991c0e35855a48bb1157295dd63ce4"
		tradeNo  = fmt.Sprint(time.Now().Unix())
		coinCode = "usdt"
		volume   = "0.1"
		payToken = "AimJY="
	)
	resp, _, err := apiV1.OpenDoUser2AppByOpenId(context.Background(), openId, tradeNo, coinCode, volume, scene, desc, device, payToken)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_OpenDoUser2AppByDragonExUid(t *testing.T) {
	var (
		dragonExUid int64 = 1000000
		tradeNo           = fmt.Sprint(time.Now().Unix())
		coinCode          = "usdt"
		volume            = "0.1"
		payToken          = "AimJY="
	)
	resp, _, err := apiV1.OpenDoUser2AppByDragonExUid(context.Background(), dragonExUid, tradeNo, coinCode, volume, scene, desc, device, payToken)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_OpenDoApp2UserByOpenId(t *testing.T) {
	var (
		openId   = "59991c0e35855a48bb1157295dd63ce4"
		tradeNo  = fmt.Sprint(time.Now().Unix())
		coinCode = "usdt"
		volume   = "0.1"
	)
	resp, _, err := apiV1.OpenDoApp2UserByOpenId(context.Background(), openId, tradeNo, coinCode, volume, scene, desc, device)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_OpenDoApp2UserByDragonExUid(t *testing.T) {
	var (
		dragonExUid int64 = 1000000
		tradeNo           = fmt.Sprint(time.Now().Unix())
		coinCode          = "usdt"
		volume            = "0.1"
	)
	resp, _, err := apiV1.OpenDoApp2UserByDragonExUid(context.Background(), dragonExUid, tradeNo, coinCode, volume, scene, desc, device)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_QueryOrderDetail(t *testing.T) {
	var (
		tradeNo = "1563964208"
	)
	resp, _, err := apiV1.QueryOrderDetail(context.Background(), tradeNo)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_ListOrders(t *testing.T) {
	var (
		dragonExUid int64 = 0
		coinCode          = ""
		direction   int64 = 0
		startTime   int64 = 0
		endTime     int64 = 0
		orderType   int64 = 1
		offset      int64 = 0
		limit       int64 = 10
	)
	resp, _, err := apiV1.ListOrdersByDragonExUid(context.Background(), dragonExUid, coinCode, direction, startTime, endTime, orderType, offset, limit)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_RedoPayCallback(t *testing.T) {
	var (
		tradeNo = "1563964208"
	)
	resp, _, err := apiV1.RedoPayCallback(context.Background(), tradeNo)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_ListUserCoinsByOpenId(t *testing.T) {
	var (
		openId = "59991c0e35855a48bb1157295dd63ce4"
	)

	resp, _, err := apiV1.ListUserCoinsByOpenId(context.Background(), openId)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_ListUserCoinsByDragonExUid(t *testing.T) {
	var (
		dragonExUid int64 = 1000000
	)

	resp, _, err := apiV1.ListUserCoinsByDragonExUid(context.Background(), dragonExUid)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_QueryUserCoinByOpenIdCoinId(t *testing.T) {
	var (
		openId       = "59991c0e35855a48bb1157295dd63ce4"
		coinId int64 = 1
	)
	resp, _, err := apiV1.QueryUserCoinByOpenIdCoinId(context.Background(), openId, coinId)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_QueryUserCoinByOpenIdCoinCode(t *testing.T) {
	var (
		openId   = "59991c0e35855a48bb1157295dd63ce4"
		coinCode = "usdt"
	)
	resp, _, err := apiV1.QueryUserCoinByOpenIdCoinCode(context.Background(), openId, coinCode)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_QueryUserCoinByDragonExUidCoinId(t *testing.T) {
	var (
		dragonExUid int64 = 1000000
		coinId      int64 = 1
	)
	resp, _, err := apiV1.QueryUserCoinByDragonExUidCoinId(context.Background(), dragonExUid, coinId)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_QueryUserCoinByDragonExUidCoinCode(t *testing.T) {
	var (
		dragonExUid int64 = 1000000
		coinCode          = "usdt"
	)
	resp, _, err := apiV1.QueryUserCoinByDragonExUidCoinCode(context.Background(), dragonExUid, coinCode)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_ListAdminOpenCoins(t *testing.T) {
	resp, _, err := apiV1.ListAdminOpenCoins(context.Background())

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_QueryAdminOpenCoinByOpenIdCoinId(t *testing.T) {
	var (
		coinId int64 = 104
	)
	resp, _, err := apiV1.QueryAdminOpenCoinByCoinId(context.Background(), coinId)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_QueryAdminOpenCoinByOpenIdCoinCode(t *testing.T) {
	var (
		coinCode = "usdt"
	)
	resp, _, err := apiV1.QueryAdminOpenCoinByCoinCode(context.Background(), coinCode)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_ListMarketReal(t *testing.T) {
	var (
		symbols = []string{"DT_USDT", "BTC_USDT", "EOS_USDT"}
	)
	resp, _, err := apiV1.ListMarketReal(context.Background(), symbols)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func displayRequestAndResponseMiddleware(ctx context.Context, oauth OAuth, req *http.Request, resp *http.Response) (err error) {
	fmt.Println(fmt.Sprintf("req.method = %+v", req.Method))
	fmt.Println(fmt.Sprintf("req.url = %+v", req.URL.String()))
	if req.Body != nil {
		reqBodyByte, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyByte))
		fmt.Println(fmt.Sprintf("req.body = %s", string(reqBodyByte)))
	}

	respBodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBodyByte))
	fmt.Println(fmt.Sprintf("resp.body = %s", string(respBodyByte)))
	return err
}
