package ctypes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// PlugmanClient is used to make requests to packages
type PlugmanClient struct {
	baseURL string
	client  http.Client
}

func NewPlugmanClient(baseURL string) *PlugmanClient {
	return &PlugmanClient{
		baseURL: baseURL,
		client: http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *PlugmanClient) ListManifest() {}

func (p *PlugmanClient) DoJSONPost(url string, body interface{}, result interface{}) error {
	return p.makeRequestWithBody("POST", fmt.Sprintf("%s%s", p.baseURL, url), body, result)
}

func (p *PlugmanClient) DoJSONGet(url string, result interface{}) error {
	return p.makeRequestWithBody("GET", fmt.Sprintf("%s%s", p.baseURL, url),nil, result)
}

func (p *PlugmanClient) makeRequestWithBody(method, url string, body interface{}, out interface{}) error {
	jsb, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(jsb))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := p.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	rsb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		var pkgErr error

		pkgErr = json.Unmarshal(rsb, &pkgErr)
		if pkgErr != nil {
			return errors.New("request pkgErr: " + string(rsb))
		} else {
			return errors.New(pkgErr.Error())
		}
	} else {
		err = json.Unmarshal(rsb, out)
		if err != nil {
			return err
		}

		return nil
	}
}