package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dubey22rohit/heyyy_yo_backend/pkg/customerror"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/helper"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/success"
)

func (app *Application) currentUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, status, err := helper.ExtractParamsFromSession(r, app.Config.Secret.SecretKey)
	if err != nil {
		switch *status {
		case http.StatusUnauthorized:
			customerror.UnauthorizedResponse(w, r, err)
		case http.StatusBadRequest:
			customerror.BadRequestResponse(w, r, err)
		case http.StatusInternalServerError:
			customerror.ServerErrorResponse(w, r, err)
		default:
			customerror.ServerErrorResponse(w, r, errors.New("something went wrong"))
		}
		return
	}

	_, err = helper.GetFromRedis(app.RedisClient, fmt.Sprintf("sessionid_%s", userID.Id))
	if err != nil {
		customerror.UnauthorizedResponse(w, r, errors.New("you are unauthorized to perform this action"))
		return
	}

	dbUser, err := app.Models.Users.Get(userID.Id)
	if err != nil {
		customerror.BadRequestResponse(w, r, err)
		return
	}

	helper.WriteJSON(w, http.StatusOK, dbUser, nil)
	if err != nil {
		customerror.ServerErrorResponse(w, r, err)
		return
	}

	success.LogSuccess(r, http.StatusOK, "user retrieved successfully")
}
