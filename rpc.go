package rpc

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"errors"
	"strings"
)

var UserAgent = "Golang h2object/rpc package"

// --------------------------------------------------------------------

func BuildHttpURL(host string, uri string, vals url.Values) *url.URL {
	u := &url.URL{
		Scheme: "http",
		Host: host,
		Path: uri,
	}
	if vals != nil {
		u.RawQuery = vals.Encode()	
	}
	return u
}

func BuildHttpsURL(host string, uri string, vals url.Values) *url.URL {
	u := &url.URL{
		Scheme: "https",
		Host: host,
		Path: uri,
	}
	if vals != nil {
		u.RawQuery = vals.Encode()	
	}
	return u	
}

// --------------------------------------------------------------------

type Client struct {
	c *http.Client
	analyzer Analyzer
}

var DefaultClient = &Client{c: http.DefaultClient, analyzer: H2OAnalyser{}}

func NewClient(c *http.Client, analyzer Analyzer) *Client {
	return &Client{
		c: c,
		analyzer: analyzer,
	}
}

// --------------------------------------------------------------------

type Logger interface {
	ReqId() string
	Xput(logs []string)
}

type Analyzer interface {
	Analyse(ret interface{}, resp *http.Response) error
}

// --------------------------------------------------------------------
func (r *Client) sent(l Logger, method string, u *url.URL, bodyType string, body io.Reader, bodyLength int) (resp *http.Response, err error) {
	var req *http.Request

	upperMethod := strings.ToUpper(method)
	switch upperMethod {
	case "GET":
		fallthrough
	case "POST":
		fallthrough
	case "PATCH":
		fallthrough
	case "PUT":
		fallthrough
	case "DELETE":
		req, err = http.NewRequest(upperMethod, u.String(), body)
		req.Header.Set("Content-Type", bodyType)
		req.ContentLength = int64(bodyLength)
		if err != nil {
			return
		}	
	default:
		err = errors.New("unsupport method: " + method)
		return
	}
	
	return r.do(l, req)
}

func (r *Client) do(l Logger,req *http.Request) (resp *http.Response, err error) {
	if l != nil {
		req.Header.Set("X-Reqid", r.logger.ReqId())
	}

	req.Header.Set("User-Agent", UserAgent)
	resp, err = r.c.Do(req)
	if err != nil {
		return
	}

	if l != nil {
		details := resp.Header["X-Log"]
		if len(details) > 0 {
			l.Xput(details)
		}
	}
	return
}

func (r *Client) Event() error {
	return nil
}

func (r *Client) Get(l Logger, u *url.URL, ret interface{}) error {
	resp, err := r.sent(l, "GET", u, "", nil, 0)
	if err != nil {
		return err
	}
	return r.analyzer.Analyse(ret, resp)
}

func (r *Client) PostBinary(l Logger, u *url.URL, rd io.Reader, length int64, ret interface{}) error {
	resp, err := r.sent(l, "POST", u, "application/octet-stream", rd, int(length))
	if err != nil {
		return err
	}
	return r.analyzer.Analyse(ret, resp)
}

func (r *Client) PutBinary(l Logger, u *url.URL, rd io.Reader, length int64, ret interface{}) error {
	resp, err := r.sent(l, "PUT", u, "application/octet-stream", rd, int(length))
	if err != nil {
		return err
	}
	return r.analyzer.Analyse(ret, resp)
}

func (r *Client) PostJson(l Logger, u *url.URL, data interface{}, ret interface{}) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := r.sent(l, "POST", u, "application/json", bytes.NewReader(msg), len(msg))
	if err != nil {
		return err
	}
	return r.analyzer.Analyse(ret, resp)
}

func (r *Client) PutJson(l Logger, u *url.URL, data interface{}, ret interface{}) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := r.sent(l, "PUT", u, "application/json", bytes.NewReader(msg), len(msg))
	if err != nil {
		return err
	}
	return r.analyzer.Analyse(ret, resp)
}

func (r *Client) PatchJson(l Logger, u *url.URL, data interface{}, ret interface{}) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := r.sent(l, "PATCH", u, "application/json", bytes.NewReader(msg), len(msg))
	if err != nil {
		return err
	}
	return r.analyzer.Analyse(ret, resp)
}

func (r *Client) PostForm(l Logger, u *url.URL, form map[string][]string, ret interface{}) error {
	msg := url.Values(form).Encode()
	resp, err := r.sent(l, "POST", u, "application/x-www-form-urlencoded", strings.NewReader(msg), len(msg))
	if err != nil {
		return err
	}
	return r.analyzer.Analyse(ret, resp)
}

func (r *Client) PutForm(l Logger, u *url.URL, form map[string][]string, ret interface{}) error {
	msg := url.Values(form).Encode()
	resp, err := r.sent(l, "PUT", u, "application/x-www-form-urlencoded", strings.NewReader(msg), len(msg))
	if err != nil {
		return err
	}
	return r.analyzer.Analyse(ret, resp)
}

func (r *Client) PatchForm(l Logger, u *url.URL, form map[string][]string, ret interface{}) error {
	msg := url.Values(form).Encode()
	resp, err := r.sent(l, "PATCH", u, "application/x-www-form-urlencoded", strings.NewReader(msg), len(msg))
	if err != nil {
		return err
	}
	return r.analyzer.Analyse(ret, resp)
}

func (r *Client) Delete(l Logger, u *url.URL, ret interface{}) error {
	resp, err := r.sent(l, "DELETE", u, "", nil, 0)
	if err != nil {
		return err
	}
	return r.analyzer.Analyse(ret, resp)
}

//! ---------------



