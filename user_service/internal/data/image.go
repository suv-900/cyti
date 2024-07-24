package data

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Image struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Size     int64
	Location string
	UserID   uint64
	User     User `gorm:"constraint:OnDelete:CASCADE;"`
}
type ImageModel struct {
	DB *gorm.DB
}

func (i ImageModel) GetProfilePicture(userid uint64) (Image, error) {
	var image Image

	ctx, cancel := context.WithTimeout(context.Background(), Context_timeout)
	defer cancel()

	t := i.DB.WithContext(ctx).Where("user_id = ?", userid).Find(&image)
	return image, t.Error
}

func (i ImageModel) UpdateProfilePicture(image *Image) error {
	ctx, cancel := context.WithTimeout(context.Background(), Context_timeout)
	defer cancel()

	t := i.DB.WithContext(ctx).Save(image)
	return t.Error
}
func (i ImageModel) RemoveProfilePicture(image *Image) error {
	ctx, cancel := context.WithTimeout(context.Background(), Context_timeout)
	defer cancel()

	t := i.DB.WithContext(ctx).Delete(image)
	return t.Error
}
