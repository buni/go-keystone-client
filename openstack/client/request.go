package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/eapache/go-resiliency/retrier"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

const authHeader = "X-Auth-Token"
const maxRetries = 3
const timeoutBetweenRetires = time.Millisecond * 100
const timeout = time.Second * 5

// MetaData per request MetaData
type metaData struct {
	ComponentName string
	OperationName string
}

// Options per request options
type options struct {
	MaxRetries            int
	TimeoutBetweenRetries time.Duration
	ReqTimeout            time.Duration
}

// request private
type request struct {
	ctx        context.Context
	method     string
	url        string
	query      string
	body       io.Reader
	authHeader string
}

// Request public type
type Request struct {
	reqClient     *client
	clientOptions options
	reqMetaData   metaData
	reqOptions    request
	mux           *sync.Mutex // just in case
}

// NewRequest prepare new request and copy client configuration

func (c *client) NewRequest(url, method string, body io.Reader) *Request {
	clientCopy := &client{}
	*clientCopy = *c                      // this is a shallow copy(surprisingly)
	clientCopy.HTTPClient = &http.Client{ // so we need a new http client
		Timeout:   timeout,
		Transport: &nethttp.Transport{}, // TODO: replicate keystone tls issue
	}
	clientCopy.mux = new(sync.Mutex) // and a new mutex

	log.Debugln(&clientCopy.mux, &c.mux, "mux")
	log.Debugln(&clientCopy.HTTPClient, &c.HTTPClient, "http client struct")
	log.Debugln(&clientCopy.HTTPClient.Transport, &c.HTTPClient.Transport, "transport ")
	return &Request{reqClient: clientCopy, clientOptions: options{MaxRetries: maxRetries, TimeoutBetweenRetries: timeoutBetweenRetires, ReqTimeout: timeout}, reqOptions: request{url: url, method: method, body: body, ctx: context.TODO(), authHeader: authHeader}, mux: new(sync.Mutex)}
}

// Context add context.Context to request for tracing purposes
func (r *Request) Context(ctx context.Context) *Request {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.reqOptions.ctx = ctx
	return r
}

// AuthHeader set authentication header key
func (r *Request) AuthHeader(name string) *Request {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.reqOptions.authHeader = name
	return r
}

// QueryKV add query key value pair
func (r *Request) QueryKV(key, value string) *Request {
	r.mux.Lock()
	defer r.mux.Unlock()
	if r.reqOptions.query == "" {
		r.reqOptions.query = "?" + key + "=" + value
		return r
	}
	r.reqOptions.query = "&" + key + "=" + value
	return r
}

// QueryString add entire query string
func (r *Request) QueryString(query string) *Request {
	r.mux.Lock()
	defer r.mux.Unlock()
	if len(query) <= 1 {
		return r
	}
	if r.reqOptions.query == "" {
		r.reqOptions.query = "?" + strings.Replace(query, "?", "", -1)
		return r
	}
	r.reqOptions.query += strings.Replace(query, "?", "&", -1)

	return r
}

// MetaData add metadata for tracing
func (r *Request) MetaData(componentName, operationName string) *Request {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.reqMetaData.ComponentName = componentName
	r.reqMetaData.OperationName = operationName
	return r
}

// MaxRetries change max retries
func (r *Request) MaxRetries(retries int) *Request {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.clientOptions.MaxRetries = retries
	return r
}

// TimeBR time between retries
func (r *Request) TimeBR(tbr time.Duration) *Request {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.clientOptions.TimeoutBetweenRetries = tbr
	return r
}

// ReqTimeout lets you set request timeout
func (r *Request) ReqTimeout(reqTimeout time.Duration) *Request {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.clientOptions.ReqTimeout = reqTimeout // cleanup this one
	r.reqClient.Timeout(timeout)
	return r
}

// Transport lets you set a custom transport
// Tracing will work with any transport but the spans will be less detailed
func (r *Request) Transport(transport http.RoundTripper) *Request {
	r.mux.Lock()
	defer r.mux.Unlock()
	log.Debugln(r.reqClient.HTTPClient.Transport, "old")
	r.reqClient.Transport(transport)
	log.Debugln(r.reqClient.HTTPClient.Transport, "new")

	return r
}

// Do execute request
func (r *Request) Do() (resp *http.Response, err error) {

	rq, err := http.NewRequest(r.reqOptions.method, r.reqOptions.url+r.reqOptions.query, r.reqOptions.body)
	if err != nil {
		return
	}

	rq.Header.Set("Content-Type", "application/json") // Fix later

	log.Debugln(rq)
	retry := retrier.New(retrier.ConstantBackoff(maxRetries, timeoutBetweenRetires), nil) // TODO: Setup a Whitelist Classifier
	rq.Header.Set(r.reqOptions.authHeader, r.reqClient.GetToken())
	rtr := 0

	err = retry.Run(func() (err error) {
		req := rq
		req = req.WithContext(r.reqOptions.ctx)
		req, ht := nethttp.TraceRequest(opentracing.GlobalTracer(), req, nethttp.ComponentName(r.reqMetaData.ComponentName), nethttp.OperationName(r.reqMetaData.OperationName))

		resp, err = r.reqClient.HTTPClient.Do(req)
		defer ht.Finish()
		rtr++
		return r.reqClient.verifyAuth(req, resp, err)
	})

	log.Debugf("attempts %v", rtr)
	return
}

// DoNonAuth do normal request
func (r *Request) DoNonAuth() (resp *http.Response, err error) {
	rq, err := http.NewRequest(r.reqOptions.method, r.reqOptions.url+r.reqOptions.query, r.reqOptions.body)
	if err != nil {
		return
	}
	retry := retrier.New(retrier.ConstantBackoff(maxRetries, timeoutBetweenRetires), nil) // TODO: Setup a Whitelist Classifier

	rtr := 0
	err = retry.Run(func() (err error) {
		req := rq
		req = req.WithContext(context.TODO())
		req, ht := nethttp.TraceRequest(opentracing.GlobalTracer(), req, nethttp.ComponentName(r.reqMetaData.ComponentName), nethttp.OperationName(r.reqMetaData.OperationName))

		resp, err = r.reqClient.HTTPClient.Do(req)
		defer ht.Finish()
		rtr++
		return r.reqClient.verifyNoAuth(req, resp, err)
	})
	log.Debugf("attempts %v", rtr)
	return
}

// func (r *Request) Query() *Request {

// 	return r
// }
