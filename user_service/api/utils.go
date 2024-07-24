package api

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// var JWTKEY = common.Config.JWTkey
var jwtKey = "JWTKEY"
var tokenExpiryTime = time.Now().Add(60 * time.Minute)
var ErrTokenInvalid = errors.New("token invalid")

type CustomPayload struct {
	ID uint64 `json:"id"`
	jwt.StandardClaims
}

func (app *application) generateToken(userid uint64) (string, error) {
	payload := CustomPayload{
		ID: userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiryTime.Unix(),
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := rawToken.SignedString(jwtKey)
	return token, err
}

func (app *application) verifyToken(token string) (uint64, error) {
	var userid uint64
	t, err := jwt.ParseWithClaims(token, &CustomPayload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		//this probably means token is expired
		return userid, err
	}
	if p, ok := t.Claims.(*CustomPayload); ok && t.Valid {
		//this probably means token is invalid
		userid = p.ID
		return userid, nil
	}
	return userid, ErrTokenInvalid
}
func (app *application) generateHashedPassword(password []byte) (string, error) {
	var hashedpass string
	hashedbytes, err := bcrypt.GenerateFromPassword(password, 3)
	if err != nil {
		return hashedpass, err
	}
	return string(hashedbytes), nil
}

// nil on success and err on fail
func (app *application) comparePassword(password []byte, dbPassword []byte) error {
	return bcrypt.CompareHashAndPassword(dbPassword, password)
}
