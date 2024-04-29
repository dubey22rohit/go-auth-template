package helper

import (
	"encoding/gob"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/dubey22rohit/heyyy_yo_backend/pkg/cookies"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/jsonlog"
	"github.com/dubey22rohit/heyyy_yo_backend/types"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func ReadIDParam(r *http.Request) (*uuid.UUID, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		return nil, errors.New("invalid id parameter")
	}
	return &id, nil
}

func ExtractParamsFromSession(r *http.Request, secretKey []byte) (*types.UserID, *int, error) {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	gobEncodedValue, err := cookies.ReadEncrypted(r, "sessionid", secretKey)

	if err != nil {
		var errorData error
		var status int
		switch {
		case errors.Is(err, http.ErrNoCookie):
			status = http.StatusUnauthorized
			errorData = errors.New("you are not authorized to access this resource")

		case errors.Is(err, cookies.ErrInvalidValue):
			logger.PrintError(err, nil, true)
			status = http.StatusBadRequest
			errorData = errors.New("invalid cookie")

		default:
			status = http.StatusInternalServerError
			errorData = errors.New("something happened getting your cookie data")

		}
		return nil, &status, errorData
	}

	var userID types.UserID

	reader := strings.NewReader(gobEncodedValue)
	if err := gob.NewDecoder(reader).Decode(&userID); err != nil {
		status := http.StatusInternalServerError
		return nil, &status, errors.New("something happened decosing your cookie data")
	}

	return &userID, nil, nil
}
