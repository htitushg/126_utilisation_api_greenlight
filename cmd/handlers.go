package main

import (
	"125_isbn_new/internal/assert"
	"125_isbn_new/internal/models"
	"125_isbn_new/internal/utils"
	"125_isbn_new/internal/validator"
	"125_isbn_new/ui"
	"errors"

	//"125_isbn_new/ui"
	"encoding/json"
	"fmt"
	"io"

	//"io/fs"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

/*
Supprimez le champ de structure FieldErrors explicite et intégrez à la place
le validateur de structure. L'intégration signifie que notre
snippetCreateForm "hérite" de tous les champs et méthodes
de notre structure Validator (y compris le champ FieldErrors).
*/
type errorForm struct {
	Numero  int
	Message string
}

// Create a new movieForm struct.
type movieForm struct {
	Id                  int64 `form:"id"`
	validator.Validator `form:"-"`
}
type tokenForm struct {
	Token               string `form:"token"`
	Name                string `form:"name"`
	validator.Validator `form:"-"`
}
type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

// Create a new userSignupForm struct.
type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// Create a new userSignupForm struct.
type userModifForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	PasswordActual      string `form:"passwordactual"`
	Password1           string `form:"password1"`
	Password2           string `form:"password2"`
	validator.Validator `form:"-"`
}

// Create a new userSignupForm struct.
type isbnForm struct {
	Isbn                string `form:"isbn"`
	validator.Validator `form:"-"`
}

/* // Create a new userLoginForm struct.
type userLoginForm struct {
	ID                  string `form:"id"`
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	Token               string `form:"token"`
	validator.Validator `form:"-"`
} */

func registerImage(url string, nomImage string) bool {

	// don't worry about errors
	response, e := http.Get(url)
	if e != nil {
		log.Println("erreur= ", e)
		return false
	}
	defer response.Body.Close()

	//open a file for writing
	//ui.Files.
	file, err := os.Create(assert.Path + "internal/img/" + nomImage + ".png")
	if err != nil {
		log.Println("erreur= ", err)
		return false
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Println("erreur= ", err)
		return false
	}
	fmt.Println("Success!")
	return true
}
func GetImageBook(book models.ApiBooks) (ok bool) {
	id := book.Items[0].ID
	if book.Items[0].VolumeInfo.ImageLinks.Thumbnail != "" {
		ok = registerImage(book.Items[0].VolumeInfo.ImageLinks.Thumbnail, id)
	}
	return ok
}
func GetImagesBooks(books []models.ApiBooks) (ok bool) {

	for _, item := range books {
		id := item.Items[0].ID
		if item.Items[0].VolumeInfo.ImageLinks.Thumbnail != "" {
			ok = registerImage(item.Items[0].VolumeInfo.ImageLinks.Thumbnail, id)
		}
	}
	return ok
}
func GetBooks() (books []models.ApiBooks, ok bool) {
	// Chemin du fichier JSON
	fileName := "livres.json"
	// Lire le fichier JSON
	data, err := ui.Files.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, false
	}

	// Décoder le fichier JSON dans une structure ApiBooks

	err = json.Unmarshal(data, &books)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, false
	}
	//_ = GetImagesBooks(books)
	return books, true
}
func TrierEtEliminerLesDoublons() {
	// Chemin du fichier JSON
	fileName := "../ui/static/livres.json"

	// Lire le fichier JSON
	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Décoder le fichier JSON dans une structure ApiBooks
	var books []models.ApiBooks
	err = json.Unmarshal(data, &books)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Fonction pour comparer les IDs
	less := func(i, j int) bool {
		return books[i].Items[0].ID < books[j].Items[0].ID
	}

	// Trier les enregistrements par ID
	sort.Slice(books, less)

	// Définir une map pour stocker les IDs uniques (clé) et leur index (valeur)
	uniqueIds := make(map[string]int)

	// Filtrer les doublons et créer un nouveau tableau
	var filteredData []models.ApiBooks

	for i, item := range books {
		if _, exists := uniqueIds[item.Items[0].ID]; !exists {
			uniqueIds[item.Items[0].ID] = i
			filteredData = append(filteredData, item)
		}
	}

	// Écrire le fichier JSON trié et sans doublons
	jsonData, err := json.MarshalIndent(filteredData, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Fichier trié et sans doublons écrit avec succès")

}
func EcrireNouvelEnregistrement(book models.ApiBooks) bool {
	var ok bool
	// Ouvrir le fichier JSON en lecture
	data, err := os.ReadFile("/home/henry/go/src/125_isbn_new/ui/static/livres.json")
	if err != nil {
		fmt.Println(err)
		return ok
	}

	// Décoder le fichier JSON dans une structure []ApiBooks
	var books []models.ApiBooks
	err = json.Unmarshal(data, &books)
	if err != nil {
		fmt.Println(err)
		return false
	}
	// Vérifier si le livre existe déjà
	resultat := false
	for _, item := range books {
		if item.Items[0].ID == book.Items[0].ID {
			resultat = true
			break
			// le livre existe déjà
		}
	}
	if !resultat { // si le livre n'existe pas resultat=false
		// Ajouter le livre en cours
		books = append(books, book)
		jsonData, err := json.MarshalIndent(books, "", "\t")
		if err != nil {
			panic(err)
		}
		// REEcrire dans le fichier json tous les livres
		err = os.WriteFile("../ui/static/livres.json", jsonData, 0666)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}
	return resultat

}
func GetBook(isbn string) (book models.ApiBooks, ok bool) {

	// instructions à exécuter pour chaque élément
	// Créer une requête HTTP GET vers l'API Google Books
	url := "https://www.googleapis.com/books/v1/volumes?q=isbn:" + isbn
	ok = true
	// Envoyer la requête et obtenir la réponse
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Lire le corps de la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	// Décoder le JSON de la réponse
	err = json.Unmarshal(body, &book)
	if err != nil {
		log.Println(err)
		return
	}
	_ = GetImageBook(book)
	return book, ok
}

func TranscrireBookVersLivre(book models.ApiBooks) (livre models.Livre, ok bool) {

	livre.Idg = book.Items[0].ID
	livre.Titre = book.Items[0].VolumeInfo.Title
	var auteur models.Auteur
	for z := 0; z < len(book.Items[0].VolumeInfo.Authors); z++ {
		livre.Auteurs = append(livre.Auteurs, auteur)
		if book.Items[0].VolumeInfo.Authors[z] != "" {
			livre.Auteurs[z].Nom = book.Items[0].VolumeInfo.Authors[z]
		} else {
			livre.Auteurs[z].Nom = "Inconnu"
		}
	}
	if book.Items[0].VolumeInfo.IndustryIdentifiers[0].Type == "OTHER" {
		livre.Isbn = book.Items[0].VolumeInfo.IndustryIdentifiers[0].Identifier
	} else {
		for z := 0; z < len(book.Items[0].VolumeInfo.IndustryIdentifiers); z++ {
			if book.Items[0].VolumeInfo.IndustryIdentifiers[z].Type == "ISBN_13" {
				livre.Isbn = book.Items[0].VolumeInfo.IndustryIdentifiers[z].Identifier
			}
		}
	}
	// Il faut tester le présence ou non d'une couverture
	if book.Items[0].VolumeInfo.ImageLinks.Thumbnail == "" {
		livre.Thumbnail = "normal.png"
	} else {
		livre.Thumbnail = book.Items[0].ID + ".png"
	}

	if book.Items[0].VolumeInfo.Publisher != "" {
		livre.Editeur.Nom = book.Items[0].VolumeInfo.Publisher
	} else {
		livre.Editeur.Nom = "Inconnu"
	}

	livre.Language = book.Items[0].VolumeInfo.Language
	livre.Publish_date = book.Items[0].VolumeInfo.PublishedDate
	livre.Nb_pages = book.Items[0].VolumeInfo.PageCount
	livre.Resume = book.Items[0].SearchInfo.TextSnippet
	livre.Description = book.Items[0].VolumeInfo.Description
	return livre, ok
}
func (app *application) IndexHandlerGet(w http.ResponseWriter, r *http.Request) {
	//app.logger.Info(models.GetCurrentFuncName())

	// Call the newTemplateData() helper to get a templateData struct containing
	// the 'default' data (which for now is just the current year), and add the
	// snippets slice to it.
	data := app.newTemplateData(r)
	data.CurrentYear = time.Now().Year()
	data.Date = time.Now().Format("02/01/2006")
	data.Message = "indexHandlerGet --- Public area! ---"
	//}
	//log.Println(r.Host)
	//app.logger.Info("Entrée dans IndexHandlerGet", "r.Host", r.Host)
	// Pass the data to the render() helper as normal.
	app.render(w, r, http.StatusOK, "index.gohtml", data)
}
func (app *application) LivreHandlerGet(w http.ResponseWriter, r *http.Request) {
	var form isbnForm

	//app.logger.Info("Entrée dans LivreHandlerGet")
	data := app.newTemplateData(r)
	data.Form = form
	data.Message = ""
	log.Println(models.GetCurrentFuncName())
	// Pass the data to the render() helper as normal.
	app.render(w, r, http.StatusOK, "livre.gohtml", data)
}

func (app *application) LivreHandlerPost(w http.ResponseWriter, r *http.Request) {
	var form isbnForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Isbn), "isbn", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Isbn, validator.IsbnRX), "isbn", "This field must contains only numérics characters")
	isvalid := form.Valid()
	if !isvalid {
		data := app.newTemplateData(r)
		data.Form = form
		data.Message = "Le code isbn doit contenir 13 chiffres"
		app.render(w, r, http.StatusUnprocessableEntity, "livre.gohtml", data)
		return
	}
	//app.logger.Info("Entrée dans LivreHandlerPost")
	data := app.newTemplateData(r)
	log.Println(models.GetCurrentFuncName())
	var livre models.Livre

	// Vérifier si le livre existe dans la base
	livre.Livre_Id = app.livres.LivreExist(form.Isbn)
	if livre.Livre_Id == 0 {
		// Acquisition du livre dan l'API
		book, ok := GetBook(form.Isbn)
		if !ok {
			log.Printf("Il n'a pas été possible d'obtenir le livre demande : %v\n", form.Isbn)
			http.Redirect(w, r, "/livre", http.StatusSeeOther)
			return
		}
		// transcrire apiBooks vers Livre
		livre, _ = TranscrireBookVersLivre(book)
		// On vérifie: L'Editeur existe-t-il dans la base ?
		livre.Editeur, ok = app.editeurs.EditeurExist(livre.Editeur.Nom) //book.Items[0].VolumeInfo.Publisher)
		if !ok {
			// L'editeur n'existe pas dans la base, on le crée
			//livre.Editeur.Nom = book.Items[0].VolumeInfo.Publisher
			livre.Editeur.Editeur_Id = app.editeurs.CreateEditeur(livre.Editeur)
		}
		livre.Livre_Id = app.livres.GetIdNewLivre()
		err = app.livres.CreateLivre(livre)
		if err != nil {
			if errors.Is(err, models.ErrDuplicateIsbn) {
				form.AddFieldError("isbn", "Isbn code or Idg code is already in use")
				data := app.newTemplateData(r)
				data.Form = form
				app.render(w, r, http.StatusUnprocessableEntity, "livre.gohtml", data)
			} else {
				app.serverError(w, r, err)
			}
			return
		}

		// Vérification de l'existence de ou des auteurs
		for z := 0; z < len(livre.Auteurs); z++ {
			// On vérifie: l'Auteur existe-t-il dans la base ?
			auteur, ok := app.auteurs.AuteurExist(livre.Auteurs[z].Nom) //book.Items[0].VolumeInfo.Authors[z])
			if !ok {
				// L'auteur n'existe pas dans la base
				auteur.Nom = livre.Auteurs[z].Nom
				// Création de l'auteur
				livre.Auteurs[z].Auteur_Id = app.auteurs.CreateAuteur(auteur)
				//créer une occurence dans la table livreauteur
				_ = app.auteurs.CreateLivreAuteur(livre.Livre_Id, livre.Auteurs[z].Auteur_Id)
			} else {
				livre.Auteurs[z].Nom = auteur.Nom
				livre.Auteurs[z].Auteur_Id = auteur.Auteur_Id
				//créer une occurence dans la table livreauteur
				_ = app.auteurs.CreateLivreAuteur(livre.Livre_Id, livre.Auteurs[z].Auteur_Id)
			}
		}
	} else {
		// Lire le livre dans la base
		livre = app.auteurs.GetLivreetEditeurAuteurs(form.Isbn)
	}
	data.Path = "/images/couverture/"
	data.Livre = livre
	// Pass the data to the render() helper as normal.
	app.render(w, r, http.StatusOK, "affichelivre.gohtml", data)
}
func (app *application) HomeHandlerGet(w http.ResponseWriter, r *http.Request) {
	//app.logger.Info("Entrée dans TestHandlerPost")
	log.Println(models.GetCurrentFuncName())
	data := app.newTemplateData(r)
	data.Username = app.username
	// Pass the data to the render() helper as normal.
	app.render(w, r, http.StatusOK, "home.gohtml", data)

}
func (app *application) LivresHandlerGet(w http.ResponseWriter, r *http.Request) {
	//app.logger.Info("Entrée dans LivresHandlerGet")
	log.Println(models.GetCurrentFuncName())
	data := app.newTemplateData(r)
	var livres []models.Livre
	livres = app.auteurs.GetLivresetEditeursAuteurs()
	data.Path = "/images/couverture/"
	data.Livres = livres
	data.Message = "Liste des Livres"
	// Pass the data to the render() helper as normal.
	app.render(w, r, http.StatusOK, "affichelivres.gohtml", data)

}
func (app *application) CouvertureLivreGet(w http.ResponseWriter, r *http.Request) {
	//app.logger.Info("Entrée dans CouvertureLivreGet")
	//log.Println(models.GetCurrentFuncName())

	nomImage, err := app.readName(r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	img, err := os.ReadFile(assert.Path + "data/img/" + nomImage)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(img)
}

// Change the signature of the home handler so it is defined as a method against
// *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//panic("oops! something went wrong") // Deliberate panic
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Call the newTemplateData() helper to get a templateData struct containing
	// the 'default' data (which for now is just the current year), and add the
	// snippets slice to it.
	data := app.newTemplateData(r)
	data.Snippets = snippets
	// Pass the data to the render() helper as normal.
	app.render(w, r, http.StatusOK, "home.gohtml", data)
}
func (app *application) index(w http.ResponseWriter, r *http.Request) {
	//panic("oops! something went wrong") // Deliberate panic
	//snippets, err := app.snippets.Latest()
	//if err != nil {
	//	app.serverError(w, r, err)
	//	return
	//}
	// Call the newTemplateData() helper to get a templateData struct containing
	// the 'default' data (which for now is just the current year), and add the
	// snippets slice to it.
	data := app.newTemplateData(r)
	//data.Snippets = snippets
	// Pass the data to the render() helper as normal.
	app.render(w, r, http.StatusOK, "index.gohtml", data)
}

// Change the signature of the snippetView handler so it is defined as a method
// against *application.
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFoundResponse(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, r, http.StatusOK, "view.gohtml", data)
}

// Change the signature of the snippetCreate handler so it is defined as a method
// against *application.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	// Initialize a new createSnippetForm instance and pass it to the template.
	// Notice how this is also a great opportunity to set any default or
	// 'initial' values for the form --- here we set the initial value for the
	// snippet expiry to 365 days.
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, r, http.StatusOK, "create.gohtml", data)
}

// Change the signature of the snippetCreatePost handler so it is defined as a method
// against *application.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")
	isvalid := form.Valid()
	if !isvalid {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.gohtml", data)
		return
	}
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Use the Put() method to add a string value ("Snippet successfully
	// created!") and the corresponding key ("flash") to the session data.
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// Update the handler so it displays the signup page.
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.gohtml", data)
}
func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	isvalid := form.Valid()
	if !isvalid {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.gohtml", data)
		return
	}
	// Tester si l'email existe déjà
	// Try to create a new user record in the database. If the email already
	// exists then add an error message to the form and re-display it.
	err = app.user.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.gohtml", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked.
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// Update the handler so it displays the login page.
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = models.UserLoginForm{}
	app.render(w, r, http.StatusOK, "login.gohtml", data)
}
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	// Decode the form data into the userLoginForm struct.// Decode the form data into the userLoginForm struct.
	var form models.UserLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Do some validation checks on the form. We check that both email and
	// password are provided, and also check the format of the email address as
	// a UX-nicety (in case the user makes a typo).
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.gohtml", data)
		return
	}
	// Check whether the credentials are valid. If they're not, add a generic
	// non-field error message and re-display the login page.

	user, err := app.user.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.gohtml", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	// Use the RenewToken() method on the current session to change the session
	// ID. It's good practice to generate a new session ID when the
	// authentication state or privilege levels changes for the user (e.g. login
	// and logout operations).
	// aller chercher le name avec l'Id
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Add the ID of the current user to the session, so that they are now
	// 'logged in'.
	app.sessionManager.Put(r.Context(), "authenticatedUserID", user.User_id)

	/* //app.logger.Info("Entrée dans MovieHandlerGet")
	data := app.newTemplateData(r)
	data.Message = ""
	log.Println(models.GetCurrentFuncName())
	// Pass the data to the render() helper as normal.
	app.render(w, r, http.StatusOK, "movie.gohtml", data)
	// #### Fin de Etablir la connection avec l'API #### */

	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) userLogoutGet(w http.ResponseWriter, r *http.Request) {
	// Use the RenewToken() method on the current session to change the session
	// ID again.
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	// Add a flash message to the session to confirm to the user that they've been
	// logged out.
	app.sessionManager.Put(r.Context(), "flash", "Vous êtes déconnecté(e)!")
	// Redirect the user to the application home page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
func (app *application) ListAuteursHandlerGet(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Entrée dans ListAuteursHandlerGet")
	log.Println(models.GetCurrentFuncName())
	data := app.newTemplateData(r)
	var auteurs []models.Auteur
	auteurs = app.auteurs.GetAuteurs()
	data.Path = "/images/couverture/"
	data.Auteurs = auteurs
	data.Message = "Liste des Auteurs"
	// Pass the data to the render() helper as normal.
	app.render(w, r, http.StatusOK, "listeauteurs.gohtml", data)
}
func (app *application) ListEditeursHandlerGet(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Entrée dans ListEditeursHandlerGet")
	log.Println(models.GetCurrentFuncName())
	data := app.newTemplateData(r)
	var editeurs []models.Editeur
	editeurs = app.editeurs.GetEditeurs()
	data.Path = "/images/couverture/"
	data.Editeurs = editeurs
	data.Message = "Liste des Editeurs"
	// Pass the data to the render() helper as normal.
	app.render(w, r, http.StatusOK, "listeediteurs.gohtml", data)
}

// Update the handler so it displays the signup page.
func (app *application) userModifGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	var user models.User
	var ok bool
	user, ok = app.user.GetUser(data.Username)
	if !ok {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	var ModifUser userModifForm
	ModifUser.Name = user.Name
	ModifUser.Email = user.Email
	data.Form = ModifUser

	app.render(w, r, http.StatusOK, "modifuser.gohtml", data)
}
func (app *application) userModifPost(w http.ResponseWriter, r *http.Request) {
	// Changement d'Email ou de mot de passe
	// Il faut ici vérifier ce que l'utilisateur a changé :
	// Vérifier que le mot de passe actuel est correct
	// si email changé ?
	//  - vérifier si l'Email entré est déjà dans la base
	//     si non -> Envoi d'un message pour validation (plus tard)
	// vérifier si un nouveau mot de passe a été entré et que les deux sont égaux,
	//  et qu'ils respectent les critères
	// Si oui prise en compte du nouveau mot de passe

	var form userModifForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password1), "password", "This field cannot be blank")
	form.CheckField(validator.EqualPwd(form.Password1, form.Password2), "pwd1&pwd2", "This fields must be equal or '' ")

	isvalid := form.Valid()
	if !isvalid {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "modifuser.gohtml", data)
		return
	}
	// Tester si l'email existe déjà
	// Try to create a new user record in the database. If the email already
	// exists then add an error message to the form and re-display it.
	err = app.user.Insert(form.Name, form.Email, form.Password1)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "modifuser.gohtml", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked.
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) indexHandlerNoMeth(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	log.Println("HTTP Error", http.StatusMethodNotAllowed)
	w.WriteHeader(http.StatusMethodNotAllowed)
	/*
		utils.Logger.Warn("indexHandlerNoMeth", slog.Int("req_id", middlewares.LogId), slog.String("req_url", r.URL.String()), slog.Int("http_status", http.StatusMethodNotAllowed))
		data := struct {
			Connect  bool
			Username string
			Date     string
			Message  string
		}{
			Date:    time.Now().Format("02/01/2006"),
			Message: fmt.Sprintf("HTTP Error %v", http.StatusMethodNotAllowed),
		}
		data.Username, data.Connect = utils.IsConnected(r)
		err := models.Tmpl["index"].ExecuteTemplate(w, "base", &data)
		if err != nil {
			log.Fatalln(err)
		}
	*/

	app.sessionManager.Put(r.Context(), "flash", "Methode non authorisée! veuillez vous connecter")

	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) indexHandlerOther(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	log.Println("HTTP Error", http.StatusNotFound)
	w.WriteHeader(http.StatusNotFound)
	/*
		utils.Logger.Warn("indexHandlerOther", slog.Int("req_id", middlewares.LogId), slog.String("req_url", r.URL.String()), slog.Int("http_status", http.StatusNotFound))
		data := struct {
			Connect  bool
			Username string
			Date     string
			Message  string
		}{
			Date:    time.Now().Format("02/01/2006"),
			Message: fmt.Sprintf("This adress is not valid : Error %v", http.StatusNotFound),
		}
		data.Username, data.Connect = utils.IsConnected(r)
		err := models.Tmpl["error404"].ExecuteTemplate(w, "base", &data)
		if err != nil {
			log.Fatalln(err)
		}
	*/
	data := app.newTemplateData(r)
	data.Message = "Adresse invalide! veuillez vous connecter"
	//app.sessionManager.Put(r.Context(), "flash", "Adresse invalide! veuillez vous connecter")
	app.render(w, r, http.StatusSeeOther, "error404.gohtml", data)
}

// ###########################################################################################
func (app *application) ConnectUserApiGet(w http.ResponseWriter, r *http.Request) {
	var tokenForm models.UserLoginForm

	data := app.newTemplateData(r)
	data.Form = tokenForm
	app.render(w, r, http.StatusOK, "signupapi.gohtml", data)
}

// ###########################################################################################
func (app *application) ConnectUserApiPost(w http.ResponseWriter, r *http.Request) {
	var form models.UserLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Do some validation checks on the form. We check that both email and
	// password are provided, and also check the format of the email address as
	// a UX-nicety (in case the user makes a typo).
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "loginapi.gohtml", data)
		return
	}
	// Il faut avoir l'Id utilisateur
	// Récupération du token envoyé par l'API, et de l'identifiant de l'utilisateur
	data := app.newTemplateData(r)
	user, ok := app.user.GetUser(data.Username)
	if !ok {
		data.Message = "Vous devez être connecté !"
	}
	// appel de l'API pour créer le compte s'il n'existe pas
	form.ID = strconv.FormatInt(int64(user.Id), 10)
	data.Form = form
	errCUAPI := app.CreateUserApi(form.Name, form.Email, form.Password, user.Id)
	if errCUAPI != nil {
		// une erreur s'est produite
		// Il faut traiter le type d'erreur (utilisateur existant...)
		fmt.Printf("Erreur : %v\n", errCUAPI)
	}

	log.Println(models.GetCurrentFuncName())
	app.render(w, r, http.StatusOK, "saisietokenapi.gohtml", data)
}

//###########################################################################################

func (app *application) ActiveUserPost(w http.ResponseWriter, r *http.Request) {
	var tokenForm models.UserLoginForm
	err := app.decodePostForm(r, &tokenForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	tokenForm.CheckField(validator.Matches(tokenForm.Token, validator.TokenRX), "token", "This field must contains 26 characters")
	if !tokenForm.Valid() {
		data := app.newTemplateData(r)
		data.Form = tokenForm
		app.render(w, r, http.StatusUnprocessableEntity, "saisietokenapi.gohtml", data)
		return
	}
	data := app.newTemplateData(r)
	// Récupération du token envoyé par l'API, et de l'identifiant de l'utilisateur
	user, ok := app.user.GetUser(data.Username)
	if !ok {
		data.Message = "Vous devez être connecté !"
	}
	tokenForm.Name = user.Name
	str_id := strconv.FormatInt(int64(user.Id), 10)
	tokenForm.ID = str_id
	log.Println(models.GetCurrentFuncName())
	// Ecrire le token dans la base de données
	var token models.AuthenticateUserApi
	token.ID = user.Id
	token.Token = tokenForm.Token
	token.Expiry = time.Now().Add(time.Hour * 24)
	errEJ := app.movies.EcrireJetonDansBase(token)
	if errEJ != nil {
		data.Message = "Impossible de sauvegarder le token !"
	}
	// renvoyer le token à l'API pour activer le compte
	var token2 models.AuthenticateUserApi
	id, _ := strconv.Atoi(tokenForm.ID)
	token2.ID = id
	token2.Token = tokenForm.Token
	cmovie, errAUM := app.ActivateUserApi(token2)
	if errAUM != nil {
		fmt.Printf("erreur = %v\n", errAUM)
	}
	if cmovie.Activated {
		data.Message = "Votre compte API greenlight est activé !"
		log.Printf("Votre compte API greenlight est activé !\n")
		http.Redirect(w, r, "/movies/authenticateuserapi", http.StatusOK)
		return
	} else {
		data.Message = " *** Votre compte API greenlight n'est pas activé ! ***"
	}

	app.render(w, r, http.StatusOK, "home.gohtml", data)

}

// ###########################################################################################
func (app *application) AuthenticateUserApiGet(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	data.Form = models.UserLoginForm{}
	// Récupération des données utilisateur (email et password) pour s'authentifier à l'API
	log.Println(models.GetCurrentFuncName())
	app.render(w, r, http.StatusOK, "loginapi.gohtml", data)

}

//###########################################################################################

func (app *application) AuthenticateUserApiPost(w http.ResponseWriter, r *http.Request) {

	var form models.UserLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Do some validation checks on the form. We check that both email and
	// password are provided, and also check the format of the email address as
	// a UX-nicety (in case the user makes a typo).
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	data := app.newTemplateData(r)
	data.Form = form
	if !form.Valid() {
		app.render(w, r, http.StatusUnprocessableEntity, "loginapi.gohtml", data)
		return
	}

	log.Println(models.GetCurrentFuncName())
	user, ok := app.user.GetUser(data.Username)
	if !ok {
		data.Message = "Vous devez être connecté !"
	}
	// Appeler la fonction d'authentification à l'API
	cmovie, errAUM := app.AuthenticateUserApi(form.Email, form.Password, user.Id)
	if errAUM != nil {
		fmt.Printf("erreur = %v\n", errAUM)
	}

	data.Message = "Vous êtes bien authentifié à l'API greenlight!"
	// Ecrire le token dans la table tokens
	var token models.AuthenticateUserApi
	token.Expiry = cmovie.Expiry
	token.ID = cmovie.User_id
	token.Token = cmovie.Token
	err = app.movies.EcrireJetonDansBase(token)
	if err != nil {
		data.Message = "Impossible d'écrire le jeton dans la base !"
		app.render(w, r, http.StatusOK, "home.gohtml", data)
	}

	app.render(w, r, http.StatusOK, "home.gohtml", data)

}

//###########################################################################################

func (app *application) MovieViewGet(w http.ResponseWriter, r *http.Request) {
	var form movieForm

	//app.logger.Info("Entrée dans MovieHandlerGet")
	data := app.newTemplateData(r)
	data.Form = form
	data.Message = ""
	log.Println(models.GetCurrentFuncName())
	// Pass the data to the render() helper as normal.
	app.render(w, r, http.StatusOK, "movie.gohtml", data)
}

func (app *application) MovieViewPost(w http.ResponseWriter, r *http.Request) {
	var form movieForm
	var err error
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}
	str_id := strconv.FormatInt(form.Id, 10)
	form.CheckField(validator.NotBlank(str_id), "id", "This field cannot be blank")
	form.CheckField(validator.Matches(str_id, validator.IdRX), "id", "This field must contains only numérics characters")
	data := app.newTemplateData(r)
	data.Form = form
	isvalid := form.Valid()
	if !isvalid {
		data.Message = "Le code id peut contenir 3 chiffres maximum"
		app.render(w, r, http.StatusUnprocessableEntity, "movie.gohtml", data)
	}
	//app.logger.Info("Entrée dans MovieHandlerPost")
	log.Println(models.GetCurrentFuncName())
	var movie models.Movie
	// Il faut vérifier qu'il existe un token valide
	user, ok := app.user.GetUser(data.Username)
	var token models.AuthenticateUserApi
	token, err = app.movies.LireJetonDansBase(user.Id)
	if err != nil {
		data.Message = "Il n'a pas été possible de lire le jeton API dans la base"
		app.render(w, r, http.StatusUnprocessableEntity, "movie.gohtml", data)
	}
	// vérifier si le token est valide
	if token.Expiry.Before(time.Now()) || (token.Token == "") {
		data.Message = "Le jeton est expiré"
		app.render(w, r, http.StatusUnprocessableEntity, "movie.gohtml", data)
	}
	// Vérifier si le livre existe dans la base
	//movie = app.movies.MovieExist(form.Id)
	//if movie.ID == 0 {
	// Acquisition du livre dans l'API
	movie, ok = app.GetMovie(str_id, token.Token)
	if !ok {
		log.Printf("Il n'a pas été possible d'obtenir le film : %v\n", form.Id)
		data.Form = form
		data.Flash = "Il n'a pas été possible d'obtenir le film !"
		app.render(w, r, http.StatusUnprocessableEntity, "movie.gohtml", data)
	}
	data.Form = movie
	app.render(w, r, http.StatusOK, "affichemovie.gohtml", data)
}

// #################################################################################
func (app *application) MoviesViewGet(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	log.Println(models.GetCurrentFuncName())
	var movies []models.Movie
	var movie movieForm
	// Il faut vérifier qu'il existe un token valide
	user, ok := app.user.GetUser(data.Username)
	var token models.AuthenticateUserApi
	var err error
	token, err = app.movies.LireJetonDansBase(user.Id)
	if err != nil {
		data.Message = "Il n'a pas été possible de lire le jeton API dans la base"
		app.render(w, r, http.StatusUnprocessableEntity, "home.gohtml", data)
	}
	// vérifier si le token est valide
	if token.Expiry.Before(time.Now()) || (token.Token == "") {
		data.Message = "Le jeton est expiré"
		app.render(w, r, http.StatusUnprocessableEntity, "home.gohtml", data)
	}
	//if movie.ID == 0 {
	// Acquisition du livre dans l'API
	movies, ok = app.GetMovies(token.Token)
	if !ok {
		log.Printf("Il n'a pas été possible d'obtenir les films demandés \n")
		data.Form = movie
		data.Flash = "Il n'a pas été possible d'obtenir les films demandés !"
		app.render(w, r, http.StatusUnprocessableEntity, "home.gohtml", data)
		data.Message = "Le jeton est expiré"
		app.render(w, r, http.StatusUnprocessableEntity, "home.gohtml", data)
	}

	data.Form = movies
	app.render(w, r, http.StatusOK, "affichemovies.gohtml", data)
}
