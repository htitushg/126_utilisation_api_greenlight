package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/justinas/nosurf"
)

func (app *application) commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: This is split across multiple lines for readability. You don't
		// need to do this in your own code.
		//app.logger.Info("Entrée dans commonHeaders")
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Server", "Go")
		next.ServeHTTP(w, r)
	})
}

var LogId = 0

// Log is a models.Middleware that writes a series of information in logs/logs_<date>.log
// in JSON format: time, client's type, request Id (incremented int),
// user's models.Session (if logged), client IP, request Method, and request URL.
func (app *application) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)
		LogId++
		//log.Println("middlewares.Log()")
		//fmt.Printf("user= %v\n", user)
		/* name := r.Context().Value("UserName")
		if name == nil { */
		//session, _ := r.Cookie("session")
		//fmt.Printf("Middleware Log app.user : %v\n", app.user)
		status := app.isAuthenticated(r)
		//fmt.Printf("Middleware Log status= %v\n", status)
		if !status {
			app.logger.Info("Visitor", slog.Int("req_id", LogId), "ip", ip, "proto", proto, "method", method, "uri", uri)
		} else {
			app.logger.Info("User", slog.Int("req_id", LogId), slog.Any("user", app.sessionManager.Cookie.Name), "ip", ip, "proto", proto, "method", method, "uri", uri)
		}

	})
}

func (app *application) requireCompteapi(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// vérifier l'existence d'un compte vers l'API greenlight
		data := app.newTemplateData(r)

		user, ok := app.user.GetUser(data.Username)
		if !ok {
			//pas de connexion
			data.Message = "Il n'a pas été possible de lire le jeton API dans la base"
			app.render(w, r, http.StatusUnprocessableEntity, "home.gohtml", data)
			return
		}
		token, err := app.movies.LireJetonDansBase(user.Id)
		if err != nil {
			data.Message = "Il n'a pas été possible de lire le jeton API dans la base"
			app.render(w, r, http.StatusUnprocessableEntity, "home.gohtml", data)
		}
		// vérifier si le token est valide
		if token.Expiry.Before(time.Now()) || (token.Token == "") {
			data.Message = "Le jeton est expiré"
			app.render(w, r, http.StatusUnprocessableEntity, "home.gohtml", data)
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//app.logger.Info("Entrée dans recoverPanic")
		// Créez une fonction différée (qui sera toujours exécutée en cas d'événement
		// panique alors que Go déroule la pile).
		defer func() {
			// Utilisez la fonction de récupération intégrée pour vérifier s'il y a eu une
			// panique ou pas. S'il y en a...
			if err := recover(); err != nil {
				// Définir un en-tête "Connection: close" sur la réponse.
				w.Header().Set("Connection", "close")
				// Appelez la méthode d'assistance app.serverError pour renvoyer une
				// Réponse 500 au serveur interne.
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//app.logger.Info("Entrée dans requireAuthentication")
		// If the user is not authenticated, redirect them to the login page and
		// return from the middleware chain so that no subsequent handlers in
		// the chain are executed.
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// Otherwise set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the users browser cache (or
		// other intermediary cache).
		w.Header().Add("Cache-Control", "no-store")
		// And call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// Create a NoSurf middleware function which uses a customized CSRF cookie with
// the Secure, Path and HttpOnly attributes set.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//app.logger.Info("Entrée dans authenticate")

		// Récupérons la valeur AuthenticatedUserID de la session à l'aide de la
		// Méthode GetInt(). Cela renvoie la valeur int zéro si
		// La valeur "authenticatedUserID" n'est pas dans la session -- auquel cas nous
		// appellons le gestionnaire(handler) suivant dans la chaîne (next.ServeHTTP(w, r))
		// comme d'habitude et "return".
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}
		// Sinon, nous vérifions si un utilisateur avec cet identifiant existe dans notre
		// base de données.
		exists, err := app.user.Exists(id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		// Si un utilisateur correspondant est trouvé, nous savons que la demande
		// provient d'un utilisateur authentifié qui existe dans notre base de données. Nous
		// créons une nouvelle copie de la requête (avec un isAuthenticatedContextKey
		// valeur à true dans le contexte de la requête) et nous l'assignons à r.
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			user, _ := app.user.SelectUserwithId(id)
			ctx = context.WithValue(ctx, "UserName", user.Name)
			r = r.WithContext(ctx)
		}
		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
