package data

import "gorm.io/gorm"

type Models struct {
	Users interface {
		AddUser(user *User) error
		GetUser(userid uint64) (User, error)
		UpdateUser(userid uint64, updates map[string]interface{}) error
		DeleteUser(user *User) error

		GetUserPassword(username string) (string, error)
		CheckUserExists(username string) (bool, error)
		UpdatePassword(userid uint64, password string) error

		GetLoginAttempts(username string) (Login, error)
		ResetLoginAttempts(username string) error
		UpdateLoginAttempts(username string) error

		FindSoftDeletedRecords() ([]User, error)
	}

	Images interface {
		UpdateProfilePicture(image *Image) error
		RemoveProfilePicture(image *Image) error
		GetProfilePicture(userid uint64) (Image, error)
	}
}

func GetModels(db *gorm.DB) Models {
	return Models{
		Users:  UserModel{DB: db},
		Images: ImageModel{DB: db},
	}
}
