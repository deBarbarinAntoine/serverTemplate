package main

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"html/template"
	"log/slog"
	"serverTemplate/internal/models"
	"serverTemplate/internal/validator"
)

type application struct {
	logger         *slog.Logger
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	users          *models.UserModel
}

type templateData struct {
	CurrentYear     int
	Message         string
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}
