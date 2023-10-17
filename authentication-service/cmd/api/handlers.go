package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil || user == nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	request, err := http.NewRequest("POST", "http://docker.for.mac.localhost:8081/log", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	if _, err := client.Do(request); err != nil {
		return err
	}
	return nil
}
