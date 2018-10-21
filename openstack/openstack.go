package openstack

import (
	"github.com/Buni/openstack-client/openstack/cinder"
	"github.com/Buni/openstack-client/openstack/client"
	"github.com/Buni/openstack-client/openstack/keystone"
)

// AuthOptions fields
type AuthOptions struct {
	Methods    []string
	TenantName string
	Domain     string
	Password   string
	Endpoint   string
}

type openstack struct {
	client   client.Client
	keystone keystone.Keystone
}

// Openstack API
type Openstack interface {
	Client() client.Client
	Authenticate() error
	Keystone() keystone.Keystone
	Cinder() cinder.Cinder
}

// NewClient instance of Authenticated client
func NewClient(ao AuthOptions) Openstack {
	keystn := keystone.New(ao.Methods, ao.TenantName, ao.Domain, ao.Password, ao.Endpoint)
	return &openstack{client: keystn.GetClient(), keystone: keystn}
}
func (o *openstack) Authenticate() error {
	return o.keystone.Authenticate()
}

// Cinder interface exposes all cinder methods
func (o *openstack) Cinder() cinder.Cinder {
	return cinder.New(o.client)
}

func (o *openstack) Keystone() keystone.Keystone {
	return o.keystone
}

func (o *openstack) Client() client.Client {
	return o.client
}
