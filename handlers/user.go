package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"qa/models/pgmodels"
	jwtservice "qa/service/jwt"

	"golang.org/x/crypto/bcrypt"
)

// TODO https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

// username | email && password
func Login(w http.ResponseWriter, r *http.Request) {
	user := &pgmodels.User{}
	json.NewDecoder(r.Body).Decode(user)
	givenPwd := user.Password

	if len(user.Password) == 0 {
		_HTTP_BAD_REQUEST(w, "you must provide password")
		return
	}
	username := user.Username
	email := user.Email
	if len(username) == 0 && len(email) == 0 {
		_HTTP_BAD_REQUEST(w, "you must provide either an email or a username for logging in")
		return
	}
	if username != "" && email != "" {
		_HTTP_BAD_REQUEST(w, "please choose either logging in with username or email, not both.")
		return
	}
	err := user.GetPassword()
	if err != nil {
		if err == sql.ErrNoRows {
			_HTTP_BAD_REQUEST(w, "there is no user with username, '"+username+"'")
			return
		}
		_HTTP_INTERNAL_SERVER_ERROR(w, err)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(givenPwd))
	if err != nil {
		_HTTP_FORBIDDEN(w, "authentication failed. wrong password")
		return
	}
	// generate access token and give it to user
	at := jwtservice.NewAccessToken(user.UserId)
	atStr, err := at.SignedString()
	if err != nil {
		_HTTP_INTERNAL_SERVER_ERROR(w, err)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": atStr,
	})
}

// get profile info
// get current user's id from r.Context
func GetProfile(w http.ResponseWriter, r *http.Request) {}

// get user info with user id (e.g. /user/17)
func GetUser(w http.ResponseWriter, r *http.Request) {
	userHandle, _ := valueFromChiUrlParam(r, "handle", STRING)
	user := &pgmodels.User{
		Handle: userHandle.(string),
	}
	err := user.GetUser()
	if err != nil {
		if err == sql.ErrNoRows {
			_HTTP_NOT_FOUND(w, "'"+userHandle.(string)+"' not found")
			return
		}
		_HTTP_INTERNAL_SERVER_ERROR(w, err)
	}
	json.NewEncoder(w).Encode(user)
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	u := &pgmodels.User{}
	json.NewDecoder(r.Body).Decode(u)

	errs, ok := u.Validate()
	if !ok {
		json.NewEncoder(w).Encode(map[string]pgmodels.ValidationErrors{
			"errors": errs,
		})
		return
	}
	u.Handle = "@" + u.Handle
	hashBx, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		_HTTP_INTERNAL_SERVER_ERROR(w, fmt.Errorf("bcrypt hash error: %s", err.Error()))
		return
	}
	err = u.Insert(string(hashBx))
	if err != nil {
		_HTTP_INTERNAL_SERVER_ERROR(w, err)
		return
	}
	newMessage("new user registered successfully").sendJSON(w)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userHandle, err := valueFromChiUrlParam(r, "handle", STRING)
	if err != nil {
		_HTTP_INTERNAL_SERVER_ERROR(w, err)
		return
	}
	u := &pgmodels.User{
		Handle: userHandle.(string),
	}
	err = u.Delete()
	if err != nil {
		if err == sql.ErrNoRows {
			_HTTP_NOT_FOUND(w, "there is not a user with handle '"+userHandle.(string)+"'")
			return
		}
		_HTTP_INTERNAL_SERVER_ERROR(w, err)
		return
	}

	newMessage("deleted user with handle '" + userHandle.(string) + "'").sendJSON(w)
}

// needs authorization
func UpdateUser(w http.ResponseWriter, r *http.Request) {

}
