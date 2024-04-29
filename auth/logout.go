package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dubey22rohit/heyyy_yo_backend/pkg/customerror"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/helper"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/success"
)

func (app *Application) logoutUserHandler(w http.ResponseWriter, r *http.Request) {
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

	//get reddis session
	_, err = helper.GetFromRedis(app.RedisClient, fmt.Sprintf("sessionid_%s", userID.Id))
	if err != nil {
		customerror.UnauthorizedResponse(w, r, errors.New("you are not authorized to perform this action"))
		return
	}

	//delete session from redis
	ctx := context.Background()
	_, err = app.RedisClient.Del(ctx, fmt.Sprintf("sessionid_%s", userID.Id)).Result()
	if err != nil {
		customerror.ServerErrorResponse(w, r, errors.New("something went wrong"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "sessionid",
		Value:   "",
		Expires: time.Now(),
	})

	success.SuccessResponse(w, r, http.StatusOK, "you have logged out successfully")
}
