package ctypes

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func newPackageClient() *PackageClient {
	return NewPackageClient(&Package{
		DBPackage:  DBPackage{
			BaseURL:        "http://localhost:5555",
			SigningKey:     "bubbles",
		},
		Nodes:      nil,
		Links:      nil,
		Events:     nil,
		Dispatches: nil,
	})
}

func TestNewPackageClient(t *testing.T) {
	type args struct {
		pkg *Package
	}
	tests := []struct {
		name string
		args args
		want *PackageClient
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPackageClient(tt.args.pkg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPackageClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackageClient_Dispatch(t *testing.T) {
	pc := newPackageClient()

	res, err := pc.Dispatch(&DispatchRequest{
		Dispatches: []DispatchCall{
			{
				RequestID:       uuid.Must(uuid.NewRandom()),
				ID:              "example_dispatch",
				ContextTree:     ContextTreeSlice{},
				MessageBody:     XMLResponse{},
				PackageSettings: MemoryContainer{},
				Sequence:        0,
			},
		},
	})
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", res)
}

func TestPackageClient_DispatchMock(t *testing.T) {
	type fields struct {
		client http.Client
		pkg    *Package
	}
	type args struct {
		request *DispatchRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *DispatchResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PackageClient{
				client: tt.fields.client,
				pkg:    tt.fields.pkg,
			}
			got, err := p.DispatchMock(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("DispatchMock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DispatchMock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackageClient_DoJSONGet(t *testing.T) {
	type fields struct {
		client http.Client
		pkg    *Package
	}
	type args struct {
		url    string
		result interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PackageClient{
				client: tt.fields.client,
				pkg:    tt.fields.pkg,
			}
			if err := p.DoJSONGet(tt.args.url, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("DoJSONGet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPackageClient_DoJSONPost(t *testing.T) {
	type fields struct {
		client http.Client
		pkg    *Package
	}
	type args struct {
		url    string
		body   interface{}
		result interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PackageClient{
				client: tt.fields.client,
				pkg:    tt.fields.pkg,
			}
			if err := p.DoJSONPost(tt.args.url, tt.args.body, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("DoJSONPost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPackageClient_ExecuteLink(t *testing.T) {
	type fields struct {
		client http.Client
		pkg    *Package
	}
	type args struct {
		request *LinkExecutionRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *LinkExecutionResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PackageClient{
				client: tt.fields.client,
				pkg:    tt.fields.pkg,
			}
			got, err := p.ExecuteLink(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExecuteLink() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackageClient_ExecuteLinkMock(t *testing.T) {
	type fields struct {
		client http.Client
		pkg    *Package
	}
	type args struct {
		request *LinkExecutionRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *LinkExecutionResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PackageClient{
				client: tt.fields.client,
				pkg:    tt.fields.pkg,
			}
			got, err := p.ExecuteLinkMock(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteLinkMock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExecuteLinkMock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackageClient_ExecuteNode(t *testing.T) {
	type fields struct {
		client http.Client
		pkg    *Package
	}
	type args struct {
		input *NodeCall
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *NodeCallResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PackageClient{
				client: tt.fields.client,
				pkg:    tt.fields.pkg,
			}
			got, err := p.ExecuteNode(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExecuteNode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackageClient_ExecuteNodeMock(t *testing.T) {
	type fields struct {
		client http.Client
		pkg    *Package
	}
	type args struct {
		input *NodeCall
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *NodeCallResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PackageClient{
				client: tt.fields.client,
				pkg:    tt.fields.pkg,
			}
			got, err := p.ExecuteNodeMock(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteNodeMock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExecuteNodeMock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackageClient_FetchManifest(t *testing.T) {
	type fields struct {
		client http.Client
		pkg    *Package
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Package
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PackageClient{
				client: tt.fields.client,
				pkg:    tt.fields.pkg,
			}
			got, err := p.FetchManifest()
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchManifest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchManifest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackageClient_GetAsset(t *testing.T) {
	type fields struct {
		client http.Client
		pkg    *Package
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    io.Reader
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PackageClient{
				client: tt.fields.client,
				pkg:    tt.fields.pkg,
			}
			got, err := p.GetAsset(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAsset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAsset() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackageClient_GetAssetBytes(t *testing.T) {
	type fields struct {
		client http.Client
		pkg    *Package
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PackageClient{
				client: tt.fields.client,
				pkg:    tt.fields.pkg,
			}
			got, err := p.GetAssetBytes(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAssetBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAssetBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackageClient_makeRequestWithBody(t *testing.T) {
	type fields struct {
		client http.Client
		pkg    *Package
	}
	type args struct {
		method       string
		url          string
		signingToken string
		body         interface{}
		out          interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PackageClient{
				client: tt.fields.client,
				pkg:    tt.fields.pkg,
			}
			if err := p.makeRequestWithBody(tt.args.method, tt.args.url, tt.args.signingToken, tt.args.body, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("makeRequestWithBody() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getSignature(t *testing.T) {
	type args struct {
		body []byte
		key  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSignature(tt.args.body, tt.args.key); got != tt.want {
				t.Errorf("getSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
