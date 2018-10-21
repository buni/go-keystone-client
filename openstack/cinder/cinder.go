package cinder

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/Buni/openstack-client/openstack/client"
	log "github.com/sirupsen/logrus"
)

// Cinder Client
type cinder struct {
	Client client.Client
}

// Cinder interface
type Cinder interface {
	GetVolume(volumeID string) (volume Volume, err error)
}

const volumePath = "/volumes/$id"

// New Cinder
func New(authClient client.Client) Cinder {
	return &cinder{Client: authClient}
}

func (c *cinder) GetVolume(volumeID string) (volume Volume, err error) {
	path := c.Client.GetEndpoint("cinderv3") + volumePath + "?limit=1&all_tenants=1"
	path = strings.Replace(path, "$id", volumeID, -1)

	resp, err := c.Client.DoAuthRequest(context.TODO(), "GET", path, nil)
	if err != nil {
		return
	}
	// c.Client.NewRequest(path, "GET", nil).Context(context.TODO()).Transport(http.DefaultTransport).Do()
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Debugln(string(respBody))

	err = json.Unmarshal(respBody, &volume)
	if err != nil {
		return
	}

	return
}
