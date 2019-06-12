package oauthgo

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

var (
	apiV1 *OAuthV1
)

func TestMain(m *testing.M) {
	const (
		appId     = "appidfortest"
		host      = "https://oauth.dragonex.io"
		accessKey = "87079e4662685c40a884baa744f571b4"
		secretKey = "a24d0648e60a5c7a9d250137c472d8f4"
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
		code   = "f38f095c79"
		scopes = []int{ScopeLogin}
	)
	resp, _, err := apiV1.DoLogin(context.Background(), code, state, device, scopes)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_RefreshToken(t *testing.T) {
	var (
		accessToken  = "3d6524isnEAoGBx5p3Z0C/3IBXesM7biMwDZKkvhj6lsXgWAb6VSlr/433k6q0hX8IG1WXFvyiZP9TBdMNm7Zwq2WdGdZ3ntdVCqBw4bZI5PMsRsN/L+aR7813rnfCKJA6GQrhlLsZV5l6hRbM/4G/4h1cI/04V9BYmd/9fKMijz0y1v2OCcx4Ge7xe8vwxMj+4JpbWIedbWPWFoZV10oLeHaQm+K51rMJ8Onxnid60="
		refreshToken = "6WKXRHI6IkpXVCJ96eyJvIjoiNTk5OTFjMGUzNTg1NWE0OGJiMTE1NzI5NWRkNjNjZTQiLCJhIjoiYXBwaWRmb3J0ZXN0IiwiZXhwIjoxNTU2MjY0ODY1LCJuYmYiOjE1NTU2NjAwNjV96TEVzAoAGhJGo94h5BqDDcd7AMkA28EZLpPkneWYn1G4"
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
		tradeNo     = fmt.Sprint(time.Now().Unix())
		coinCode    = "usdt"
		volume      = "0.1"
		redirectUrl = "https://www.google.com"
	)
	resp, _, err := apiV1.PreUser2App(context.Background(), tradeNo, coinCode, volume, scene, desc, device, state, redirectUrl)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_DoApp2UserByOpenId(t *testing.T) {
	var (
		openId   = "9b442f4561905244a02f1cbc4bd4db5c"
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

func TestOAuthV1_QueryOrderDetail(t *testing.T) {
	var (
		tradeNo = ""
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
		offset      int64 = 0
		limit       int64 = 10
	)
	resp, _, err := apiV1.ListOrdersByDragonExUid(context.Background(), dragonExUid, coinCode, direction, startTime, endTime, offset, limit)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_RedoPayCallback(t *testing.T) {
	var (
		tradeNo = ""
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

func displayRequestAndResponseMiddleware(ctx context.Context, oauth OAuth, req *http.Request, resp *http.Response) (err error) {
	reqBodyByte, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyByte))

	fmt.Println(fmt.Sprintf("req.method = %+v", req.Method))
	fmt.Println(fmt.Sprintf("req.url = %+v", req.URL.String()))
	fmt.Println(fmt.Sprintf("req.body = %s", string(reqBodyByte)))

	respBodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBodyByte))

	fmt.Println(fmt.Sprintf("resp.body = %s", string(respBodyByte)))
	return
}
