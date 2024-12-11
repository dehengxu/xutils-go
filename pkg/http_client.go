package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPClientParams struct {
	AutoRedirect bool          `json:"auto_redirect"`
	Timeout      time.Duration `json:"timeout"`
}

type HttpInfo struct {
	Headers    map[string][]string
	Status     string
	StatusCode int
	Location   string
	Body       []byte
	Error      error
}

type HttpError struct {
	// error
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e HttpError) Error() string {
	// return fmt.Sprintf("code: %v, message: %v", e.Code, e.Message)
	return e.Message
}

type WebClient struct {
	Cookie    string `json:"cookie"`
	UserAgent string `json:"user-agent"`
}

func (c *WebClient) GetCookie() string {
	return c.Cookie
}

func (c *WebClient) GetUserAgent() string {
	return c.UserAgent
}

func (c *WebClient) SetCookie(cookie string) {
	c.Cookie = cookie
}

func (c *WebClient) SetUserAgent(userAgent string) {
	c.UserAgent = userAgent
}

func (c *WebClient) GetHeaders() map[string]string {
	headers := make(map[string]string)
	headers["Cookie"] = c.Cookie
	headers["User-Agent"] = c.UserAgent
	return headers
}

func NewWebClient(params map[string]string) *WebClient {
	userAgent := params["user-agent"]
	if userAgent == "" {
		userAgent = "Mozilla"
	}
	return &WebClient{
		Cookie:    params["cookie"],
		UserAgent: userAgent,
	}
}

type HttpResponse struct {
	Code           string          `json:"code,omitempty"`
	Message        string          `json:"message,omitempty"`
	Result         bool            `json:"result,omitempty"`
	Data           json.RawMessage `json:"data,omitempty"`
	RequestID      string          `json:"requestID,omitempty"`
	WebAnnotations interface{}     `json:"web_annotations,omitempty"`
	HttpInfo       *HttpInfo       //`json:"httpInfo,omitempty"`
}

func NewHttpClient(param HTTPClientParams) *http.Client {
	client := &http.Client{
		Timeout: func() time.Duration {
			if param.Timeout == 0 {
				return 10 * time.Second
			} else {
				return param.Timeout
			}
		}(),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !param.AutoRedirect {
				return http.ErrUseLastResponse
			} else {
				return nil
			}
		},
	}
	return client
}

func _HttpRequest(c *WebClient, method string, url string, postBodyBytes *[]byte, headersIn map[string]string) (*HttpInfo, error) {
	client := NewHttpClient(HTTPClientParams{Timeout: 30 * time.Second})

	// fmt.Printf("postBody: %v\n", string(postBodyBytes))
	var req *http.Request
	var err error
	if postBodyBytes != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(*postBodyBytes))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		fmt.Printf("_Request create request failed: %v\n", err)
		return nil, err
	}
	headers := c.GetHeaders()
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	for k, v := range headersIn {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("_Request client do request failed: %v\n", err)
		return nil, HttpError{Code: -1, Message: fmt.Sprintf("client do request failed: %v", err)}
	}
	httpInfo := HttpInfo{StatusCode: res.StatusCode, Status: res.Status}
	defer res.Body.Close()
	if body, err := io.ReadAll(res.Body); err != nil {
		httpInfo.Error = err
		httpInfo.Body = []byte{}
	} else {
		httpInfo.Body = body
	}

	return &httpInfo, nil
}

type HttpRequestParamsInterface interface {
	GetParams() []interface{}
}
type HttpEmptyRequestParams struct {
	HttpRequestParamsInterface
}

func (r HttpEmptyRequestParams) GetParams() []interface{} {
	return []interface{}{}
}

type HttpRequestBodyInterface interface {
}
type HttpEmptyRequestBody struct {
	HttpRequestBodyInterface
}

type HttpResponseInterface interface {
}

type HttpRedirectResonse struct {
	HttpResponseInterface
	HttpInfo *HttpInfo
}

type APIRequestParams struct {
	BaseURL    string
	URLPattern string
	Cookie     string
	Headers    map[string]string
}

func URLBind(baseUrl, urlPattern string, queries ...interface{}) string {
	// url := fmt.Sprintf(urlPattern, queries...)
	url := XSprintf(urlPattern, queries...)
	url = fmt.Sprintf("%s%s", baseUrl, url)
	return url
}

type HttpAPIError struct {
	error
	// HTTPCode  int       `json:"httpCode"`
	// Location  string    `json:"location"`
	HTTPInfo  *HttpInfo `json:"httpInfo"`
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	RequestID string    `json:"requestID"`
}

func (e HttpAPIError) Error() string {
	// fmt.Printf("BCSAPIError type: %v, error: %v\n", reflect.TypeOf(e), e.error)
	return e.error.Error()
}

type XHttpResponse struct {
	HttpInfo *HttpInfo       `json:"httpInfo"`
	Data     json.RawMessage `json:"data"`
}

func HttpRequest[R1 HttpRequestParamsInterface, R2 HttpRequestBodyInterface](method string, param APIRequestParams, reqParam *R1, reqBody *R2) (*XHttpResponse, error) {
	fmt.Printf("args : %v\n", (*reqParam).GetParams())
	api_url := ""
	var args []interface{} = (*reqParam).GetParams()
	fmt.Printf("BCS_RequestAPI args: %v\n", args)
	api_url = URLBind(param.BaseURL, param.URLPattern, args[:]...)
	fmt.Printf("[%v] BCS_RequestAPI api_url: %v\n", method, api_url)

	var err error
	var reqBodyBytes []byte
	if reqBody != nil {
		reqBodyBytes, err = json.Marshal(reqBody)
		if err != nil {
			fmt.Printf("RequestAPI parse req body Error: %v\n", err)
			e := HttpAPIError{error: err, Code: 0, Message: fmt.Sprintf("PostAPI parse req body Error: %v\n", err), RequestID: ""}
			return nil, e //err
		}
	}

	b := NewWebClient(func() map[string]string {
		if param.Cookie != "" {
			return map[string]string{"cookie": param.Cookie}
		} else {
			return map[string]string{}
		}
	}())

	var httpRes XHttpResponse
	// var res HttpResponse
	var httpInfo *HttpInfo = nil

	if param.Headers == nil {
		httpInfo, err = _HttpRequest(b, method, api_url, &reqBodyBytes, nil)
	} else {
		httpInfo, err = _HttpRequest(b, method, api_url, &reqBodyBytes, param.Headers)
	}
	// res.HttpInfo = httpInfo
	if err != nil {
		fmt.Printf("PostAPI request Error: %v\n", err)
		e := HttpAPIError{error: err, Code: 0, Message: fmt.Sprintf("PostAPI request Error: %v\n", err), RequestID: ""}
		return nil, e
	}
	fmt.Printf("res.Date: %v\n", string(httpInfo.Body))

	httpRes = XHttpResponse{HttpInfo: httpInfo, Data: httpInfo.Body}
	return &httpRes, nil
}

// 解析业务 response
func RequestAPI[R1 HttpRequestParamsInterface, R2 HttpRequestBodyInterface, T1 HttpResponseInterface](method string, param APIRequestParams, reqParam *R1, reqBody *R2, response *T1) (*T1, error) {

	if httpRes, err := HttpRequest(method, param, reqParam, reqBody); err != nil {
		return nil, err
	} else {
		fmt.Printf("get json data: %v\n", string(httpRes.Data))

		err = json.Unmarshal(httpRes.Data, &response)
		if err != nil {
			return nil, err
		}

		return response, err
	}
}
