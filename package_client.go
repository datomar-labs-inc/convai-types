package ctypes

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// PackageClient is used to make requests to packages
type PackageClient struct {
	client http.Client
	pkg    *Package
}

func NewPackageClient(pkg *Package) *PackageClient {
	return &PackageClient{
		pkg: pkg,
		client: http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *PackageClient) DoJSONPost(url string, body interface{}, result interface{}) error {
	return p.makeRequestWithBody("POST", fmt.Sprintf("%s%s", p.pkg.BaseURL, url), p.pkg.SigningKey, body, result)
}

func (p *PackageClient) DoJSONGet(url string, result interface{}) error {
	return p.makeRequestWithBody("GET", fmt.Sprintf("%s%s", p.pkg.BaseURL, url), p.pkg.SigningKey, nil, result)
}

func (p *PackageClient) FetchManifest() (*Package, error) {
	var result Package

	err := p.DoJSONGet("/manifest", &result)
	if err != nil {
		return nil, err
	}

	p.pkg.Dispatches = result.Dispatches
	p.pkg.Links = result.Links
	p.pkg.Nodes = result.Nodes
	p.pkg.Events = result.Events

	return &result, nil
}

func (p *PackageClient) ExecuteNode(input *NodeCall) (*NodeCallResult, error) {
	var result NodeExecutionResponse

	err := p.DoJSONPost("/nodes/execute", NodeExecutionRequest{
		Calls: []NodeCall{*input},
	}, &result)
	if err != nil {
		return nil, err
	}

	return &result.Results[0], nil
}

func (p *PackageClient) ExecuteNodeMock(input *NodeCall) (*NodeCallResult, error) {
	var result NodeExecutionResponse

	err := p.DoJSONPost("/nodes/execute-mock", NodeExecutionRequest{
		Calls: []NodeCall{*input},
	}, &result)
	if err != nil {
		return nil, err
	}

	return &result.Results[0], nil
}

func (p *PackageClient) ExecuteLink(request *LinkExecutionRequest) (*LinkExecutionResponse, error) {
	var result LinkExecutionResponse

	err := p.DoJSONPost("/nodes/execute", request, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (p *PackageClient) ExecuteLinkMock(request *LinkExecutionRequest) (*LinkExecutionResponse, error) {
	var result LinkExecutionResponse

	err := p.DoJSONPost("/nodes/execute-mock", request, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (p *PackageClient) Dispatch(request *DispatchRequest) (*DispatchResponse, error) {
	var result DispatchResponse

	err := p.DoJSONPost("/dispatch/execute", request, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (p *PackageClient) DispatchMock(request *DispatchRequest) (*DispatchResponse, error) {
	var result DispatchResponse

	err := p.DoJSONPost("/dispatch/execute-mock", request, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// TODO
func (p *PackageClient) GetAsset(filename string) (io.Reader, error) {
	return nil, nil
}

// TODO
func (p *PackageClient) GetAssetBytes(filename string) ([]byte, error) {
	return nil, nil
}

func (p *PackageClient) makeRequestWithBody(method, url, signingToken string, body interface{}, out interface{}) error {
	jsb, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(jsb))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Convai-Signature", getSignature(jsb, signingToken))

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

func getSignature(body []byte, key string) string {
	mac := hmac.New(sha512.New, []byte(key))
	mac.Write(body)
	hash := mac.Sum(nil)

	return base64.StdEncoding.EncodeToString(hash)
}
