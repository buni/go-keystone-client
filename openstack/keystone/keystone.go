package keystone

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"
	"time"

	"github.com/Buni/openstack-client/openstack/client"
	log "github.com/sirupsen/logrus"
)

// Keystone interface
type Keystone interface {
	Authenticate() (err error)
	ReAuthenticate() (err error)
	GetToken() string
	GetEndpoint(name string) string
	GetClient() client.Client
}

// keystone Auth
type keystone struct {
	Auth      authScoped        `json:"auth"`
	Token     token             `json:"-"`
	Endpoint  string            `json:"-"`
	Endpoints map[string]string `json:"-"`
	client    client.Client
	updated   time.Time
	mux       *sync.Mutex
}

// New instance of keystone
func New(methods []string, name, dom, pass, endpoint string) Keystone {
	auth := authScoped{Identity: identity{Methods: methods, Password: password{user{Name: name, Domain: domain{ID: dom}, Password: pass}}}, Scope: scope{Project: project{Name: name, Domain: domain{ID: dom}}}}
	return newClient(&keystone{Auth: auth, Token: token{}, Endpoint: endpoint, Endpoints: make(map[string]string), mux: new(sync.Mutex)})
}

func newClient(k *keystone) *keystone {
	k.client = client.New(k)
	return k
}

func (k *keystone) GetClient() client.Client {
	return k.client
}

// Authenticate Authenticate
func (k *keystone) Authenticate() (err error) {
	k.mux.Lock()

	defer k.mux.Unlock()
	var jsonResponse resp

	payload, err := json.Marshal(k)
	if err != nil {
		return
	}

	resp, err := k.client.DoRequest(context.TODO(), "POST", k.Endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	log.Debugln(string(respBody))

	err = json.Unmarshal(respBody, &jsonResponse)
	if err != nil {
		return
	}
	// c.parseEndpoints(jsonResponse, "public")
	// c.parseEndpoints(jsonResponse, "internal")
	k.Token.Value = resp.Header.Get("X-Subject-Token")
	k.Token.ExperiesAt = jsonResponse.Token.ExpiresAt
	k.Token.ProjectID = jsonResponse.Token.Project.ID

	for _, service := range jsonResponse.Token.Catalog {
		for _, endpoint := range service.Endpoints {
			if endpoint.Interface == "public" {
				k.Endpoints[service.Name] = endpoint.URL
			}
		}
	}

	return
}

// ReAuthenticate Authenticate
func (k *keystone) ReAuthenticate() (err error) {
	k.mux.Lock()

	defer k.mux.Unlock()

	if time.Since(k.updated) < time.Minute {
		return errors.New("token was updated less than a minute ago")
	}

	err = k.Authenticate()
	if err != nil {
		return
	}

	k.updated = time.Now()
	return
}

// GetToken returns the currently set keystone token
func (k *keystone) GetToken() string {
	k.mux.Lock()
	defer k.mux.Unlock()
	return k.Token.Value
}

// GetEndpoint returns openstack public endpoints
func (k *keystone) GetEndpoint(name string) string {
	k.mux.Lock()
	defer k.mux.Unlock()
	return k.Endpoints[name]
}
