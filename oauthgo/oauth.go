package oauthgo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type BeforeHandler func(ctx context.Context, oauth OAuth, req *http.Request) (err error)

type AfterHandler func(ctx context.Context, oauth OAuth, req *http.Request, resp *http.Response) (err error)

type OAuth interface {
	// return AppId
	GetAppId() (appId string)

	// return host
	GetHost() (host string)

	// return access key
	GetAccessKey() (accessKey string)

	// return secret key
	GetSecretKey() (secretKey string)

	// return sign key
	GetSignKey() (checkKey string)

	// make http.Request
	NewRequest(method, url string, values map[string]interface{}, header http.Header) (req *http.Request, err error)

	// do http request
	Do(ctx context.Context, req *http.Request, i interface{}) (hResp *http.Response, err error)

	// add middleware before doing request
	Before(handlers ...BeforeHandler)

	// add middleware after doing request
	After(handlers ...AfterHandler)
}

const (
	HeaderAuth          = "Auth"
	HeaderAppId         = "App-Id"
	HeaderDate          = "Date"
	HeaderContentType   = "Content-Type"
	GmtFormat           = "Mon, 02 Jan 2006 15:04:05 GMT"
	MIMEApplicationJSON = "application/json"
)

var (
	ErrDoNotSupportThisMethod = errors.New("do not support this method")
	ErrNot200StatusCode       = errors.New("not 200 http status code")
)

type DefaultOAuth struct {
	appId          string
	host           string
	accessKey      string
	secretKey      string
	signKey        string
	beforeHandlers []BeforeHandler
	afterHandlers  []AfterHandler
}

func NewDefaultOAuth(appId, host, accessKey, secretKey, signKey string) (client *DefaultOAuth) {
	return &DefaultOAuth{
		appId:     appId,
		host:      host,
		accessKey: accessKey,
		secretKey: secretKey,
		signKey:   signKey,
	}
}

func (oauth *DefaultOAuth) GetAppId() string {
	return oauth.appId
}

func (oauth *DefaultOAuth) GetHost() string {
	return oauth.host
}

func (oauth *DefaultOAuth) GetAccessKey() string {
	return oauth.accessKey
}

func (oauth *DefaultOAuth) GetSecretKey() string {
	return oauth.secretKey
}

func (oauth *DefaultOAuth) GetSignKey() string {
	return oauth.signKey
}

func (oauth *DefaultOAuth) NewRequest(method, path string, values map[string]interface{}, header http.Header) (req *http.Request, err error) {
	switch strings.ToUpper(method) {
	case http.MethodGet, "":
		return oauth.makeGetRequest(path, values, header)

	case http.MethodPost:
		return oauth.makePostRequest(path, values, header)

	default:
		return req, ErrDoNotSupportThisMethod
	}
}

func (oauth *DefaultOAuth) Do(ctx context.Context, req *http.Request, i interface{}) (hResp *http.Response, err error) {
	var reqBodyByte []byte
	if req.Body != nil {
		if reqBodyByte, err = ioutil.ReadAll(req.Body); err != nil {
			return hResp, err
		}
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyByte))

	// do something before doing request
	for _, handler := range oauth.beforeHandlers {
		if err = handler(ctx, oauth, req); err != nil {
			return hResp, err
		}
	}

	// do request
	client := http.Client{}
	if hResp, err = client.Do(req); err != nil {
		return hResp, err
	}
	defer func() {
		_ = hResp.Body.Close()
	}()
	req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyByte))

	// do something after doing request
	for _, handler := range oauth.afterHandlers {
		if err = handler(ctx, oauth, req, hResp); err != nil {
			return hResp, err
		}
	}

	if hResp.StatusCode != http.StatusOK {
		return hResp, ErrNot200StatusCode
	}

	body, err := ioutil.ReadAll(hResp.Body)
	if err != nil {
		return hResp, err
	}
	hResp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if err = json.Unmarshal(body, i); err != nil {
		return hResp, err
	}
	return
}

func (oauth *DefaultOAuth) Before(handlers ...BeforeHandler) {
	if oauth.beforeHandlers == nil {
		oauth.beforeHandlers = make([]BeforeHandler, 0, len(handlers))
	}

	oauth.beforeHandlers = append(oauth.beforeHandlers, handlers...)
}

func (oauth *DefaultOAuth) After(handlers ...AfterHandler) {
	if oauth.afterHandlers == nil {
		oauth.afterHandlers = make([]AfterHandler, 0, len(handlers))
	}

	oauth.afterHandlers = append(oauth.afterHandlers, handlers...)
}

func (oauth *DefaultOAuth) makeGetRequest(path string, values map[string]interface{}, header http.Header) (req *http.Request, err error) {
	method := http.MethodGet

	if strings.HasSuffix(oauth.host, "/") {
		oauth.host = oauth.host[:len(oauth.host)-1]
	}

	valuesStringSlice := make([]string, 0, len(values))
	for k, v := range values {
		valuesStringSlice = append(valuesStringSlice, fmt.Sprintf("%s=%s", k, fmt.Sprint(v)))
	}

	reqUrl := fmt.Sprintf("%s%s?%s", oauth.host, path, strings.Join(valuesStringSlice, "&"))
	if req, err = http.NewRequest(method, reqUrl, nil); err != nil {
		return req, err
	}

	for k, v := range header {
		req.Header.Add(k, strings.Join(v, ""))
	}
	return
}

func (oauth *DefaultOAuth) makePostRequest(path string, values map[string]interface{}, header http.Header) (req *http.Request, err error) {
	method := http.MethodPost

	if strings.HasSuffix(oauth.host, "/") {
		oauth.host = oauth.host[:len(oauth.host)-1]
	}
	reqUrl := fmt.Sprintf("%s%s", oauth.host, path)

	bodyByte, err := json.Marshal(values)
	if err != nil {
		return req, err
	}

	if req, err = http.NewRequest(method, reqUrl, strings.NewReader(string(bodyByte))); err != nil {
		return req, err
	}

	req.Header.Add(HeaderDate, time.Now().UTC().Format(GmtFormat))
	req.Header.Add(HeaderContentType, MIMEApplicationJSON)
	req.Header.Add(HeaderAppId, oauth.appId)

	sign, err := sign(method, path, req.Header, oauth.secretKey)
	if err != nil {
		return req, err
	}
	auth := fmt.Sprintf("%s:%s", oauth.accessKey, sign)
	req.Header.Add(HeaderAuth, auth)

	for k, v := range header {
		req.Header.Add(k, strings.Join(v, ""))
	}
	return
}

type BaseResponse struct {
	Ok   bool   `json:"ok"`
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

func (r BaseResponse) IsOk() bool {
	return r.Code == CodeOk
}
