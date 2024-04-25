package main

import (
	"github.com/justinas/alice"
	"net/http"
	"serverTemplate/ui"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /about", dynamic.ThenFunc(app.about))
	mux.Handle("GET /login", dynamic.ThenFunc(app.login))
	mux.Handle("POST /login", dynamic.ThenFunc(app.loginPost))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /protected", protected.ThenFunc(app.restrictedSomething))
	mux.Handle("POST /logout", protected.ThenFunc(app.logoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
