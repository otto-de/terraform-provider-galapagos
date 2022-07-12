package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Client encodes semantic logic from Galapagos REST
// Client also takes inspiration from k8s's RESTClient.
type Client struct {
	cfg    *RESTConfig
	Client *http.Client
}

type CreateOptions struct {
	IsIdempotent bool
	Environment  string
}

func NewClient(cfg *RESTConfig, client *http.Client) *Client {
	return &Client{
		cfg:    cfg,
		Client: client,
	}
}

func (c *Client) SendCreate(ctx context.Context, opts CreateOptions, createRequest interface{}, createResponse interface{}) error {

	jsonReader, jsonWriter := io.Pipe()
	defer jsonReader.Close()

	go func() {

		defer jsonWriter.Close()
		je := json.NewEncoder(jsonWriter)
		err := je.Encode(createRequest)
		if err != nil {
			tflog.Error(ctx, "Async JSON encoding failed", map[string]interface{}{
				"err": err,
			})
		}
	}()

	var method string
	if opts.IsIdempotent {
		method = http.MethodPut

	} else {
		method = http.MethodPost
	}

	var url string
	if opts.Environment == "" {
		url = fmt.Sprintf("%s/api/%s", c.cfg.BaseUrl, c.cfg.Type.Plural)
	} else {
		url = fmt.Sprintf("%s/api/%s/%s", c.cfg.BaseUrl, opts.Environment, c.cfg.Type.Plural)
	}

	httpReq, err := http.NewRequest(method, url, jsonReader)
	if err != nil {
		return fmt.Errorf("Building Request to %s failed: %w", url, err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := c.Client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("Put call to %s failed: %w", url, err)
	}

	jd := json.NewDecoder(httpResp.Body)
	err = jd.Decode(&createResponse)
	if err != nil {
		return fmt.Errorf("JSON decoding failed: %w", err)
	}

	return nil
}

func (c *Client) SendDelete(ctx context.Context, deleteRequest fmt.Stringer, deleteResponse interface{}) error {

	escaped := url.PathEscape(deleteRequest.String())

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/%s/%s", c.cfg.BaseUrl, c.cfg.Type.Singular, escaped), nil)
	if err != nil {
		return fmt.Errorf("Building request to %s failed: %w", c.cfg.BaseUrl, err)
	}
	_, err = c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("Delete call to %s failed: %w", c.cfg.BaseUrl, err)
	}

	return nil
}

func (c *Client) SendDescribe(ctx context.Context, describeRequest fmt.Stringer, describeResponse interface{}) error {
	escaped := url.PathEscape(describeRequest.String())

	resp, err := c.Client.Get(fmt.Sprintf("%s/api/%s/%s", c.cfg.BaseUrl, c.cfg.Type.Singular, escaped))
	if err != nil {
		return fmt.Errorf("Describe call to %s failed: %w", c.cfg.BaseUrl, err)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return ErrNotExist
	case http.StatusOK:
		break
	default:
		panic(fmt.Sprintf("Status code %d unexpected", resp.StatusCode))
	}

	jd := json.NewDecoder(resp.Body)
	err = jd.Decode(describeResponse)
	if err != nil {
		return fmt.Errorf("JSON decoding failed: %w", err)
	}

	return nil
}
