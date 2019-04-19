package oauthgo

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	ErrDifferentResponseSign = errors.New("different sign")
)

func CheckResponseMiddleware(oauth OAuth, req *http.Request, resp *http.Response) (err error) {
	dragonExTs := resp.Header.Get("dragonex-ts")
	dragonExSign := resp.Header.Get("dragonex-sign")

	respBodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBodyByte))
	strToSign := fmt.Sprintf("%s%s%s", string(respBodyByte), dragonExTs, oauth.GetCheckKey())
	h := md5.New()
	_, err = io.WriteString(h, strToSign)
	if err != nil {
		return err
	}

	sign := fmt.Sprintf("%x", h.Sum(nil))
	if sign[:8] != dragonExSign {
		return ErrDifferentResponseSign
	}

	return nil
}
