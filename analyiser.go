package rpc

import (
	"net/http"
	"encoding/json"
	"strings"
)

type H2OAnalyser struct{
} 

func (analyser H2OAnalyser) Analyse(ret interface{}, resp *http.Response) (err error) {
	defer resp.Body.Close()

	if resp.StatusCode/100 == 2 {
		if ret != nil && resp.ContentLength != 0 {
			err = json.NewDecoder(resp.Body).Decode(ret)
			if err != nil {
				return
			}
		}
		if resp.StatusCode == 200 {
			return nil
		}
	}
	return ResponseError(resp)
}

// --------------------------------------------------------------------

type ErrorInfo struct {
	Err     string   `json:"error"`
	Reqid   string   `json:"reqid"`
	Details []string `json:"details"`
	Code    int      `json:"code"`
}

func (r *ErrorInfo) Error() string {
	msg, _ := json.Marshal(r)
	return string(msg)
}

// --------------------------------------------------------------------

type ErrorRet struct {
	Error string `json:"error"`
}

func ResponseError(resp *http.Response) (err error) {

	e := &ErrorInfo{
		Details: resp.Header["X-Log"],
		Reqid:   resp.Header.Get("X-Reqid"),
		Code:    resp.StatusCode,
	}
	if resp.StatusCode > 299 {
		if resp.ContentLength != 0 {
			if ct, ok := resp.Header["Content-Type"]; ok && strings.Contains(ct[0], "application/json") {
				var ret1 ErrorRet
				json.NewDecoder(resp.Body).Decode(&ret1)
				e.Err = ret1.Error
			}
		}
	}
	return e
}

