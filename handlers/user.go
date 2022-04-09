package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"qa/models/pgmodels"

	"golang.org/x/crypto/bcrypt"
)

// TODO https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

func Login(w http.ResponseWriter, r *http.Request) {

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

/*
	need a json object containing all the new fields of user with {handle}
*/
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userHandle, err := valueFromChiUrlParam(r, "handle", STRING)
	if err != nil {
		_HTTP_INTERNAL_SERVER_ERROR(w, err)
		return
	}

	user := &pgmodels.User{
		Handle: userHandle.(string),
	}

	payload := &pgmodels.User{}
	json.NewDecoder(r.Body).Decode(payload)

	/* check every field in r.Body */
	/* as soon as any field is there, we update the-user-with-handle's that field with the new value */
	/* precedences: username > email > password > handle */

	if len(payload.Username) > 0 {
		user.UpdateUsername(payload.Username)
		newMessage("updated username of user with handle '" + user.Handle + "' to '" + payload.Username + "'").sendJSON(w)
		return
	} else if len(payload.Email) > 0 {
		user.UpdateEmail(payload.Email)
		newMessage("updated email of user with handle '" + user.Handle + "' to '" + payload.Email + "'").sendJSON(w)
		return
	} else if len(payload.Password) > 0 {
		user.UpdatePassword(payload.Password)
		newMessage("updated password of user with handle '" + user.Handle + "'").sendJSON(w)
		return
	} else if len(payload.Handle) > 0 {
		user.UpdateHandle(payload.Handle)
		newMessage("updated handle of user with handle '" + user.Handle + "' to '" + payload.Handle + "'").sendJSON(w)
		return
	} else {
		_HTTP_BAD_REQUEST(w, "empty payload. nothing to update")
	}
}
