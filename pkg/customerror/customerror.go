package customerror

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dubey22rohit/heyyy_yo_backend/pkg/helper"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/jsonlog"
)

type envelope map[string]interface{}

func LogError(r *http.Request, err error) {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	}, true)
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := helper.WriteJSON(w, status, env, nil)
	if err != nil {
		LogError(r, err)
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	LogError(r, err)
	message := "the server encountered a problem and could not process your request"
	ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	ErrorResponse(w, r, http.StatusNotFound, message)
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	LogError(r, err)
	ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}
func UnauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	LogError(r, err)
	ErrorResponse(w, r, http.StatusUnauthorized, err.Error())
}

func FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
