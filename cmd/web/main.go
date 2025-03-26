package main

import (
	"github.com/alexedwards/scs/v2"
	"github.com/packden/bookings/pkg/config"
	"github.com/packden/bookings/pkg/handlers"
	"github.com/packden/bookings/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Println("Error creating template cache:", err)
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	log.Printf("Listening on port %s", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)

}
