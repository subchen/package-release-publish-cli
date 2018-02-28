package bintray

import (
	"errors"
	"net/http"

	"github.com/subchen/go-curl"
	"github.com/subchen/go-stack/config/json"
)

const BINTRAY_API_PREFIX = "https://api.bintray.com"

//const BINTRAY_API_PREFIX = "http://127.0.0.1"

type bintrayClient struct {
	subject string
	apikey  string
}

func newBintrayClient(subject, apikey string) *bintrayClient {
	return &bintrayClient{
		subject: subject,
		apikey:  apikey,
	}
}

func (c *bintrayClient) newRequest() *curl.Request {
	req := curl.NewRequest()
	req.SetBasicAuth(c.subject, c.apikey)
	return req
}

// POST /repos/:subject/:repo
func (c *bintrayClient) getRespErr(resp *curl.Response, err error, forceCreate bool) error {
	if err != nil {
		return err
	}
	if !resp.OK() {
		if resp.StatusCode == 409 && forceCreate {
			return nil
		}
		data, err := resp.JSON()
		if err != nil {
			return errors.New(resp.Status)
		}
		return errors.New(json.NewQuery(data).Query("message").AsString())
	}
	return nil
}
