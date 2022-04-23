package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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

// give http.StatusInternalServerError to w, and log error to os.Stdout
func _HTTP_INTERNAL_SERVER_ERROR(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println("HTTP INTERNAL SERVER ERROR 500 -> ", err)
}

func _HTTP_BAD_REQUEST(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
	log.Println("HTTP BAD REQUEST 400 -> " + message)
}

func _HTTP_NOT_FOUND(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(message))
	log.Println("HTTP NOT FOUND 404 -> " + message)
}

func _HTTP_FORBIDDEN(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(message))
	log.Println("HTTP FORBIDDEN 403 -> " + message)
}

type urlParamConvertible uint

const (
	INT urlParamConvertible = iota
	BOOL
	STRING
)

// get url param value convert it into desired type and return it
func valueFromChiUrlParam(r *http.Request, paramName string, as urlParamConvertible) (interface{}, error) {
	urlParam := chi.URLParam(r, paramName)

	switch as {
	case INT:
		res, err := strconv.Atoi(urlParam)
		if err != nil {
			return nil, err
		}
		return res, nil
	//string
	default:
		return urlParam, nil
	}
}
