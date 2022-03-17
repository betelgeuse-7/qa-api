package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type message struct {
	Msg string `json:"message"`
}

func newMessage(msg interface{}) *message {
	return &message{
		Msg: fmt.Sprintf("%v", msg),
	}
}

func (m *message) sendJSON(dest io.Writer) {
	json.NewEncoder(dest).Encode(m)
}

func (m *message) json() string {
	bx, err := json.Marshal(m)
	if err != nil {
		return fmt.Sprintf("ERR: %s", err.Error())
	}
	return string(bx)
}

// give http.StatusInternalServerError to w, and log error to os.Stdout
func http500(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(err)
}
