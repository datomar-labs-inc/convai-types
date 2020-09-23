package ctypes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// PackmanClient is used to make requests to packages
type PackmanClient struct {
	baseURL string
	client  http.Client
}

type ManifestResponse struct {
	Page      PageInfo  `json:"page"`
	Manifests []Package `json:"data"`
}

type PageInfo struct {
	PreviousCursor string `json:"previousCursor"`
	NextCursor     string `json:"nextCursor"`
	Count          int    `json:"count"`
}

func NewPackmanClient(baseURL string) *PackmanClient {
	return &PackmanClient{
		baseURL: baseURL,
		client: http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *PackmanClient) ListManifest(ids []uuid.UUID) (*ManifestResponse, error) {
	var stringFormIds []string

	for _, id := range ids {
		stringFormIds = append(stringFormIds, id.String())
	}

	idString := strings.Join(stringFormIds, ",")

	var manifests ManifestResponse

	err := p.DoJSONGet(fmt.Sprintf("/v1/manifest?ids=%s", idString), &manifests)
	if err != nil {
		return nil, err
	}

	return &manifests, nil
}

func (p *PackmanClient) DoJSONPost(url string, body interface{}, result interface{}) error {
	return p.makeRequestWithBody("POST", fmt.Sprintf("%s%s", p.baseURL, url), body, result)
}

func (p *PackmanClient) DoJSONGet(url string, result interface{}) error {
	return p.makeRequestWithBody("GET", fmt.Sprintf("%s%s", p.baseURL, url), nil, result)
}

func (p *PackmanClient) makeRequestWithBody(method, url string, body interface{}, out interface{}) error {
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
