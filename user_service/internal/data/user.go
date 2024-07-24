package data

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	//	_ "github.com/lib/pq"
)

// root:Core@123@/blogweb?
// postgres://core:12345678@localhost:5432/cloud

type User struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Username string `gorm:"<-:false"`
	Email    string
	Password string

	Bio       string
	BirthDate time.Time

	IsDel soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
}

type Login struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	FailedLoginAttempts uint
	FailedLoginTime     time.Time

	UserID uint64
	User   User `gorm:"constraint:OnDelete:CASCADE;"`
}

type UserModel struct {
	DB *gorm.DB
}

var (
	ErrRecordNotFound = errors.New("user not found")
	ErrConflict       = errors.New("user exists")
	ErrUnknown        = errors.New("unknown error occured")
)

const Context_timeout = 5 * time.Second

var AnonymousUser = &User{}

func (u *User) IsAnonymousUser() bool {
	return *u == *AnonymousUser
}

func (u UserModel) AddUser(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), Context_timeout)
	defer cancel()

	err := u.DB.WithContext(ctx).Create(user).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

func (u UserModel) GetUser(userid uint64) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Context_timeout)
	defer cancel()

	var user User

	err := u.DB.WithContext(ctx).First(&user, userid).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return user, ErrRecordNotFound
		default:
			return user, err
		}
	}

	return user, nil
}

func (u UserModel) UpdatePassword(userid uint64, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), Context_timeout)
	defer cancel()

	err := u.DB.WithContext(ctx).Raw("UPDATE users SET password = ? WHERE user_id = ?", password, userid).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

func (u UserModel) GetUserPassword(username string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Context_timeout)
	defer cancel()

	var db_password string

	err := u.DB.WithContext(ctx).Raw("SELECT password FROM users WHERE username = ?", username).Scan(&db_password).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return db_password, ErrRecordNotFound
		default:
			return db_password, err
		}
	}

	return db_password, nil
}

func (u UserModel) UpdateUser(userid uint64, updates map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), Context_timeout)
	defer cancel()

	err := u.DB.WithContext(ctx).Model(&User{ID: userid}).Omit("password").Updates(updates).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (u UserModel) CheckUserExists(username string) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), Context_timeout)
	defer cancel()

	r := u.DB.WithContext(ctx).Where(&User{Username: username})

	return r.RowsAffected != 0, r.Error
}

func (u UserModel) DeleteUser(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), Context_timeout)
	defer cancel()

	err := u.DB.WithContext(ctx).Delete(user).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}
func (u UserModel) FindSoftDeletedRecords() ([]User, error) {
	var users []User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := u.DB.WithContext(ctx).Where("is_del = 1").Find(&users)
	return users, t.Error
}

func (u UserModel) GetLoginAttempts(username string) (Login, error) {
	var result Login

	ctx, cancel := context.WithTimeout(context.Background(), Context_timeout)
	defer cancel()

	err := u.DB.WithContext(ctx).Raw(`SELECT 
	failed_login_attempts,failed_login_time 
	FROM users WHERE username = ?`, username).Scan(&result).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return result, ErrRecordNotFound
		default:
			return result, err
		}
	}

	return result, nil
}

func (u UserModel) UpdateLoginAttempts(username string) error {

	ctx, cancel := context.WithTimeout(context.Background(), Context_timeout)
	defer cancel()

	err := u.DB.WithContext(ctx).Raw(`UPDATE users SET
	failed_login_attempts = failed_login_attempts + 1,
	failed_login_time = ? WHERE username = ?`, time.Now(), username).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil

}
func (u UserModel) ResetLoginAttempts(username string) error {

	ctx, cancel := context.WithTimeout(context.Background(), Context_timeout)
	defer cancel()

	err := u.DB.WithContext(ctx).Raw(`UPDATE users SET
	failed_login_attempts = 0 WHERE username = ?`, username).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil

}
