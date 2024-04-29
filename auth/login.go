package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"net/http"

	"github.com/dubey22rohit/heyyy_yo_backend/auth/internal/data"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/cookies"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/customerror"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/helper"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/success"
)

func (app *Application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	//Data sent by user
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//reading user input to JSON
	err := helper.ReadJSON(w, r, &input)
	if err != nil {
		customerror.BadRequestResponse(w, r, err)
		return
	}

	dbUser, err := app.Models.Users.GetByEmail(input.Email)
	if err != nil {
		customerror.BadRequestResponse(w, r, err)
		return
	}

	match, err := dbUser.Password.Matches(input.Password)
	if err != nil {
		return
	}

	if !match {
		customerror.BadRequestResponse(w, r, errors.New("email and password combination does not match"))
		return
	}

	var userID = data.UserID{
		Id: dbUser.ID,
	}

	//gob encoding the user data, storing the encoded output in the buffer
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(&userID)
	if err != nil {
		// TODO: this error message is for debugging, will change it in PROD
		customerror.ServerErrorResponse(w, r, errors.New("error encoding data"))
		return
	}

	session := buf.String()

	//store session in redis
	err = helper.StoreInRedis(app.RedisClient, "sessionid_", session, userID.Id, app.Config.Secret.SessionExpiration)
	if err != nil {
		customerror.LogError(r, err)
	}

	cookie := http.Cookie{
		Name:     "sessionid",
		Value:    session,
		Path:     "/",
		MaxAge:   int(app.Config.Secret.SessionExpiration.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	// write an encrypted cookie containing gob encoded data as normal
	err = cookies.WriteEncrypted(w, cookie, app.Config.Secret.SecretKey)
	if err != nil {
		customerror.ServerErrorResponse(w, r, errors.New("something happened setting your cookie data"))
		return
	}

	helper.WriteJSON(w, http.StatusOK, dbUser, nil)
	if err != nil {
		customerror.ServerErrorResponse(w, r, err)
		return
	}
	success.LogSuccess(r, http.StatusOK, "logged in successfully")
}
