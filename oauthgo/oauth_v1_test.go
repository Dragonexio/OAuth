package oauthgo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	apiV1 *OAuthV1
)

func TestMain(m *testing.M) {
	const (
		appId     = "appidfortest"
		host      = "http://127.0.0.1:9101"
		accessKey = "87079e4662685c40a884baa744f571b4"
		secretKey = "a24d0648e60a5c7a9d250137c472d8f4"
		checkKey  = "testKey"
	)

	apiV1 = NewOAuthV1(appId, host, accessKey, secretKey, checkKey)

	apiV1.After(DisplayRequestAndRespponseMiddleware)
	apiV1.After(CheckResponseMiddleware)

	fmt.Println(fmt.Sprintf("%+v", apiV1))
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
		code   = "fbbd3724dc"
		scopes = []int{ScopeLogin}
	)
	resp, _, err := apiV1.DoLogin(code, state, device, scopes)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_RefreshToken(t *testing.T) {
	var (
		accessToken  = "sXGOiJIUzI1NiIsInR5cCI6IkpXVCJ9GeyJhIjoiYXBwaWRmb3J0ZXN0IiwibyI6IjliNDQyZjQ1NjE5MDUyNDRhMDJmMWNiYzRiZDRkYjVjIiwiZXhwIjoxNTUyNDczMzkxLCJuYmYiOjE1NTIzODY5OTF9GRwYMXjx1o7kzg_PW6u1f5MvUtMVyPADp7tlsnUHx4ew"
		refreshToken = "hNySYqfEOiJIUzI1NiIsInR5cCI6IkpXVCJ9FeyJvIjoiOWI0NDJmNDU2MTkwNTI0NGEwMmYxY2JjNGJkNGRiNWMiLCJhIjoiYXBwaWRmb3J0ZXN0IiwiZXhwIjoxNTUyOTkxNzkxLCJuYmYiOjE1NTIzODY5OTF9F01vE7GZjosaJrQDU-Gj6ZqG5Eo9vuvJ-kIqOaF"
	)
	resp, _, err := apiV1.RefreshToken(accessToken, refreshToken)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_LogoutToken(t *testing.T) {
	var (
		accessToken = "hDANbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9deyJhIjoiYXBwaWRmb3J0ZXN0IiwibyI6IjliNDQyZjQ1NjE5MDUyNDRhMDJmMWNiYzRiZDRkYjVjIiwiZXhwIjoxNTUyNDc0MDEyLCJuYmYiOjE1NTIzODc2MTJ9dR_qg7ytCLiSi7KzLOZt-HdckbZDxSCFrp1B789"
	)
	resp, _, err := apiV1.LogoutToken(accessToken)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_QueryUserDetail(t *testing.T) {
	var (
		openId = "59991c0e35855a48bb1157295dd63ce4"
	)
	resp, _, err := apiV1.QueryUserDetail(openId)

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
	resp, _, err := apiV1.PreUser2App(tradeNo, coinCode, volume, scene, desc, device, state, redirectUrl)

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
	resp, _, err := apiV1.DoApp2UserByOpenId(openId, tradeNo, coinCode, volume, scene, desc, device)

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
	resp, _, err := apiV1.DoApp2UserByDragonExUid(dragonExUid, tradeNo, coinCode, volume, scene, desc, device)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_QueryOrderDetail(t *testing.T) {
	var (
		tradeNo = "1552289753"
	)
	resp, _, err := apiV1.QueryOrderDetail(tradeNo)

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
	resp, _, err := apiV1.ListOrdersByDragonExUid(dragonExUid, coinCode, direction, startTime, endTime, offset, limit)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}

func TestOAuthV1_RedoPayCallback(t *testing.T) {
	var (
		tradeNo = ""
	)
	resp, _, err := apiV1.RedoPayCallback(tradeNo)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.Ok)
}
