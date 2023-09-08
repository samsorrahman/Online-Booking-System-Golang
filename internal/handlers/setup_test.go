package handlers

import (
	// Import necessary packages and libraries
	// ...

	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/samsorrahman/Online-Booking-System-Golang/internal/config"
	"github.com/samsorrahman/Online-Booking-System-Golang/internal/models"
	"github.com/samsorrahman/Online-Booking-System-Golang/internal/render"
)

// Define global variables
var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

// getRoutes returns an http.Handler that defines all the application routes and middleware.
func getRoutes() http.Handler {
	// Register the Reservation model for use with the gob package
	gob.Register(models.Reservation{})

	// Set the production flag (change this to true in production)
	app.InProduction = false

	// Initialize and configure the session manager
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	// Store the session manager in the app configuration
	app.Session = session

	// Create the template cache and configure the application
	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = true

	// Create a new repository and handlers for the application
	repo := NewRepo(&app)
	NewHandlers(repo)

	// Initialize templates for rendering
	render.NewTemplates(&app)

	// Create a new chi router for routing requests
	mux := chi.NewRouter()

	// Middleware: Recoverer - recovers from panics and logs them
	mux.Use(middleware.Recoverer)

	// Middleware: NoSurf - adds CSRF protection to routes
	// mux.Use(NoSurf)

	// Middleware: SessionLoad - loads and saves session data for the current request
	mux.Use(SessionLoad)

	// Define the application routes
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	// Serve static files from the "static" directory
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// NoSurf is a middleware function that provides CSRF protection for routes.
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	// Configure the CSRF protection cookie
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad is a middleware function that loads and saves session data for the current request.
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTestTemplateCache creates a template cache as a map.
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// Find all ".page.tmpl" files in the templates directory
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// Iterate over each page template and parse it
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// Find and parse any layout templates for the page
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
