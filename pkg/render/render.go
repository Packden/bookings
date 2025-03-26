package render

import (
	"bytes"
	"github.com/packden/bookings/pkg/config"
	"github.com/packden/bookings/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData can be used append global application data to td
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate processes html template files disk read no cache being used
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		// get template cache from app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("template not found")
	}
	buf := new(bytes.Buffer)
	td = AddDefaultData(td)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println("Error executing template:", err)
	}
	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing template:", err)
	}
}
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	// get all file pages
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		log.Println("Error finding template pages:", err)
		return myCache, err
	}
	for i, page := range pages {
		name := filepath.Base(page)
		log.Println("Loading template:", i)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
		}
		if err != nil {
			return myCache, err
		}
		myCache[name] = ts
	}
	return myCache, nil
}
