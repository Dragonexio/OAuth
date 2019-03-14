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

type BeforeHandler func(req *http.Request) (err error)

type AfterHandler func(req *http.Request, resp *http.Response) (err error)

type OAuth interface {
	// make http.Request
	NewRequest(method, url string, values map[string]interface{}, header http.Header) (req *http.Request, err error)

	// do http request
	Do(req *http.Request, i interface{}) (hResp *http.Response, err error)

	// add middleware before doing request
	Before(handlers ...BeforeHandler)

	// add middleware after doing request
	After(handlers ...AfterHandler)

	// return AppId
	GetAppId() (appId string)
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
	beforeHandlers []BeforeHandler
	afterHandlers  []AfterHandler
}

func NewDefaultOAuth(appId, host, accessKey, secretKey string) (client *DefaultOAuth) {
	return &DefaultOAuth{
		appId:     appId,
		host:      host,
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

func (d *DefaultOAuth) Do(req *http.Request, i interface{}) (hResp *http.Response, err error) {
	// before doing request
	for _, handler := range d.beforeHandlers {
		err = handler(req)
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

	// after doing request
	for _, handler := range d.afterHandlers {
		err = handler(req, hResp)
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

func (d *DefaultOAuth) NewRequest(method, path string, values map[string]interface{}, header http.Header) (req *http.Request, err error) {
	switch strings.ToUpper(method) {
	case http.MethodGet, "":
		return d.makeGetRequest(path, values, header)

	case http.MethodPost:
		return d.makePostRequest(path, values, header)

	default:
		return req, ErrNotSupportMethod
	}
}

func (d *DefaultOAuth) Before(handlers ...BeforeHandler) {
	if d.beforeHandlers == nil {
		d.beforeHandlers = make([]BeforeHandler, 0, len(handlers))
	}

	d.beforeHandlers = append(d.beforeHandlers, handlers...)
}

func (d *DefaultOAuth) After(handlers ...AfterHandler) {
	if d.afterHandlers == nil {
		d.afterHandlers = make([]AfterHandler, 0, len(handlers))
	}

	d.afterHandlers = append(d.afterHandlers, handlers...)
}

func (d *DefaultOAuth) GetAppId() string {
	return d.appId
}
func (d *DefaultOAuth) makeGetRequest(path string, values map[string]interface{}, header http.Header) (req *http.Request, err error) {
	method := http.MethodGet

	if strings.HasSuffix(d.host, "/") {
		d.host = d.host[:len(d.host)-1]
	}

	valuesStringSlice := make([]string, 0, len(values))
	for k, v := range values {
		valuesStringSlice = append(valuesStringSlice, fmt.Sprintf("%s=%s", k, fmt.Sprint(v)))
	}

	reqUrl := fmt.Sprintf("%s%s?%s", d.host, path, strings.Join(valuesStringSlice, "&"))
	req, err = http.NewRequest(method, reqUrl, nil)
	if err != nil {
		return
	}

	for k, v := range header {
		req.Header.Add(k, strings.Join(v, ""))
	}
	return
}

func (d *DefaultOAuth) makePostRequest(path string, values map[string]interface{}, header http.Header) (req *http.Request, err error) {
	method := http.MethodPost

	if strings.HasSuffix(d.host, "/") {
		d.host = d.host[:len(d.host)-1]
	}
	reqUrl := fmt.Sprintf("%s%s", d.host, path)

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
	req.Header.Add(HeaderAppId, d.appId)

	sign, err := sign(method, path, req.Header, d.secretKey)
	if err != nil {
		return
	}
	auth := fmt.Sprintf("%s:%s", d.accessKey, sign)
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
