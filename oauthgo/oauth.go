package oauthgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type BeforeHandler func(oauth OAuth, req *http.Request) (err error)

type AfterHandler func(oauth OAuth, req *http.Request, resp *http.Response) (err error)

type OAuth interface {
	// return AppId
	GetAppId() (appId string)

	// return host
	GetHost() (host string)

	// return access key
	GetAccessKey() (accessKey string)

	// return secret key
	GetSecretKey() (secretKey string)

	// return check key
	GetCheckKey() (checkKey string)

	// make http.Request
	NewRequest(method, url string, values map[string]interface{}, header http.Header) (req *http.Request, err error)

	// do http request
	Do(req *http.Request, i interface{}) (hResp *http.Response, err error)

	// add middleware beforeHandlers doing request
	Before(handlers ...BeforeHandler)

	// add middleware afterHandlers doing request
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
	ErrNotSupportMethod = errors.New("not support this method")
	ErrNot200StatusCode = errors.New("not 200 http status code")
)

type DefaultOAuth struct {
	appId          string
	host           string
	accessKey      string
	secretKey      string
	checkKey       string
	beforeHandlers []BeforeHandler
	afterHandlers  []AfterHandler
}

func NewDefaultOAuth(appId, host, accessKey, secretKey, checkKey string) (client *DefaultOAuth) {
	return &DefaultOAuth{
		appId:     appId,
		host:      host,
		accessKey: accessKey,
		secretKey: secretKey,
		checkKey:  checkKey,
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

func (oauth *DefaultOAuth) GetCheckKey() string {
	return oauth.checkKey
}

func (oauth *DefaultOAuth) NewRequest(method, path string, values map[string]interface{}, header http.Header) (req *http.Request, err error) {
	switch strings.ToUpper(method) {
	case http.MethodGet, "":
		return oauth.makeGetRequest(path, values, header)

	case http.MethodPost:
		return oauth.makePostRequest(path, values, header)

	default:
		return req, ErrNotSupportMethod
	}
}

func (oauth *DefaultOAuth) Do(req *http.Request, i interface{}) (hResp *http.Response, err error) {
	// do something before doing request
	for _, handler := range oauth.beforeHandlers {
		err = handler(oauth, req)
		if err != nil {
			return
		}
	}

	client := http.Client{}
	hResp, err = client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = hResp.Body.Close()
	}()

	// do something after doing request
	for _, handler := range oauth.afterHandlers {
		err = handler(oauth, req, hResp)
		if err != nil {
			return
		}
	}

	if hResp.StatusCode != http.StatusOK {
		return hResp, ErrNot200StatusCode
	}

	body, err := ioutil.ReadAll(hResp.Body)
	if err != nil {
		return
	}
	hResp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	err = json.Unmarshal(body, i)
	if err != nil {
		return
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
	req, err = http.NewRequest(method, reqUrl, nil)
	if err != nil {
		return
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
		return
	}

	req, err = http.NewRequest(method, reqUrl, strings.NewReader(string(bodyByte)))
	if err != nil {
		return
	}

	req.Header.Add(HeaderDate, time.Now().UTC().Format(GmtFormat))
	req.Header.Add(HeaderContentType, MIMEApplicationJSON)
	req.Header.Add(HeaderAppId, oauth.appId)

	sign, err := sign(method, path, req.Header, oauth.secretKey)
	if err != nil {
		return
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
