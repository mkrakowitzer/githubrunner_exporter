package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	ApiURL = "https://api.github.com/"
)

var (
	etags = map[string]Etag{}
)

type Etag struct {
	Key  string
	Date string
	Body []byte
}

type Client struct {
	http *http.Client
}

// ClientOption represents an argument to NewClient
type ClientOption = func(http.RoundTripper) http.RoundTripper

// NewClient initializes a Client
func NewClient(opts ...ClientOption) *Client {
	tr := http.DefaultTransport
	for _, opt := range opts {
		tr = opt(tr)
	}

	http := &http.Client{Transport: tr}
	client := &Client{http: http}
	return client
}

// Build a request with Etag data.
// When using Etags request to the github API do not count against the rate limit quota
// See https://docs.github.com/en/rest/overview/resources-in-the-rest-api#rate-limiting
// and https://docs.github.com/en/rest/reference/rate-limit
func buildRequest(url string, method string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if _, present := etags[url]; present {
		req.Header.Add("If-None-Match", etags[url].Key)
		req.Header.Add("If-Modified-Since", etags[url].Date)
	}
	return req, err
}

// Update Etag cache.
func updateEtagCache(url string, resp *http.Response) ([]byte, error) {
	var b []byte
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	etags[url] = Etag{
		resp.Header.Get("ETag"),
		resp.Header.Get("Date"),
		b,
	}
	return b, nil
}

// REST performs a REST request and parses the response.
func (c Client) REST(method string, p string, body io.Reader, data interface{}) error {

	var b []byte

	url := "https://api.github.com/" + p

	req, err := buildRequest(url, method, body)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !success && resp.StatusCode != 304 {
		return handleHTTPError(resp)
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	if resp.StatusCode != 304 {
		b, _ = updateEtagCache(url, resp)
	} else {
		b = etags[url].Body
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	return nil
}

// handleHTTPError handles http errors
func handleHTTPError(resp *http.Response) error {
	var message string
	var parsedBody struct {
		Message string `json:"message"`
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		message = string(body)
	} else {
		message = parsedBody.Message
	}

	return fmt.Errorf("http error, '%s' failed (%d): '%s'", resp.Request.URL, resp.StatusCode, message)
}

// AddHeader turns a RoundTripper into one that adds a request header
func AddHeader(name, value string) ClientOption {
	return func(tr http.RoundTripper) http.RoundTripper {
		return &funcTripper{roundTrip: func(req *http.Request) (*http.Response, error) {
			req.Header.Add(name, value)
			return tr.RoundTrip(req)
		}}
	}
}

// AddHeaderFunc is an AddHeader that gets the string value from a function
func AddHeaderFunc(name string, value func() string) ClientOption {
	return func(tr http.RoundTripper) http.RoundTripper {
		return &funcTripper{roundTrip: func(req *http.Request) (*http.Response, error) {
			req.Header.Add(name, value())
			return tr.RoundTrip(req)
		}}
	}
}

type funcTripper struct {
	roundTrip func(*http.Request) (*http.Response, error)
}

func (tr funcTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return tr.roundTrip(req)
}
