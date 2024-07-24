package api

import (
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"user_service/internal/data"
)

var imagesDir = "/home/core/go/kl/user_service/store/images/"

func (app *application) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user data.User

	err := app.readJSON(r, w, user)
	if err != nil {
		app.sendErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	user.Password, err = app.generateHashedPassword([]byte(user.Password))
	if err != nil {
		app.internalServerError(w, r)
		return
	}

	err = app.models.Users.AddUser(&user)
	if err != nil {
		app.internalServerError(w, r)
		return
	}

	token, err := app.generateToken(user.ID)
	if err != nil {
		app.internalServerError(w, r)
		return
	}

	w.Header().Add("Authentication-Token", token)
	w.WriteHeader(http.StatusOK)
}

func (app *application) LoginUser(w http.ResponseWriter, r *http.Request) {

	var userLogin struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := app.readJSON(r, w, userLogin)
	if err != nil {
		app.sendErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	//check login attempts
	loginAttempts, err := app.models.Users.GetLoginAttempts(userLogin.Username)
	if err != nil {
		app.internalServerError(w, r)
		return
	}

	//avoid hardcoding
	if loginAttempts.FailedLoginAttempts > 5 {
		timeNow := time.Now()
		timeDiff := loginAttempts.FailedLoginTime.Sub(timeNow)

		if timeDiff.Abs().Hours() < 2 {
			//try again after this hours
			return
		} else {
			if err := app.models.Users.ResetLoginAttempts(userLogin.Username); err != nil {
				log.Error("error while reseting login attempts ", err)
				app.internalServerError(w, r)
				return
			}
		}
	}

	dbpassword, err := app.models.Users.GetUserPassword(userLogin.Username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.userNotFound(w, r)
		default:
			app.internalServerError(w, r)
		}
		return
	}

	if mismatch := app.comparePassword([]byte(userLogin.Password), []byte(dbpassword)); mismatch != nil {
		if err := app.models.Users.UpdateLoginAttempts(userLogin.Username); err != nil {
			log.Error("error while updating login attempts ", err)
			app.internalServerError(w, r)
			return
		}
		app.wrongcredentials(w, r)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (app *application) CheckUserExists(w http.ResponseWriter, r *http.Request) {
	var username string

	err := app.readJSON(r, w, username)
	if err != nil {
		app.sendErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	exists, err := app.models.Users.CheckUserExists(username)
	if err != nil {
		log.Error(err)
		app.internalServerError(w, r)
		return
	}

	if exists {
		app.userExists(w, r)
	} else {
		w.WriteHeader(http.StatusOK)
	}

}

func (app *application) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	//READJSON
	var password string

	err := app.readJSON(r, w, password)
	if err != nil {
		app.sendErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	user := app.contextGetUser(r)

	hashedpassword, err := app.generateHashedPassword([]byte(password))
	if err != nil {
		app.internalServerError(w, r)
		return
	}

	if err := app.models.Users.UpdatePassword(user.ID, hashedpassword); err != nil {
		app.internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (app *application) UpdateUserDetails(w http.ResponseWriter, r *http.Request) {

	var userDetails struct {
		Bio       string    `json:"bio"`
		BirthDate time.Time `json:"birthdate"`
	}

	err := app.readJSON(r, w, userDetails)
	if err != nil {
		app.sendErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	user := app.contextGetUser(r)

	user.Bio = userDetails.Bio
	user.BirthDate = userDetails.BirthDate

	if err := app.models.Users.UpdateUser(user); err != nil {
		app.internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) GetUserProfilePicture(w http.ResponseWriter, r *http.Request) {
	var userid uint64

	userid, err := app.readParamID(r)
	if err != nil {
		app.sendErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	image, err := app.models.Images.GetProfilePicture(userid)
	if err != nil {
		app.internalServerError(w, r)
		return
	}

	http.ServeFile(w, r, image.Location)
}

func (app *application) UpdateProfilePicture(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		app.fileSizeLimitExceded(w, r)
		return
	}

	file, h, err := r.FormFile("image")
	if err != nil {
		app.internalServerError(w, r)
		return
	}

	user := app.contextGetUser(r)

	image := &data.Image{
		UserID:   user.ID,
		Size:     h.Size,
		Location: imagesDir + h.Filename,
	}

	if err := app.models.Images.UpdateProfilePicture(image); err != nil {
		app.internalServerError(w, r)
		return
	}

	destfile, err := os.Create(imagesDir + h.Filename)
	if err != nil {
		app.internalServerError(w, r)
		return
	}
	defer destfile.Close()

	_, err = io.Copy(destfile, file)
	if err != nil {
		app.internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) GetDeletedUsers(w http.ResponseWriter, r *http.Request) {
	allDeletedUsers, err := app.models.Users.FindSoftDeletedRecords()
	if err != nil {
		app.internalServerError(w, r)
		return
	}

	app.writeJSON(w, envelope{"deleted_users": allDeletedUsers}, http.StatusOK)
}

// // probably huge
// func (app *application) GetAllActiveUsers(w http.ResponseWriter, r *http.Request) {
// 	activeUsers, err := app.models.Users.FindActiveUsers()
// 	if err != nil {
// 		app.internalServerError(w, r)
// 		return
// 	}

// 	app.writeJSON(w, envelope{"deleted_users": activeUsers}, http.StatusOK)
// }
