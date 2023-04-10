package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	logging "github.com/ipfs/go-log/v2"

	"fil-kms/app/global/http_response"
	"fil-kms/app/global/variables"
)

var log = logging.Logger("client")

func Invoke(method string, params []byte) ([]byte, error) {
	clientOpt := NewClientOpt(fmt.Sprintf("%s:%v", "http://127.0.0.1", variables.ServerPort))
	return clientOpt.Invoke(method, params)
}

type clientOpt struct {
	host string
}

func NewClientOpt(host string) *clientOpt {
	return &clientOpt{host: host}
}

const localInvokeUrl = "/local"

type InvokeParams struct {
	Method string `json:"method"`
	Params []byte `json:"params"`
}

func (fo *clientOpt) Invoke(method string, params []byte) ([]byte, error) {
	invokeUrl, err := url.Parse(fo.host)
	if err != nil {
		log.Errorf("parse url failed url: %v,err: %v", fo.host, err)
		return nil, err
	}
	invokeUrl, err = invokeUrl.Parse(localInvokeUrl)
	if err != nil {
		log.Errorf("parse url failed url: %v,err: %v", invokeUrl, err)
		return nil, err
	}
	if params == nil {
		params = []byte{}
	}
	invokeParams, err := json.Marshal(InvokeParams{
		method,
		params,
	})
	if err != nil {
		log.Errorf("parse params failed err: %v", err)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, invokeUrl.String(), bytes.NewReader(invokeParams))
	if err != nil {
		log.Errorf("NewRequest failed err: %v", err)
		return nil, err
	}
	req.Close = true
	req.Header = http.Header{}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	cli := http.Client{}
	cli.Timeout = time.Hour

	resp, err := cli.Do(req)
	if err != nil {
		log.Errorf("request http failed err: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		if resp.StatusCode < 500 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("fetch host %v:%v", resp.Status, resp.StatusCode)
			}
			return nil, fmt.Errorf(string(body))
		}
		return nil, fmt.Errorf("fetch host %v:%v", resp.Status, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read body failed err: %v", err)
		return nil, err
	}

	responseInfo := http_response.ResponseInfo{}
	err = json.Unmarshal(body, &responseInfo)
	if err != nil {
		log.Errorf("unmarshal response failed err: %v", err)
		return nil, err
	}
	if responseInfo.Data == nil {
		return nil, nil
	}
	var data []byte
	err = json.Unmarshal(responseInfo.Data, &data)
	if err != nil {
		log.Errorf("unmarshal response Data failed err: %v", err)
		return nil, err
	}
	return data, nil
}
