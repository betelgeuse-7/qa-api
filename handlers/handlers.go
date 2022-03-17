package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"qa/models"

	"golang.org/x/crypto/bcrypt"
)

/* -------------------- USERS -------------------- */

func GetUserProfile(w http.ResponseWriter, r *http.Request) {}

func NewUser(w http.ResponseWriter, r *http.Request) {
	u := &models.User{}
	json.NewDecoder(r.Body).Decode(u)

	hashBx, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		http500(w, fmt.Errorf("bcrypt hash error: %s", err.Error()))
		return
	}
	err = u.Insert(string(hashBx))
	if err != nil {
		http500(w, err)
		return
	}
	newMessage("new user registered successfully").sendJSON(w)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

}

/* -------------------- QUESTIONS -------------------- */

func NewQuestion(w http.ResponseWriter, r *http.Request) {

}

func NewAnswer(w http.ResponseWriter, r *http.Request) {

}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {

}

/* -------------------- ANSWERS -------------------- */

func DeleteAnswer(w http.ResponseWriter, r *http.Request) {

}

func UpdateQuestion(w http.ResponseWriter, r *http.Request) {

}

func UpdateAnswer(w http.ResponseWriter, r *http.Request) {

}
