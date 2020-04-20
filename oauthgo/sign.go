package oauthgo

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

func sign(method, path string, header http.Header, secretKey string) (sign string, err error) {
	newHeaders := map[string]string{}
	for k, v := range header {
		newHeaders[strings.ToLower(k)] = strings.Join(v, "")
	}

	contentSha1 := newHeaders["content-sha1"]
	contentType := newHeaders["content-type"]
	date := newHeaders["date"]

	dragonHeaders := make([]string, 0, len(newHeaders))
	canonicalizedDragonExHeaders := ""
	for k, v := range newHeaders {
		if strings.HasPrefix(k, "dragonex-") {
			dragonHeaders = append(dragonHeaders, fmt.Sprintf("%s:%s", k, v))
		}
	}
	sort.Strings(dragonHeaders)
	if len(dragonHeaders) != 0 {
		canonicalizedDragonExHeaders = strings.Join(dragonHeaders, "\n")
		canonicalizedDragonExHeaders += "\n"
	}

	stringsToSignSlice := []string{
		strings.ToUpper(method),
		contentSha1,
		contentType,
		date,
		canonicalizedDragonExHeaders,
	}
	stringToSign := strings.Join(stringsToSignSlice, "\n")
	stringToSign += path

	sha1Hash := hmac.New(sha1.New, []byte(secretKey))
	if _, err = sha1Hash.Write([]byte(stringToSign)); err != nil {
		return
	}

	sign = base64.StdEncoding.EncodeToString(sha1Hash.Sum(nil))
	return
}
