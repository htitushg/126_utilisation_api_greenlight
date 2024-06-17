package main

import (
	"125_isbn_new/internal/models"
	"125_isbn_new/ui"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the build progresses.
type templateData struct {
	CurrentYear     int
	Date            string
	Message         string
	Path            string
	Book            models.ApiBooks
	Livre           models.Livre
	Books           []models.ApiBooks
	Livres          []models.Livre
	Auteurs         []models.Auteur
	Editeurs        []models.Editeur
	Movies          []models.Movie
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string // Add a Flash field to the templateData struct.
	IsAuthenticated bool   // Add an IsAuthenticated field to the templateData struct.
	Username        string
	CSRFToken       string // Add a CSRFToken field.
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
	// Return the empty string if time has the zero value.
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 January 2006 à 15:04")
}

/*
Initialise un objet template.FuncMap et le stocke dans une variable globale.
C'est essentiellement une carte à clé de chaîne qui agit comme une recherche
entre les noms de nos fonctions de modèle personnalisées
et les fonctions elles-mêmes.
*/
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	// Use fs.Glob() to get a slice of all filepaths in the ui.Files embedded
	// filesystem which match the pattern 'html/pages/*.tmpl'. This essentially
	// gives us a slice of all the 'page' templates for the application, just
	// like before.
	pages, err := fs.Glob(ui.Files, "templates/*.gohtml")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		// Create a slice containing the filepath patterns for the templates we
		// want to parse.
		patterns := []string{
			"templates/layouts/base.gohtml",
			"templates/partials/*.gohtml",
			page,
		}
		/*
			Utiliser ParseFS() à la place de ParseFiles() pour analyser( to parse)
			les fichiers template venant du système de fichiers ui.Files intégré(embedded) .

			Funcs ajoute les éléments de la carte(map) d'arguments à la map de fonctions du modèle.
			Il doit être appelé avant que le modèle ne soit analysé.
			Cela panique si une valeur entrée dans la carte n'est pas une fonction avec un retour approprié.
			Cependant, il est légal d’écraser des éléments de la carte. La valeur de retour
			est le modèle, donc les appels peuvent être chaînés.
		*/
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
