package pkg

import "net/http"

type HTTPClientParams struct {
	AutoRedirect bool `json:"auto_redirect"`
}

func NewHttpClient(param HTTPClientParams) *http.Client {
	client := &http.Client{
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
