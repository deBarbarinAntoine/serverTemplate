package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
	"log/slog"
	"net/http"
	"runtime/debug"
	"serverTemplate/internal/validator"
	"time"
)

func (app *application) decodePostForm(r *http.Request, dst any) error {

	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		status = http.StatusInternalServerError
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), slog.String("method", method), slog.String("URI", uri), slog.String("trace", trace))
	http.Error(w, http.StatusText(status), status)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		//debug
		app.logger.Debug("User is not authenticated", slog.Any("isAuthenticated", isAuthenticated))
		return false
	}

	return isAuthenticated
}

func newUserLoginForm() *userLoginForm {
	return &userLoginForm{
		Validator: *validator.NewValidator(),
	}
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {

	// retrieving the appropriate set of templates
	ts, ok := app.templateCache[page]
	if !ok {
		app.serverError(w, r, fmt.Errorf("the template %s does not exist", page))
		return
	}

	// creating a bytes Buffer
	buf := new(bytes.Buffer)

	// executing the template in the buffer to catch any possible parsing error,
	// so that the user doesn't see a half-empty page
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// if it's all okay, write the status in the header and write the buffer in the ResponseWriter
	w.WriteHeader(status)

	buf.WriteTo(w)
}
