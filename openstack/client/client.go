package client

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/eapache/go-resiliency/retrier"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

type client struct {
	HTTPClient *http.Client
	Keystone   Keystone
	maxRetries int
	mux        *sync.Mutex
}

// Keystone import cycle prevention
type Keystone interface {
	Authenticate() (err error)
	ReAuthenticate() (err error)
	GetToken() string
	GetEndpoint(name string) string
}

// Client interface
type Client interface {
	DoAuthRequest(ctx context.Context, method, url string, body io.Reader) (resp *http.Response, err error)
	DoRequest(ctx context.Context, method, url string, body io.Reader) (resp *http.Response, err error)
	Transport(transport http.RoundTripper)
	Timeout(timeout time.Duration)
	MaxRetries(maxRetries int)
	NewRequest(url, method string, body io.Reader) *Request
	Keystone
}

// New auth client
func New(k Keystone) Client {
	// http.Transport{}
	return &client{HTTPClient: &http.Client{
		Timeout:   timeout,
		Transport: &nethttp.Transport{}, // TODO: replicate keystone tls issue with this transport
	},
		Keystone:   k,
		maxRetries: maxRetries,
		mux:        new(sync.Mutex)}
}

// DoAuthRequest prepare and do Request with retry
func (c *client) DoAuthRequest(ctx context.Context, method, url string, body io.Reader) (resp *http.Response, err error) {
	rq, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	r := retrier.New(retrier.ConstantBackoff(c.maxRetries, timeoutBetweenRetires), nil) // TODO: Setup a Whitelist Classifier
	rq.Header.Set(authHeader, c.Keystone.GetToken())
	rtr := 0

	err = r.Run(func() (err error) {
		req := rq // copy the original request to split the retry spans
		req = req.WithContext(ctx)
		req, ht := nethttp.TraceRequest(opentracing.GlobalTracer(), req, nethttp.ComponentName(url), nethttp.OperationName(method))

		resp, err = c.HTTPClient.Do(req)
		defer ht.Finish()
		rtr++
		return c.verifyAuth(req, resp, err)
	})

	log.Debugf("attempts %v", rtr)

	return
}

// DoRequest do normal request
func (c *client) DoRequest(ctx context.Context, method, url string, body io.Reader) (resp *http.Response, err error) {
	rq, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	r := retrier.New(retrier.ConstantBackoff(c.maxRetries, timeoutBetweenRetires), nil) // TODO: Setup a Whitelist Classifier
	rtr := 0

	err = r.Run(func() (err error) {
		req := rq // copy the original request to split the retry spans
		req = req.WithContext(ctx)
		req, ht := nethttp.TraceRequest(opentracing.GlobalTracer(), req, nethttp.ComponentName(url), nethttp.OperationName(method))

		resp, err = c.HTTPClient.Do(req)
		defer ht.Finish()
		rtr++
		return c.verifyNoAuth(req, resp, err)
	})

	log.Debugf("attempts %v", rtr)
	return
}

func (c *client) verifyAuth(req *http.Request, resp *http.Response, err error) error {
	switch {
	case err != nil:
		return err
	case resp.StatusCode == 401 || resp.StatusCode == 403:
		resp.Body.Close()
		err = c.ReAuthenticate()
		if err != nil {
			log.Errorln("Another Goroutine is refreshing the Token")
			return fmt.Errorf("%s", "Another Goroutine is refreshing the Token")
		}
		req.Header.Set(authHeader, c.GetToken())
		return fmt.Errorf("Code %v  %s", resp.StatusCode, "Access Denied token Expired")
	case resp.StatusCode > 299:
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		return fmt.Errorf("Code %v  %s", resp.StatusCode, string(respBody))
	case resp.StatusCode < 299:
		log.Debugln(resp.StatusCode)
	}
	return err
}

func (c *client) verifyNoAuth(req *http.Request, resp *http.Response, err error) error {
	switch {
	case err != nil:
		return err
	case resp.StatusCode > 299:
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		return fmt.Errorf("Code %v  %s", resp.StatusCode, string(respBody))
	case resp.StatusCode < 299:
		log.Debugln(resp.StatusCode)
	}
	return err
}

func (c *client) Transport(transport http.RoundTripper) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.HTTPClient.Transport = transport
}

func (c *client) Timeout(timeout time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.HTTPClient.Timeout = timeout
}
func (c *client) MaxRetries(maxRetries int) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.maxRetries = maxRetries
}

// Authenticate Authenticate
func (c *client) Authenticate() (err error) {
	return c.Keystone.Authenticate()
}

// GetToken get current ks token
func (c *client) GetToken() string {
	return c.Keystone.GetToken()
}

// GetEndpoint returns public openstack endpoins by name
func (c *client) GetEndpoint(name string) string {
	return c.Keystone.GetEndpoint(name)
}

// ReAuthenticate same as Authenticate but implements a timeout incase multiple request fail at once and require new token
func (c *client) ReAuthenticate() (err error) {
	return c.Keystone.ReAuthenticate()
}
