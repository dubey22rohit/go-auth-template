package success

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/dubey22rohit/heyyy_yo_backend/pkg/customerror"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/helper"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/jsonlog"
)

type envelope map[string]interface{}

func LogSuccess(r *http.Request, status int, message interface{}) {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	logger.PrintInfo(fmt.Sprintf("%v", message), map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
		"status":         strconv.Itoa(status),
	}, true)
}

func SuccessResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"message": message}

	err := helper.WriteJSON(w, status, env, nil)
	if err != nil {
		customerror.LogError(r, err)
		w.WriteHeader(500)
	}
	LogSuccess(r, status, message)
}
