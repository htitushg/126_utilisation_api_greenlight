package main

import (
	"125_isbn_new/internal/assert"
	"bufio"
	"bytes"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-playground/form/v4"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/nosurf"

	//"errors"
	"fmt"
	//"github.com/go-playground/form/v4"
	//"github.com/justinas/nosurf"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		// Use debug.Stack() to get the stack trace. This returns a byte slice, which
		// we need to convert to a string so that it's readable in the log entry.
		trace = string(debug.Stack())
	)
	// Include the trace in the log entry.
	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Create an newTemplateData() helper, which returns a pointer to a templateData
// struct initialized with the current year. Note that we're not using the
// *http.Request parameter here at the moment, but we will do later in the book.
func (app *application) newTemplateData(r *http.Request) templateData {
	var name string
	if r.Context().Value("UserName") != nil {
		name = r.Context().Value("UserName").(string)
	}
	return templateData{
		CurrentYear:     time.Now().Year(),
		Date:            time.Now().Format("02/01/2006"),
		Message:         "",
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		Username:        name,
		CSRFToken:       nosurf.Token(r), // Add the CSRF token.
	}
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}

// Create a new decodePostForm() helper method. The second parameter here, dst,
// is the target destination that we want to decode the form data into.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// Call ParseForm() on the request, in the same way that we did in our
	// snippetCreatePost handler.
	err := r.ParseForm()
	if err != nil {
		return err
	}
	// Call Decode() on our decoder instance, passing the target destination as
	// the first parameter.
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If we try to use an invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError.We use
		// errors.As() to check for this and raise a panic rather than returning
		// the error.
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		// For all other errors, we return them as normal.
		return err
	}
	return nil
}

// Return true if the current request is from an authenticated user, otherwise
// return false.
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	fmt.Printf("isAuthenticated: %#v\n", isAuthenticated)
	fmt.Printf("ok: %#v\n", ok)
	if !ok {
		return false
	}
	return isAuthenticated
}

// Retrieve the "id" URL parameter from the current request context, then convert it to
// an integer and return it. If the operation isn't successful, return 0 and an error.
func (app *application) readIDParam(r *http.Request) (int, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return int(id), nil
}

// Retrieve the "id" URL parameter from the current request context, then convert it to
// an integer and return it. If the operation isn't successful, return 0 and an error.
func (app *application) readName(r *http.Request) (string, error) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")
	return name, nil
}
func getEnvrc() (string, error) {
	var value string
	var key string
	// Lire le fichier .envrc
	data, err := os.ReadFile(assert.Path + ".envrc")
	if err != nil {
		log.Fatal("Error reading .envrc file:", err)
	}

	// Décomposer le contenu du fichier en lignes
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()

		// Ignorer les lignes vides et les commentaires
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Décomposer la ligne en clé et valeur
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("Invalid line in .envrc file: %s\n", line)
			continue
		}

		key = parts[0]
		value = parts[1]

		// Définir la variable d'environnement
		err := os.Setenv(key, value)
		if err != nil {
			log.Printf("Error setting environment variable %s: %s\n", key, err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error scanning .envrc file:", err)
	}
	return value, err
}
