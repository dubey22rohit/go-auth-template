package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"

	"github.com/dubey22rohit/heyyy_yo_backend/auth/internal/tokens"
	"github.com/dubey22rohit/heyyy_yo_backend/auth/internal/validator"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/customerror"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/helper"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/success"
)

func (app *Application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helper.ReadIDParam(r)

	if err != nil {
		customerror.BadRequestResponse(w, r, err)
		return
	}

	var input struct {
		Secret string `json:"token"`
	}

	err = helper.ReadJSON(w, r, &input)
	if err != nil {
		customerror.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if tokens.ValidateSecret(v, input.Secret); !v.Valid() {
		customerror.FailedValidationResponse(w, r, v.Errors)
		return
	}

	hash, err := helper.GetFromRedis(app.RedisClient, fmt.Sprintf("activation_%s", id))
	if err != nil {
		customerror.BadRequestResponse(w, r, err)
		return
	}

	tokenHash := fmt.Sprintf("%x\n", sha256.Sum256([]byte(input.Secret)))

	if *hash != tokenHash {
		app.Logger.PrintError(errors.New("the supplied token is invalid"), nil, app.Config.Debug)
		customerror.FailedValidationResponse(w, r, map[string]string{
			"token": "The supplied token is invalid",
		})
		return
	}

	_, err = app.Models.Users.ActivateUser(*id)
	if err != nil {
		customerror.ServerErrorResponse(w, r, err)
		return
	}

	ctx := context.Background()
	deleted, err := app.RedisClient.Del(ctx, fmt.Sprintf("activation_%s", id)).Result()
	if err != nil {
		app.Logger.PrintError(err, map[string]string{
			"key": fmt.Sprintf("activation_%s", id),
		}, app.Config.Debug)

	}

	app.Logger.PrintInfo(fmt.Sprintf("Token hash was deleted successfully :activation_%d", deleted), nil, app.Config.Debug)

	success.SuccessResponse(w, r, http.StatusOK, "Account activated successfully.")
}
