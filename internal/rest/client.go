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

func NewClient(cfg *RESTConfig, client *http.Client) *Client {
	return &Client{
		cfg:    cfg,
		Client: client,
	}
}

func (c *Client) SendCreate(ctx context.Context, createRequest interface{}, createResponse interface{}) error {

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

	postResp, err := c.Client.Post(fmt.Sprintf("%s/%s", c.cfg.BaseUrl, c.cfg.Type.Plural), "application/json", jsonReader)
	if err != nil {
		return fmt.Errorf("Post call to %s failed: %w", c.cfg.BaseUrl, err)
	}

	jd := json.NewDecoder(postResp.Body)
	err = jd.Decode(&createResponse)
	if err != nil {
		return fmt.Errorf("JSON decoding failed: %w", err)
	}

	return nil
}

func (c *Client) SendDelete(ctx context.Context, deleteRequest fmt.Stringer, deleteResponse interface{}) error {

	escaped := url.PathEscape(deleteRequest.String())

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s/%s", c.cfg.BaseUrl, c.cfg.Type.Singular, escaped), nil)
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

	resp, err := c.Client.Get(fmt.Sprintf("%s/%s/%s", c.cfg.BaseUrl, c.cfg.Type.Singular, escaped))
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
