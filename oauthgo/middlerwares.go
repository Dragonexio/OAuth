package oauthgo

import (
	"bytes"
	"context"
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

func CheckResponseMiddleware(ctx context.Context, oauth OAuth, req *http.Request, resp *http.Response) (err error) {
	dragonExTs := resp.Header.Get("dragonex-ts")
	dragonExSign := resp.Header.Get("dragonex-sign")

	respBodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.WithStack(err)
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBodyByte))
	strToSign := fmt.Sprintf("%s%s%s", string(respBodyByte), dragonExTs, oauth.GetSignKey())
	h := md5.New()
	_, err = io.WriteString(h, strToSign)
	if err != nil {
		return errors.WithStack(err)
	}

	sign := fmt.Sprintf("%x", h.Sum(nil))
	if sign[:8] != dragonExSign {
		return errors.WithStack(ErrDifferentResponseSign)
	}

	return errors.WithStack(err)
}
