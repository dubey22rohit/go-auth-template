package main

import (
	"net/http"
	"strconv"

	"github.com/dubey22rohit/heyyy_yo_backend/pkg/customerror"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/helper"
)

func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "available",
		"debug":   strconv.FormatBool(app.Config.Debug),
		"version": version,
	}

	err := helper.WriteJSON(w, http.StatusOK, data, nil)

	if err != nil {
		customerror.ServerErrorResponse(w, r, err)
	}
}
