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

func BuildHttpURL(addr string, uri string, vals url.Values) *url.URL {
	u := &url.URL{
		Scheme: "http",
		Host: addr,
		Path: uri,
	}
	if vals != nil {
		u.RawQuery = vals.Encode()	
	}
	return u
}

func BuildHttpsURL(addr string, uri string, vals url.Values) *url.URL {
	u := &url.URL{
		Scheme: "https",
		Host: addr,
		Path: uri,
	}
	if vals != nil {
		u.RawQuery = vals.Encode()	
	}
	return u	
}

// package's interface definitions

//! interface before request sending
type PreRequest interface {
	Do(*http.Request) *http.Request
}

//! interface for the logger server side log info
type Logger interface {
	ReqId() string
	Xput(logs []string)
}

//! interface for the response analyser
type Analyzer interface {
	Analyse(ret interface{}, resp *http.Response) error
}

type Client struct {
	conn     	*http.Client
	prepare   	PreRequest
	analyzer 	Analyzer
}

func NewClient(analyzer Analyzer) *Client {
	return &Client{
		conn: &http.Client{},
		analyzer: analyzer,
	}
}

func NewClient2(conn *http.Client, analyzer Analyzer) *Client {
	return &Client{
		conn: conn,
		analyzer: analyzer,
	}
}

func (c *Client) Prepare(prepare PreRequest) {
	c.prepare = prepare
}

func (c *Client) sent(l Logger, method string, u *url.URL, bodyType string, body io.Reader, bodyLength int) (resp *http.Response, err error) {
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
		if bodyLength > 0 {
			req.ContentLength = int64(bodyLength)	
		}		
		if err != nil {
			return
		}	
	default:
		err = errors.New("unsupport method: " + method)
		return
	}
	
	return c.do(l, req)
}

func (c *Client) do(l Logger,req *http.Request) (resp *http.Response, err error) {
	if l != nil {
		req.Header.Set("X-Reqid", l.ReqId())
	}

	var real *http.Request
	if c.prepare != nil {
		real = c.prepare.Do(req)
	} else {
		real = req
	}
	real.Header.Set("User-Agent", UserAgent)

	resp, err = c.conn.Do(real)
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

func (c *Client) Get(l Logger, u *url.URL, ret interface{}) error {
	resp, err := c.sent(l, "GET", u, "", nil, 0)
	if err != nil {
		return err
	}
	return c.analyzer.Analyse(ret, resp)
}

func (c *Client) GetResponse(l Logger, u *url.URL) (*http.Response, error) {
	return c.sent(l, "GET", u, "", nil, 0)
}

func (c *Client) Post(l Logger, u *url.URL, bodyType string, body io.Reader, length int64, ret interface{}) error {
	resp, err := c.sent(l, "POST", u, bodyType, body, int(length))
	if err != nil {
		return err
	}
	return c.analyzer.Analyse(ret, resp)
}

func (c *Client) Put(l Logger, u *url.URL, bodyType string, body io.Reader, length int64, ret interface{}) error {
	resp, err := c.sent(l, "PUT", u, bodyType, body, int(length))
	if err != nil {
		return err
	}
	return c.analyzer.Analyse(ret, resp)
}

func (c *Client) Patch(l Logger, u *url.URL, bodyType string, body io.Reader, length int64, ret interface{}) error {
	resp, err := c.sent(l, "PATCH", u, bodyType, body, int(length))
	if err != nil {
		return err
	}
	return c.analyzer.Analyse(ret, resp)
}

func (c *Client) Delete(l Logger, u *url.URL, ret interface{}) error {
	resp, err := c.sent(l, "DELETE", u, "", nil, 0)
	if err != nil {
		return err
	}
	return c.analyzer.Analyse(ret, resp)
}

func (c *Client) PostJson(l Logger, u *url.URL, data interface{}, ret interface{}) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := c.sent(l, "POST", u, "application/json", bytes.NewReader(msg), len(msg))
	if err != nil {
		return err
	}
	return c.analyzer.Analyse(ret, resp)
}

func (c *Client) PutJson(l Logger, u *url.URL, data interface{}, ret interface{}) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := c.sent(l, "PUT", u, "application/json", bytes.NewReader(msg), len(msg))
	if err != nil {
		return err
	}
	return c.analyzer.Analyse(ret, resp)
}

func (c *Client) PatchJson(l Logger, u *url.URL, data interface{}, ret interface{}) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := c.sent(l, "PATCH", u, "application/json", bytes.NewReader(msg), len(msg))
	if err != nil {
		return err
	}
	return c.analyzer.Analyse(ret, resp)
}

func (c *Client) PostForm(l Logger, u *url.URL, form map[string][]string, ret interface{}) error {
	msg := url.Values(form).Encode()
	resp, err := c.sent(l, "POST", u, "application/x-www-form-urlencoded", strings.NewReader(msg), len(msg))
	if err != nil {
		return err
	}
	return c.analyzer.Analyse(ret, resp)
}

func (c *Client) PutForm(l Logger, u *url.URL, form map[string][]string, ret interface{}) error {
	msg := url.Values(form).Encode()
	resp, err := c.sent(l, "PUT", u, "application/x-www-form-urlencoded", strings.NewReader(msg), len(msg))
	if err != nil {
		return err
	}
	return c.analyzer.Analyse(ret, resp)
}

func (c *Client) PatchForm(l Logger, u *url.URL, form map[string][]string, ret interface{}) error {
	msg := url.Values(form).Encode()
	resp, err := c.sent(l, "PATCH", u, "application/x-www-form-urlencoded", strings.NewReader(msg), len(msg))
	if err != nil {
		return err
	}
	return c.analyzer.Analyse(ret, resp)
}

func (c *Client) PostMultiPartForm(l Logger, u *url.URL, multipart *MultipartForm, ret interface{}) error {
	ct, err := multipart.ContentType()
	if err != nil {
		return err
	}

	rd, err := multipart.Reader()
	if err != nil {
		return err
	}

	sz, err := multipart.Size()
	if err != nil {
		return err
	}

	resp, err := c.sent(l, "POST", u, ct, rd, sz)
	if err != nil {
		return err
	}
	return c.analyzer.Analyse(ret, resp)
}

func (c *Client) PutMultiPartForm(l Logger, u *url.URL, multipart *MultipartForm, ret interface{}) error {
		ct, err := multipart.ContentType()
	if err != nil {
		return err
	}

	rd, err := multipart.Reader()
	if err != nil {
		return err
	}

	sz, err := multipart.Size()
	if err != nil {
		return err
	}

	resp, err := c.sent(l, "PUT", u, ct, rd, sz)
	if err != nil {
		return err
	}
	return c.analyzer.Analyse(ret, resp)
}

func (c *Client) PatchMultiPartForm(l Logger, u *url.URL, multipart *MultipartForm, ret interface{}) error {
	ct, err := multipart.ContentType()
	if err != nil {
		return err
	}

	rd, err := multipart.Reader()
	if err != nil {
		return err
	}

	sz, err := multipart.Size()
	if err != nil {
		return err
	}

	resp, err := c.sent(l, "PATCH", u, ct, rd, sz)
	if err != nil {
		return err
	}
	return c.analyzer.Analyse(ret, resp)
}


