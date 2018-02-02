package bintray

import (
	"errors"
	"net/http"

	"github.com/mozillazg/request"
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

func (c *bintrayClient) newRequest() *request.Request {
	req := request.NewRequest(new(http.Client))
	req.BasicAuth = request.BasicAuth{c.subject, c.apikey}
	return req
}

// POST /repos/:subject/:repo
func (c *bintrayClient) getRespErr(resp *request.Response, err error, forceCreate bool) error {
	if err != nil {
		return err
	}
	if !resp.OK() {
		if resp.StatusCode == 409 && forceCreate {
			return nil
		}
		json, err := resp.Json()
		if err != nil {
			return errors.New(resp.Reason())
		}
		return errors.New(json.Get("message").MustString())
	}
	return nil
}
