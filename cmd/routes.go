package main

import (
	"125_isbn_new/ui"
	"io/fs"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/justinas/alice"
)

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
	app.logger.Info("Entrée dans routes.go")
	f := fs.FS(ui.Files)
	v, _ := fs.Sub(f, "static")
	router := httprouter.New()

	// Convertit l'assistant notFoundResponse() en http.Handler en utilisant
	// l’Adaptateur http.HandlerFunc(), puis la définit comme gestionnaire
	// d'erreurs personnalisé pour 404 : Réponses introuvables.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// De même, convertit l'assistant methodNotAllowedResponse() en http.Handler et définit
	// en tant que gestionnaire d'erreurs personnalisé pour les réponses 405 Method Not Allowed.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Utilisez la fonction http.FileServerFS() pour créer un gestionnaire HTTP qui
	// sert les fichiers intégrés dans ui.Files. Il est important de noter que nos
	// fichiers statiques sont contenus dans le dossier système de fichiers intégré
	//"static" de ui.Files. Ainsi, par exemple, notre feuille de style CSS se trouve à l'adresse
	// "static/css/main.css". Cela signifie que nous n'avons plus besoin d'indiquer le
	// préfixe de l'URL de la requête -- toutes les requêtes commençant par /static/ peuvent
	// transmis directement au serveur de fichiers et au fichier statique correspondant
	// le fichier sera servi (tant qu'il existe).
	//router.Handler(http.MethodGet, "/static/", http.FileServerFS(ui.Files))
	router.ServeFiles("/static/*filepath", http.FS(v))

	// Sending the assets to the clients: remplacé par la ligne au dessus
	// fs := http.FileServer(http.Dir(models.Path + "assets"))
	// mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	// Add a new GET /ping route.
	router.HandlerFunc(http.MethodGet, "/ping", ping)
	router.HandlerFunc(http.MethodGet, "/images/couverture/:name", app.CouvertureLivreGet)

	// Add the authenticate() middleware to the chain.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.index))
	router.Handler(http.MethodGet, "/index", dynamic.ThenFunc(app.IndexHandlerGet))
	router.Handler(http.MethodGet, "/affichelivres", dynamic.ThenFunc(app.LivresHandlerGet))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/view", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodGet, "/user/modif", dynamic.ThenFunc(app.userModifGet))
	router.Handler(http.MethodPost, "/user/modif", dynamic.ThenFunc(app.userModifPost))

	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/home", protected.ThenFunc(app.HomeHandlerGet))
	router.Handler(http.MethodGet, "/livre", protected.ThenFunc(app.LivreHandlerGet))
	router.Handler(http.MethodPost, "/affichelivre", protected.ThenFunc(app.LivreHandlerPost))
	router.Handler(http.MethodGet, "/afficheauteurs", protected.ThenFunc(app.ListAuteursHandlerGet))
	router.Handler(http.MethodGet, "/afficheediteurs", protected.ThenFunc(app.ListEditeursHandlerGet))
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodGet, "/user/logout", protected.ThenFunc(app.userLogoutGet))
	router.Handler(http.MethodGet, "/movies/connectuserapi", protected.ThenFunc(app.ConnectUserApiGet))
	router.Handler(http.MethodPost, "/movies/connectuserapi", protected.ThenFunc(app.ConnectUserApiPost))
	router.Handler(http.MethodPost, "/movies/activeuserapi", protected.ThenFunc(app.ActiveUserPost))
	router.Handler(http.MethodGet, "/movies/authenticateuserapi", protected.ThenFunc(app.AuthenticateUserApiGet))
	router.Handler(http.MethodPost, "/movies/authenticateuserapi", protected.ThenFunc(app.AuthenticateUserApiPost))
	router.Handler(http.MethodGet, "/movies/refreshtokensuserapi", protected.ThenFunc(app.refreshTokensHandlerApiGet))

	//protectedmovie := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate, app.requireAuthentication, app.requireCompteapi)
	router.Handler(http.MethodGet, "/movie/view", protected.ThenFunc(app.MovieViewGet))
	router.Handler(http.MethodPost, "/movie/view", protected.ThenFunc(app.MovieViewPost))
	router.Handler(http.MethodGet, "/movies/view", protected.ThenFunc(app.MoviesViewGet))

	standard := alice.New(app.recoverPanic, app.Log, app.commonHeaders)
	return standard.Then(router)

}
