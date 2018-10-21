package openstack

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

var clientAuth Openstack

const mockURL = "http://mock.api"
const keystoneURI = "/identity/v3/auth/tokens"
const keystoneURL = mockURL + keystoneURI
const keystoneResponse = `{"token": {"is_domain": false, "methods": ["password"], "roles": [{"id": "ebc6c937b13044579fb58d04d777d1d0", "name": "member"}, {"id": "706726bcf1674d16af7703745ec983e1", "name": "reader"}, {"id": "b5abb3602f584ccbb30e6914d36bc491", "name": "admin"}], "expires_at": "2018-08-13T15:39:29.000000Z", "project": {"domain": {"id": "default", "name": "Default"}, "id": "31ae23a9a786499f82bc5bb18bc9ac9f", "name": "admin"}, "catalog": [{"endpoints": [{"region_id": "RegionOne", "url": "http://mock.api/volume/v3/31ae23a9a786499f82bc5bb18bc9ac9f", "region": "RegionOne", "interface": "public", "id": "5926ae88c7e140428f57c74839b4ca3b"}], "type": "volumev3", "id": "0f228f7c8c90462eb2c29539fe468044", "name": "cinderv3"}, {"endpoints": [{"region_id": "RegionOne", "url": "http://mock.api:8889/", "region": "RegionOne", "interface": "public", "id": "965fd5581d7343309c9c5c1f9dccf313"}, {"region_id": "RegionOne", "url": "http://mock.api:8889/", "region": "RegionOne", "interface": "internal", "id": "dc9dfe5158374803bdf72f9fe84f3e45"}, {"region_id": "RegionOne", "url": "http://mock.api:8889/", "region": "RegionOne", "interface": "admin", "id": "e140965e6c794ceb8db1c3b7a0454df3"}], "type": "rating", "id": "1a675f90fccd43c8ac9921ed7f992f39", "name": "cloudkitty"}, {"endpoints": [{"region_id": "RegionOne", "url": "http://mock.api/compute/v2/31ae23a9a786499f82bc5bb18bc9ac9f", "region": "RegionOne", "interface": "public", "id": "593749abdff74ee4997599f8adee61ac"}], "type": "compute_legacy", "id": "4fba4ebb476348c18eecbb1bc9ead053", "name": "nova_legacy"}, {"endpoints": [{"region_id": "RegionOne", "url": "http://mock.api:9696/", "region": "RegionOne", "interface": "public", "id": "04bbac3a95bb4f96beba20f6f7cded72"}], "type": "network", "id": "6f1f97a258dd44b7b4dafde57cce5a78", "name": "neutron"}, {"endpoints": [{"region_id": "RegionOne", "url": "http://mock.api/volume/v1/31ae23a9a786499f82bc5bb18bc9ac9f", "region": "RegionOne", "interface": "public", "id": "5c8bfe78dc384001b3448e775fc78c9c"}], "type": "volume", "id": "7468255b397d4d078525bcaabdd01cd6", "name": "cinder"}, {"endpoints": [{"region_id": "RegionOne", "url": "http://mock.api/image", "region": "RegionOne", "interface": "public", "id": "875548cbd5804309a8720d59631153e3"}], "type": "image", "id": "87e9a7096e2e42e7be44f1966c422d96", "name": "glance"}, {"endpoints": [{"region_id": "RegionOne", "url": "http://mock.api/placement", "region": "RegionOne", "interface": "public", "id": "bfe3edf02efa4488a1222a55f6969352"}], "type": "placement", "id": "8de0991fb7c14864bcaa2135d44158de", "name": "placement"}, {"endpoints": [{"region_id": "RegionOne", "url": "http://mock.api/identity", "region": "RegionOne", "interface": "public", "id": "5d06aa9562674edfa8046f6ca643c47f"}, {"region_id": "RegionOne", "url": "http://mock.api/identity", "region": "RegionOne", "interface": "admin", "id": "a1f510bd49c34f83bc609fe131d1750c"}], "type": "identity", "id": "92298c23c7b648fa834f55fdcb135345", "name": "keystone"}, {"endpoints": [{"region_id": "RegionOne", "url": "http://mock.api/metric", "region": "RegionOne", "interface": "public", "id": "20a3ebe9a4b943e4b4e86e1e91ba6e3e"}, {"region_id": "RegionOne", "url": "http://mock.api/metric", "region": "RegionOne", "interface": "internal", "id": "c2ef804a23954ab8866f5718bc1ad038"}, {"region_id": "RegionOne", "url": "http://mock.api/metric", "region": "RegionOne", "interface": "admin", "id": "d9b7cefb958144939932e6bf28811d5b"}], "type": "metric", "id": "a216aa47eefc404c86a08745f1c1db1e", "name": "gnocchi"}, {"endpoints": [{"region_id": "RegionOne", "url": "http://mock.api/compute/v2.1", "region": "RegionOne", "interface": "public", "id": "8f5cd5de48f9434b915b7ba7d06f3d3c"}], "type": "compute", "id": "a760511b36bb469482809b4230c86e63", "name": "nova"}, {"endpoints": [{"region_id": "RegionOne", "url": "http://mock.api/volume/v2/31ae23a9a786499f82bc5bb18bc9ac9f", "region": "RegionOne", "interface": "public", "id": "2966f04ea90242e1ad93f18a1aa98a1b"}], "type": "volumev2", "id": "bd97b39a81ce425e9d2b5cddd8805240", "name": "cinderv2"}, {"endpoints": [{"region_id": "RegionOne", "url": "http://mock.api/volume/v3/31ae23a9a786499f82bc5bb18bc9ac9f", "region": "RegionOne", "interface": "public", "id": "07fdfb2ea76d426c8378830cb5565ea3"}], "type": "block-storage", "id": "f7fc9240f43741d29ece002d53e182b4", "name": "cinder"}], "user": {"password_expires_at": null, "domain": {"id": "default", "name": "Default"}, "id": "6aa23f0e0e464250ab99a34946d50c17", "name": "admin"}, "audit_ids": ["UFC1o1BdSyG-Ox6YidvNHQ"], "issued_at": "2018-08-13T14:39:29.000000Z"}}`
const cinderURI = `/volume/v3/31ae23a9a786499f82bc5bb18bc9ac9f/volumes/7a66eb97-9cd0-46b7-9ecf-9be6c4b8dac3`
const cinderGetVolumeResponse = `{"volume": {"migration_status": null, "attachments": [{"server_id": "3441b857-59b0-4908-8236-fdc48aed8084", "attachment_id": "b584d421-ef6f-4ca5-bd16-ddf667138378", "attached_at": "2018-08-06T08:47:16.000000", "host_name": null, "volume_id": "7a66eb97-9cd0-46b7-9ecf-9be6c4b8dac3", "device": "/dev/vda", "id": "7a66eb97-9cd0-46b7-9ecf-9be6c4b8dac3"}], "links": [{"href": "http://1.1.1.249/volume/v3/31ae23a9a786499f82bc5bb18bc9ac9f/volumes/7a66eb97-9cd0-46b7-9ecf-9be6c4b8dac3", "rel": "self"}, {"href": "http://1.1.1.249/volume/31ae23a9a786499f82bc5bb18bc9ac9f/volumes/7a66eb97-9cd0-46b7-9ecf-9be6c4b8dac3", "rel": "bookmark"}], "availability_zone": "nova", "os-vol-host-attr:host": "mitaka-gnocchi@lvmdriver-1#lvmdriver-1", "encrypted": false, "updated_at": "2018-08-06T08:47:17.000000", "replication_status": null, "snapshot_id": null, "id": "7a66eb97-9cd0-46b7-9ecf-9be6c4b8dac3", "size": 1, "user_id": "6aa23f0e0e464250ab99a34946d50c17", "os-vol-tenant-attr:tenant_id": "31ae23a9a786499f82bc5bb18bc9ac9f", "os-vol-mig-status-attr:migstat": null, "metadata": {"attached_mode": "rw"}, "status": "in-use", "volume_image_metadata": {"checksum": "f8ab98ff5e73ebab884d80c9dc9c7290", "min_ram": "0", "disk_format": "qcow2", "image_name": "cirros-0.3.5-x86_64-disk", "image_id": "993a6a61-a144-4595-961a-3b92d774e9d6", "container_format": "bare", "min_disk": "0", "size": "13267968"}, "description": "", "multiattach": false, "source_volid": null, "consistencygroup_id": null, "os-vol-mig-status-attr:name_id": null, "name": "", "bootable": "true", "created_at": "2018-08-06T08:46:51.000000", "volume_type": "lvmdriver-1"}}`
const novaServersDetailedResponse = `{"servers": [{"OS-EXT-STS:task_state": null, "addresses": {"public": [{"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:51:86:ff", "version": 6, "addr": "2001:db8::3", "OS-EXT-IPS:type": "fixed"}, {"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:51:86:ff", "version": 4, "addr": "1.1.1.2", "OS-EXT-IPS:type": "fixed"}]}, "links": [{"href": "http://1.1.1.249/compute/v2.1/servers/3441b857-59b0-4908-8236-fdc48aed8084", "rel": "self"}, {"href": "http://1.1.1.249/compute/servers/3441b857-59b0-4908-8236-fdc48aed8084", "rel": "bookmark"}], "image": "", "OS-EXT-STS:vm_state": "active", "OS-EXT-SRV-ATTR:instance_name": "instance-00000002", "OS-SRV-USG:launched_at": "2018-08-06T08:47:35.000000", "flavor": {"id": "1", "links": [{"href": "http://1.1.1.249/compute/flavors/1", "rel": "bookmark"}]}, "id": "3441b857-59b0-4908-8236-fdc48aed8084", "security_groups": [{"name": "default"}], "user_id": "6aa23f0e0e464250ab99a34946d50c17", "OS-DCF:diskConfig": "AUTO", "accessIPv4": "", "accessIPv6": "", "progress": 0, "OS-EXT-STS:power_state": 1, "OS-EXT-AZ:availability_zone": "nova", "config_drive": "", "status": "ACTIVE", "updated": "2018-08-06T08:47:36Z", "hostId": "b0fc6c1e8f3f11ac77f8de36b344a7a542cdbc889a1b0b5aa4adc78b", "OS-EXT-SRV-ATTR:host": "mitaka-gnocchi", "OS-SRV-USG:terminated_at": null, "key_name": null, "OS-EXT-SRV-ATTR:hypervisor_hostname": "mitaka-gnocchi", "name": "Test", "created": "2018-08-06T08:46:30Z", "tenant_id": "31ae23a9a786499f82bc5bb18bc9ac9f", "os-extended-volumes:volumes_attached": [{"id": "7a66eb97-9cd0-46b7-9ecf-9be6c4b8dac3"}], "metadata": {}}, {"OS-EXT-STS:task_state": null, "addresses": {"private": [{"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:14:c0:72", "version": 4, "addr": "10.0.0.17", "OS-EXT-IPS:type": "fixed"}, {"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:14:c0:72", "version": 6, "addr": "fd9c:ea38:ad0d:0:f816:3eff:fe14:c072", "OS-EXT-IPS:type": "fixed"}]}, "links": [{"href": "http://1.1.1.249/compute/v2.1/servers/9e5ec46a-2104-4f9b-8221-4ca56bc5c687", "rel": "self"}, {"href": "http://1.1.1.249/compute/servers/9e5ec46a-2104-4f9b-8221-4ca56bc5c687", "rel": "bookmark"}], "image": "", "OS-EXT-STS:vm_state": "active","OS-EXT-SRV-ATTR:instance_name": "instance-00000001", "OS-SRV-USG:launched_at": "2018-08-02T09:53:55.000000", "flavor": {"id": "1", "links": [{"href": "http://1.1.1.249/compute/flavors/1", "rel": "bookmark"}]}, "id": "9e5ec46a-2104-4f9b-8221-4ca56bc5c687", "security_groups": [{"name": "default"}], "user_id": "6aa23f0e0e464250ab99a34946d50c17", "OS-DCF:diskConfig": "AUTO", "accessIPv4": "", "accessIPv6": "", "progress": 0, "OS-EXT-STS:power_state": 1, "OS-EXT-AZ:availability_zone": "nova", "config_drive": "", "status": "ACTIVE", "updated": "2018-08-08T14:04:43Z", "hostId": "8f0341c333840841c8fc45d7d25b683dff9e4a96a4e86cc3000659c9", "OS-EXT-SRV-ATTR:host": "mitaka-gnocchi", "OS-SRV-USG:terminated_at": null, "key_name": null, "OS-EXT-SRV-ATTR:hypervisor_hostname": "mitaka-gnocchi", "name": "test", "created": "2018-08-02T09:52:53Z", "tenant_id": "d2690bd20b7d4bc8b2085bbc585d83fb", "os-extended-volumes:volumes_attached": [{"id": "a0887dce-9b8c-41b3-b631-2e8e8044eb79"}], "metadata": {}}]}`

func TestOpenStackNewClient(t *testing.T) {
	defer gock.Off()

	var err error

	authOptions := AuthOptions{[]string{"password"}, "admin", "default", "secret", keystoneURL}
	gock.New(mockURL).
		Post(keystoneURI).
		Reply(200).
		JSON(keystoneResponse)

	clientAuth = NewClient(authOptions)
	clientAuth.Client().MaxRetries(0)
	err = clientAuth.Authenticate()
	assert.Nil(t, err)

	gock.New(mockURL).
		Post(keystoneURI).
		Reply(500).
		JSON(keystoneResponse)
	failClient := NewClient(authOptions)
	failClient.Client().MaxRetries(0)
	err = failClient.Authenticate()
	assert.NotNil(t, err)

	// Verify that we don't have pending mocks
	assert.Equal(t, gock.IsDone(), true)
}

func TestCinderGetVolume(t *testing.T) {
	defer gock.Off()
	gock.New(mockURL).
		Get(cinderURI).
		Reply(200).
		JSON(cinderGetVolumeResponse)

	cinderClient := clientAuth.Cinder()
	_, err := cinderClient.GetVolume("7a66eb97-9cd0-46b7-9ecf-9be6c4b8dac3")
	// log.Debugln(err)

	assert.Nil(t, err)

	gock.New(mockURL).
		Get(cinderURI).
		Reply(500).
		JSON(cinderGetVolumeResponse)

	_, err = cinderClient.GetVolume("7a66eb97-9cd0-46b7-9ecf-9be6c4b8dac3")
	assert.NotNil(t, err)

	assert.Equal(t, gock.IsDone(), true)
}
