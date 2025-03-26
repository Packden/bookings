package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
)

// this file must not import any functions from the application. Limit libraries to standard go packages.

// AppConfig store application state values used across all packages
type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	InProduction  bool
	Session       *scs.SessionManager
}
