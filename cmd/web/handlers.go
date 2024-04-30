package main

import (
	"errors"
	"net/http"
	"serverTemplate/internal/models"
	"serverTemplate/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	data.Message = "Welcome to my website!"

	app.sessionManager.Put(r.Context(), "flash", "Welcome home!")

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the About page!"))
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	data.Form = newUserLoginForm()

	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {

	form := newUserLoginForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	var id int
	id, err = app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Invalid credentials")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/protected", http.StatusSeeOther)
}

func (app *application) logoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserId")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) restrictedSomething(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the restricted page!"))
}
