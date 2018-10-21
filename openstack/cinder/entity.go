package cinder

// Volume cinder volume response type
type Volume struct {
	Volume volume `json:"volume"`
}
type volume struct {
	Status              string              `json:"status"`
	VolumeImageMetadata volumeImageMetadata `json:"volume_image_metadata"`
}

type volumeImageMetadata struct {
	Checksum          string `json:"checksum"`
	MinRAM            string `json:"min_ram"`
	DiskFormat        string `json:"disk_format"`
	ImageName         string `json:"image_name"`
	ImageID           string `json:"image_id"`
	SignatureVerified string `json:"signature_verified"`
	ContainerFormat   string `json:"container_format"`
	MinDisk           string `json:"min_disk"`
	Size              string `json:"size"`
}
