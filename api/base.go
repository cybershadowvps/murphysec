package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ErrTokenInvalid = errors.New("Token invalid")
var ErrServerRequest = errors.New("Send request failed")
var ErrParseErrMsg = errors.New("Parse error message failed")
var ErrTimeout = errors.New("API request timeout")

var C *Client

type Client struct {
	client  *http.Client
	baseUrl string
	Token   string
}

func NewClient(baseUrl string) *Client {
	c := new(http.Client)
	p := regexp.MustCompile("/*$")
	baseUrl = p.ReplaceAllString(strings.TrimSpace(baseUrl), "")
	c.Timeout = time.Second * 300
	i, e := strconv.Atoi(os.Getenv("API_TIMEOUT"))
	if e == nil && i > 0 {
		c.Timeout = time.Duration(int64(time.Second) * int64(i))
	}
	cl := &Client{client: c, baseUrl: baseUrl}
	return cl
}

func (c *Client) POST(relUri string, body io.Reader) *http.Request {
	u, e := http.NewRequest(http.MethodPost, c.baseUrl+relUri, body)
	if e != nil {
		panic(e)
	}
	return u
}

func (c *Client) PostJson(relUri string, a interface{}) *http.Request {
	u := c.POST(relUri, bytes.NewReader(must.Byte(json.Marshal(a))))
	u.Header.Set("Content-Type", "application/json")
	return u
}

func (c *Client) GET(relUri string) *http.Request {
	u, e := http.NewRequest(http.MethodGet, c.baseUrl+relUri, nil)
	if e != nil {
		panic(e)
	}
	return u
}

func (c *Client) Do(req *http.Request, resBody interface{}) error {
	var noBody bool
	if t := reflect.TypeOf(resBody); t == nil {
		noBody = true
	} else {
		if t.Kind() != reflect.Ptr {
			panic("resBody must be a pointer")
		}
	}
	res, e := c.client.Do(req)
	if e != nil {
		e := e.(*url.Error)
		if e.Timeout() {
			return ErrTimeout
		}
		return errors.Wrap(ErrServerRequest, e.Error())
	}
	data, e := io.ReadAll(res.Body)
	if e != nil {
		return errors.Wrap(ErrServerRequest, "read response body failed:"+e.Error())
	}
	defer res.Body.Close()
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		if noBody {
			return nil
		}
		if e := json.Unmarshal(data, resBody); e != nil {
			return errors.Wrap(ErrServerRequest, "parse response body as json failed")
		}
		return nil
	}
	if res.StatusCode >= 400 {
		var m CommonApiErr
		if e := json.Unmarshal(data, &m); e != nil {
			return errors.Wrap(ErrServerRequest, "parse error message failed.")
		}
		return &m
	}
	return errors.Wrap(ErrServerRequest, fmt.Sprintf("http code %d - %s", res.StatusCode, res.Status))
}

type CommonApiErr struct {
	EError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Details string `json:"details"`
	} `json:"error"`
}

func (c *CommonApiErr) Error() string {
	return fmt.Sprintf("[%d]%s", c.EError.Code, c.EError.Details)
}

func readCommonErr(data []byte, statusCode int) error {
	panic(1)
}

func readHttpBody(res *http.Response) ([]byte, error) {
	data, e := io.ReadAll(res.Body)
	if e != nil {
		logger.Warn.Println("read body failed.", e.Error())
		return nil, e
	}
	logger.Debug.Println(string(data))
	logger.Debug.Println("body size", len(data), "bytes")
	_ = res.Body.Close()
	return data, e
}
