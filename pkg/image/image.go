package image

type Image struct {
	ID           uint64 `json:"id,omitempty"`
	Image        []byte `json:"image,omitempty"`
	ImageType    string `json:"image_type,omitempty"`
	EncodedImage string `gorm:"-"`
}
