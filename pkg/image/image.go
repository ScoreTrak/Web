package image

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Image struct {
	ID           uuid.UUID `json:"id,omitempty"`
	Image        []byte    `json:"image,omitempty"`
	ImageType    string    `json:"image_type,omitempty"`
	EncodedImage string    `gorm:"-"`
}

func (p *Image) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		u, err := uuid.NewV4()
		if err != nil {
			return err
		}
		p.ID = u
	}
	return nil
}
