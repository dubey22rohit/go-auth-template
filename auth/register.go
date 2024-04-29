package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/dubey22rohit/heyyy_yo_backend/auth/internal/data"
	"github.com/dubey22rohit/heyyy_yo_backend/auth/internal/tokens"
	"github.com/dubey22rohit/heyyy_yo_backend/auth/internal/validator"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/customerror"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/helper"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/success"
)

func (app *Application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := helper.ReadJSON(w, r, &input)
	if err != nil {
		customerror.BadRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Email:    input.Email,
		Username: input.Username,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		customerror.ServerErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		customerror.FailedValidationResponse(w, r, v.Errors)
		return
	}

	userID, err := app.Models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email already exists")
			customerror.FailedValidationResponse(w, r, v.Errors)
		default:
			customerror.ServerErrorResponse(w, r, err)
		}
		return
	}

	otp, err := tokens.GenerateOTP()
	if err != nil {
		customerror.LogError(r, err)
	}

	err = helper.StoreInRedis(app.RedisClient, "activation_", otp.Hash, userID.Id, app.Config.TokenExpiration.Duration)
	if err != nil {
		customerror.LogError(r, err)
	}

	now := time.Now()
	expiration := now.Add(app.Config.TokenExpiration.Duration)
	exact := expiration.Format(time.RFC1123)

	helper.RunBackgroundTask(func() {
		data := map[string]interface{}{
			"token":       tokens.FormatOTP(otp.Secret),
			"userID":      userID.Id,
			"frontendURL": app.Config.FrontendURL,
			"expiration":  app.Config.TokenExpiration.DurationString,
			"exact":       exact,
		}
		err = app.Mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			customerror.LogError(r, err)
		}
		app.Logger.PrintInfo("Email successfully sent.", nil, app.Config.Debug)
	})

	success.SuccessResponse(w, r, http.StatusAccepted, "Your account creation was accepted successfully. Check your email address and follow the instruction to activate your account. Ensure you activate your account before the token expires")
}
